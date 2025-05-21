/*
Copyright

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"log"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/mod/semver"

	"github.com/helm/helm-mapkubeapis/pkg/mapping"
)

// KubeConfig are the Kubernetes configuration settings
type KubeConfig struct {
	Context string
	File    string
}

// MapOptions are the options for mapping deprecated APIs in a release
type MapOptions struct {
	DryRun           bool
	KubeConfig       KubeConfig
	MapFile          string
	ReleaseName      string
	ReleaseNamespace string
}

// UpgradeDescription is description of why release was upgraded
const UpgradeDescription = "Kubernetes deprecated API upgrade - DO NOT rollback from this version"

// ReplaceManifestUnSupportedAPIs returns a release manifest with deprecated or removed
// Kubernetes APIs updated to supported APIs
func ReplaceManifestUnSupportedAPIs(origManifest, mapFile string, kubeConfig KubeConfig, additionalMappings ...*mapping.Mapping) (string, error) {
	var modifiedManifest = origManifest
	var err error
	var mapMetadata *mapping.Metadata

	// Load the mapping data
	if mapMetadata, err = mapping.LoadMapfile(mapFile); err != nil {
		return "", errors.Wrapf(err, "Failed to load mapping file: %s", mapFile)
	}

	mapMetadata.Mappings = append(mapMetadata.Mappings, additionalMappings...)

	// get the Kubernetes server version
	kubeVersionStr, err := getKubernetesServerVersion(kubeConfig)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(kubeVersionStr) {
		return "", errors.Errorf("Failed to get Kubernetes server version")
	}

	// Check for deprecated or removed APIs and map accordingly to supported versions
	modifiedManifest, err = ReplaceManifestData(mapMetadata, modifiedManifest, kubeVersionStr)
	if err != nil {
		return "", err
	}

	return modifiedManifest, nil
}

// ReplaceManifestData scans the release manifest string for deprecated APIs in a given Kubernetes version and replaces
// their groups and versions if there is a successor, or fully removes the manifest for that specific resource if no
// successors exist (such as the PodSecurityPolicy API).
func ReplaceManifestData(mapMetadata *mapping.Metadata, modifiedManifest string, kubeVersionStr string) (string, error) {
	for _, mapping := range mapMetadata.Mappings {
		deprecatedAPI := mapping.DeprecatedAPI
		supportedAPI := mapping.NewAPI
		var apiVersionStr string
		if mapping.DeprecatedInVersion != "" {
			apiVersionStr = mapping.DeprecatedInVersion
		} else {
			apiVersionStr = mapping.RemovedInVersion
		}

		if !semver.IsValid(apiVersionStr) {
			return "", errors.Errorf("Failed to get the deprecated or removed Kubernetes version for API: %s", strings.ReplaceAll(deprecatedAPI, "\n", " "))
		}

		if count := strings.Count(modifiedManifest, deprecatedAPI); count > 0 {
			if semver.Compare(apiVersionStr, kubeVersionStr) > 0 {
				log.Printf("The following API:\n\"%s\" does not require mapping as the "+
					"API is not deprecated or removed in Kubernetes \"%s\"\n", deprecatedAPI, kubeVersionStr)
				// skip to next mapping
				continue
			}
			if supportedAPI == "" {
				log.Printf("Found %d instances of deprecated or removed Kubernetes API:\n\"%s\"\nNo supported API equivalent\n", count, deprecatedAPI)
				modifiedManifest = removeDeprecatedAPIWithoutSuccessor(count, deprecatedAPI, modifiedManifest)
			} else {
				log.Printf("Found %d instances of deprecated or removed Kubernetes API:\n\"%s\"\nSupported API equivalent:\n\"%s\"\n", count, deprecatedAPI, supportedAPI)
				modifiedManifest = strings.ReplaceAll(modifiedManifest, deprecatedAPI, supportedAPI)
			}
		}
	}
	return modifiedManifest, nil
}

// removeDeprecatedAPIWithoutSuccessor removes a deprecated API that has no successor specified in the mapping file.
func removeDeprecatedAPIWithoutSuccessor(count int, deprecatedAPI string, modifiedManifest string) string {
	for repl := 0; repl < count; repl++ {
		// find the position where the API header is
		apiIndex := strings.Index(modifiedManifest, deprecatedAPI)

		// find the next separator index
		separatorIndex := strings.Index(modifiedManifest[apiIndex:], "---\n")

		// find the previous separator index
		previousSeparatorIndex := strings.LastIndex(modifiedManifest[:apiIndex], "---\n")

		/*
		 * if no previous separator index was found, it means the resource is at the beginning and not
		 * prefixed by ---
		 */
		if previousSeparatorIndex == -1 {
			previousSeparatorIndex = 0
		}

		if separatorIndex == -1 { // this means we reached the end of input
			modifiedManifest = modifiedManifest[:previousSeparatorIndex]
		} else {
			modifiedManifest = modifiedManifest[:previousSeparatorIndex] + modifiedManifest[separatorIndex+apiIndex:]
		}
	}

	modifiedManifest = strings.Trim(modifiedManifest, "\n")
	return modifiedManifest
}

func getKubernetesServerVersion(kubeConfig KubeConfig) (string, error) {
	clientSet := GetClientSetWithKubeConfig(kubeConfig.File, kubeConfig.Context)
	if clientSet == nil {
		return "", errors.Errorf("kubernetes cluster unreachable")
	}
	kubeVersion, err := clientSet.ServerVersion()
	if err != nil {
		return "", errors.Wrap(err, "kubernetes cluster unreachable")
	}
	return kubeVersion.GitVersion, nil
}

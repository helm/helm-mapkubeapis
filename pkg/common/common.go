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

	utils "github.com/maorfr/helm-plugin-utils/pkg"
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
func ReplaceManifestUnSupportedAPIs(origManifest, mapFile string, kubeConfig KubeConfig) (string, error) {
	var err error
	var mapMetadata *mapping.Metadata

	// Load the mapping data
	if mapMetadata, err = mapping.LoadMapfile(mapFile); err != nil {
		return "", errors.Wrapf(err, "Failed to load mapping file: %s", mapFile)
	}

	// get the Kubernetes server version
	kubeVersionStr, err := getKubernetesServerVersion(kubeConfig)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(kubeVersionStr) {
		return "", errors.Errorf("Failed to get Kubernetes server version")
	}

	updatedDocuments := []string{}
	for _, originalDocument := range SeparateDocuments(origManifest) {
		newDocument := originalDocument
		// Check for deprecated or removed APIs and map accordingly to supported versions
		for _, mapping := range mapMetadata.Mappings {
			var apiVersionStr string
			if mapping.DeprecatedInVersion != "" {
				apiVersionStr = mapping.DeprecatedInVersion
			} else {
				apiVersionStr = mapping.RemovedInVersion
			}
			if !semver.IsValid(apiVersionStr) {
				return "", errors.Errorf("Failed to get the deprecated or removed Kubernetes version for API: %s %s", mapping.DeprecatedAPI.APIVersion, mapping.DeprecatedAPI.Kind)
			}

			if semver.Compare(apiVersionStr, kubeVersionStr) > 0 {
				log.Printf("The resource %s/%s does not require mapping as the "+
					"API is not deprecated or removed in Kubernetes '%s'\n",
					mapping.DeprecatedAPI.APIVersion,
					mapping.DeprecatedAPI.Kind,
					kubeVersionStr,
				)
			} else {
				newDocument, err = CheckForOldAPI(originalDocument, mapping)
				if err != nil {
					log.Printf("could not update old api: %s\n", err)
					newDocument = originalDocument
				}
			}
			updatedDocuments = append(updatedDocuments, newDocument)
		}
	}

	modifiedManifest := strings.Join(updatedDocuments, "\n")
	return modifiedManifest, nil
}

func getKubernetesServerVersion(kubeConfig KubeConfig) (string, error) {
	clientSet := utils.GetClientSetWithKubeConfig(kubeConfig.File, kubeConfig.Context)
	if clientSet == nil {
		return "", errors.Errorf("kubernetes cluster unreachable")
	}
	kubeVersion, err := clientSet.ServerVersion()
	if err != nil {
		return "", errors.Wrap(err, "kubernetes cluster unreachable")
	}
	return kubeVersion.GitVersion, nil
}

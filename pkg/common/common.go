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
	"fmt"
	"github.com/helm/helm-mapkubeapis/pkg/mapping"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"strings"
)

// KubeConfig are the Kubernetes configuration settings
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
	var modifiedManifest string
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

	// Check for deprecated or removed APIs and map accordingly to supported versions
	modifiedManifest, err = ReplaceManifestData(mapMetadata, origManifest, kubeVersionStr)
	if err != nil {
		return "", err
	}

	return modifiedManifest, nil
}

type genericYAML map[string]any

func parseYAML(s string) []genericYAML {
	decoder := yaml.NewDecoder(strings.NewReader(s))
	var docs []genericYAML
	for {
		var y genericYAML
		err := decoder.Decode(&y)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("Failed to decode YAML: %v\n", err)
			continue
		}
		docs = append(docs, y)
	}
	return docs
}

// ReplaceManifestData scans the release manifest string for deprecated APIs in a given Kubernetes version and replaces
// their groups and versions if there is a successor, or fully removes the manifest for that specific resource if no
// successors exist (such as the PodSecurityPolicy API).
func ReplaceManifestData(mapMetadata *mapping.Metadata, origManifest string, kubeVersionStr string) (string, error) {
	yamlDocs := parseYAML(origManifest)
	for _, m := range mapMetadata.Mappings {
		deprecatedAPI := m.DeprecatedAPI
		supportedAPI := m.NewAPI
		deprecatedAPIVersion, deprecatedAPIKind := mapping.ParseAPIString(deprecatedAPI)
		supportedAPIVersion, supportedAPIKind := mapping.ParseAPIString(supportedAPI)

		//fmt.Printf("deprecatedAPI: %s\nsupportedAPI: %s\n", deprecatedAPI, supportedAPI)
		//fmt.Printf("deprecatedAPIVersion: %s\ndeprecatedAPIKind: %s\n", deprecatedAPIVersion, deprecatedAPIKind)
		//fmt.Printf("supportedAPIVersion: %s\nsupportedAPIKind: %s\n", supportedAPIVersion, supportedAPIKind)

		var apiVersionStr = m.RemovedInVersion
		if m.DeprecatedInVersion != "" {
			apiVersionStr = m.DeprecatedInVersion
		}

		if !semver.IsValid(apiVersionStr) {
			return "", errors.Errorf("Failed to get the deprecated or removed Kubernetes version for API: %s", strings.ReplaceAll(deprecatedAPI, "\n", " "))
		}

		var count = 0

	docLoop:
		for idx, doc := range yamlDocs {
			version, _ := doc["apiVersion"].(string)
			kind, _ := doc["kind"].(string)
			if version == deprecatedAPIVersion && kind == deprecatedAPIKind {
				fmt.Printf("Found deprecated or removed Kubernetes version for API: %s %s\n", deprecatedAPIVersion, deprecatedAPIKind)
				fmt.Println("original: ", doc)
				if semver.Compare(apiVersionStr, kubeVersionStr) > 0 {
					log.Printf("The following API:\n\"%s\" does not require mapping as the "+
						"API is not deprecated or removed in Kubernetes \"%s\"\n", deprecatedAPI, kubeVersionStr)
					// skip to next mapping
					break docLoop
				}
				count++
				if supportedAPI != "" {
					doc["apiVersion"] = supportedAPIVersion
					doc["kind"] = supportedAPIKind
					fmt.Println("modified: ", doc)
				} else {
					yamlDocs = append(yamlDocs[:idx], yamlDocs[idx+1:]...)
					fmt.Println("deleted doc without replacement")
				}
			}
		}
		if count > 0 {
			if supportedAPI == "" {
				log.Printf("Found %d instances of deprecated or removed Kubernetes API:\n\"%s\"\nNo supported API equivalent\n", count, deprecatedAPI)
			} else {
				log.Printf("Found %d instances of deprecated or removed Kubernetes API:\n\"%s\"\nSupported API equivalent:\n\"%s\"\n", count, deprecatedAPI, supportedAPI)
			}
		}
	}

	var sb strings.Builder
	encoder := yaml.NewEncoder(&sb)
	encoder.SetIndent(2) // match test cases
	for _, doc := range yamlDocs {
		if err := encoder.Encode(doc); err != nil {
			log.Printf("Failed to encode document: %v\n", err)
		}
	}
	if err := encoder.Close(); err != nil {
		log.Printf("Failed to close encoder: %v\n", err)
	}
	return "---\n" + sb.String(), nil // always start with a document separator
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

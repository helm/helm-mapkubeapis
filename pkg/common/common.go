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

	"github.com/hickeyma/helm-mapkubeapis/pkg/mapping"
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
	StorageType      string
	TillerOutCluster bool
}

// UpgradeDescription is description of why release was upgraded
const UpgradeDescription = "Kubernetes deprecated API upgrade - DO NOT rollback from this version"

// ReplaceManifestUnSupportedAPIs returns a release manifest with deprecated or removed
// Kubernetes APIs updated to supported APIs
func ReplaceManifestUnSupportedAPIs(origManifest, mapFile, kubeVersion string) (string, error) {
	var modifiedManifest string
	var err error
	var mapMetadata *mapping.Metadata

	// Load the mapping data
	if mapMetadata, err = mapping.LoadMapfile(mapFile); err != nil {
		return "", errors.Wrapf(err, "Failed to load mapping file: %s", mapFile)
	}

	// Check for deprecated or removed APIs and map accordingly to supported versions
	for _, mapping := range mapMetadata.Mappings {
		deprecatedAPI := mapping.DeprecatedAPI
		supportedAPI := mapping.NewAPI
		if kubeVersion != "" {
			var version string
			if mapping.DeprecatedInVersion != "" {
				version = mapping.DeprecatedInVersion
			} else {
				version = mapping.RemovedInVersion
			}
			if version > kubeVersion {
				log.Printf("The following API does not required mapping now as it is not deprecated till Kubernetes '%s':\n\"%s\"\n", mapping.DeprecatedInVersion,
					deprecatedAPI)
				continue
			}
		}
		var modManifestForAPI string
		if len(modifiedManifest) <= 0 {
			modManifestForAPI = strings.ReplaceAll(origManifest, deprecatedAPI, supportedAPI)
			if modManifestForAPI != origManifest {
				log.Printf("Found deprecated or removed Kubernetes API:\n\"%s\"\nSupported API equivalent:\n\"%s\"\n", deprecatedAPI, supportedAPI)
			}

		} else {
			modManifestForAPI = strings.ReplaceAll(modifiedManifest, deprecatedAPI, supportedAPI)
			if modManifestForAPI != modifiedManifest {
				log.Printf("Found deprecated or removed Kubernetes API:\n\"%s\"\nSupported API equivalent:\n\"%s\"\n", deprecatedAPI, supportedAPI)
			}
		}
		modifiedManifest = modManifestForAPI
	}

	return modifiedManifest, nil
}

/*
Copyright The Helm Authors.

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

package v3

import (
	"fmt"
	"log"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"

	common "github.com/hickeyma/helm-mapkubeapis/pkg/common"
)

// MapOptions are the options for mapping deprecated apis in a release
type MapOptions struct {
	DryRun           bool
	KubeConfig       common.KubeConfig
	ReleaseName      string
	ReleaseNamespace string
}

var mappedAPIs = map[string]string{
	"apiVersion: extensions/v1beta1\nkind: NetworkPolicy":     "apiVersion: networking.k8s.io/v1\nkind: NetworkPolicy",
	"apiVersion: extensions/v1beta1\nkind: PodSecurityPolicy": "apiVersion:  policy/v1beta1\nkind: PodSecurityPolicy",
	"apiVersion: extensions/v1beta1\nkind: DaemonSet":         "apiVersion: apps/v1\nkind: DaemonSet",
	"apiVersion: apps/v1beta2\nkind: DaemonSet":               "apiVersion: apps/v1\nkind: DaemonSet",
	"apiVersion: extensions/v1beta1\nkind: Deployment":        "apiVersion: apps/v1\nkind: Deployment",
	"apiVersion: apps/v1beta1\nkind: Deployment":              "apiVersion: apps/v1\nkind: Deployment",
	"apiVersion: apps/v1beta12\nkind: Deployment":             "apiVersion: apps/v1\nkind: Deployment",
	"apiVersion: apps/v1beta1\nkind: StatefulSet":             "apiVersion: apps/v1\nkind: StatefulSet",
	"apiVersion: apps/v1beta12\nkind: StatefulSet":            "apiVersion: apps/v1\nkind: StatefulSet",
	"apiVersion: extensions/v1beta1\nkind: ReplicaSet":        "apiVersion: apps/v1\nkind: ReplicaSet",
	"apiVersion: apps/v1beta1\nkind: ReplicaSet":              "apiVersion: apps/v1\nkind: ReplicaSet",
	"apiVersion: apps/v1beta12\nkind: ReplicaSet":             "apiVersion: apps/v1\nkind: ReplicaSet",
	"apiVersion: extensions/v1beta1\nkind: Ingress":           "apiVersion: networking.k8s.io/v1beta1\nkind: Ingress"}

// MapReleaseWithDeprecatedAPIs checks the latest release version for any deprecated APIs in its metadata
// If it finds any, it will create a new release version with the APIs mapped to the supported versions
func MapReleaseWithDeprecatedAPIs(mapOptions MapOptions) error {
	cfg, err := GetActionConfig(mapOptions.ReleaseNamespace, mapOptions.KubeConfig)
	if err != nil {
		return fmt.Errorf("Failed to get Helm action configuration due to the following error: %s", err)
	}

	var releaseName = mapOptions.ReleaseName
	log.Printf("Get release '%s' latest version.\n", releaseName)
	releaseToMap, err := getLatestRelease(releaseName, cfg)
	if err != nil {
		return fmt.Errorf("Failed to get release '%s' latest version due to the following error: %s", mapOptions.ReleaseName, err)
	}

	log.Printf("Check release '%s' for deprecated APIs.\n", releaseName)
	var origManifest = releaseToMap.Manifest
	modifiedManifest := replaceManifestDeprecatedAPIs(origManifest)
	if modifiedManifest == origManifest {
		log.Printf("Release '%s' has no deprecated APIs.\n", releaseName)
		return nil
	}

	log.Printf("Deprecated APIs exist, updating release: %s.\n", releaseName)
	if !mapOptions.DryRun {
		if err := updateRelease(releaseToMap, modifiedManifest, cfg); err != nil {
			return fmt.Errorf("Failed to update release '%s' due to the following error: %s", releaseName, err)
		}
		log.Printf("Release '%s' with deprecated APIs updated successfully to new version.\n", releaseName)
	}

	return nil
}

func replaceManifestDeprecatedAPIs(origManifest string) string {
	var modifiedManifest string

	// Check for deprecated APIs and map accordingly to supported versions
	for deprecatedAPI, supportedAPI := range mappedAPIs {
		if len(modifiedManifest) <= 0 {
			modifiedManifest = strings.ReplaceAll(origManifest, deprecatedAPI, supportedAPI)
		} else {
			modifiedManifest = strings.ReplaceAll(modifiedManifest, deprecatedAPI, supportedAPI)
		}
	}

	return modifiedManifest
}

func updateRelease(origRelease *release.Release, modifiedManifest string, cfg *action.Configuration) error {
	// Update current release version to be superseded
	origRelease.Info.Status = release.StatusSuperseded
	if err := cfg.Releases.Update(origRelease); err != nil {
		return fmt.Errorf("failed to update current release version: %s", err)
	}

	// Using a shallow copy of  current release version to update the object with the modification
	// and then store this new version
	var newRelease = origRelease
	newRelease.Manifest = modifiedManifest
	newRelease.Info.Description = "Kubernetes deprecated API upgrade - DO NOT rollback from this version"
	newRelease.Info.LastDeployed = cfg.Now()
	newRelease.Version = origRelease.Version + 1
	newRelease.Info.Status = release.StatusDeployed
	if err := cfg.Releases.Create(newRelease); err != nil {
		return fmt.Errorf("failed to create new release version: %s", err)
	}
	return nil
}

func getLatestRelease(releaseName string, cfg *action.Configuration) (*release.Release, error) {
	return cfg.Releases.Last(releaseName)
}

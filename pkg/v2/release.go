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

package v2

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"

	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/storage"
	"k8s.io/helm/pkg/timeconv"

	common "github.com/helm/helm-mapkubeapis/pkg/common"
)

// MapReleaseWithUnSupportedAPIs checks the latest release version for any deprecated or removed APIs in its metadata
// If it finds any, it will create a new release version with the APIs mapped to the supported versions
func MapReleaseWithUnSupportedAPIs(mapOptions common.MapOptions) error {
	var releaseName = mapOptions.ReleaseName
	log.Printf("Get release '%s' latest version.\n", releaseName)
	storageDriver, err := GetStorageDriver(mapOptions)
	if err != nil {
		return errors.Wrapf(err, "Failed to get release '%s' latest version", mapOptions.ReleaseName)
	}
	releaseToMap, err := getLatestRelease(releaseName, storageDriver)
	if err != nil {
		return errors.Wrapf(err, "Failed to get release '%s' latest version", mapOptions.ReleaseName)
	}

	log.Printf("Check release '%s' for deprecated or removed APIs...\n", releaseName)
	var origManifest = releaseToMap.Manifest
	modifiedManifest, err := common.ReplaceManifestUnSupportedAPIs(origManifest, mapOptions.MapFile, mapOptions.KubeConfig)
	if err != nil {
		return err
	}
	log.Printf("Finished checking release '%s' for deprecated or removed APIs.\n", releaseName)
	if modifiedManifest == origManifest {
		log.Printf("Release '%s' has no deprecated or removed APIs.\n", releaseName)
		return nil
	}

	if mapOptions.DryRun {
		log.Printf("Deprecated or removed APIs exist, no changes will be made to release: %s.\n", releaseName)
	} else {
		log.Printf("Deprecated or removed APIs exist, updating release: %s.\n", releaseName)
		if err := updateRelease(releaseToMap, modifiedManifest, storageDriver); err != nil {
			return errors.Wrapf(err, "Failed to update release '%s'", releaseName)
		}
		log.Printf("Release '%s' with deprecated or removed APIs updated successfully to new version.\n", releaseName)
	}

	return nil
}

func getLatestRelease(releaseName string, storageDriver *storage.Storage) (*release.Release, error) {
	return storageDriver.Last(releaseName)
}

func updateRelease(origRelease *release.Release, modifiedManifest string, storageDriver *storage.Storage) error {
	// Update current release version to be superseded
	log.Printf("Set status of release version '%s' to 'superseded'.\n", getReleaseVersionName(origRelease))
	origRelease.Info.Status.Code = release.Status_SUPERSEDED
	if err := storageDriver.Update(origRelease); err != nil {
		return errors.Wrapf(err, "failed to update release version '%s'", getReleaseVersionName(origRelease))
	}
	log.Printf("Release version '%s' updated successfully.\n", getReleaseVersionName(origRelease))

	// Using a shallow copy of  current release version to update the object with the modification
	// and then store this new version
	var newRelease = origRelease
	newRelease.Manifest = modifiedManifest
	newRelease.Info.Description = common.UpgradeDescription
	newRelease.Info.LastDeployed = timeconv.Timestamp(time.Now())
	newRelease.Version = origRelease.Version + 1
	newRelease.Info.Status.Code = release.Status_DEPLOYED
	log.Printf("Add release version '%s' with updated supported APIs.\n", getReleaseVersionName(origRelease))
	if err := storageDriver.Create(newRelease); err != nil {
		return errors.Wrapf(err, "failed to create new release version '%s'", getReleaseVersionName(origRelease))
	}
	log.Printf("Release version '%s' added successfully.\n", getReleaseVersionName(origRelease))
	return nil
}

func getReleaseVersionName(rel *release.Release) string {
	return fmt.Sprintf("%s.v%d", rel.Name, rel.Version)
}

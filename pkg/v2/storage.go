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
	"os"

	utils "github.com/maorfr/helm-plugin-utils/pkg"
	"github.com/pkg/errors"

	"k8s.io/helm/pkg/storage"
	"k8s.io/helm/pkg/storage/driver"

	common "github.com/hickeyma/helm-mapkubeapis/pkg/common"
)

// GetStorageDriver return handle to Helm v2 backend storage driver
func GetStorageDriver(mapOptions common.MapOptions) (*storage.Storage, error) {
	clientSet := utils.GetClientSetWithKubeConfig(mapOptions.KubeConfig.File, mapOptions.KubeConfig.Context)
	if clientSet == nil {
		return nil, errors.Errorf("kubernetes cluster unreachable")
	}
	namespace := mapOptions.ReleaseNamespace
	storageType := getStorageType(mapOptions)

	switch storageType {
	case "configmap", "configmaps", "":
		cfgMaps := driver.NewConfigMaps(clientSet.CoreV1().ConfigMaps(namespace))
		cfgMaps.Log = newLogger("storage/driver").Printf
		return storage.Init(cfgMaps), nil
	case "secret", "secrets":
		secrets := driver.NewSecrets(clientSet.CoreV1().Secrets(namespace))
		secrets.Log = newLogger("storage/driver").Printf
		return storage.Init(secrets), nil
	default:
		// Not sure what to do here.
		panic("Unknown storage driver")
	}
}

func getStorageType(mapOptions common.MapOptions) string {
	var storage string
	if !mapOptions.TillerOutCluster {
		storage = utils.GetTillerStorageWithKubeConfig(mapOptions.ReleaseNamespace,
			mapOptions.KubeConfig.File, mapOptions.KubeConfig.Context)
	} else {
		storage = mapOptions.StorageType
	}
	return storage
}

func newLogger(prefix string) *log.Logger {
	if len(prefix) > 0 {
		prefix = fmt.Sprintf("[%s] ", prefix)
	}
	return log.New(os.Stderr, prefix, log.Flags())
}

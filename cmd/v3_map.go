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

package cmd

import (
	"errors"
	"io"
	"log"

	"github.com/spf13/cobra"

	"helm-mapkubeapis/pkg/common"
	v3 "helm-mapkubeapis/pkg/v3"
)

// MapOptions sontains the options for Map operation
type MapOptions struct {
	DryRun           bool
	ReleaseName      string
	ReleaseNamespace string
}

func newV3MapCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "v3map [flags] RELEASE RELEASE_NAMESPACE",
		Short: "map v3 release deprecated Kubernetes APIs in-place",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("name of release to be mapped and the release namespace need to be passed")
			}
			return nil
		},

		RunE: runMap,
	}

	flags := cmd.Flags()
	settings.AddFlags(flags)

	return cmd

}

func runMap(cmd *cobra.Command, args []string) error {
	releaseName := args[0]
	releaseNamespace := args[1]
	mapOptions := MapOptions{
		DryRun:           settings.DryRun,
		ReleaseName:      releaseName,
		ReleaseNamespace: releaseNamespace,
	}
	kubeConfig := common.KubeConfig{
		Context: settings.KubeContext,
		File:    settings.KubeConfigFile,
	}

	return V3Map(mapOptions, kubeConfig)
}

// V3Map checks for Kubernetes deprectaed APIs in the manifest of the last deployed release version
// and maps those deprecated APIs to the supported versions. It then adds a new release version with
// the updated APIs and supersedes the version with the deprecated APIs.
func V3Map(mapOptions MapOptions, kubeConfig common.KubeConfig) error {
	if mapOptions.DryRun {
		log.Println("NOTE: This is in dry-run mode, the following actions will not be executed.")
		log.Println("Run without --dry-run to take the actions described below:")
		log.Println()
	}

	log.Printf("Release '%s' will be checked for deprecated APIs and will be updated if necessary to supported API versions.\n", mapOptions.ReleaseName)

	v3MapOptions := v3.MapOptions{
		DryRun:           mapOptions.DryRun,
		KubeConfig:       kubeConfig,
		ReleaseName:      mapOptions.ReleaseName,
		ReleaseNamespace: mapOptions.ReleaseNamespace,
	}

	if err := v3.MapReleaseWithDeprecatedAPIs(v3MapOptions); err != nil {
		return err
	}

	log.Printf("Map of release '%s' deprecated APIs, completed successfully.\n", mapOptions.ReleaseName)

	return nil
}

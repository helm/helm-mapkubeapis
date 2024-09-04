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

package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/helm/helm-mapkubeapis/pkg/common"
	v3 "github.com/helm/helm-mapkubeapis/pkg/v3"
)

// MapOptions contains the options for Map operation
type MapOptions struct {
	DryRun           bool
	MapFile          string
	ReleaseName      string
	ReleaseNamespace string
}

var (
	settings *EnvSettings
)

func newMapCmd(_ io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "mapkubeapis [flags] RELEASE",
		Short:        "Map release deprecated or removed Kubernetes APIs in-place",
		Long:         "Map release deprecated or removed Kubernetes APIs in-place",
		SilenceUsage: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				err := cmd.Help()
				if err != nil {
					return err
				}
				os.Exit(1)
			} else if len(args) > 1 {
				return errors.New("only one release name may be passed at a time")
			}
			return nil
		},

		RunE: runMap,
	}

	flags := cmd.PersistentFlags()
	flags.ParseErrorsWhitelist.UnknownFlags = true

	settings = new(EnvSettings)

	// Get the default mapping file
	if ctx := os.Getenv("HELM_PLUGIN_DIR"); ctx != "" {
		settings.MapFile = filepath.Join(ctx, "config", "Map.yaml")
	} else {
		settings.MapFile = filepath.Join("config", "Map.yaml")
	}

	// When run with the Helm plugin framework, Helm plugins are not passed the
	// plugin flags that correspond to Helm global flags e.g. helm mapkubeapis v3map --kube-context ...
	// The flag values are set to corresponding environment variables instead.
	// The flags are passed as expected when run directly using the binary.
	// The below allows to use Helm's --kube-context global flag.
	if ctx := os.Getenv("HELM_KUBECONTEXT"); ctx != "" {
		settings.KubeContext = ctx
	}

	// Note that the plugin's --kubeconfig flag is set by the Helm plugin framework to
	// the KUBECONFIG environment variable instead of being passed into the plugin.

	settings.AddFlags(flags)

	return cmd
}

func runMap(_ *cobra.Command, args []string) error {
	releaseName := args[0]
	mapOptions := MapOptions{
		DryRun:           settings.DryRun,
		MapFile:          settings.MapFile,
		ReleaseName:      releaseName,
		ReleaseNamespace: settings.Namespace,
	}
	kubeConfig := common.KubeConfig{
		Context: settings.KubeContext,
		File:    settings.KubeConfigFile,
	}

	return Map(mapOptions, kubeConfig)
}

// Map checks for Kubernetes deprecated or removed APIs in the manifest of the last deployed release version
// and maps those API versions to supported versions. It then adds a new release version with
// the updated APIs and supersedes the version with the unsupported APIs.
func Map(mapOptions MapOptions, kubeConfig common.KubeConfig) error {
	if mapOptions.DryRun {
		log.Println("NOTE: This is in dry-run mode, the following actions will not be executed.")
		log.Println("Run without --dry-run to take the actions described below:")
		log.Println()
	}

	log.Printf("Release '%s' will be checked for deprecated or removed Kubernetes APIs and will be updated if necessary to supported API versions.\n", mapOptions.ReleaseName)

	options := common.MapOptions{
		DryRun:           mapOptions.DryRun,
		KubeConfig:       kubeConfig,
		MapFile:          mapOptions.MapFile,
		ReleaseName:      mapOptions.ReleaseName,
		ReleaseNamespace: mapOptions.ReleaseNamespace,
	}

	if err := v3.MapReleaseWithUnSupportedAPIs(options); err != nil {
		return err
	}

	log.Printf("Map of release '%s' deprecated or removed APIs to supported versions, completed successfully.\n", mapOptions.ReleaseName)

	return nil
}

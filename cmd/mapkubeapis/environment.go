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
	"github.com/hickeyma/helm-mapkubeapis/pkg/common"
	"github.com/spf13/pflag"
)

// EnvSettings defined settings
type EnvSettings struct {
	DryRun           bool
	KubeConfigFile   string
	KubeContext      string
	MapFile          string
	Namespace        string
	RunV2            bool
	StorageType      string
	TillerOutCluster bool
}

// New returns default env settings
func New() *EnvSettings {
	envSettings := EnvSettings{}
	return &envSettings
}

// AddBaseFlags binds base flags to the given flagset.
func (s *EnvSettings) AddBaseFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.DryRun, "dry-run", false, "simulate a command")
}

// AddFlags binds flags to the given flagset.
func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	s.AddBaseFlags(fs)
	fs.StringVar(&s.KubeConfigFile, "kubeconfig", "", "path to the kubeconfig file")
	fs.StringVar(&s.KubeContext, "kube-context", s.KubeContext, "name of the kubeconfig context to use")
	fs.StringVar(&s.MapFile, "mapfile", common.DefaultMappingFile, "path to the API mapping file")
	fs.StringVar(&s.Namespace, "namespace", s.Namespace, "namespace scope of the release. For Helm v2, this is the Tiller namespace e.g. kube-system")
	fs.BoolVar(&s.RunV2, "v2", false, "run for Helm v2 release. The default is Helm v3.")
	fs.BoolVar(&s.TillerOutCluster, "tiller-out-cluster", false, "for Helm v2 only - when Tiller is not running in the cluster e.g. Tillerless")
	fs.StringVarP(&s.StorageType, "release-storage", "s", "secrets", "for Helm v2 only - release storage type/object. It can be 'secrets' or 'configmaps'. This is only used with the 'tiller-out-cluster' flag")
}

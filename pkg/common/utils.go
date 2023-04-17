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

/*
This functionality is copied from
https://github.com/maorfr/helm-plugin-utils/blob/master/pkg/utils.go.
The reason for duplicating it is because it is no longer possible to keep maintaining the
source repo in sync with Kubernetes Go client changes as required by helm-mapkubeapis.
It has been copied in co-operation with author maorfr.
*/

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClientSetWithKubeConfig returns a kubernetes ClientSet
func GetClientSetWithKubeConfig(kubeConfigFile, context string) *kubernetes.Clientset {
	var kubeConfigFiles []string
	if kubeConfigFile != "" {
		kubeConfigFiles = append(kubeConfigFiles, kubeConfigFile)
	} else if kubeConfigPath := os.Getenv("KUBECONFIG"); kubeConfigPath != "" {
		// The KUBECONFIG environment variable holds a list of kubeconfig files.
		// For Linux and Mac, the list is colon-delimited. For Windows, the list
		// is semicolon-delimited. Ref:
		// https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable
		var separator string
		if runtime.GOOS == "windows" {
			separator = ";"
		} else {
			separator = ":"
		}
		kubeConfigFiles = strings.Split(kubeConfigPath, separator)
	} else {
		kubeConfigFiles = append(kubeConfigFiles, filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	}

	config, err := buildConfigFromFlags(context, kubeConfigFiles)
	if err != nil {
		log.Fatal(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return clientset
}

func buildConfigFromFlags(context string, kubeConfigFiles []string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{Precedence: kubeConfigFiles},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

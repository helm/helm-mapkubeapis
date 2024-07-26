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

package v3

import (
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"

	common "github.com/helm/helm-mapkubeapis/pkg/common"
)

var (
	settings = cli.New()
)

// GetActionConfig returns action configuration based on Helm env
func GetActionConfig(namespace string, kubeConfig common.KubeConfig) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)

	// Add kube config settings passed by user
	settings.KubeConfig = kubeConfig.File
	settings.KubeContext = kubeConfig.Context

	// check if the namespace is passed by the user. If not get Helm to return the current namespace
	if namespace == "" {
		namespace = settings.Namespace()
	}

	err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debug)
	if err != nil {
		return nil, err
	}

	return actionConfig, err
}

func debug(format string, v ...interface{}) {
	if settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		err := log.Output(2, fmt.Sprintf(format, v...))
		if err != nil {
			return
		}
	}
}

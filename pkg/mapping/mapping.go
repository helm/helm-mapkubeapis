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

package mapping

import (
	"strings"
)

// Mapping describes mappings which defines the Kubernetes
// API deprecations and the new replacement API
type Mapping struct {
	// From is the API looking to be mapped
	DeprecatedAPI string `json:"deprecatedAPI"`

	// To is the API to be mapped to
	NewAPI string `json:"newAPI"`

	// Kubernetes version API is deprecated in
	DeprecatedInVersion string `json:"deprecatedInVersion,omitempty"`

	// Kubernetes version API is removed in
	RemovedInVersion string `json:"removedInVersion,omitempty"`
}

const (
	apiVersionLabel     = "apiVersion:"
	kindLabel           = "kind:"
	apiVersionLabelSize = len(apiVersionLabel)
	kindLabelSize       = len(kindLabel)
)

// ParseAPIString parses the API string into version and kind components
func ParseAPIString(apiString string) (version, kind string) {
	idx := strings.Index(apiString, apiVersionLabel)
	if idx != -1 {
		temps := apiString[idx+apiVersionLabelSize:]
		idx = strings.Index(temps, "\n")
		if idx != -1 {
			temps = temps[:idx]
			version = strings.TrimSpace(temps)
		}
	}
	idx = strings.Index(apiString, kindLabel)
	if idx != -1 {
		temps := apiString[idx+kindLabelSize:]
		idx = strings.Index(temps, "\n")
		if idx != -1 {
			temps = temps[:idx]
			kind = strings.TrimSpace(temps)
		}
	}
	return version, kind
}

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

// Mapping describes mappings which defines the Kubernetes
// API deprecations and the new replacement API
type Mapping struct {
	// From is the API looking to be mapped
	DeprecatedAPI API `json:"deprecatedAPI"`

	// To is the API to be mapped to
	NewAPI API `json:"newAPI"`

	// Kubernetes version API is deprecated in
	DeprecatedInVersion string `json:"deprecatedInVersion,omitempty"`

	// Kubernetes version API is removed in
	RemovedInVersion string `json:"removedInVersion,omitempty"`
}

// API the apiVersion and kind uniquely identify the resources being
// updated.
type API struct {
	// Kind is a Kubernetes resource such as Deployment, PriorityClass, Pod that's
	// found in an API version.
	Kind string `json:"kind"`

	// APIVersion is the version of the Kubernetes resources. This usually takes the form
	// scheduling.k8s.io/v1beta1, extensions/v1beta1, or networking.k8s.io/v1beta1
	APIVersion string `json:"apiVersion"`
}

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

import (
	"log"
	"strings"
)

// KubeConfig are the Kubernetes configuration settings
type KubeConfig struct {
	Context string
	File    string
}

// MapOptions are the options for mapping deprecated apis in a release
type MapOptions struct {
	DryRun           bool
	KubeConfig       KubeConfig
	ReleaseName      string
	ReleaseNamespace string
	StorageType      string
	TillerOutCluster bool
}

// UpgradeDescription is description of why release was upgraded
const UpgradeDescription = "Kubernetes deprecated API upgrade - DO NOT rollback from this version"

var mappedAPIs = map[string]string{
	"apiVersion: extensions/v1beta1\nkind: Deployment":                             "apiVersion: apps/v1\nkind: Deployment",
	"apiVersion: apps/v1beta1\nkind: Deployment":                                   "apiVersion: apps/v1\nkind: Deployment",
	"apiVersion: apps/v1beta2\nkind: Deployment":                                   "apiVersion: apps/v1\nkind: Deployment",
	"apiVersion: apps/v1beta1\nkind: StatefulSet":                                  "apiVersion: apps/v1\nkind: StatefulSet",
	"apiVersion: apps/v1beta2\nkind: StatefulSet":                                  "apiVersion: apps/v1\nkind: StatefulSet",
	"apiVersion: extensions/v1beta1\nkind: DaemonSet":                              "apiVersion: apps/v1\nkind: DaemonSet",
	"apiVersion: apps/v1beta2\nkind: DaemonSet":                                    "apiVersion: apps/v1\nkind: DaemonSet",
	"apiVersion: extensions/v1beta1\nkind: ReplicaSet":                             "apiVersion: apps/v1\nkind: ReplicaSet",
	"apiVersion: apps/v1beta1\nkind: ReplicaSet":                                   "apiVersion: apps/v1\nkind: ReplicaSet",
	"apiVersion: apps/v1beta2\nkind: ReplicaSet":                                   "apiVersion: apps/v1\nkind: ReplicaSet",
	"apiVersion: extensions/v1beta1\nkind: NetworkPolicy":                          "apiVersion: networking.k8s.io/v1\nkind: NetworkPolicy",
	"apiVersion: extensions/v1beta1\nkind: PodSecurityPolicy":                      "apiVersion: policy/v1beta1\nkind: PodSecurityPolicy",
	"apiVersion: apiextensions.k8s.io/v1beta1\nkind: CustomResourceDefinition":     "apiVersion: apiextensions.k8s.io/v1\nkind: CustomResourceDefinition",
	"apiVersion: extensions/v1beta1\nkind: Ingress":                                "apiVersion: networking.k8s.io/v1beta1\nkind: Ingress",
	"apiVersion: rbac.authorization.k8s.io/v1alpha1\nkind: ClusterRole":            "apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRole",
	"apiVersion: rbac.authorization.k8s.io/v1alpha1\nkind: ClusterRoleList":        "apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRoleList",
	"apiVersion: rbac.authorization.k8s.io/v1alpha1\nkind: ClusterRoleBinding":     "apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRoleBinding",
	"apiVersion: rbac.authorization.k8s.io/v1alpha1\nkind: ClusterRoleBindingList": "apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRoleBindingList",
	"apiVersion: rbac.authorization.k8s.io/v1alpha1\nkind: Role":                   "apiVersion: rbac.authorization.k8s.io/v1\nkind: Role",
	"apiVersion: rbac.authorization.k8s.io/v1alpha1\nkind: RoleList":               "apiVersion: rbac.authorization.k8s.io/v1\nkind: RoleList",
	"apiVersion: rbac.authorization.k8s.io/v1alpha1\nkind: RoleBinding":            "apiVersion: rbac.authorization.k8s.io/v1\nkind: RoleBinding",
	"apiVersion: rbac.authorization.k8s.io/v1alpha1\nkind: RoleBindingList":        "apiVersion: rbac.authorization.k8s.io/v1\nkind: RoleBindingList",
	"apiVersion: rbac.authorization.k8s.io/v1beta1\nkind: ClusterRole":             "apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRole",
	"apiVersion: rbac.authorization.k8s.io/v1beta1\nkind: ClusterRoleList":         "apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRoleList",
	"apiVersion: rbac.authorization.k8s.io/v1beta1\nkind: ClusterRoleBinding":      "apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRoleBinding",
	"apiVersion: rbac.authorization.k8s.io/v1beta1\nkind: ClusterRoleBindingList":  "apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRoleBindingList",
	"apiVersion: rbac.authorization.k8s.io/v1beta1\nkind: Role":                    "apiVersion: rbac.authorization.k8s.io/v1\nkind: Role",
	"apiVersion: rbac.authorization.k8s.io/v1beta1\nkind: RoleList":                "apiVersion: rbac.authorization.k8s.io/v1\nkind: RoleList",
	"apiVersion: rbac.authorization.k8s.io/v1beta1\nkind: RoleBinding":             "apiVersion: rbac.authorization.k8s.io/v1\nkind: RoleBinding",
	"apiVersion: rbac.authorization.k8s.io/v1beta1\nkind: RoleBindingList":         "apiVersion: rbac.authorization.k8s.io/v1\nkind: RoleBindingList",
}

// ReplaceManifestUnSupportedAPIs returns a release manifest with deprecated or removed
// Kubernetes APIs updated to supported APIs
func ReplaceManifestUnSupportedAPIs(origManifest string) string {
	var modifiedManifest string

	// Check for deprecated or removed APIs and map accordingly to supported versions
	for deprecatedAPI, supportedAPI := range mappedAPIs {
		var modManifestForAPI string
		if len(modifiedManifest) <= 0 {
			modManifestForAPI = strings.ReplaceAll(origManifest, deprecatedAPI, supportedAPI)
			if modManifestForAPI != origManifest {
				log.Printf("Found deprecated or removed Kubernetes API:\n\"%s\"\nSupported API equivalent:\n\"%s\"\n", deprecatedAPI, supportedAPI)
			}

		} else {
			modManifestForAPI = strings.ReplaceAll(modifiedManifest, deprecatedAPI, supportedAPI)
			if modManifestForAPI != modifiedManifest {
				log.Printf("Found deprecated or removed Kubernetes API:\n\"%s\"\nSupported API equivalent:\n\"%s\"\n", deprecatedAPI, supportedAPI)
			}
		}
		modifiedManifest = modManifestForAPI
	}

	return modifiedManifest
}

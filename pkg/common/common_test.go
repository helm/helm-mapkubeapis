package common_test

import (
	"testing"

	"github.com/helm/helm-mapkubeapis/pkg/common"
	"github.com/helm/helm-mapkubeapis/pkg/mapping"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestCommon(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Deprecated APIs replacement suite")
}

// CheckDecode verifies that the passed YAML is parsing correctly
// It doesn't check semantic correctness
func CheckDecode(manifest string) error {
	_, err := common.ParseYAML(manifest)
	return err
}

var _ = ginkgo.Describe("replacing deprecated APIs", ginkgo.Ordered, func() {
	var mapFile *mapping.Metadata

	var deprecatedPodDisruptionBudget mapping.APIVersionKind
	var newPodDisruptionBudget mapping.APIVersionKind

	var deprecatedDeployment mapping.APIVersionKind
	var newDeployment mapping.APIVersionKind

	var deprecatedPodSecurityPolicy mapping.APIVersionKind

	ginkgo.BeforeAll(func() {
		// "apiVersion: policy/v1beta1\nkind: PodDisruptionBudget\n"
		deprecatedPodDisruptionBudget = mapping.APIVersionKind{
			APIVersion: "policy/v1beta1",
			Kind:       "PodDisruptionBudget",
		}
		// "apiVersion: policy/v1\nkind: PodDisruptionBudget\n"
		newPodDisruptionBudget = mapping.APIVersionKind{
			APIVersion: "policy/v1",
			Kind:       "PodDisruptionBudget",
		}

		// "apiVersion: apps/v1beta2\nkind: Deployment\n"
		deprecatedDeployment = mapping.APIVersionKind{
			APIVersion: "apps/v1beta2",
			Kind:       "Deployment",
		}
		// "apiVersion: apps/v1\nkind: Deployment\n"
		newDeployment = mapping.APIVersionKind{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		}

		// "apiVersion: policy/v1beta1\nkind: PodSecurityPolicy\n"
		deprecatedPodSecurityPolicy = mapping.APIVersionKind{
			APIVersion: "policy/v1beta1",
			Kind:       "PodSecurityPolicy",
		}

		mapFile = &mapping.Metadata{
			Mappings: []*mapping.Mapping{
				{
					// - deprecatedAPI: "apiVersion: policy/v1beta1\nkind: PodDisruptionBudget\n"
					//   newAPI: "apiVersion: policy/v1\nkind: PodDisruptionBudget\n"
					//   deprecatedInVersion: "v1.21"
					//   removedInVersion: "v1.25"
					DeprecatedAPI:       deprecatedPodDisruptionBudget,
					NewAPI:              newPodDisruptionBudget,
					DeprecatedInVersion: "v1.21",
					RemovedInVersion:    "v1.25",
				},
				{
					// - deprecatedAPI: "apiVersion: apps/v1beta2\nkind: Deployment\n"
					//   newAPI: "apiVersion: apps/v1\nkind: Deployment\n"
					//   deprecatedInVersion: "v1.9"
					//   removedInVersion: "v1.16"
					DeprecatedAPI:       deprecatedDeployment,
					NewAPI:              newDeployment,
					DeprecatedInVersion: "v1.9",
					RemovedInVersion:    "v1.16",
				},
				{
					// - deprecatedAPI: "apiVersion: policy/v1beta1\nkind: PodSecurityPolicy"
					//   deprecatedInVersion: "v1.21"
					//   removedInVersion: "v1.25"
					DeprecatedAPI:    deprecatedPodSecurityPolicy,
					RemovedInVersion: "v1.25",
				},
			},
		}
	})

	ginkgo.When("a deprecated API exists in the manifest", func() {
		ginkgo.When("it is a superseded API", func() {
			var (
				deploymentManifest                           string
				expectedResultingDeploymentManifest          string
				podDisruptionBudgetManifest                  string
				expectedResultingPodDisruptionBudgetManifest string
			)

			ginkgo.BeforeAll(func() {
				deploymentManifest = `---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: test
  namespace: test-ns
spec:
  template:
    containers:
      - image: test-image
        name: test-container
`

				expectedResultingDeploymentManifest = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
spec:
  template:
    containers:
      - image: test-image
        name: test-container
`

				podDisruptionBudgetManifest = `---
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: pdb-test
  namespace: test-ns
`

				expectedResultingPodDisruptionBudgetManifest = `---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: pdb-test
  namespace: test-ns
`
			})

			ginkgo.It("replaces deprecated resources with a new version in Kubernetes v1.25", func() {
				kubeVersion125 := "v1.25"
				modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, deploymentManifest, kubeVersion125)

				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(modifiedDeploymentManifest).To(gomega.Equal(expectedResultingDeploymentManifest))

				modifiedPdbManifest, err := common.ReplaceManifestData(mapFile, podDisruptionBudgetManifest, kubeVersion125)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(modifiedPdbManifest).To(gomega.Equal(expectedResultingPodDisruptionBudgetManifest))

				err = CheckDecode(modifiedDeploymentManifest)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.When("it is a removed API", func() {
			var kubeVersion125 = "v1.25"
			var expectedResultManifest = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: test-ns
`

			ginkgo.When("it is in the beginning of the manifest", func() {
				var podSecurityPolicyManifest = `---
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: test-psp
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: test-ns
`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.APIVersion))
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.Kind))
					gomega.Expect(expectedResultManifest).To(gomega.Equal(expectedResultManifest))

					err = CheckDecode(modifiedDeploymentManifest)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})

			ginkgo.When("it is at the end of the manifest", func() {
				var podSecurityPolicyManifest = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: test-ns
---
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: test-psp
`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.APIVersion))
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.Kind))
					gomega.Expect(modifiedDeploymentManifest).To(gomega.Equal(expectedResultManifest))

					err = CheckDecode(modifiedDeploymentManifest)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})

			ginkgo.When("it is in the middle of other manifests", func() {
				var podSecurityPolicyManifest = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
---
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: test-psp
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: test-ns
`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.APIVersion))
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.Kind))
					gomega.Expect(modifiedDeploymentManifest).To(gomega.Equal(expectedResultManifest))

					err = CheckDecode(modifiedDeploymentManifest)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})

			ginkgo.When("a three-dash is missing at the beginning", func() {
				var podSecurityPolicyManifest = `apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: test-psp
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: test-ns
`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.APIVersion))
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.Kind))
					gomega.Expect(modifiedDeploymentManifest).To(gomega.Equal(expectedResultManifest))

					err = CheckDecode(modifiedDeploymentManifest)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})

			ginkgo.When("apiVersion is not the first field", func() {
				var podSecurityPolicyManifest = `---
metadata:
  name: test-psp
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: test-ns
`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.APIVersion))
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.Kind))
					gomega.Expect(modifiedDeploymentManifest).To(gomega.Equal(expectedResultManifest))

					err = CheckDecode(modifiedDeploymentManifest)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})

			ginkgo.When("apiVersion is not the first field and a three-dash is missing at the beginning of the manifest", func() {
				var podSecurityPolicyManifest = `metadata:
  name: test-psp
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: test-ns
`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.APIVersion))
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.Kind))
					gomega.Expect(modifiedDeploymentManifest).To(gomega.Equal(expectedResultManifest))

					err = CheckDecode(modifiedDeploymentManifest)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})

			ginkgo.When("apiVersion is not the first field and the resource is in the middle of the manifest", func() {
				var podSecurityPolicyManifest = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
---
metadata:
  name: test-psp
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: test-ns
`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.APIVersion))
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.Kind))
					gomega.Expect(modifiedDeploymentManifest).To(gomega.Equal(expectedResultManifest))

					err = CheckDecode(modifiedDeploymentManifest)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})

			ginkgo.When("apiVersion is not the first field and the resource is at the end of the manifest", func() {
				var podSecurityPolicyManifest = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: test-ns
---
metadata:
  name: test-psp
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
spec:
  allowPrivilegeEscalation: true
`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.APIVersion))
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy.Kind))
					gomega.Expect(modifiedDeploymentManifest).To(gomega.Equal(expectedResultManifest))

					err = CheckDecode(modifiedDeploymentManifest)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})
		})
	})
})

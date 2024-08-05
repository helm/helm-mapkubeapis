package common_test

import (
	"bytes"
	"github.com/helm/helm-mapkubeapis/pkg/common"
	"github.com/helm/helm-mapkubeapis/pkg/mapping"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"testing"

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
	decoder := yaml.NewDecoder(bytes.NewBufferString(manifest))

	for {
		var value interface{}

		err := decoder.Decode(&value)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

var _ = ginkgo.Describe("replacing deprecated APIs", ginkgo.Ordered, func() {
	var mapFile *mapping.Metadata

	var deprecatedPodDisruptionBudget string
	var newPodDisruptionBudget string

	var deprecatedDeployment string
	var newDeployment string

	var deprecatedPodSecurityPolicy string

	ginkgo.BeforeAll(func() {
		deprecatedPodDisruptionBudget = "apiVersion: policy/v1beta1\nkind: PodDisruptionBudget\n"
		newPodDisruptionBudget = "apiVersion: policy/v1\nkind: PodDisruptionBudget\n"

		deprecatedDeployment = "apiVersion: apps/v1beta2\nkind: Deployment\n"
		newDeployment = "apiVersion: apps/v1\nkind: Deployment\n"

		deprecatedPodSecurityPolicy = "apiVersion: policy/v1beta1\nkind: PodSecurityPolicy\n"

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
    - name: test-container
      image: test-image`

				expectedResultingDeploymentManifest = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test-ns
spec:
  template:
    containers:
    - name: test-container
      image: test-image`

				podDisruptionBudgetManifest = `---
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: pdb-test
  namespace: test-ns`

				expectedResultingPodDisruptionBudgetManifest = `---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: pdb-test
  namespace: test-ns`
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
  namespace: test-ns`

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
  namespace: test-ns`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy))
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
  name: test-psp`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy))
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
  namespace: test-ns`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy))
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
  namespace: test-ns`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy))
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
  namespace: test-ns`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy))
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
  namespace: test-ns`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy))
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
  namespace: test-ns`

				ginkgo.It("removes the deprecated API manifest and leaves a valid YAML", func() {
					modifiedDeploymentManifest, err := common.ReplaceManifestData(mapFile, podSecurityPolicyManifest, kubeVersion125)

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy))
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
					gomega.Expect(modifiedDeploymentManifest).ToNot(gomega.ContainSubstring(deprecatedPodSecurityPolicy))
					gomega.Expect(modifiedDeploymentManifest).To(gomega.Equal(expectedResultManifest))

					err = CheckDecode(modifiedDeploymentManifest)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})
		})
	})
})

package common_test

import (
	"testing"

	"github.com/helm/helm-mapkubeapis/pkg/common"
	"github.com/helm/helm-mapkubeapis/pkg/mapping"
)

const (
	basicPriorityClass = `
---
apiVersion: scheduling.k8s.io/v1beta1
description: important-priorityclass critical
globalDefault: false
kind: PriorityClass
metadata:
  labels:
    team: mapkubeapis
  name: important-priorityclass-critical
value: 1000000000
`

	multidocYamlOne = `---
app: mapkubeapis
description: updates kube apis
`
	multidocYamlTwo = `---
app: helm
description: deploys kube resources
`
)

func TestCheckForOldAPI(t *testing.T) {
	manifest, err := common.UnmarshallManifest(basicPriorityClass)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	updated, err := common.CheckForOldAPI(
		manifest,
		&mapping.Mapping{
			DeprecatedAPI: mapping.API{
				Kind:       "PriorityClass",
				APIVersion: "scheduling.k8s.io/v1beta1",
			},
			NewAPI: mapping.API{
				Kind:       "PriorityClass",
				APIVersion: "scheduling.k8s.io/v1",
			},
			RemovedInVersion: "v1.22",
		},
		"v1.22",
	)
	if err != nil {
		t.Logf("failed to check manifest: %s", err)
		t.Fail()
	}

	if !updated {
		t.Logf("found the old API version in the new manifest")
		t.Fail()
	}
}

func TestSeparateDocuments(t *testing.T) {
	docs := common.SeparateDocuments(multidocYamlOne + multidocYamlTwo)
	if docs[0] != multidocYamlOne {
		t.Logf(docs[0])
		t.Fail()
	}

	if docs[1] != multidocYamlTwo {
		t.Logf(docs[1])
		t.Fail()
	}
}

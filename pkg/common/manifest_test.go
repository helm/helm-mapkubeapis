package common_test

import (
	"strings"
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
)

func TestCheckForOldAPI(t *testing.T) {
	manifestWithNewAPIVersion, err := common.CheckForOldAPI(
		basicPriorityClass,
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
	)
	if err != nil {
		t.Logf("failed to check manifest: %s", err)
		t.Fail()
	}

	if strings.Contains(manifestWithNewAPIVersion, "scheduling.k8s.io/v1beta1") {
		t.Logf("found the old API version in the new manifest")
		t.Fail()
	}
}

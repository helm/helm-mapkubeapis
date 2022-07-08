package common_test

import (
	"strings"
	"testing"

	"github.com/helm/helm-mapkubeapis/pkg/common"
)

const (
	basicDeployment = `
---
apiVersion: extensions/v1beta1
description: "the basic nginx-deployment from the kube docs"
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
  app: nginx
spec:
  replicas: 3
  selector:
  matchLabels:
    app: nginx
  template:
  metadata:
    labels:
    app: nginx
  spec:
    containers:
    - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
`
)

func TestReplaceManifestUnSupportedAPIs(t *testing.T) {
	multidocManifest := basicDeployment + basicPriorityClass
	newManifest, err := common.ReplaceManifestUnSupportedAPIs(multidocManifest, "../../fixtures/test-map.yaml", "v1.22")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if strings.Contains(newManifest, "extensions/v1beta1") || strings.Contains(newManifest, "scheduling.k8s.io/v1beta1") {
		t.Fail()
		t.Log(newManifest)
	}
}

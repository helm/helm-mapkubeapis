# Helm mapdepapis Plugin

## Overview

Helm chart templates uses `API version` and `Kind` properties when defining Kuberentes resources, similar to  manifest files. Kubernetes can deprecate `API versions` bentween minor releases. An example of such [deprecation is in Kubernetes 1.16](https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/).

If you upgrade a Kubernetes cluster to a later release which deprecates API versions of deployed Helm releases then a Helm upgrade iof such a release will fail. The `helm upgrade` command fails as it attempts to create a diff patch between the current deployed release which uses the Kubernetes APIs that are deprecated against the chart you are passing with the updated API versions. This is because when Kubernetes removes an API version, the go libraries can no longer parse the deprecated objects and Helm therfore fails calling the libraries.

Errors similar to the following can be seen:

```
Error: UPGRADE FAILED: unable to build kubernetes objects from current release manifest: unable to recognize "": no matches for kind "Deployment" in version "apps/v1beta1"
```

The `mapdepapis` plugin fixes the issue by mapping the deprecated Kubernetes APIs inline in the release metadata.

> Note: It currently support Helm v3 only.

## Usage

### Map Helm v3 deprecated Kubernetes APIs

Map Helm v3 release deprecated Kubernetes APIs in-place

Usage:
  mapkubapis [command]

Available Commands:
  help        Help about any command
  v3map       map v3 release deprecated Kubernetes APIs in-place

Flags:
  -h, --help   help for mapkubapis

## Developer (From Source) Install

If you would like to handle the build yourself, this is the recommended way to do it.

You must first have [Go v1.13](http://golang.org) installed, and then you run:

```console
$ mkdir -p ${GOPATH}/src/github.com
$ cd $_
$ git clone git@github.com:hickeyma/helm-mapkubeapis.git
$ cd helm-mapkubeapis
$ make build
$ helm plugin install <your_path>/helm-mapkubeapis
```

That last command will use the binary that you built.

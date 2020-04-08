# Helm mapkubeapis Plugin

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/hickeyma/helm-mapkubeapis)](https://goreportcard.com/report/github.com/hickeyma/helm-mapkubeapis)
[![CircleCI](https://circleci.com/gh/hickeyma/helm-mapkubeapis/tree/master.svg?style=svg)](https://circleci.com/gh/hickeyma/helm-mapkubeapis/tree/master)
[![Release](https://img.shields.io/github/release/hickeyma/helm-mapkubeapis.svg?style=flat-square)](https://github.com/hickeyma/helm-mapkubeapis/releases/latest)

mapkubeapis is a simple Helm plugin which is designed to update Helm release metadata that contains deprecated Kubernetes APIs to a new instance with supported Kubernetes APIs. Jump to [background to the issue](#background-to-the-issue) for more details on the problem space that the plugin solves.

> Note: It currently supports Helm v3 only.

## Prerequisite

- Helm v3 client with `mapkubeapis` plugin installed on the same system
- Access to the cluster(s) that Helm v3 manages. This access is similar to `kubectl` access using [kubeconfig files](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/).
  The `--kubeconfig`, `--kube-context` and `--namespace` flags can be used to set the kubeconfig path, kube context and namespace context to override the environment configuration.

## Install

Based on the version in `plugin.yaml`, release binary will be downloaded from GitHub:

```console
$ helm plugin install https://github.com/hickeyma/helm-mapkubeapis
Downloading and installing helm-mapkubeapis v0.0.1 ...
https://github.com/hickeyma/helm-mapkubeapis/releases/download/v0.0.1/helm-mapkubeapis_0.0.1_darwin_amd64.tar.gz
Installed plugin: mapkubeapis
```

### For Windows (using WSL)

Helm's plugin install hook system relies on `/bin/sh`, regardless of the operating system present. Windows users can work around this by using Helm under [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10).
```
$ wget https://get.helm.sh/helm-v3.0.0-linux-amd64.tar.gz
$ tar xzf helm-v3.0.0-linux-amd64.tar.gz
$ ./linux-amd64/helm plugin install https://github.com/hickeyma/helm-mapkubeapis
```

## Usage

### Map Helm deprecated Kubernetes APIs

Map release deprecated Kubernetes APIs in-place:

```console
$ helm mapkubeapis [flags] RELEASE 

Flags:
      --dry-run               simulate a command
  -h, --help                  help for v3map
      --kube-context string   name of the kubeconfig context to use
      --kubeconfig string     path to the kubeconfig file
      --namespace string      namespace scope of the release
```

Example output:

```console
$ helm mapkubeapis oldapi-chrt
2020/04/08 16:15:30 Release 'oldapi-chrt' will be checked for deprecated Kubernetes APIs and will be updated if necessary to supported API versions.
2020/04/08 16:15:30 Get release 'oldapi-chrt' latest version.
2020/04/08 16:15:30 Check release 'oldapi-chrt' for deprecated APIs...
2020/04/08 16:15:30 Found deprecated Kubernetes API:
"apiVersion: extensions/v1beta1
kind: Ingress"
Supported API equivalent:
"apiVersion: networking.k8s.io/v1beta1
kind: Ingress"
2020/04/08 16:15:30 Found deprecated Kubernetes API:
"apiVersion: apps/v1beta1
kind: Deployment"
Supported API equivalent:
"apiVersion: apps/v1
kind: Deployment"
2020/04/08 16:15:30 Finished checking release 'oldapi-chrt' for deprecated APIs.
2020/04/08 16:15:30 Deprecated APIs exist, updating release: oldapi-chrt.
2020/04/08 16:15:30 Set status of release version 'oldapi-chrt.v1' to 'superseded'.
2020/04/08 16:15:30 Release version 'oldapi-chrt.v1' updated successfully.
2020/04/08 16:15:30 Add release version 'oldapi-chrt.v2' with updated supported APIs.
2020/04/08 16:15:30 Release version 'oldapi-chrt.v2' added successfully.
2020/04/08 16:15:30 Release 'oldapi-chrt' with deprecated APIs updated successfully to new version.
2020/04/08 16:15:30 Map of release 'oldapi-chrt' deprecated APIs to supported APIs, completed successfully.
```

## Background to the issue

Helm chart templates uses `API version` and `Kind` properties when defining Kuberentes resources, similar to  manifest files. Kubernetes can deprecate `API versions` between minor releases. An example of such [deprecation is in Kubernetes 1.16](https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/).

There is a Kubernetes API deprecation policy and all chart maintainers needs to be cognisant of this and update chart Kubernetes APIs appropriately. This is not such an issue when installing a chart as it will just fail if the chart API versions are not fully compliant. You then need to get the latest chart version or update the chart yourself.

This is however become a problem for Helm releases that are already deployed with APIs that are within the deprecation time period. If the Kubernetes cluster (containing such releases) is updated to a later version where the APIs become deprecated, then Helm becomes unable to manage such releases anymore. It does not matter is the chart being passed in the upgrade contains the supported API.
 
An example of this is the `helm upgrade` command. It fails with an error similar to the following:

```
Error: UPGRADE FAILED: unable to build kubernetes objects from current release manifest: unable to recognize "": no matches for kind "Deployment" in version "apps/v1beta1"
```

Helm fails because it attempts to create a diff patch between the current deployed release which contains the Kubernetes APIs that are deprecated against the chart you are passing with the updated/supported API versions. The underlying reason for failure is due to when Kubernetes removes an API version, the Go libraries can no longer parse the deprecated objects and Helm therefore fails calling the libraries.

The `mapkubeapis` plugin fixes the issue by mapping releases which contain deprecated Kubernetes APIs to supported APIs. This is performed inline in the release metadata where the existing release is `superseded` and a new release (metadata only) is added. The deployed Kubernetes resources are updated automatically by Kubernetes during upgrade of its version. Once this operation is completed, you can then upgrade using the chart with supported APIs.

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

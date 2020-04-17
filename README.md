# Helm mapkubeapis Plugin

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/hickeyma/helm-mapkubeapis)](https://goreportcard.com/report/github.com/hickeyma/helm-mapkubeapis)
[![CircleCI](https://circleci.com/gh/hickeyma/helm-mapkubeapis/tree/master.svg?style=svg)](https://circleci.com/gh/hickeyma/helm-mapkubeapis/tree/master)
[![Release](https://img.shields.io/github/release/hickeyma/helm-mapkubeapis.svg?style=flat-square)](https://github.com/hickeyma/helm-mapkubeapis/releases/latest)

mapkubeapis is a simple Helm plugin which is designed to update Helm release metadata that contains deprecated or removed Kubernetes APIs to a new instance with supported Kubernetes APIs. Jump to [background to the issue](#background-to-the-issue) for more details on the problem space that the plugin solves.

## Prerequisite

- Helm client with `mapkubeapis` plugin installed on the same system
- Access to the cluster(s) that Helm manages. This access is similar to `kubectl` access using [kubeconfig files](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/).
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

### Map Helm deprecated or removed Kubernetes APIs

Map release deprecated or removed Kubernetes APIs in-place:

```console
$ helm mapkubeapis [flags] RELEASE 

Flags:
      --dry-run                  simulate a command
  -h, --help                     help for mapkubeapis
      --kube-context string      name of the kubeconfig context to use
      --kubeconfig string        path to the kubeconfig file
      --namespace string         namespace scope of the release. For Helm v2, this is the Tiller namespace (e.g. kube-system)
  -s, --release-storage string   for Helm v2 only - release storage type/object. It can be 'secrets' or 'configmaps'. This is only used with the 'tiller-out-cluster' flag (default "secrets")
      --tiller-out-cluster       for Helm v2 only - when Tiller is not running in the cluster e.g. Tillerless
      --v2                       run for Helm v2 release (default is Helm v3)
```

Example output:

```console
$ helm mapkubeapis oldapi-chrt --namespace test
2020/04/17 13:05:45 Release 'v2-oldapi' will be checked for deprecated or removed Kubernetes APIs and will be updated if necessary to supported API versions.
2020/04/17 13:05:45 Get release 'v2-oldapi' latest version.
2020/04/17 13:05:45 Check release 'v2-oldapi' for deprecated or removed APIs...
2020/04/17 13:05:45 Found deprecated or removed Kubernetes API:
"apiVersion: apps/v1beta1
kind: Deployment"
Supported API equivalent:
"apiVersion: apps/v1
kind: Deployment"
2020/04/17 13:05:45 Found deprecated or removed Kubernetes API:
"apiVersion: extensions/v1beta1
kind: Ingress"
Supported API equivalent:
"apiVersion: networking.k8s.io/v1beta1
kind: Ingress"
2020/04/17 13:05:45 Finished checking release 'v2-oldapi' for deprecated or removed APIs.
2020/04/17 13:05:45 Deprecated or removed APIs exist, updating release: v2-oldapi.
2020/04/17 13:05:45 Set status of release version 'v2-oldapi.v1' to 'superseded'.
2020/04/17 13:05:45 Release version 'v2-oldapi.v1' updated successfully.
2020/04/17 13:05:45 Add release version 'v2-oldapi.v2' with updated supported APIs.
2020/04/17 13:05:45 Release version 'v2-oldapi.v2' added successfully.
2020/04/17 13:05:45 Release 'v2-oldapi' with deprecated or removed APIs updated successfully to new version.
2020/04/17 13:05:45 Map of release 'v2-oldapi' deprecated or removed APIs to supported versions, completed successfully.
```

## Background to the issue

Kubernetes is an API-driven system and the API evolves over time to reflect the evolving understanding of the problem space. This is common practice across systems and their APIs. An important part of evolving APIs is a good deprecation policy and process to inform users of how changes to APIs are implemented. In other words, consumers of your API need to know in advance in what release an API will be removed or changed. This removes the element of surprise and unexpected breaking changes to the consumers. 

The [Kubernetes deprecation policy](https://kubernetes.io/docs/reference/using-api/deprecation-policy/) documents how Kubernetes handles the changes to its API versions. The policy for deprecation states the timeframe that API versions will be supported following a deprecation announcement. It is therefore important to be aware of deprecation announcements to mimimize the effect to you once an API version goes out of support.
This is an example of [the removal of deprecated API versions in Kubernetes 1.16](https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/). 

Helm chart templates uses Kubernetes `API version` and `Kind` properties when defining Kuberentes resources, similar to  manifest files. This therefore means that Helm users and chart maintainers need to be aware when Kubvernetes API versions have been deprecated and in what Kubernetes version they will removed.

This is not such a big issue when installing a chart as it will just fail if the chart API versions are no longer supported. In this situation, you then need to get the latest chart version or update the chart yourself.

This does however become a problem for Helm releases that are already deployed with APIs that are no longer supported. If the Kubernetes cluster (containing such releases) is updated to a version where the APIs are removed, then Helm becomes unable to manage such releases anymore. It does not matter if the chart being passed in the upgrade contains the supported API versions.
 
It fails with an error similar to the following:

```
Error: UPGRADE FAILED: unable to build kubernetes objects from current release manifest: unable to recognize "": no matches for kind "Deployment" in version "apps/v1beta1"
```

Helm fails because it attempts to create a diff patch between the current deployed release which contains the Kubernetes APIs that are removed against the chart you are passing with the updated/supported API versions. The underlying reason for failure is due because when Kubernetes removes an API version, its Go libraries can no longer parse the removed objects and Helm therefore fails calling the libraries.

The `mapkubeapis` plugin fixes the issue by mapping releases which contain deprecated or removed Kubernetes APIs to supported APIs. This is performed inline in the release metadata where the existing release is `superseded` and a new release (metadata only) is added. The deployed Kubernetes resources are updated automatically by Kubernetes during upgrade of its version. Once this operation is completed, you can then upgrade using the chart with supported APIs.

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

# Helm mapkubeapis Plugin

## Overview

Helm chart templates uses `API version` and `Kind` properties when defining Kuberentes resources, similar to  manifest files. Kubernetes can deprecate `API versions` between minor releases. An example of such [deprecation is in Kubernetes 1.16](https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/).

There is an Kubernetes API deprecation policy in place and all chart maintainers needs to be cognisant of this and update chart Kubernetes APIs appropriately. This is not such an issue when installing a chart as it will just fail if the chart API versions are not fully compliant. You then need to get the latest chart version or update the chart yourself.

This can however become a problem for Helm releases already deployed which use APIs within derprecated policy. If the Kubernetes cluster (containing such releases) is updated to a later version where the Kubernetes API deprecation is applied, then Helm becomes unable to manage such releases anymore.
 
An example of this is the `helm upgrade` command . It fails as it will attempt to create a diff patch between the current deployed release which contains the Kubernetes APIs that are deprecated against the chart you are passing with the updated/supported API versions. This is because when Kubernetes removes an API version, the Go libraries can no longer parse the deprecated objects and Helm therefore fails calling the libraries.

Errors similar to the following can be seen:

```
Error: UPGRADE FAILED: unable to build kubernetes objects from current release manifest: unable to recognize "": no matches for kind "Deployment" in version "apps/v1beta1"
```

The `mapkubeapis` plugin fixes the issue by mapping releases which contain deprecated Kubernetes APIs to supported APIs. This is performed inline in the release metadata. Once this operation is completed, you can then upgrade using the chart with supported APIs.

> Note: It currently support Helm v3 only.

## Install

Based on the version in `plugin.yaml`, release binary will be downloaded from GitHub:

```console
$ helm plugin install https://github.com/hickeyma/helm-mapkubeapis.git
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

### Map Helm v3 deprecated Kubernetes APIs

Map v3 release deprecated Kubernetes APIs in-place:

```console
$ helm mapkubeapis v3map RELEASE RELEASE_NAMESPACE [flags]

Flags:
      --dry-run               simulate a command
  -h, --help                  help for v3map
      --kube-context string   name of the kubeconfig context to use
      --kubeconfig string     path to the kubeconfig file
```

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

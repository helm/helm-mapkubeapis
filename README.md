# Helm mapkubeapis Plugin

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/helm/helm-mapkubeapis)](https://goreportcard.com/report/github.com/helm/helm-mapkubeapis)
[![CircleCI](https://circleci.com/gh/helm/helm-mapkubeapis/tree/master.svg?style=svg)](https://circleci.com/gh/helm/helm-mapkubeapis/tree/master)
[![Release](https://img.shields.io/github/release/helm/helm-mapkubeapis.svg?style=flat-square)](https://github.com/helm/helm-mapkubeapis/releases/latest)

`mapkubeapis` is a Helm v3 plugin which updates in-place Helm release metadata that contains deprecated or removed Kubernetes APIs to a new instance with supported Kubernetes APIs. Jump to [background to the issue](#background-to-the-issue) for more details on the problem space that the plugin solves.

> Note: Charts need to be updated also to supported Kubernetes APIs to avoid failure during deployment in a Kubernetes version. This is a separate task to the plugin.

## Prerequisite

- Helm client with `mapkubeapis` plugin installed on the same system
- Access to the cluster(s) that Helm manages. This access is similar to `kubectl` access using [kubeconfig files](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/).
  The `--kubeconfig`, `--kube-context` and `--namespace` flags can be used to set the kubeconfig path, kube context and namespace context to override the environment configuration.
- If you try and upgrade a release with unsupported APIs then the upgrade will fail. This is ok in Helm v3 as it will not generate a failed release for Helm.
- The plugin updates the lastest release version. The latest release version should be in a `deployed` state as you want to update a successful deployment. If it is not then you need to delete the latest release version. The command to remove a release version is:
    - Helm v3: `kubectl delete configmap/secret sh.helm.release.v1.<release_name>.v<latest_version_number> --namespace <release_namespace>`

## Install

Based on the version in `plugin.yaml`, release binary will be downloaded from GitHub:

```console
$ helm plugin install https://github.com/helm/helm-mapkubeapis
Downloading and installing helm-mapkubeapis v0.1.0 ...
https://github.com/helm/helm-mapkubeapis/releases/download/v0.1.0/helm-mapkubeapis_0.1.0_darwin_amd64.tar.gz
Installed plugin: mapkubeapis
```

### For Windows (using WSL)

Helm's plugin install hook system relies on `/bin/sh`, regardless of the operating system present. Windows users can work around this by using Helm under [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10).
```
$ wget https://get.helm.sh/helm-v3.0.0-linux-amd64.tar.gz
$ tar xzf helm-v3.0.0-linux-amd64.tar.gz
$ ./linux-amd64/helm plugin install https://github.com/helm/helm-mapkubeapis
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
      --mapfile string           path to the API mapping file (default "config/Map.yaml")
      --namespace string         namespace scope of the release
```

Example output:

```console
$ helm mapkubeapis cluster-role-example --namespace test-cluster-role-example         
2022/02/07 18:48:49 Release 'cluster-role-example' will be checked for deprecated or removed Kubernetes APIs and will be updated if necessary to supported API versions.
2022/02/07 18:48:49 Get release 'cluster-role-example' latest version.
2022/02/07 18:48:49 Check release 'cluster-role-example' for deprecated or removed APIs...
2022/02/07 18:48:49 Found 1 instances of deprecated or removed Kubernetes API:
"apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
"
Supported API equivalent:
"apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
"
2022/02/07 18:48:49 Found 1 instances of deprecated or removed Kubernetes API:
"apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
"
Supported API equivalent:
"apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
"
2022/02/07 18:48:49 Finished checking release 'cluster-role-example' for deprecated or removed APIs.
2022/02/07 18:48:49 Deprecated or removed APIs exist, updating release: cluster-role-example.
2022/02/07 18:48:49 Set status of release version 'cluster-role-example.v1' to 'superseded'.
2022/02/07 18:48:49 Release version 'cluster-role-example.v1' updated successfully.
2022/02/07 18:48:49 Add release version 'cluster-role-example.v2' with updated supported APIs.
2022/02/07 18:48:49 Release version 'cluster-role-example.v2' added successfully.
2022/02/07 18:48:49 Release 'cluster-role-example' with deprecated or removed APIs updated successfully to new version.
2022/02/07 18:48:49 Map of release 'cluster-role-example' deprecated or removed APIs to supported versions, completed successfully.
```

## API Mapping

The mapping information of deprecated or removed APIs to supported APIs is configured in the [Map.yaml](https://github.com/helm/helm-mapkubeapis/blob/master/config/Map.yaml) file. The file is a list of entries similar to the following:

```yaml
 - deprecatedAPI: "apiVersion: extensions/v1beta1\nkind: Deployment"
    newAPI: "apiVersion: apps/v1\nkind: Deployment"
    deprecatedInVersion: "v1.9"
    removedInVersion: "v1.16"
```

The plugin when performing update of a Helm release metadata first loads the map file from the `config` directory where the plugin is run from. If the map file is a different name or in a different location, you can use the `--mapfile` flag to specify the different mapping file.

The OOTB mapping file is configured as follows:

- The search and replace strings are in order with `apiVersion` first and then `kind`. This should be changed if the Helm release metadata is rendered with different search/replace string.
- The strings contain UNIX/Linux line feeds. This means that `\n` is used to signify line separation between properties in the strings. This should be changed if the Helm release metadata is rendered in Windows or Mac.
- Each mapping contains the Kubernetes version that the API is deprecated and removed in. This information is important as the plugin checks that the deprecated version (uses removed if deprecated unset) is later than the Kubernetes version that it is running against. If it is then no mapping occurs for this API as it not yet deprecated in this Kubernetes version and hence the new API is not yet supported. Otherwise, the mapping can proceed.

> Note: The Helm release metadata can be checked by following the steps in:
- Helm v3: [Updating API Versions of a Release Manifest](https://helm.sh/docs/topics/kubernetes_apis/#updating-api-versions-of-a-release-manifest)

## Background to the issue

For details on the background to this issue, it is recommended to read the docs appropriate to your Helm version. The docs can be accessed as follows:
- Helm v3: [Deprecated Kubernetes APIs](https://helm.sh/docs/topics/kubernetes_apis)

The Helm documentation describes the problem when Helm releases that are already deployed with APIs that are no longer supported. If the Kubernetes cluster (containing such releases) is updated to a version where the APIs are removed, then Helm becomes unable to manage such releases anymore. It does not matter if the chart being passed in the upgrade contains the supported API versions or not.

This is what the `mapkubeapis` plugin resolves. It fixes the issue by mapping releases which contain deprecated or removed Kubernetes APIs to supported APIs. This is performed inline in the release metadata where the existing release is `superseded` and a new release (metadata only) is added. The deployed Kubernetes resources are updated automatically by Kubernetes during upgrade of its version. Once this operation is completed, you can then upgrade using the chart with supported APIs.

## Helm v2 Support

Helm [v2.17.0](https://github.com/helm/helm/releases/tag/v2.17.0) was the final release of Helm v2 in October 2020. Helm v2 is unsupported since November 2020, as detailed in [Helm 2 and the Charts Project Are Now Unsupported](https://helm.sh/blog/helm-2-becomes-unsupported/). `mapkubeapis` Helm v2 support finished in [release v0.2.0](https://github.com/helm/helm-mapkubeapis/releases/tag/v0.2.0).

## Developer (From Source) Install

If you would like to handle the build yourself, this is the recommended way to do it.

You must first have [Go v1.18+](http://golang.org) installed, and then you run:

```console
$ mkdir -p ${GOPATH}/src/github.com
$ cd $_
$ git clone git@github.com:helm/helm-mapkubeapis.git
$ cd helm-mapkubeapis
$ make
$ export HELM_LINTER_PLUGIN_NO_INSTALL_HOOK=true
$ helm plugin install <your_path>/helm-mapkubeapis
```

That last command will use the binary that you built.

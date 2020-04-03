module helm-mapkubeapis

go 1.13

require (
	github.com/helm/helm-2to3 v0.5.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v2 v2.2.4
	helm.sh/helm/v3 v3.1.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)

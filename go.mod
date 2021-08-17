module github.com/hickeyma/helm-mapkubeapis

go 1.13

require (
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/maorfr/helm-plugin-utils v0.0.0-20200216074820-36d2fcf6ae86
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/mod v0.3.1-0.20200828183125-ce943fd02449
	helm.sh/helm/v3 v3.6.1
	k8s.io/helm v2.17.0+incompatible
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)

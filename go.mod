module github.com/helm/helm-mapkubeapis

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0 // indirect
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/maorfr/helm-plugin-utils v0.0.0-20200216074820-36d2fcf6ae86
	github.com/pkg/errors v0.9.1
	github.com/rubenv/sql-migrate v0.0.0-20200402132117-435005d389bc // indirect
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/mod v0.5.0
	helm.sh/helm/v3 v3.1.3
	k8s.io/helm v2.17.0+incompatible
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)

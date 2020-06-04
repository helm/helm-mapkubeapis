module github.com/hickeyma/helm-mapkubeapis

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1 // indirect
	github.com/golang/protobuf v1.4.0 // indirect
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/lib/pq v1.3.0 // indirect
	github.com/maorfr/helm-plugin-utils v0.0.0-20200216074820-36d2fcf6ae86
	github.com/pkg/errors v0.9.1
	github.com/rubenv/sql-migrate v0.0.0-20200402132117-435005d389bc // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	golang.org/x/mod v0.3.0
	helm.sh/helm/v3 v3.1.2
	k8s.io/helm v2.16.6+incompatible
	sigs.k8s.io/yaml v1.1.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)

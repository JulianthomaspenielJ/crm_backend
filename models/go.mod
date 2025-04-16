module crm/models

go 1.23.1

toolchain go1.23.3

replace crm/routes => ../routes

require golang.org/x/oauth2 v0.23.0

require (
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.35.1 // indirect
	golang.org/x/sys v0.26.0 // indirect
)

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.8.1
	github.com/joho/godotenv v1.5.1
//github.com/pelletier/go-toml v1.9.5
)

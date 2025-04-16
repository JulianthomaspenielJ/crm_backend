module crm/queue

go 1.23.1

//require github.com/go-redis/redis/v8 v8.11.5

replace crm/utils => ../utils

require gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df

require (
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/joho/godotenv v1.5.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)

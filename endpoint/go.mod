module github.com/RedeployAB/container-apps-dapr/endpoint

go 1.21.4

replace github.com/RedeployAB/container-apps-dapr/common => ../common

require (
	github.com/caarlos0/env/v10 v10.0.0
	github.com/dapr/go-sdk v1.9.1
	github.com/google/go-cmp v0.6.0
)

require (
	github.com/dapr/dapr v1.12.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231120223509-83a465c0220f // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

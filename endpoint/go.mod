module github.com/RedeployAB/container-apps-dapr/endpoint

go 1.20

replace github.com/RedeployAB/container-apps-dapr/common => ../common

require (
	github.com/RedeployAB/container-apps-dapr/common v0.0.0-00010101000000-000000000000
	github.com/caarlos0/env/v8 v8.0.0
	github.com/dapr/go-sdk v1.7.0
	github.com/go-logr/logr v1.2.4
	github.com/google/go-cmp v0.5.9
)

require (
	github.com/go-logr/zerologr v1.2.3 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/rs/zerolog v1.29.1 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.56.3 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

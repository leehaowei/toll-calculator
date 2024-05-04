# toll-calculator

`https://hub.docker.com/r/bitnami/kafka`

```
docker compose up -d
```

## Installing protobuf & gRPC
```
brew install protobuf
```

https://grpc.io/docs/languages/go/quickstart/
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

dependenncies
```
go get google.golang.org/protobuf
go get google.golang.org/grpc
```

prometheus golang client
```
go get github.com/prometheus/client_golang/prometheus
```

```
go mod edit -go=1.20
```
https://go.dev/dl/


### grafana
`http://host.docker.internal:9090`
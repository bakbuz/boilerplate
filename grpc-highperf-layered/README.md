# task-tracker

## Description

This is a simple task tracker application. It allows you to create lists and tasks.

## Hot to start (local)

Execute next commands to start application:

```bash
go mod download
go run ./cmd/app/main.go
```

## Hot to start (docker)

Execute next commands to start application:

```bash
docker-compose up -d
```


# grpc-highperf-layered


# buf.yaml

go install github.com/bufbuild/buf/cmd/buf@latest
buf config init
buf lint api/protos


# golangci-lint

go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint --version
golangci-lint run

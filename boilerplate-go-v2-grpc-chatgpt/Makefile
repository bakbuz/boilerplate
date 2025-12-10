.PHONY: proto build docker test

proto:
	protoc --go_out=./api/gen --go-grpc_out=./api/gen api/proto/shop.proto

build: proto
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/shop ./cmd/server

docker: build
	docker build -t yourorg/shop:latest .

test:
	go test ./... -v

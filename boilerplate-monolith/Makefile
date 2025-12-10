.PHONY: proto build docker test

proto:
	cls
	protoc --go_out=./api/pb --go-grpc_out=./api/pb ./api/protos/*.proto

build: proto
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/catalog ./cmd/server

docker: build
	docker build -t maydere/catalog:latest .

test:
	go test ./... -v

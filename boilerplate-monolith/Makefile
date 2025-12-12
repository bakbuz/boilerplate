.PHONY: proto build docker test u

proto:
	cls
	protoc --proto_path=./api/protos --go_out=./api/pb --go-grpc_out=./api/pb ./api/protos/*.proto

build: proto
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/catalog ./cmd/server

docker: build
	docker build -t maydere/catalog:latest .

test:
	go test ./... -v

u:
	cls
	go get -u all
	go mod tidy
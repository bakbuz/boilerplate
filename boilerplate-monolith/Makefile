.PHONY: proto build docker test u r

proto:
	cls
	protoc 
    	--proto_path=./api/protos 
		--go_out=./api/pb 
		--go-grpc_out=./api/pb 
		./api/protos/*.proto

build: proto
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/catalog ./cmd/server

build_sm: proto
	# -s -w flagleri binary boyutunu düşürür (debug bilgilerini siler)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/server cmd/server/main.go

docker: build
	docker build -t maydere/catalog:latest .

test:
	go test ./... -v

u:
	cls
	go get -u all
	go mod tidy

r:
	cls
	go run main.go
.PHONY: proto build docker test u r

# Variables
PROTO_PATH=api/protos/v1
OUTPUT_PATH=api/pb/v1

# Generate gRPC code
proto:
	@echo $(OUTPUT_PATH)
	protoc \
    	--proto_path=./${PROTO_PATH} \
		--go_out=./${OUTPUT_PATH} \
		--go-grpc_out=./${OUTPUT_PATH} \
		./${PROTO_PATH}/*.proto

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

pb-lint:
	buf lint ./api/protos/common/v1/common.proto
	buf lint ./api/protos/demo/v1/demo.proto
	buf lint ./api/protos/catalog/v1/brand.proto
	buf lint ./api/protos/catalog/v1/product.proto

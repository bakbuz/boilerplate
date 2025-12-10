.PHONY: proto build run clean

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/catalog.proto

build:
	# -s -w flagleri binary boyutunu düşürür (debug bilgilerini siler)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

clean:
	rm -rf bin/
.PHONY: run swag fix tidy build build-l

run:
	go run main.go

swag:
	swag init

fix1:
	go mod tidy
	go fmt ./...
	go fix
	go clean

fix:
	go mod tidy
	gofmt -s -w .
	make swag
	make tidy

tidy:
	go get -u all
	go mod tidy

build:
	set GOARCH=amd64
	set GOOS=linux
	go build -o bin/app main.go

build-l:
	GOOS=linux GOARCH=amd64 go build -o bin/app main.go

val1:
	go run _examples/validate/main.go

val2:
	go run _examples/validate2/main.go	
.PHONY: buf-gen buf-lint build build_sm build-prod fmt lint security test test-integration clean docker docker-run deploy profile u r

# Variables
BINARY_OUTPUT=.bin/server
OUTPUT_PATH=api/gen

# Generate gRPC code
buf-gen:
	buf generate

buf-lint:
	buf lint ./api/protos/common/v1/common.proto
	buf lint ./api/protos/demo/v1/demo.proto
	buf lint ./api/protos/catalog/v1/brand.proto
	buf lint ./api/protos/catalog/v1/product.proto

build: buf-gen
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/catalog ./cmd/server

# Build with optimizations
build_sm: buf-gen
	# -s -w flagleri binary boyutunu düşürür (debug bilgilerini siler)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/server cmd/server/main.go

# Build for production with PGO
build-prod: buf-gen
	@echo "Building with Profile Guided Optimization..."
	go build -pgo=auto -o bin/$(BINARY_NAME)-prod ./cmd/server

# Run formatter
fmt:
	@echo "Running formatter..."
	go fmt ./...

# Run linter
lint:
	cls
	@echo "linter version"
	golangci-lint --version
	@echo "Running linter..."
	golangci-lint run

# Security scan
security:
	gosec ./...	

# Run tests
test:
	@echo "Running tests..."
	go test ./...

test-repo:
	@echo "Running repository integration tests..."
	go test ./internal/repository/...

test-svc:
	@echo "Running service unit tests..."
	go test ./internal/service/...

# Integration tests
test-integration:
	@echo "Running integration tests..."
	go test -tags=integration ./tests/integration/... -v

test-e2e:
	@echo "Running end-to-end tests..."
	go test ./tests/e2e -v

test-benchmark:
	@echo "Running benchmark tests..."
	go test ./tests/benchmark/service -v	

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf $(OUTPUT_PATH)
	rm -f coverage.out	

# Docker build
docker: build
	@echo "Building Docker image..."
	docker build -t grpc-backend:latest .

# Run locally with docker-compose
docker-run:
	docker-compose up --build

# Deploy to AWS EC2
deploy: build
	@echo "Deploying to EC2..."
	scp bin/$(BINARY_NAME) deployment/ user@ec2-instance:/opt/grpc-backend/
	ssh user@ec2-instance "sudo systemctl restart grpc-backend"

# Performance profiling
profile:
	@echo "Starting CPU profiling..."
	go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...

u:
	cls
	go get -u all
	go mod tidy

r:
	cls
	go run main.go



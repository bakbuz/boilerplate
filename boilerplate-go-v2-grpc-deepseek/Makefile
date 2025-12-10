.PHONY: all build test clean proto docker deploy

# Variables
BINARY_NAME=server
PROTO_PATH=api/v1
OUTPUT_PATH=internal/gen

all: build

# Generate gRPC code
proto:
	@mkdir -p $(OUTPUT_PATH)
	protoc -I=$(PROTO_PATH) \
		--go_out=$(OUTPUT_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUTPUT_PATH) --go-grpc_opt=paths=source_relative \
		$(PROTO_PATH)/*.proto

# Build with optimizations
build: proto
	@echo "Building binary with optimizations..."
	CGO_ENABLED=0 go build \
		-ldflags="-w -s" \
		-trimpath \
		-o bin/$(BINARY_NAME) \
		./cmd/server

# Build for production with PGO
build-prod: proto
	@echo "Building with Profile Guided Optimization..."
	go build -pgo=auto -o bin/$(BINARY_NAME)-prod ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v -race -coverprofile=coverage.out

# Integration tests
test-integration:
	@echo "Running integration tests..."
	go test -tags=integration ./tests/integration/... -v

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf $(OUTPUT_PATH)
	rm -f coverage.out

# Docker build
docker:
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

# Lint
lint:
	golangci-lint run ./...

# Security scan
security:
	gosec ./...
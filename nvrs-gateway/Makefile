# Specify the Go version if needed
GOVERSION := 1.23.1

# Default target
.PHONY: all
all: test

# Download dependencies
.PHONY: deps
deps:
	go mod download

# Generate gRPC code from .proto files
.PHONY: generate
generate: deps
	protoc --go_out=. --go-grpc_out=. protos/*.proto

# Run tests
.PHONY: test
test: deps generate
	go test ./handlers -v

# Start the gRPC server only
.PHONY: run-grpc
run-grpc: deps generate
	go run grpc/grpc_server.go

# Start the entire application
.PHONY: run
run: deps generate
	go run cmd/nvrs-gateway/main.go

# Clean build artifacts (if any)
.PHONY: clean
clean:
	go clean
	rm -f protos/*.pb.go  # Remove generated .pb.go files if needed

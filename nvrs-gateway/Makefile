# Specify the Go version if needed
GOVERSION := 1.19

# Default target
.PHONY: all
all: test

# Download dependencies
.PHONY: deps
deps:
	go mod download

# Run tests
.PHONY: test
test: deps
	go test ./handlers -v

# Start the application
.PHONY: run
run: deps
	go run cmd/nvrs-gateway/main.go

# Clean build artifacts (if any)
.PHONY: clean
clean:
	go clean

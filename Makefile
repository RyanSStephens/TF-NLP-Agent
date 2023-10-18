.PHONY: build test clean run serve install deps

# Variables
BINARY_NAME=tf-nlp-agent
VERSION=1.0.0
BUILD_DIR=build

# Build the application
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/agent

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean

# Run the CLI application
run:
	go run ./cmd/agent

# Start the web server
serve:
	go run ./cmd/agent serve

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/agent
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/agent
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/agent

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate example configurations
examples:
	mkdir -p examples
	./$(BUILD_DIR)/$(BINARY_NAME) generate "Create an AWS VPC with public and private subnets" -o examples/vpc.tf
	./$(BUILD_DIR)/$(BINARY_NAME) generate "Deploy a web application with load balancer and database" -o examples/webapp.tf

# Install the binary
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/ 
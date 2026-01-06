.PHONY: build install test clean run help

# Binary name
BINARY_NAME=gitai

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

# Install to system
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	go clean
	@echo "Clean complete!"

# Run the application (for development)
run: build
	./$(BINARY_NAME)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .
	@echo "Multi-platform build complete!"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  install       - Build and install to /usr/local/bin"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Remove build artifacts"
	@echo "  run           - Build and run the application"
	@echo "  deps          - Download dependencies"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  help          - Show this help message"

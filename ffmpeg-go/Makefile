# Makefile for ffmpeg-go

.PHONY: build examples clean test install-deps

# Default target
all: build

# Build all examples
build: examples

# Build examples
examples:
	@echo "Building examples..."
	go build -o bin/info ./examples/info
	go build -o bin/decode ./examples/decode
	@echo "Examples built successfully!"

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Dependencies installed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Install the package
install:
	@echo "Installing ffmpeg-go..."
	go install .
	@echo "Installation complete!"
#!/bin/bash

# Test script for ffmpeg-go

set -e

echo "Testing ffmpeg-go project..."

# Run Go tests
echo "Running Go tests..."
go test ./ffmpeg

echo "All tests passed!"

echo "Building project..."
make build

echo "Build successful!"
echo "Binaries created in bin/ directory:"
ls -la bin/
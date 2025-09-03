#!/bin/bash

# Build script for ffmpeg-go examples

set -e

echo "Building ffmpeg-go examples..."

# Build info example
echo "Building info example..."
cd examples/info
go build -o ../../bin/info .
cd ../..

# Build decode example
echo "Building decode example..."
cd examples/decode
go build -o ../../bin/decode .
cd ../..

echo "Build complete! Binaries are in the bin/ directory."
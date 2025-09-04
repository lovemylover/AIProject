#!/bin/bash

# macOS specific build script for FLV Player C++

echo "Building FLV Player C++ for macOS..."

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "Homebrew not found. Please install Homebrew first."
    exit 1
fi

# Check if FFmpeg is installed
if ! brew list ffmpeg &> /dev/null; then
    echo "FFmpeg not found. Installing FFmpeg..."
    brew install ffmpeg
fi

# Get Homebrew prefix
HOMEBREW_PREFIX=$(brew --prefix)
echo "Homebrew prefix: $HOMEBREW_PREFIX"

# Create build directory if it doesn't exist
if [ ! -d "build" ]; then
    mkdir build
fi

# Change to build directory
cd build

# Run CMake with explicit paths
cmake .. -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_CXX_FLAGS="-I$HOMEBREW_PREFIX/include" \
    -DCMAKE_EXE_LINKER_FLAGS="-L$HOMEBREW_PREFIX/lib -lavformat -lavcodec -lavutil -lswscale"

# Check if cmake succeeded
if [ $? -ne 0 ]; then
    echo "CMake configuration failed!"
    exit 1
fi

# Build the project
make

# Check if make succeeded
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Build successful!"
echo "Executable is located at: build/flvcpp"
echo "Run with: ./flvcpp <path_to_flv_file>"
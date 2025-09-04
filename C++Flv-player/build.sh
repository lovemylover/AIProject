#!/bin/bash

# Build script for FLV Player C++

echo "Building FLV Player C++..."

# Create build directory if it doesn't exist
if [ ! -d "build" ]; then
    mkdir build
fi

# Change to build directory
cd build

# Run CMake
cmake ..

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
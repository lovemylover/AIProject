#!/bin/bash

# Complete setup and build script for FLV Player C++

echo "Setting up and building FLV Player C++..."

# Set execute permissions
chmod +x build.sh
chmod +x build.bat
chmod +x build_macos.sh
chmod +x set_permissions.sh

echo "Permissions set. Now building..."

# Use the macOS specific build script
./build_macos.sh
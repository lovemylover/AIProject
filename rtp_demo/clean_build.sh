#!/bin/bash

# Clean and rebuild the RTP demo applications
cd "$(dirname "$0")"

echo "Cleaning previous builds..."
rm -f server client

echo "Making scripts executable..."
chmod +x *.sh

echo "Building applications..."
./build.sh

if [ -f "server" ] && [ -f "client" ]; then
    echo "Build successful!"
    echo "Server size: $(du -h server | cut -f1)"
    echo "Client size: $(du -h client | cut -f1)"
else
    echo "Build failed!"
    exit 1
fi
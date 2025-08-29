#!/bin/bash

# Test building the RTP demo applications
cd "$(dirname "$0")"

echo "Testing build process..."

# Clean any existing binaries
rm -f server client

# Make scripts executable
chmod +x *.sh

# Run build script
./build.sh

# Check if binaries were created
if [ -f "server" ] && [ -f "client" ]; then
    echo "Build successful!"
    echo "Server size: $(du -h server | cut -f1)"
    echo "Client size: $(du -h client | cut -f1)"
    
    # Show binary information
    echo ""
    echo "Server info:"
    file server
    
    echo ""
    echo "Client info:"
    file client
    
    # Clean up
    rm -f server client
else
    echo "Build failed!"
    exit 1
fi
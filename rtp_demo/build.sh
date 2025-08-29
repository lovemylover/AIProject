#!/bin/bash

# Make sure we're in the right directory
cd "$(dirname "$0")"

echo "Building RTP Server..."
go build -o server server.go
if [ $? -ne 0 ]; then
    echo "Failed to build server"
    exit 1
fi

echo "Building RTP Client..."
go build -o client client.go
if [ $? -ne 0 ]; then
    echo "Failed to build client"
    exit 1
fi

echo "Build complete!"
echo "Run server with: ./server"
echo "Run client with: ./client 127.0.0.1:5004 /path/to/video.mp4"
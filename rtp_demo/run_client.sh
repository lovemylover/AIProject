#!/bin/bash

# Make sure we're in the right directory
cd "$(dirname "$0")"

if [ $# -ne 2 ]; then
    echo "Usage: $0 <server_address:port> <mp4_file>"
    echo "Example: $0 127.0.0.1:5004 test.mp4"
    exit 1
fi

SERVER_ADDR=$1
MP4_FILE=$2

if [ ! -f "$MP4_FILE" ]; then
    echo "Error: File '$MP4_FILE' not found!"
    exit 1
fi

echo "Starting RTP Client..."
echo "Streaming '$MP4_FILE' to $SERVER_ADDR"
./client $SERVER_ADDR $MP4_FILE
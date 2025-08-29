# RTP Client-Server Demo

This is a simple RTP (Real-time Transport Protocol) implementation in Go that demonstrates how to stream video data from an MP4 file over UDP.

## Features

- RTP client that reads MP4 files and streams them over UDP
- RTP server that receives and processes RTP packets
- H.264 payload handling with NALU parsing
- Support for different RTP H.264 packetization modes:
  - Single NAL Unit Packets
  - STAP-A (Single-Time Aggregation Packet type A)
  - FU-A (Fragmentation Unit type A)
- Detailed H.264 NAL Unit type identification
- Sequence number and timestamp management

## Prerequisites

- Go 1.15 or higher
- An MP4 video file for testing

## Building the Applications

1. Navigate to the project directory:
   ```
   cd rtp_demo
   ```

2. Build the applications using the build script:
   ```
   chmod +x build.sh
   ./build.sh
   ```

   Or build manually:
   ```
   go build -o server server.go
   go build -o client client.go
   ```

## Running the Applications

1. Start the RTP server:
   ```
   chmod +x run_server.sh
   ./run_server.sh
   ```
   
   By default, the server listens on port 5004. You can specify a different address:
   ```
   ./server :5005
   ```

2. In another terminal, run the RTP client:
   ```
   chmod +x run_client.sh
   ./run_client.sh 127.0.0.1:5004 /path/to/your/video.mp4
   ```

   Or run directly:
   ```
   ./client 127.0.0.1:5004 /path/to/your/video.mp4
   ```

## How It Works

### Client Side
1. The client opens the specified MP4 file
2. It reads the file in chunks (simulating video frames)
3. Each chunk is packaged into an RTP packet with appropriate headers
4. RTP packets are sent to the server via UDP

### Server Side
1. The server listens for UDP packets on the specified port
2. When a packet arrives, it parses the RTP header
3. It extracts the payload and processes it based on the payload type
4. For H.264 payloads (payload type 96), it performs detailed NALU parsing:
   - Identifies Single NAL Unit packets
   - Parses STAP-A aggregation packets
   - Handles FU-A fragmentation units
   - Displays detailed information about each NAL Unit type
5. Information about each received packet and NAL Unit is printed to the console

## RTP Header Structure

The implementation uses the standard RTP header format:

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|V=2|P|X|  CC   |M|     PT      |       sequence number         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           timestamp                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           synchronization source (SSRC) identifier            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|            contributing source (CSRC) identifiers             |
|                             ....                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

## H.264 NAL Unit Types

When processing H.264 payloads, the server identifies the following NAL Unit types:

- 1: Coded slice of a non-IDR picture
- 5: Coded slice of an IDR picture
- 6: Supplemental enhancement information (SEI)
- 7: Sequence parameter set (SPS)
- 8: Picture parameter set (PPS)
- 9: Access unit delimiter
- 10: End of sequence
- 11: End of stream
- 12: Filler data
- 24: STAP-A (Single-time aggregation packet)
- 25: STAP-B
- 26: MTAP16
- 27: MTAP24
- 28: FU-A (Fragmentation unit)
- 29: FU-B

## Limitations

This is a simplified demonstration implementation with the following limitations:

1. The MP4 parsing is very basic and doesn't properly extract H.264 NAL units
2. No support for RTCP (RTP Control Protocol)
3. No error correction or packet retransmission
4. No support for multiple streams or synchronization
5. No proper H.264 decoder (only parsing NAL Unit structure)
6. Does not reconstruct fragmented FU-A packets into complete NAL Units

## Possible Improvements

1. Implement proper MP4 parsing to extract H.264 NAL units
2. Add H.264 decoder using a library like FFmpeg
3. Implement RTCP for control and feedback
4. Add error handling and packet retransmission
5. Support for multiple simultaneous clients
6. Add support for audio streams (e.g., AAC)

## License

This project is licensed under the MIT License.
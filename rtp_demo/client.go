package main

import (
	_ "bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// RTPHeader represents the RTP header
type RTPHeader struct {
	Version        uint8  // 2 bits
	Padding        bool   // 1 bit
	Extension      bool   // 1 bit
	CSRCCount      uint8  // 4 bits
	Marker         bool   // 1 bit
	PayloadType    uint8  // 7 bits
	SequenceNumber uint16 // 16 bits
	Timestamp      uint32 // 32 bits
	SSRC           uint32 // 32 bits
}

// RTPPacket represents an RTP packet
type RTPPacket struct {
	Header  RTPHeader
	Payload []byte
}

// MP4Reader reads MP4 file and extracts video data
type MP4Reader struct {
	file      *os.File
	buffer    []byte
	pos       int64
	chunkSize int
}

// NewMP4Reader creates a new MP4 reader
func NewMP4Reader(filename string) (*MP4Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &MP4Reader{
		file:      file,
		buffer:    make([]byte, 1024),
		chunkSize: 1024,
	}, nil
}

// ReadNextChunk reads the next chunk of data from the MP4 file
func (r *MP4Reader) ReadNextChunk() ([]byte, error) {
	// For demonstration purposes, we'll just read chunks of data
	n, err := r.file.Read(r.buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if n == 0 {
		return nil, io.EOF
	}

	// Return a copy of the buffer with actual data read
	data := make([]byte, n)
	copy(data, r.buffer[:n])

	return data, nil
}

// Close closes the MP4 file
func (r *MP4Reader) Close() error {
	return r.file.Close()
}

// RTPClient represents an RTP client
type RTPClient struct {
	conn       *net.UDPConn
	remoteAddr *net.UDPAddr
	seqNum     uint16
	timestamp  uint32
	ssrc       uint32
}

// NewRTPClient creates a new RTP client
func NewRTPClient(remoteAddr string) (*RTPClient, error) {
	addr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &RTPClient{
		conn:       conn,
		remoteAddr: addr,
		seqNum:     1,
		timestamp:  0,
		ssrc:       12345, // Random SSRC
	}, nil
}

// MarshalHeader marshals the RTP header into bytes
func (h *RTPHeader) MarshalHeader() []byte {
	buf := make([]byte, 12)

	// First byte: V(2) P(1) X(1) CC(4)
	buf[0] = (h.Version << 6) | (btou8(h.Padding) << 5) | (btou8(h.Extension) << 4) | (h.CSRCCount & 0x0F)

	// Second byte: M(1) PT(7)
	buf[1] = (btou8(h.Marker) << 7) | (h.PayloadType & 0x7F)

	// Sequence number
	binary.BigEndian.PutUint16(buf[2:4], h.SequenceNumber)

	// Timestamp
	binary.BigEndian.PutUint32(buf[4:8], h.Timestamp)

	// SSRC
	binary.BigEndian.PutUint32(buf[8:12], h.SSRC)

	return buf
}

// btou8 converts bool to uint8
func btou8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// SendPacket sends an RTP packet
func (c *RTPClient) SendPacket(payload []byte) error {
	header := RTPHeader{
		Version:        2,
		Padding:        false,
		Extension:      false,
		CSRCCount:      0,
		Marker:         false,
		PayloadType:    96, // Dynamic type for H.264
		SequenceNumber: c.seqNum,
		Timestamp:      c.timestamp,
		SSRC:           c.ssrc,
	}

	headerBytes := header.MarshalHeader()
	packet := append(headerBytes, payload...)

	_, err := c.conn.Write(packet)
	if err != nil {
		return err
	}

	fmt.Printf("Sent RTP packet: Seq=%d, TS=%d, Size=%d\n", c.seqNum, c.timestamp, len(payload))

	// Update sequence number and timestamp
	c.seqNum++
	c.timestamp += 3000 // Assuming 90kHz clock rate and ~33ms per frame

	return nil
}

// Close closes the RTP client
func (c *RTPClient) Close() error {
	return c.conn.Close()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: client <server_address:port> <mp4_file>")
		os.Exit(1)
	}

	serverAddr := os.Args[1]
	mp4File := os.Args[2]

	// Create RTP client
	client, err := NewRTPClient(serverAddr)
	if err != nil {
		fmt.Printf("Failed to create RTP client: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// Open MP4 file
	reader, err := NewMP4Reader(mp4File)
	if err != nil {
		fmt.Printf("Failed to open MP4 file: %v\n", err)
		os.Exit(1)
	}
	defer reader.Close()

	fmt.Printf("Sending video stream to %s\n", serverAddr)

	// Send video frames
	frameInterval := time.Second / 30 // 30 FPS
	ticker := time.NewTicker(frameInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Read next frame/NAL unit
			chunk, err := reader.ReadNextChunk()
			if err == io.EOF {
				fmt.Println("End of video stream")
				return
			} else if err != nil {
				fmt.Printf("Error reading chunk: %v\n", err)
				continue
			}

			// Send RTP packet
			err = client.SendPacket(chunk)
			if err != nil {
				fmt.Printf("Error sending RTP packet: %v\n", err)
			}
		}
	}
}

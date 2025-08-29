package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// RTPPacketHeader represents the RTP header   12字节
type RTPPacketHeader struct {
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

// RTPServer represents an RTP server
type RTPServer struct {
	conn     *net.UDPConn
	addr     *net.UDPAddr
	received int
}

// NewRTPServer creates a new RTP server
func NewRTPServer(listenAddr string) (*RTPServer, error) {
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &RTPServer{
		conn: conn,
		addr: addr,
	}, nil
}

// UnmarshalHeader unmarshals the RTP header from bytes
func (h *RTPPacketHeader) UnmarshalHeader(data []byte) error {
	if len(data) < 12 {
		return fmt.Errorf("RTP header too short")
	}

	// First byte: V(2) P(1) X(1) CC(4)
	h.Version = (data[0] >> 6) & 0x03
	h.Padding = (data[0] >> 5 & 0x01) == 1
	h.Extension = (data[0] >> 4 & 0x01) == 1
	h.CSRCCount = data[0] & 0x0F

	// Second byte: M(1) PT(7)
	h.Marker = (data[1] >> 7) == 1
	h.PayloadType = data[1] & 0x7F

	// Sequence number
	h.SequenceNumber = binary.BigEndian.Uint16(data[2:4])

	// Timestamp
	h.Timestamp = binary.BigEndian.Uint32(data[4:8])

	// SSRC
	h.SSRC = binary.BigEndian.Uint32(data[8:12])

	return nil
}

// Start starts the RTP server
func (s *RTPServer) Start() {
	fmt.Printf("RTP server listening on %s\n", s.addr.String())

	buffer := make([]byte, 65536) // Max UDP packet size

	for {
		n, clientAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading UDP message: %v\n", err)
			continue
		}

		s.received++

		// Parse RTP header
		header := &RTPPacketHeader{}
		err = header.UnmarshalHeader(buffer[:n])
		if err != nil {
			fmt.Printf("Error parsing RTP header: %v\n", err)
			continue
		}

		// Extract payload (skip 12-byte header)
		payload := buffer[12:n]

		// Print packet info
		fmt.Printf("Received RTP packet #%d from %s: Seq=%d, TS=%d, PT=%d, Size=%d\n",
			s.received, clientAddr.String(), header.SequenceNumber, header.Timestamp, header.PayloadType, len(payload))

		// Process payload based on payload type
		s.processPayload(header, payload)
	}
}

// processPayload processes the RTP payload based on its type
func (s *RTPServer) processPayload(header *RTPPacketHeader, payload []byte) {
	switch header.PayloadType {
	case 96: // Dynamic type, assuming H.264
		fmt.Printf("  -> H.264 video frame, size: %d bytes\n", len(payload))
		// Parse H.264 NAL Units
		s.parseH264NALUs(payload)
	default:
		fmt.Printf("  -> Unknown payload type: %d\n", header.PayloadType)
	}
}

// parseH264NALUs parses H.264 NAL Units from the payload
func (s *RTPServer) parseH264NALUs(payload []byte) {
	if len(payload) == 0 {
		return
	}

	// Handle different RTP H.264 payload formats
	// Check for STAP-A (Single-Time Aggregation Packet type A)
	if payload[0] == 24 {
		s.parseSTAPA(payload)
		return
	}

	// Check for FU-A (Fragmentation Unit type A)
	if payload[0] == 28 {
		s.parseFUA(payload)
		return
	}

	// Single NAL unit packet
	s.parseSingleNALU(payload)
}

// parseSingleNALU parses a single NAL unit
func (s *RTPServer) parseSingleNALU(payload []byte) {
	if len(payload) < 1 {
		return
	}

	// First byte contains NAL unit header
	nalHeader := payload[0]
	nalType := nalHeader & 0x1F

	// Get NAL unit name
	nalTypeName := getNALUnitName(nalType)

	fmt.Printf("    -> Single NAL Unit - Type: %d (%s), Size: %d bytes\n", nalType, nalTypeName, len(payload))

	// For SPS/PPS, print additional info
	switch nalType {
	case 7: // SPS
		s.parseSPS(payload[1:])
	case 8: // PPS
		s.parsePPS(payload[1:])
	}
}

// parseSTAPA parses STAP-A packets
func (s *RTPServer) parseSTAPA(payload []byte) {
	fmt.Printf("    -> STAP-A Packet\n")

	if len(payload) < 1 {
		return
	}

	// Skip the STAP-A header (first byte)
	data := payload[1:]
	offset := 0

	for offset < len(data) {
		// Check if we have enough bytes for the length field
		if offset+2 > len(data) {
			break
		}

		// Get NAL unit length (16-bit big endian)
		nalLength := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += 2

		// Check if we have enough bytes for the NAL unit
		if offset+nalLength > len(data) {
			break
		}

		// Extract NAL unit
		nalUnit := data[offset : offset+nalLength]
		offset += nalLength

		if len(nalUnit) > 0 {
			nalHeader := nalUnit[0]
			nalType := nalHeader & 0x1F
			nalTypeName := getNALUnitName(nalType)

			fmt.Printf("      -> NAL Unit - Type: %d (%s), Size: %d bytes\n", nalType, nalTypeName, len(nalUnit))
		}
	}
}

// parseFUA parses FU-A packets
func (s *RTPServer) parseFUA(payload []byte) {
	if len(payload) < 2 {
		return
	}

	// FU-A header is the second byte
	fuHeader := payload[1]
	startBit := (fuHeader >> 7) & 0x01
	endBit := (fuHeader >> 6) & 0x01
	nalType := fuHeader & 0x1F

	nalTypeName := getNALUnitName(nalType)

	// Reconstruct the NAL header from the FU indicator (first byte) and FU header
	fuIndicator := payload[0]
	nalHeader := (fuIndicator & 0xE0) | nalType // Keep F, NRI from FU indicator, use type from FU header

	fmt.Printf("    -> FU-A Packet - NAL Type: %d (%s), Start: %d, End: %d, Size: %d bytes\n",
		nalType, nalTypeName, startBit, endBit, len(payload))

	// If this is the start of a fragmented NAL unit, show the reconstructed header
	if startBit == 1 {
		fmt.Printf("      -> Reconstructed NAL Header: 0x%02X\n", nalHeader)
	}
}

// getNALUnitName returns the name of the NAL unit type
func getNALUnitName(nalType uint8) string {
	switch nalType {
	case 1:
		return "Coded slice of a non-IDR picture"
	case 5:
		return "Coded slice of an IDR picture"
	case 6:
		return "Supplemental enhancement information (SEI)"
	case 7:
		return "Sequence parameter set (SPS)"
	case 8:
		return "Picture parameter set (PPS)"
	case 9:
		return "Access unit delimiter"
	case 10:
		return "End of sequence"
	case 11:
		return "End of stream"
	case 12:
		return "Filler data"
	case 24:
		return "STAP-A (Single-time aggregation packet)"
	case 25:
		return "STAP-B"
	case 26:
		return "MTAP16"
	case 27:
		return "MTAP24"
	case 28:
		return "FU-A (Fragmentation unit)"
	case 29:
		return "FU-B"
	default:
		return "Reserved"
	}
}

// parseSPS parses Sequence Parameter Set
func (s *RTPServer) parseSPS(spsData []byte) {
	fmt.Printf("      -> SPS Data: %d bytes\n", len(spsData))
	// In a full implementation, you would parse the SPS fields here
	// For now, we just acknowledge that we received SPS data
}

// parsePPS parses Picture Parameter Set
func (s *RTPServer) parsePPS(ppsData []byte) {
	fmt.Printf("      -> PPS Data: %d bytes\n", len(ppsData))
	// In a full implementation, you would parse the PPS fields here
	// For now, we just acknowledge that we received PPS data
}

// Close closes the RTP server
func (s *RTPServer) Close() error {
	return s.conn.Close()
}

func main() {
	listenAddr := ":5004"
	if len(os.Args) > 1 {
		listenAddr = os.Args[1]
	}

	server, err := NewRTPServer(listenAddr)
	if err != nil {
		fmt.Printf("Failed to create RTP server: %v\n", err)
		os.Exit(1)
	}
	defer server.Close()

	fmt.Printf("Starting RTP server on %s\n", listenAddr)
	server.Start()
}

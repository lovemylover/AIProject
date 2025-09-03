package main

import (
	"fmt"
	"log"
	"os"

	"ffmpeg-go/ffmpeg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <video_file>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Initialize FFmpeg
	ffmpeg.InitializeFFmpeg()
	defer ffmpeg.CleanupFFmpeg()

	// Create decoder
	decoder, err := ffmpeg.NewDecoder(filename)
	if err != nil {
		log.Fatalf("Error creating decoder: %v", err)
	}
	defer decoder.Close()

	// Decode a few frames
	for i := 0; i < 5; i++ {
		frame, err := decoder.DecodeNextFrame()
		if err != nil {
			log.Printf("Error decoding frame %d: %v", i, err)
			break
		}

		fmt.Printf("Decoded frame %d: %dx%d\n", i, frame.Width, frame.Height)
		fmt.Printf("Frame data size: %d bytes\n", len(frame.Data))
	}
}
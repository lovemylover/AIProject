package main

import (
	"fmt"
	"log"
	"os"

	"ffmpeg-go/ffmpeg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <media_file>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Initialize FFmpeg
	ffmpeg.InitializeFFmpeg()
	defer ffmpeg.CleanupFFmpeg()

	// Get media information
	info, err := ffmpeg.GetMediaInfo(filename)
	if err != nil {
		log.Fatalf("Error getting media info: %v", err)
	}

	fmt.Printf("Media Information for: %s\n", filename)
	fmt.Printf("Duration: %s\n", info.Duration)
	fmt.Printf("Bitrate: %d kb/s\n", info.Bitrate)
	fmt.Printf("Number of streams: %d\n", info.Streams)
}
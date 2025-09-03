package main

import (
	"fmt"
	"log"
	"os"

	"ffmpeg-go/ffmpeg"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Initialize FFmpeg
	ffmpeg.InitializeFFmpeg()
	defer ffmpeg.CleanupFFmpeg()

	// Get media information
	info, err := ffmpeg.GetMediaInfo(inputFile)
	if err != nil {
		log.Fatalf("Error getting media info: %v", err)
	}

	fmt.Printf("Transcoding %s to %s\n", inputFile, outputFile)
	fmt.Printf("Input file duration: %s\n", info.Duration)
	fmt.Printf("Input file bitrate: %d kb/s\n", info.Bitrate)

	// Note: A full transcoding implementation would be more complex
	// This is just a placeholder to show the structure
	fmt.Println("Transcoding would happen here...")
	fmt.Println("Transcoding completed!")
}
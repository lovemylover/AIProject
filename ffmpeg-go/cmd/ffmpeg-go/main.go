package main

import (
	"fmt"
	"log"
	"os"

	"ffmpeg-go/ffmpeg"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Initialize FFmpeg
	ffmpeg.InitializeFFmpeg()
	defer ffmpeg.CleanupFFmpeg()

	switch command {
	case "info":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ffmpeg-go info <media_file>")
			os.Exit(1)
		}
		handleInfoCommand(os.Args[2])
	case "version":
		handleVersionCommand()
	case "validate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ffmpeg-go validate <media_file>")
			os.Exit(1)
		}
		handleValidateCommand(os.Args[2])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("ffmpeg-go - FFmpeg Go bindings CLI tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  ffmpeg-go info <media_file>     Get media file information")
	fmt.Println("  ffmpeg-go validate <media_file> Check if file is valid media")
	fmt.Println("  ffmpeg-go version               Show FFmpeg version")
	fmt.Println("")
}

func handleInfoCommand(filename string) {
	info, err := ffmpeg.GetMediaInfo(filename)
	if err != nil {
		log.Fatalf("Error getting media info: %v", err)
	}

	fmt.Printf("Media File: %s\n", filename)
	fmt.Printf("Duration: %s\n", info.Duration)
	fmt.Printf("Bitrate: %d kb/s\n", info.Bitrate)
	fmt.Printf("Number of streams: %d\n", info.Streams)
	
	if len(info.StreamInfo) > 0 {
		fmt.Println("\nStream details:")
		for _, stream := range info.StreamInfo {
			fmt.Printf("  Stream #%d: %s (%s)\n", stream.Index, stream.CodecType, stream.CodecName)
		}
	}
}

func handleVersionCommand() {
	version := ffmpeg.GetFFmpegVersion()
	fmt.Printf("FFmpeg Version: %s\n", version)
}

func handleValidateCommand(filename string) {
	if ffmpeg.IsValidMediaFile(filename) {
		fmt.Printf("'%s' is a valid media file.\n", filename)
	} else {
		fmt.Printf("'%s' is NOT a valid media file.\n", filename)
	}
}
package ffmpeg

import (
	"testing"
)

func TestInitializeFFmpeg(t *testing.T) {
	// This test just verifies that the FFmpeg initialization functions exist
	// and can be called without panicking
	InitializeFFmpeg()
	CleanupFFmpeg()
}

func TestMediaInfoStruct(t *testing.T) {
	info := &MediaInfo{
		Duration: "10.5 seconds",
		Bitrate:  128,
		Streams:  2,
	}

	if info.Duration != "10.5 seconds" {
		t.Errorf("Expected duration '10.5 seconds', got '%s'", info.Duration)
	}

	if info.Bitrate != 128 {
		t.Errorf("Expected bitrate 128, got %d", info.Bitrate)
	}

	if info.Streams != 2 {
		t.Errorf("Expected 2 streams, got %d", info.Streams)
	}
}
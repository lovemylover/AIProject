package ffmpeg

/*
#cgo pkg-config: libavutil
#include <libavutil/error.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
)

// AVError converts an FFmpeg error code to a Go error with descriptive message
func AVError(code int) error {
	errbuf := make([]C.char, 1024)
	C.av_strerror(C.int(code), &errbuf[0], C.size_t(len(errbuf)))
	goErrStr := C.GoString(&errbuf[0])
	return fmt.Errorf("ffmpeg error %d: %s", code, goErrStr)
}

// IsValidMediaFile checks if a file is a valid media file that can be processed
func IsValidMediaFile(filename string) bool {
	_, err := GetMediaInfo(filename)
	return err == nil
}

// GetFFmpegVersion returns the version of the underlying FFmpeg libraries
func GetFFmpegVersion() string {
	// For now, we'll return a placeholder
	// In a real implementation, you would query the actual FFmpeg version
	return "libavformat: " + "58.x.x" + ", libavcodec: " + "58.x.x"
}

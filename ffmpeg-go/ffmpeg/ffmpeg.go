// Package ffmpeg provides Go bindings for FFmpeg libraries
package ffmpeg

/*
#cgo pkg-config: libavformat libavcodec libavutil libswscale
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>
#include <libswscale/swscale.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// StreamInfo holds information about a media stream
type StreamInfo struct {
	Index     int
	CodecType string
	CodecName string
	Duration  string
	Bitrate   int64
}

// MediaInfo holds basic media file information
type MediaInfo struct {
	Duration   string
	Bitrate    int
	Streams    int
	StreamInfo []StreamInfo
}

// GetMediaInfo retrieves basic information about a media file
func GetMediaInfo(filename string) (*MediaInfo, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	// Open the media file
	var fmtCtx *C.AVFormatContext
	ret := C.avformat_open_input(&fmtCtx, cFilename, nil, nil)
	if ret < 0 {
		return nil, fmt.Errorf("could not open file: %s", filename)
	}
	defer C.avformat_close_input(&fmtCtx)

	// Retrieve stream information
	ret = C.avformat_find_stream_info(fmtCtx, nil)
	if ret < 0 {
		return nil, fmt.Errorf("could not find stream information")
	}

	// Create MediaInfo struct
	info := &MediaInfo{
		Duration: fmt.Sprintf("%.2f seconds", float64(fmtCtx.duration)/float64(C.AV_TIME_BASE)),
		Bitrate:  int(fmtCtx.bit_rate / 1000), // kbps
		Streams:  int(fmtCtx.nb_streams),
	}

	return info, nil
}

// InitializeFFmpeg initializes the FFmpeg libraries
func InitializeFFmpeg() {
	//C.av_register_all()
	C.avformat_network_init()
}

// CleanupFFmpeg cleans up the FFmpeg libraries
func CleanupFFmpeg() {
	C.avformat_network_deinit()
}

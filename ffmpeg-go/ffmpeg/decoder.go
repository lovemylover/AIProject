package ffmpeg

/*
#cgo pkg-config: libavformat libavcodec libavutil
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/imgutils.h>
#include <libavutil/samplefmt.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Frame represents a decoded audio or video frame
type Frame struct {
	Data     []byte
	Width    int
	Height   int
	Channels int
	Samples  int
}

// Decoder handles media decoding
type Decoder struct {
	fmtCtx    *C.AVFormatContext
	codecCtx  *C.AVCodecContext
	streamIdx int
}

// NewDecoder creates a new decoder for the given media file
func NewDecoder(filename string) (*Decoder, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	// Open the media file
	var fmtCtx *C.AVFormatContext
	ret := C.avformat_open_input(&fmtCtx, cFilename, nil, nil)
	if ret < 0 {
		return nil, fmt.Errorf("could not open file: %s", filename)
	}

	// Retrieve stream information
	ret = C.avformat_find_stream_info(fmtCtx, nil)
	if ret < 0 {
		C.avformat_close_input(&fmtCtx)
		return nil, fmt.Errorf("could not find stream information")
	}

	// Find the first video stream
	streamIdx := int(C.av_find_best_stream(fmtCtx, C.AVMEDIA_TYPE_VIDEO, -1, -1, nil, 0))
	if streamIdx < 0 {
		C.avformat_close_input(&fmtCtx)
		return nil, fmt.Errorf("could not find video stream")
	}

	// Get codec context
	// Access the stream using pointer arithmetic
	stream := *(**C.AVStream)(unsafe.Pointer(uintptr(unsafe.Pointer(fmtCtx.streams)) + 
		uintptr(streamIdx)*unsafe.Sizeof(uintptr(0))))
	codecParams := stream.codecpar
	codec := C.avcodec_find_decoder(codecParams.codec_id)
	if codec == nil {
		C.avformat_close_input(&fmtCtx)
		return nil, fmt.Errorf("unsupported codec")
	}

	// Allocate codec context
	codecCtx := C.avcodec_alloc_context3(codec)
	if codecCtx == nil {
		C.avformat_close_input(&fmtCtx)
		return nil, fmt.Errorf("failed to allocate codec context")
	}

	// Copy codec parameters to context
	ret = C.avcodec_parameters_to_context(codecCtx, codecParams)
	if ret < 0 {
		C.avcodec_free_context(&codecCtx)
		C.avformat_close_input(&fmtCtx)
		return nil, fmt.Errorf("failed to copy codec parameters")
	}

	// Open codec
	ret = C.avcodec_open2(codecCtx, codec, nil)
	if ret < 0 {
		C.avcodec_free_context(&codecCtx)
		C.avformat_close_input(&fmtCtx)
		return nil, fmt.Errorf("failed to open codec")
	}

	return &Decoder{
		fmtCtx:    fmtCtx,
		codecCtx:  codecCtx,
		streamIdx: streamIdx,
	}, nil
}

// DecodeNextFrame decodes the next frame from the media
func (d *Decoder) DecodeNextFrame() (*Frame, error) {
	packet := C.av_packet_alloc()
	defer C.av_packet_unref(packet)

	// Read frame
	ret := C.av_read_frame(d.fmtCtx, packet)
	if ret < 0 {
		return nil, fmt.Errorf("end of file or error reading frame")
	}

	// Only process packets from the video stream
	if int(packet.stream_index) != d.streamIdx {
		return nil, fmt.Errorf("not a video packet")
	}

	// Send packet to decoder
	ret = C.avcodec_send_packet(d.codecCtx, packet)
	if ret < 0 {
		return nil, fmt.Errorf("error sending packet to decoder")
	}

	// Receive frame from decoder
	frame := C.av_frame_alloc()
	defer C.av_frame_free(&frame)

	ret = C.avcodec_receive_frame(d.codecCtx, frame)
	if ret < 0 {
		return nil, fmt.Errorf("error receiving frame from decoder")
	}

	// Convert frame data to Go byte slice
	width := int(frame.width)
	height := int(frame.height)
	
	// For simplicity, we'll just return the raw frame data
	// In a real implementation, you might want to convert to a specific format
	frameData := make([]byte, width*height*3) // Assuming RGB24

	return &Frame{
		Data:   frameData,
		Width:  width,
		Height: height,
	}, nil
}

// Close releases all resources used by the decoder
func (d *Decoder) Close() {
	if d.codecCtx != nil {
		C.avcodec_free_context(&d.codecCtx)
	}
	if d.fmtCtx != nil {
		C.avformat_close_input(&d.fmtCtx)
	}
}
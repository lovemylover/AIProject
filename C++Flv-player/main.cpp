#include <iostream>
#include <string>
#include <thread>
#include <chrono>

extern "C" {
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libswscale/swscale.h>
#include <libavutil/imgutils.h>
}

#include <SDL2/SDL.h>

#ifdef _WIN32
#include <windows.h>
#else
#include <unistd.h>
#endif

class FLVPlayer {
private:
    AVFormatContext* formatContext;
    AVCodecContext* codecContext;
    const AVCodec* decoder;
    AVFrame* frame;
    AVFrame* frameRGB;
    AVPacket* packet;
    SwsContext* swsContext;
    uint8_t* buffer;
    int videoStreamIndex;
    
    // SDL2 variables for video rendering
    SDL_Window* window;
    SDL_Renderer* renderer;
    SDL_Texture* texture;
    bool sdlInitialized;
    
    // Flag to control playback
    bool shouldStop;
    
public:
    FLVPlayer() : formatContext(nullptr), codecContext(nullptr), decoder(nullptr),
                  frame(nullptr), frameRGB(nullptr), packet(nullptr), 
                  swsContext(nullptr), buffer(nullptr), videoStreamIndex(-1),
                  window(nullptr), renderer(nullptr), texture(nullptr), sdlInitialized(false),
                  shouldStop(false) {}
    
    ~FLVPlayer() {
        cleanup();
        cleanupSDL();
    }
    
    bool initialize(const std::string& filename) {
        // Open the file
        if (avformat_open_input(&formatContext, filename.c_str(), nullptr, nullptr) != 0) {
            std::cerr << "Could not open file: " << filename << std::endl;
            return false;
        }
        
        // Retrieve stream information
        if (avformat_find_stream_info(formatContext, nullptr) < 0) {
            std::cerr << "Could not find stream information" << std::endl;
            return false;
        }
        
        // Find the first video stream
        const AVCodec **decoderPtr = &decoder;
        videoStreamIndex = av_find_best_stream(formatContext, AVMEDIA_TYPE_VIDEO, -1, -1, decoderPtr, 0);
        if (videoStreamIndex < 0) {
            std::cerr << "Could not find video stream" << std::endl;
            return false;
        }
        
        // Get codec context
        AVStream* stream = formatContext->streams[videoStreamIndex];
        codecContext = avcodec_alloc_context3(decoder);
        if (!codecContext) {
            std::cerr << "Could not allocate codec context" << std::endl;
            return false;
        }
        
        // Copy codec parameters to context
        if (avcodec_parameters_to_context(codecContext, stream->codecpar) < 0) {
            std::cerr << "Could not copy codec parameters to context" << std::endl;
            return false;
        }
        
        // Open codec
        if (avcodec_open2(codecContext, decoder, nullptr) < 0) {
            std::cerr << "Could not open codec" << std::endl;
            return false;
        }
        
        // Allocate frames
        frame = av_frame_alloc();
        frameRGB = av_frame_alloc();
        if (!frame || !frameRGB) {
            std::cerr << "Could not allocate frames" << std::endl;
            return false;
        }
        
        // Determine required buffer size and allocate buffer
        int numBytes = av_image_get_buffer_size(AV_PIX_FMT_RGB24, codecContext->width, codecContext->height, 1);
        buffer = (uint8_t*)av_malloc(numBytes * sizeof(uint8_t));
        
        // Assign appropriate parts of buffer to image planes in frameRGB
        av_image_fill_arrays(frameRGB->data, frameRGB->linesize, buffer, AV_PIX_FMT_RGB24, 
                            codecContext->width, codecContext->height, 1);
        
        // Initialize SWS context for software scaling
        swsContext = sws_getContext(codecContext->width, codecContext->height, codecContext->pix_fmt,
                                   codecContext->width, codecContext->height, AV_PIX_FMT_RGB24,
                                   SWS_BILINEAR, nullptr, nullptr, nullptr);
        
        // Allocate packet
        packet = av_packet_alloc();
        if (!packet) {
            std::cerr << "Could not allocate packet" << std::endl;
            return false;
        }
        
        std::cout << "Successfully initialized player for file: " << filename << std::endl;
        std::cout << "Video resolution: " << codecContext->width << "x" << codecContext->height << std::endl;
        std::cout << "Frame rate: " << av_q2d(stream->avg_frame_rate) << " fps" << std::endl;
        
        // Initialize SDL for video rendering
        if (!initializeSDL()) {
            std::cerr << "Warning: Could not initialize SDL for video rendering" << std::endl;
            // We can still play without video rendering
        }
        
        return true;
    }
    
    bool initializeSDL() {
        // Initialize SDL
        if (SDL_Init(SDL_INIT_VIDEO) < 0) {
            std::cerr << "SDL could not initialize! SDL_Error: " << SDL_GetError() << std::endl;
            return false;
        }
        
        // Create window
        window = SDL_CreateWindow("FLV Player", 
                                 SDL_WINDOWPOS_UNDEFINED, 
                                 SDL_WINDOWPOS_UNDEFINED, 
                                 codecContext->width, 
                                 codecContext->height, 
                                 SDL_WINDOW_SHOWN);
        if (window == nullptr) {
            std::cerr << "Window could not be created! SDL_Error: " << SDL_GetError() << std::endl;
            SDL_Quit();
            return false;
        }
        
        // Create renderer
        renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED);
        if (renderer == nullptr) {
            std::cerr << "Renderer could not be created! SDL_Error: " << SDL_GetError() << std::endl;
            SDL_DestroyWindow(window);
            SDL_Quit();
            return false;
        }
        
        // Create texture for rendering
        texture = SDL_CreateTexture(renderer, 
                                   SDL_PIXELFORMAT_RGB24, 
                                   SDL_TEXTUREACCESS_STREAMING, 
                                   codecContext->width, 
                                   codecContext->height);
        if (texture == nullptr) {
            std::cerr << "Texture could not be created! SDL_Error: " << SDL_GetError() << std::endl;
            SDL_DestroyRenderer(renderer);
            SDL_DestroyWindow(window);
            SDL_Quit();
            return false;
        }
        
        sdlInitialized = true;
        std::cout << "SDL initialized successfully for video rendering" << std::endl;
        return true;
    }
    
    bool play() {
        if (!formatContext || !codecContext) {
            std::cerr << "Player not initialized" << std::endl;
            return false;
        }
        
        std::cout << "Starting playback..." << std::endl;
        
        int frameCount = 0;
        while (av_read_frame(formatContext, packet) >= 0 && !shouldStop) {
            if (packet->stream_index == videoStreamIndex) {
                // Decode video frame
                int response = avcodec_send_packet(codecContext, packet);
                if (response < 0) {
                    std::cerr << "Error sending packet for decoding" << std::endl;
                    av_packet_unref(packet);
                    continue;
                }
                
                while (response >= 0 && !shouldStop) {
                    response = avcodec_receive_frame(codecContext, frame);
                    if (response == AVERROR(EAGAIN) || response == AVERROR_EOF) {
                        break;
                    } else if (response < 0) {
                        std::cerr << "Error during decoding" << std::endl;
                        break;
                    }
                    
                    // Convert the image from its native format to RGB
                    sws_scale(swsContext, (uint8_t const* const*)frame->data, frame->linesize, 0,
                             codecContext->height, frameRGB->data, frameRGB->linesize);
                    
                    // Process the decoded frame (in a real application, you would render it)
                    processFrame(frameRGB, frameCount);
                    
                    frameCount++;
                    
                    // Control playback speed (simulate real-time playback)
                    AVRational frameRate = formatContext->streams[videoStreamIndex]->avg_frame_rate;
                    if (frameRate.num && frameRate.den) {
                        int delayMs = (1000 * frameRate.den) / frameRate.num;
                        #ifdef _WIN32
                        Sleep(delayMs);
                        #else
                        usleep(delayMs * 1000);
                        #endif
                    }
                }
            }
            av_packet_unref(packet);
        }
        
        if (shouldStop) {
            std::cout << "Playback stopped by user" << std::endl;
        } else {
            std::cout << "Playback finished. Total frames: " << frameCount << std::endl;
        }
        return true;
    }
    
    void processFrame(AVFrame* frame, int frameNumber) {
        // Render frame to SDL window if SDL is initialized
        if (sdlInitialized) {
            // Update texture with new frame data
            SDL_UpdateTexture(texture, nullptr, frame->data[0], frame->linesize[0]);
            
            // Clear renderer
            SDL_RenderClear(renderer);
            
            // Copy texture to renderer
            SDL_RenderCopy(renderer, texture, nullptr, nullptr);
            
            // Present renderer
            SDL_RenderPresent(renderer);
            
            // Handle SDL events (for window closing)
            SDL_Event event;
            while (SDL_PollEvent(&event)) {
                if (event.type == SDL_QUIT) {
                    std::cout << "Window closed by user"  << std::endl;
                    shouldStop = true;
                }
            }
        }
        
        // Print frame information periodically
        if (frameNumber % 30 == 0) {  // Print every 30 frames
            std::cout << "Processing frame " << frameNumber << std::endl;
        }
    }
    
    void cleanup() {
        if (packet) {
            av_packet_free(&packet);
        }
        if (swsContext) {
            sws_freeContext(swsContext);
            swsContext = nullptr;
        }
        if (buffer) {
            av_free(buffer);
            buffer = nullptr;
        }
        if (frameRGB) {
            av_frame_free(&frameRGB);
        }
        if (frame) {
            av_frame_free(&frame);
        }
        if (codecContext) {
            avcodec_free_context(&codecContext);
        }
        if (formatContext) {
            avformat_close_input(&formatContext);
        }
    }
    
    void cleanupSDL() {
        if (texture) {
            SDL_DestroyTexture(texture);
            texture = nullptr;
        }
        if (renderer) {
            SDL_DestroyRenderer(renderer);
            renderer = nullptr;
        }
        if (window) {
            SDL_DestroyWindow(window);
            window = nullptr;
        }
        if (sdlInitialized) {
            SDL_Quit();
            sdlInitialized = false;
        }
    }
};

int main(int argc, char* argv[]) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <flv_file>" << std::endl;
        std::cerr << "Example: " << argv[0] << " sample.flv" << std::endl;
        return 1;
    }
    
    std::string filename = argv[1];
    
    // Initialize FFmpeg libraries
    avformat_network_init();
    
    // Create player instance
    FLVPlayer player;
    
    // Initialize player with the specified file
    if (!player.initialize(filename)) {
        std::cerr << "Failed to initialize player" << std::endl;
        return 1;
    }
    
    // Start playback
    if (!player.play()) {
        std::cerr << "Failed to play video" << std::endl;
        return 1;
    }
    
    std::cout << "Program finished successfully" << std::endl;
    return 0;
}
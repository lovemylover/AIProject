cmake_minimum_required(VERSION 3.10)
project(FLVPlayer)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Set build type to Release if not specified
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Release)
endif()

# Special handling for macOS with Homebrew
if(APPLE)
    # Check if Homebrew is installed
    find_program(HOMEBREW brew)
    if(HOMEBREW)
        # Get Homebrew prefix
        execute_process(
            COMMAND brew --prefix
            OUTPUT_VARIABLE HOMEBREW_PREFIX
            OUTPUT_STRIP_TRAILING_WHITESPACE
        )
        
        if(HOMEBREW_PREFIX)
            # Add Homebrew include and lib paths
            include_directories(${HOMEBREW_PREFIX}/include)
            link_directories(${HOMEBREW_PREFIX}/lib)
        endif()
    endif()
endif()

# Find required packages
find_package(PkgConfig REQUIRED)

# Try to find SDL2
pkg_check_modules(SDL2 REQUIRED sdl2)

# Try to find FFmpeg libraries
# First try pkg-config
if(PKG_CONFIG_FOUND)
    pkg_check_modules(AVCODEC libavcodec)
    pkg_check_modules(AVFORMAT libavformat)
    pkg_check_modules(AVUTIL libavutil)
    pkg_check_modules(SWSCALE libswscale)
endif()

# If pkg-config didn't work, try manual search with more paths
if(NOT AVCODEC_FOUND)
    find_path(AVCODEC_INCLUDE_DIR 
        NAMES libavcodec/avcodec.h
        PATHS /usr/local/include /opt/local/include /opt/homebrew/include /usr/include
        NO_DEFAULT_PATH
    )
    
    find_library(AVCODEC_LIBRARY 
        NAMES avcodec
        PATHS /usr/local/lib /opt/local/lib /opt/homebrew/lib /usr/lib
        NO_DEFAULT_PATH
    )
    
    if(AVCODEC_INCLUDE_DIR AND AVCODEC_LIBRARY)
        set(AVCODEC_FOUND TRUE)
        set(AVCODEC_LIBRARIES ${AVCODEC_LIBRARY})
        set(AVCODEC_INCLUDE_DIRS ${AVCODEC_INCLUDE_DIR})
    endif()
endif()

if(NOT AVFORMAT_FOUND)
    find_path(AVFORMAT_INCLUDE_DIR 
        NAMES libavformat/avformat.h
        PATHS /usr/local/include /opt/local/include /opt/homebrew/include /usr/include
        NO_DEFAULT_PATH
    )
    
    find_library(AVFORMAT_LIBRARY 
        NAMES avformat
        PATHS /usr/local/lib /opt/local/lib /opt/homebrew/lib /usr/lib
        NO_DEFAULT_PATH
    )
    
    if(AVFORMAT_INCLUDE_DIR AND AVFORMAT_LIBRARY)
        set(AVFORMAT_FOUND TRUE)
        set(AVFORMAT_LIBRARIES ${AVFORMAT_LIBRARY})
        set(AVFORMAT_INCLUDE_DIRS ${AVFORMAT_INCLUDE_DIR})
    endif()
endif()

if(NOT AVUTIL_FOUND)
    find_path(AVUTIL_INCLUDE_DIR 
        NAMES libavutil/avutil.h
        PATHS /usr/local/include /opt/local/include /opt/homebrew/include /usr/include
        NO_DEFAULT_PATH
    )
    
    find_library(AVUTIL_LIBRARY 
        NAMES avutil
        PATHS /usr/local/lib /opt/local/lib /opt/homebrew/lib /usr/lib
        NO_DEFAULT_PATH
    )
    
    if(AVUTIL_INCLUDE_DIR AND AVUTIL_LIBRARY)
        set(AVUTIL_FOUND TRUE)
        set(AVUTIL_LIBRARIES ${AVUTIL_LIBRARY})
        set(AVUTIL_INCLUDE_DIRS ${AVUTIL_INCLUDE_DIR})
    endif()
endif()

if(NOT SWSCALE_FOUND)
    find_path(SWSCALE_INCLUDE_DIR 
        NAMES libswscale/swscale.h
        PATHS /usr/local/include /opt/local/include /opt/homebrew/include /usr/include
        NO_DEFAULT_PATH
    )
    
    find_library(SWSCALE_LIBRARY 
        NAMES swscale
        PATHS /usr/local/lib /opt/local/lib /opt/homebrew/lib /usr/lib
        NO_DEFAULT_PATH
    )
    
    if(SWSCALE_INCLUDE_DIR AND SWSCALE_LIBRARY)
        set(SWSCALE_FOUND TRUE)
        set(SWSCALE_LIBRARIES ${SWSCALE_LIBRARY})
        set(SWSCALE_INCLUDE_DIRS ${SWSCALE_INCLUDE_DIR})
    endif()
endif()

# Check if all required libraries were found
if(NOT AVCODEC_FOUND OR NOT AVFORMAT_FOUND OR NOT AVUTIL_FOUND OR NOT SWSCALE_FOUND)
    message(STATUS "FFmpeg libraries not found through standard methods, trying alternative approach...")
    
    # Try linking with full paths on macOS
    if(APPLE)
        set(FFMPEG_LIBS_FOUND TRUE)
        set(AVCODEC_LIBRARIES avcodec)
        set(AVFORMAT_LIBRARIES avformat)
        set(AVUTIL_LIBRARIES avutil)
        set(SWSCALE_LIBRARIES swscale)
    else()
        message(FATAL_ERROR "FFmpeg libraries not found. Please install FFmpeg development packages.")
    endif()
else()
    message(STATUS "Found FFmpeg libraries:")
    message(STATUS "  avcodec: ${AVCODEC_LIBRARIES}")
    message(STATUS "  avformat: ${AVFORMAT_LIBRARIES}")
    message(STATUS "  avutil: ${AVUTIL_LIBRARIES}")
    message(STATUS "  swscale: ${SWSCALE_LIBRARIES}")
endif()

# Include directories
if(AVCODEC_INCLUDE_DIRS)
    include_directories(${AVCODEC_INCLUDE_DIRS})
endif()

if(AVFORMAT_INCLUDE_DIRS)
    include_directories(${AVFORMAT_INCLUDE_DIRS})
endif()

if(AVUTIL_INCLUDE_DIRS)
    include_directories(${AVUTIL_INCLUDE_DIRS})
endif()

if(SWSCALE_INCLUDE_DIRS)
    include_directories(${SWSCALE_INCLUDE_DIRS})
endif()

# Include SDL2 directories
if(SDL2_INCLUDE_DIRS)
    include_directories(${SDL2_INCLUDE_DIRS})
endif()

# Add executable
add_executable(flvcpp main.cpp)

# Link libraries
target_link_libraries(flvcpp 
    ${AVCODEC_LIBRARIES}
    ${AVFORMAT_LIBRARIES}
    ${AVUTIL_LIBRARIES}
    ${SWSCALE_LIBRARIES}
    ${SDL2_LIBRARIES}
)

# Additional linker flags for macOS
if(APPLE)
    target_link_libraries(flvcpp "-lavformat -lavcodec -lavutil -lswscale -lSDL2")
endif()

# On Windows, we might need additional libraries
if(WIN32)
    target_link_libraries(flvcpp ws2_32)
endif()

# Compiler-specific options
if(MSVC)
    # Force to always compile with W4
    if(CMAKE_CXX_FLAGS MATCHES "/W[0-4]")
        string(REGEX REPLACE "/W[0-4]" "/W4" CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS}")
    else()
        set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /W4")
    endif()
elseif(CMAKE_COMPILER_IS_GNUCC OR CMAKE_COMPILER_IS_GNUCXX)
    # Update if necessary
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wno-long-long -pedantic")
endif()
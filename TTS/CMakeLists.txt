cmake_minimum_required(VERSION 3.10)
project(TTS)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Find espeak-ng library
find_library(ESPEAK_NG_LIBRARY NAMES espeak-ng)
find_path(ESPEAK_NG_INCLUDE_DIR NAMES espeak-ng/espeak_ng.h)

# If not found, try alternative names
if(NOT ESPEAK_NG_LIBRARY)
    find_library(ESPEAK_NG_LIBRARY NAMES espeak_ng espeak-ng libespeak-ng)
endif()

if(NOT ESPEAK_NG_INCLUDE_DIR)
    find_path(ESPEAK_NG_INCLUDE_DIR NAMES espeak-ng/espeak_ng.h PATHS /usr/local/include /usr/include)
endif()

# Check if espeak-ng is found
if(ESPEAK_NG_LIBRARY AND ESPEAK_NG_INCLUDE_DIR)
    message(STATUS "Found espeak-ng: ${ESPEAK_NG_LIBRARY}")
    message(STATUS "Found espeak-ng headers: ${ESPEAK_NG_INCLUDE_DIR}")
else()
    message(WARNING "espeak-ng not found, trying espeak instead")
    
    # Try to find espeak as fallback
    find_library(ESPEAK_LIBRARY NAMES espeak)
    find_path(ESPEAK_INCLUDE_DIR NAMES espeak/speak_lib.h)
    
    if(ESPEAK_LIBRARY AND ESPEAK_INCLUDE_DIR)
        message(STATUS "Found espeak: ${ESPEAK_LIBRARY}")
        message(STATUS "Found espeak headers: ${ESPEAK_INCLUDE_DIR}")
        set(USE_ESPEAK TRUE)
    else()
        message(FATAL_ERROR "Neither espeak-ng nor espeak found. Please install one of them:
            Ubuntu/Debian: sudo apt-get install libespeak-ng1 libespeak-ng-dev
            macOS: brew install espeak-ng
            Windows: Download from https://github.com/espeak-ng/espeak-ng")
    endif()
endif()

# Include directories
if(USE_ESPEAK)
    include_directories(${ESPEAK_INCLUDE_DIR})
else()
    include_directories(${ESPEAK_NG_INCLUDE_DIR})
endif()

# Source files
set(SOURCES
    src/main.cpp
    src/TTSProcessor.cpp
)

# Header files
set(HEADERS
    include/TTSProcessor.h
)

# Create executable
add_executable(${PROJECT_NAME} ${SOURCES} ${HEADERS})

# Link libraries
if(USE_ESPEAK)
    target_link_libraries(${PROJECT_NAME} ${ESPEAK_LIBRARY})
    target_compile_definitions(${PROJECT_NAME} PRIVATE USE_ESPEAK)
else()
    target_link_libraries(${PROJECT_NAME} ${ESPEAK_NG_LIBRARY})
endif()

# Install target
install(TARGETS ${PROJECT_NAME} DESTINATION bin)

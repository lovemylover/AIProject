# Text-to-Speech (TTS) Program in C++

A C++ implementation of a text-to-speech program using eSpeak-ng library. This program can convert text to speech and save the output as WAV files.

## Features

- Convert text to speech using eSpeak-ng
- Save output as WAV audio files
- Adjustable speech parameters (rate, pitch, volume)
- Multiple voice support
- Command-line interface
- Cross-platform compatibility (Linux, macOS, Windows)

## Prerequisites

Before building the project, ensure you have the following dependencies installed:

### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install cmake libespeak-ng1 libespeak-ng-dev build-essential
```

### macOS
```bash
# Install Homebrew if not already installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install dependencies
brew install cmake espeak-ng
```

### Windows
1. Install CMake from https://cmake.org/download/
2. Download eSpeak-ng from https://github.com/espeak-ng/espeak-ng
3. Follow the installation instructions for eSpeak-ng on Windows

## Building the Project

1. Clone or download this repository
2. Navigate to the TTS directory
3. Run the setup script:
   ```bash
   ./setup.sh
   ```
4. Build the project:
   ```bash
   ./build.sh
   ```

Alternatively, you can build manually:
```bash
mkdir build
cd build
cmake ..
make
```

## Usage

After building, you can run the TTS program with various options:

### Basic usage
```bash
./run.sh "Hello, world!"
```

### Save to WAV file
```bash
./run.sh --output hello.wav "Hello, world!"
```

### Adjust speech parameters
```bash
./run.sh --rate 200 --pitch 75 --volume 150 "Welcome to TTS program"
```

### Read text from file
```bash
./run.sh --file input.txt
```

### List available voices
```bash
./run.sh --list-voices
```

### Set specific voice
```bash
./run.sh --voice en-us "Hello from American English voice"
```

## Command Line Options

```
--help, -h          Show help message
--rate <wpm>        Set speech rate (words per minute, default: 175)
--pitch <value>     Set pitch (0-100, default: 50)
--volume <value>    Set volume (0-200, default: 100)
--voice <name>      Set voice (default: default)
--list-voices       List available voices
--file <filename>   Read text from file
--output <filename> Save output to WAV file
```

## Project Structure

```
TTS/
├── CMakeLists.txt     # Build configuration
├── README.md          # This file
├── build.sh           # Build script
├── run.sh             # Run script
├── setup.sh           # Setup script
├── include/           # Header files
│   └── TTSProcessor.h # TTS processor class
└── src/               # Source files
    ├── main.cpp       # Main program
    └── TTSProcessor.cpp # TTS processor implementation
```

## How It Works

1. The program uses the eSpeak-ng library for text-to-speech synthesis
2. Audio data is captured through a callback function during synthesis
3. The captured audio data is formatted as a WAV file
4. The WAV file can be saved to disk or processed further

## Troubleshooting

### "eSpeak-ng not found" error
Make sure you have installed the eSpeak-ng development libraries:
- Ubuntu/Debian: `sudo apt-get install libespeak-ng-dev`
- macOS: `brew install espeak-ng`

### No audio output
When running without the `--output` flag, the program generates WAV data but doesn't play it through speakers. To hear the audio, either:
1. Save to a file and play it with an audio player
2. Modify the code to integrate with an audio playback library

### Voice not found
Use `--list-voices` to see available voices on your system. The default voice is "default".

## Extending the Program

You can extend this program by:
1. Adding support for different audio formats (MP3, OGG, etc.)
2. Integrating with audio playback libraries for real-time playback
3. Adding support for SSML (Speech Synthesis Markup Language)
4. Implementing a GUI interface
5. Adding support for phoneme generation

## License

This project is licensed under the MIT License.
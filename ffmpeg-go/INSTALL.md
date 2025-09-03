# Installation Guide

## Prerequisites

Before using this library, you need to install FFmpeg development libraries.

### Ubuntu/Debian

```bash
sudo apt update
sudo apt install ffmpeg libavformat-dev libavcodec-dev libavutil-dev libswscale-dev pkg-config
```

### CentOS/RHEL/Fedora

```bash
# For CentOS/RHEL
sudo yum install epel-release
sudo yum install ffmpeg-devel pkgconfig

# For Fedora
sudo dnf install ffmpeg-devel pkgconfig
```

### macOS

Using Homebrew:

```bash
brew install ffmpeg pkg-config
```

### Windows

For Windows, you'll need to download FFmpeg development libraries from the official website and set up the appropriate environment variables.

## Go Dependencies

This project uses CGO to interface with FFmpeg C libraries, so you need to have a C compiler installed:

### Linux
GCC is usually pre-installed. If not:
```bash
sudo apt install build-essential  # Ubuntu/Debian
```

### macOS
Install Xcode command line tools:
```bash
xcode-select --install
```

### Windows
Install MinGW or use Visual Studio with C++ support.

## Building the Project

After installing the prerequisites:

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd ffmpeg-go
   ```

2. Install Go dependencies:
   ```bash
   make install-deps
   ```

3. Build the examples:
   ```bash
   make build
   ```

4. Run the examples:
   ```bash
   ./bin/info <path_to_media_file>
   ./bin/decode <path_to_video_file>
   ```

## Troubleshooting

### pkg-config errors

If you encounter pkg-config errors, make sure the PKG_CONFIG_PATH environment variable is set correctly:

```bash
export PKG_CONFIG_PATH="/usr/local/lib/pkgconfig:/usr/lib/pkgconfig"
```

### Library linking issues

If you get linking errors, you might need to set the CGO_LDFLAGS:

```bash
export CGO_LDFLAGS="-L/usr/local/lib"
```
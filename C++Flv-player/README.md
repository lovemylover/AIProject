# C++ FLV Player

一个使用FFmpeg库实现的简单C++ FLV播放器。

## 架构说明

这个FLV播放器采用了模块化的架构设计：

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   FLV File      │───▶│  FFmpeg Library  │───▶│  Frame Processor│
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              ▲                       ▲
                              │                       │
                              ▼                       ▼
                      ┌──────────────────┐    ┌─────────────────┐
                      │  AVFormat Module │    │  Display System │
                      └──────────────────┘    └─────────────────┘
```

### 核心组件

1. **FLVParser**: 解析FLV文件格式
2. **FFmpeg Integration**: 使用FFmpeg库进行解码
3. **FrameProcessor**: 处理解码后的视频帧
4. **Main Controller**: 控制播放流程

## 技术特点

- 基于FFmpeg库进行FLV解析和解码
- 支持H.264视频编码和AAC音频编码
- 跨平台支持（Windows/Linux/macOS）
- CMake构建系统
- 面向对象设计

## 文件结构

```
C++Flv-player/
├── main.cpp         # 主程序文件
├── CMakeLists.txt   # CMake构建配置
├── build.sh         # Linux/macOS构建脚本
├── build.bat        # Windows构建脚本
└── README.md        # 说明文档
```

## 依赖项

### 必需的库

1. **FFmpeg libraries**:
   - libavformat: 处理各种多媒体容器格式
   - libavcodec: 提供编解码功能
   - libavutil: 提供各种实用函数
   - libswscale: 提供图像缩放和色彩空间转换

2. **SDL2 library**:
   - 用于创建窗口和渲染视频画面

### 安装依赖

#### Ubuntu/Debian:
```bash
sudo apt-get update
sudo apt-get install build-essential cmake
sudo apt-get install libavformat-dev libavcodec-dev libavutil-dev libswscale-dev
sudo apt-get install libsdl2-dev
```

#### CentOS/RHEL/Fedora:
```bash
sudo yum install gcc gcc-c++ make cmake
sudo yum install ffmpeg-devel SDL2-devel
# 或对于较新的版本:
sudo dnf install gcc gcc-c++ make cmake ffmpeg-devel SDL2-devel
```

#### macOS (使用Homebrew):
```bash
brew install cmake ffmpeg sdl2
```

**注意**: 在macOS上，您可能需要设置一些环境变量才能正确链接FFmpeg库：
```bash
export PKG_CONFIG_PATH="/opt/homebrew/lib/pkgconfig:/usr/local/lib/pkgconfig:$PKG_CONFIG_PATH"
export LDFLAGS="-L/opt/homebrew/lib"
export CPPFLAGS="-I/opt/homebrew/include"
```

#### Windows:
1. 下载并安装Visual Studio
2. 下载FFmpeg开发包: https://www.ffmpeg.org/download.html
3. 下载SDL2开发包: https://www.libsdl.org/download-2.0.php
4. 解压并设置环境变量

## 构建项目

### Linux/macOS:

```bash
# 添加执行权限
chmod +x build.sh

# 运行构建脚本
./build.sh
```

### macOS (推荐):

在macOS上，我们提供了专门的构建脚本来处理Homebrew安装的FFmpeg库：

```bash
# 添加执行权限
chmod +x build_macos.sh

# 运行macOS专用构建脚本
./build_macos.sh
```

这个脚本会自动检测Homebrew安装的FFmpeg库并正确配置链接路径。

### Windows:

```cmd
# 运行构建脚本
build.bat
```

### 手动构建:

```bash
# 创建构建目录
mkdir build
cd build

# 运行CMake
cmake ..

# 构建项目
make  # Linux/macOS
# 或
cmake --build . --config Release  # Windows
```

## 使用方法

构建完成后，运行播放器：

```bash
# Linux/macOS
./build/flvcpp path/to/your/video.flv

# Windows
build\Release\flvcpp.exe path\to\your\video.flv
```

## 代码结构说明

### FLVPlayer类

封装了所有播放功能：

1. **initialize()**: 初始化播放器，打开文件并设置解码器
2. **play()**: 开始播放，解码并处理视频帧
3. **processFrame()**: 处理解码后的帧（可扩展为渲染到窗口）
4. **cleanup()**: 清理资源

### 关键技术点

1. **FFmpeg集成**: 使用libavformat, libavcodec, libswscale等库
2. **内存管理**: 正确分配和释放FFmpeg资源
3. **错误处理**: 检查每个FFmpeg函数的返回值
4. **跨平台支持**: 处理不同操作系统的差异

## 扩展建议

1. **音频支持**: 添加音频解码和播放功能
2. **播放控制**: 实现暂停、快进、快退等功能
3. **网络流支持**: 支持RTMP、HTTP等网络流播放
4. **多格式支持**: 扩展支持更多视频格式
5. **硬件加速**: 利用GPU进行解码加速
6. **用户界面**: 添加播放控制按钮和进度条

## 常见问题

### 编译错误

如果遇到FFmpeg相关错误，请确保：
1. FFmpeg开发包已正确安装
2. CMake能找到FFmpeg库和头文件
3. 环境变量设置正确

### 运行时错误

如果播放失败，请检查：
1. FLV文件路径是否正确
2. FLV文件是否损坏
3. 视频编码格式是否支持

## 许可证

本项目仅供学习和参考使用。
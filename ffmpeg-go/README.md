# FFmpeg-Go

FFmpeg API接口的Go语言封装库。

## 功能

- 封装FFmpeg的核心功能
- 提供易于使用的Go接口
- 支持音视频解码、编码、转码等功能

## 安装

```bash
go get github.com/yourusername/ffmpeg-go
```

## 使用示例

```go
package main

import (
    "fmt"
    "ffmpeg-go/ffmpeg"
)

func main() {
    // 示例：获取媒体文件信息
    info, err := ffmpeg.GetMediaInfo("input.mp4")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Duration: %s\n", info.Duration)
    fmt.Printf("Bitrate: %d kb/s\n", info.Bitrate)
}
```
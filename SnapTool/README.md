# SnapTool - 视频截图工具

SnapTool 是一个基于Go语言和FFmpeg的快速视频截图工具，可以按照25fps的频率对视频进行截图。

## 功能特点

- 按照25fps的频率对视频进行截图
- 自动计算视频总帧数并逐帧截图
- 显示截图进度
- 支持自定义输出目录
- 并发处理提高截图速度
- 自动检查依赖环境

## 依赖要求

- Go 1.19 或更高版本
- FFmpeg 和 ffprobe (需要添加到系统PATH中)

## 安装步骤

1. 确保已安装Go环境和FFmpeg工具包
2. 克隆或下载本项目
3. 进入项目目录并运行:

```bash
go build
```

## 使用方法

编译后会生成可执行文件，使用方法如下：

```bash
./SnapTool <video_file> [output_directory]
```

参数说明:
- `<video_file>`: 要截图的视频文件路径 (必需)
- `[output_directory]`: 截图输出目录，默认为当前目录下的`screenshots`文件夹 (可选)

### 示例

```bash
# 对test.mp4视频按25fps截图，保存到默认目录
./SnapTool test.mp4

# 对test.mp4视频按25fps截图，保存到指定目录
./SnapTool test.mp4 ./my_screenshots
```

## 输出文件

截图文件将保存为JPEG格式，命名为:
- frame_000000.jpg
- frame_000001.jpg
- frame_000002.jpg
- ...

## 注意事项

1. 确保FFmpeg和ffprobe命令可以在终端中直接运行
2. 视频文件路径需要正确无误
3. 程序会自动创建输出目录（如果不存在）
4. 截图质量设置为较高水平(-q:v 2)，可以根据需要调整
5. 程序默认使用5个并发进程来加速截图处理

## 技术实现

程序通过以下步骤实现截图功能：

1. 使用ffprobe获取视频总时长
2. 计算按25fps需要截图的总帧数
3. 使用goroutine并发处理截图任务
4. 使用ffmpeg按时间点逐帧截图
5. 显示进度信息直到完成

## 性能优化

- 使用并发处理提高截图速度
- 合理设置ffmpeg参数优化截图质量与速度
- 预先检查依赖环境确保程序正常运行

## 许可证

MIT License
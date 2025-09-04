# AIProject
每天一个用 AI 开发一个小工具
###### 1. 项目介绍
AIProject 是一个致力于每天开发一个基于人工智能技术的小工具的项目。无论是自动化任务、数据分析还是创意生成，我们都希望通过 AI 技术为用户提供实用且有趣的解决方案。
###### 2. 项目目标
- 每天开发并发布一个新的 AI 工具。
- 提供详细的文档和使用指南。
- 鼓励社区参与和反馈。
- 关注 AI 技术的最新发展和趋势。
- 与 AI 技术社区合作，分享经验和知识。
###### 3. 项目进度
- 20250827 开发一个去水印工具
  - 功能：去除图片中的水印
    用法1 (自动检测水印): python watermark_remover.py <输入图像路径> <输出图像路径>
    用法2 (手动指定水印): python watermark_remover.py <输入图像路径> <输出图像路径> <x> <y> <宽度> <高度>
    示例: python watermark_remover.py input.jpg output.jpg
    技术栈：Python, OpenCV, NumPy
    需要本地安装 Python, OpenCV, NumPy

- 20280828 开发一个文本转语音的工具
  - 功能：将文本转换为语音
    用法：./TTS -f test.txt
    需要技术栈：C++
    需要本地安装：eSpeak-ng 库

- 20280829 开发一个rtp解析服务端与发送客户端demo
  - 功能：解析rtp包与发送rtp包
    用法：./rtp_server
    需要技术栈：go
    需要本地安装：go语言环境

 -20250903 开发一个ffmpeg的cgo库
   - 功能：调用ffmpeg的api
    - 用法：./ffmpeg_cgo
    - 需要技术栈：cgo
    - 需要本地安装：ffmpeg库

 -20250904开发一个基于C++与Ffmpeg的FLV播放器
  - 功能：播放flv视频
  - 用法：./flv_player <视频路径>
  - 需要技术栈：C++
  - 需要本地安装：ffmpeg库 
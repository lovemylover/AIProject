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

- 20280812 
#!/bin/bash

# Stream Pusher 构建脚本

set -e

echo "🚀 开始构建 Stream Pusher..."

# 进入源码目录
cd "$(dirname "$0")/src"

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go环境未安装，请先安装Go"
    exit 1
fi

# 检查FFmpeg
if ! command -v ffmpeg &> /dev/null; then
    echo "❌ FFmpeg未安装，请先安装FFmpeg"
    echo "   macOS: brew install ffmpeg"
    echo "   Ubuntu: sudo apt install ffmpeg"
    exit 1
fi

# 构建
echo "📦 正在编译..."
go build -o ../bin/stream-pusher main.go

echo "✅ 构建完成！"
echo "📍 可执行文件位置: $(pwd)/../bin/stream-pusher"
echo ""
echo "🎯 使用方法:"
echo "   ./bin/stream-pusher -input \"http://camera-url\" -output \"rtmp://server-url\""
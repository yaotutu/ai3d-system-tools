# Stream Pusher

一个简单的Go程序，用于将stream流推送到B站直播服务器。

## 功能特性

- 接收各种格式的输入流（HTTP流、本地文件等）
- 使用FFmpeg进行流处理和转码
- 推送到B站直播服务器（RTMP协议）
- 支持优雅退出
- 高性能流复制（无重编码）

## 依赖要求

- Go 1.16+
- FFmpeg（需要安装在系统中）

## 安装FFmpeg

### macOS
```bash
brew install ffmpeg
```

### Ubuntu/Debian
```bash
sudo apt update
sudo apt install ffmpeg
```

### Windows
下载FFmpeg并添加到系统PATH中

## 使用方法

### 基本用法
```bash
go run main.go -input <输入流地址> -output <B站推流地址>
```

### 示例
```bash
# 从HTTP流推送到B站
go run main.go -input "http://example.com/live.m3u8" -output "rtmp://live-push.bilivideo.com/live-bvc/YOUR_STREAM_KEY"

# 从本地文件推送到B站
go run main.go -input "video.mp4" -output "rtmp://live-push.bilivideo.com/live-bvc/YOUR_STREAM_KEY"

# 从摄像头推送到B站（Linux/macOS）
go run main.go -input "/dev/video0" -output "rtmp://live-push.bilivideo.com/live-bvc/YOUR_STREAM_KEY"
```

## B站直播推流地址获取

1. 登录B站，进入直播间
2. 开启直播
3. 获取推流地址，格式为：`rtmp://live-push.bilivideo.com/live-bvc/YOUR_STREAM_KEY`

## 支持的输入格式

- HTTP Live Streaming (HLS) - `.m3u8`
- DASH流
- 本地视频文件 - `.mp4`, `.avi`, `.mkv` 等
- 网络摄像头设备
- UDP/TCP流

## 注意事项

- 确保网络连接稳定
- B站推流码请妥善保管，不要泄露
- 推流前确认B站直播间已开启
- 程序使用Ctrl+C优雅退出
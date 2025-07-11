# AI3D System Tools

一个专业的直播和流媒体工具集合，包含多个高性能的推流和处理工具。

## 📁 项目结构

```
ai3d-system-tools/
├── stream-pusher/          # Go实现的智能流推送工具
│   ├── src/               # 源代码
│   ├── docs/              # 项目文档和测试报告
│   ├── configs/           # 配置文件示例
│   ├── tests/             # 测试文件
│   └── README.md          # 项目说明
├── [future-tools]/        # 未来的其他工具
└── README.md              # 总体说明
```

## 🛠️ 工具列表

### ✅ Stream Pusher (已完成)
- **功能**: 将摄像头流推送到RTMP直播服务器
- **技术栈**: Go + FFmpeg
- **特性**: 
  - 智能参数优化
  - 自动重连机制
  - 详细错误诊断
  - 75分钟+稳定运行验证
- **状态**: 🏆 生产可用

### 🔄 规划中的工具
- **Multi-Stream Pusher**: 多路流同时推送
- **Stream Recorder**: 直播流录制工具
- **Stream Monitor**: 流状态监控工具
- **Stream Transcoder**: 实时转码工具

## 🚀 快速开始

### Stream Pusher
```bash
cd stream-pusher/src
go run main.go -input "http://your-camera-url" -output "rtmp://your-rtmp-server"
```

详细使用说明请查看各工具的README文档。

## 📊 项目特色

- **高性能**: 优化的参数配置，确保实时处理
- **高稳定性**: 经过长时间验证的稳定性
- **智能诊断**: 详细的错误分析和解决建议
- **易于扩展**: 模块化设计，便于添加新工具

## 🤝 贡献

欢迎提交Issue和Pull Request来改进项目。

## 📄 许可证

MIT License
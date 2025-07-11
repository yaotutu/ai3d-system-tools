package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// StreamPusher 流推送器结构体
type StreamPusher struct {
	inputURL     string // 输入流地址
	outputURL    string // 输出RTMP地址
	cmd          *exec.Cmd
	maxRetries   int    // 最大重试次数
	retryDelay   time.Duration // 重试间隔
	currentRetry int    // 当前重试次数
}

// ErrorType 错误类型枚举
type ErrorType int

const (
	ErrorTypeInputStream ErrorType = iota  // 输入流错误
	ErrorTypeRTMPOutput                     // RTMP输出错误
	ErrorTypeEncoding                       // 编码错误
	ErrorTypeNetwork                        // 网络错误
	ErrorTypeUnknown                        // 未知错误
)

// NewStreamPusher 创建新的流推送器
func NewStreamPusher(inputURL, outputURL string) *StreamPusher {
	return &StreamPusher{
		inputURL:     inputURL,
		outputURL:    outputURL,
		maxRetries:   5,
		retryDelay:   5 * time.Second,
		currentRetry: 0,
	}
}

// CheckInputStream 检查输入流状态
func (sp *StreamPusher) CheckInputStream() error {
	log.Printf("🔍 检查输入流: %s", sp.inputURL)
	
	// 解析URL
	parsedURL, err := url.Parse(sp.inputURL)
	if err != nil {
		log.Printf("❌ 输入流URL格式错误: %v", err)
		return fmt.Errorf("输入流URL格式错误: %v", err)
	}
	
	// 如果是HTTP/HTTPS流，检查连接性
	if parsedURL.Scheme == "http" || parsedURL.Scheme == "https" {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Head(sp.inputURL)
		if err != nil {
			log.Printf("❌ 无法连接到输入流: %v", err)
			return fmt.Errorf("无法连接到输入流: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode >= 400 {
			log.Printf("❌ 输入流返回错误状态: %d %s", resp.StatusCode, resp.Status)
			return fmt.Errorf("输入流返回错误状态: %d %s", resp.StatusCode, resp.Status)
		}
		
		log.Printf("✅ 输入流连接正常 (状态: %d, 类型: %s)", resp.StatusCode, resp.Header.Get("Content-Type"))
	} else {
		log.Printf("ℹ️  检测到非HTTP流: %s", parsedURL.Scheme)
	}
	
	return nil
}

// CheckRTMPOutput 检查RTMP输出地址
func (sp *StreamPusher) CheckRTMPOutput() error {
	log.Printf("🔍 检查RTMP输出地址: %s", sp.outputURL)
	
	// 解析RTMP URL
	parsedURL, err := url.Parse(sp.outputURL)
	if err != nil {
		log.Printf("❌ RTMP URL格式错误: %v", err)
		return fmt.Errorf("RTMP URL格式错误: %v", err)
	}
	
	if parsedURL.Scheme != "rtmp" {
		log.Printf("❌ 不是有效的RTMP地址，应该以rtmp://开头")
		return fmt.Errorf("不是有效的RTMP地址，应该以rtmp://开头")
	}
	
	log.Printf("✅ RTMP地址格式正确")
	log.Printf("   服务器: %s", parsedURL.Host)
	log.Printf("   推流码: %s", strings.TrimPrefix(parsedURL.Path, "/"))
	
	return nil
}

// AnalyzeError 分析FFmpeg错误输出
func (sp *StreamPusher) AnalyzeError(exitCode int, output string) ErrorType {
	// 输入流相关错误
	inputErrors := []string{
		"HTTP error 502", "HTTP error 404", "HTTP error 500",
		"Connection refused", "No route to host",
		"Error opening input file", "Server returned 5XX",
		"Invalid data found when processing input",
	}
	
	// RTMP输出相关错误
	rtmpErrors := []string{
		"Broken pipe", "Connection reset by peer",
		"RTMP_SendPacket", "Failed to connect",
		"Handshake failed", "Authentication failed",
	}
	
	// 编码相关错误
	encodingErrors := []string{
		"not compatible with flv", "codec not supported",
		"Conversion failed", "Encoder not found",
	}
	
	for _, errStr := range inputErrors {
		if strings.Contains(output, errStr) {
			return ErrorTypeInputStream
		}
	}
	
	for _, errStr := range rtmpErrors {
		if strings.Contains(output, errStr) {
			return ErrorTypeRTMPOutput
		}
	}
	
	for _, errStr := range encodingErrors {
		if strings.Contains(output, errStr) {
			return ErrorTypeEncoding
		}
	}
	
	return ErrorTypeUnknown
}

// PrintErrorSuggestion 根据错误类型打印解决建议
func (sp *StreamPusher) PrintErrorSuggestion(errorType ErrorType, output string) {
	switch errorType {
	case ErrorTypeInputStream:
		log.Println("")
		log.Println("🔧 输入流问题解决建议:")
		log.Println("   1. 检查摄像头设备是否在线: ping 192.168.201.124")
		log.Println("   2. 确认摄像头服务是否运行")
		log.Println("   3. 尝试在浏览器中打开流地址")
		log.Println("   4. 检查网络连接")
		if strings.Contains(output, "502") {
			log.Println("   5. 502错误通常表示摄像头服务内部错误，需要重启设备")
		}
		
	case ErrorTypeRTMPOutput:
		log.Println("")
		log.Println("🔧 RTMP推流问题解决建议:")
		log.Println("   1. 检查B站直播间是否已开启")
		log.Println("   2. 确认推流码是否正确")
		log.Println("   3. 检查RTMP服务器连接: telnet 192.168.200.68 1935")
		log.Println("   4. 推流码可能已过期，请重新获取")
		if strings.Contains(output, "Broken pipe") {
			log.Println("   5. 'Broken pipe'通常表示服务器主动断开连接")
			log.Println("   6. 可能是推流时间限制或推流码问题")
		}
		
	case ErrorTypeEncoding:
		log.Println("")
		log.Println("🔧 编码问题解决建议:")
		log.Println("   1. 尝试不同的编码参数")
		log.Println("   2. 检查FFmpeg是否支持所需的编解码器")
		log.Println("   3. 输入流格式可能不兼容")
		
	default:
		log.Println("")
		log.Println("🔧 通用解决建议:")
		log.Println("   1. 检查网络连接")
		log.Println("   2. 确认所有服务正常运行")
		log.Println("   3. 查看完整的FFmpeg输出日志")
	}
	log.Println("")
}

// Start 开始推流
func (sp *StreamPusher) Start() error {
	// 构建FFmpeg命令
	// -i: 输入源
	// -c:v libx264: 视频编码为H.264（RTMP兼容）
	// -c:a aac: 音频编码为AAC（RTMP兼容）
	// -preset ultrafast: 快速编码预设
	// -f flv: 输出格式为FLV（RTMP协议要求）
	args := []string{
		"-re",                       // 实时读取输入（重要！）
		"-i", sp.inputURL,           // 输入流
		"-f", "lavfi",               // 使用lavfi生成静音音频
		"-i", "anullsrc=channel_layout=stereo:sample_rate=48000", // 生成静音音频源(48kHz匹配OBS)
		
		// 视频编码参数（参考OBS设置）
		"-c:v", "libx264",           // 视频编码为H.264
		"-preset", "veryfast",       // 编码预设（匹配OBS）
		"-tune", "zerolatency",      // 零延迟调优
		"-x264-params", "nal-hrd=cbr", // CBR模式
		"-b:v", "1000k",             // 视频码率（匹配640x480分辨率）
		"-maxrate", "1000k",         // 最大码率
		"-bufsize", "1000k",         // 缓冲区大小
		"-g", "250",                 // 关键帧间隔（匹配OBS keyint）
		"-keyint_min", "25",         // 最小关键帧间隔
		"-r", "30",                  // 帧率30fps（保持标准）
		"-s", "640x480",             // 分辨率（保持原始分辨率）
		
		// 音频编码参数（参考OBS设置）
		"-c:a", "aac",               // 音频编码为AAC
		"-b:a", "160k",              // 音频码率（匹配OBS）
		"-ar", "48000",              // 采样率48kHz（匹配OBS）
		"-ac", "2",                  // 双声道
		
		"-shortest",                 // 以最短流为准
		"-f", "flv",                 // 输出格式
		"-flvflags", "no_duration_filesize", // FLV优化参数
		sp.outputURL,                // 输出RTMP地址
	}

	sp.cmd = exec.Command("ffmpeg", args...)
	sp.cmd.Stdout = os.Stdout
	sp.cmd.Stderr = os.Stderr

	log.Printf("🚀 开始推流: %s -> %s", sp.inputURL, sp.outputURL)
	log.Printf("📋 FFmpeg命令: ffmpeg %v", args)
	log.Println("💡 提示: 使用 Ctrl+C 停止推流")
	log.Println("")

	return sp.cmd.Start()
}

// StartWithRetry 带重连的推流启动
func (sp *StreamPusher) StartWithRetry() error {
	for sp.currentRetry <= sp.maxRetries {
		if sp.currentRetry > 0 {
			log.Printf("🔄 第 %d/%d 次重连，等待 %v...", sp.currentRetry, sp.maxRetries, sp.retryDelay)
			time.Sleep(sp.retryDelay)
		}

		log.Printf("🚀 启动推流 (尝试 %d/%d)", sp.currentRetry+1, sp.maxRetries+1)
		
		err := sp.Start()
		if err != nil {
			log.Printf("❌ 启动失败: %v", err)
			sp.currentRetry++
			continue
		}

		// 等待推流结束
		waitErr := sp.Wait()
		if waitErr == nil {
			log.Println("✅ 推流正常结束")
			return nil
		}

		// 分析错误类型
		errorType := sp.AnalyzeError(1, waitErr.Error())
		
		// 如果是输入流错误，不重连
		if errorType == ErrorTypeInputStream {
			log.Printf("❌ 输入流错误，停止重连: %v", waitErr)
			return waitErr
		}

		log.Printf("⚠️ 推流中断: %v", waitErr)
		sp.currentRetry++
		
		if sp.currentRetry <= sp.maxRetries {
			log.Printf("📡 检测到网络中断，将自动重连...")
		}
	}

	return fmt.Errorf("达到最大重试次数 (%d)，停止重连", sp.maxRetries)
}

// Stop 停止推流
func (sp *StreamPusher) Stop() error {
	if sp.cmd != nil && sp.cmd.Process != nil {
		log.Println("正在停止推流...")
		return sp.cmd.Process.Kill()
	}
	return nil
}

// Wait 等待推流结束
func (sp *StreamPusher) Wait() error {
	if sp.cmd != nil {
		return sp.cmd.Wait()
	}
	return nil
}

func main() {
	// 命令行参数解析
	var (
		inputURL  = flag.String("input", "", "输入流地址 (例如: http://example.com/stream.m3u8)")
		outputURL = flag.String("output", "", "B站直播推流地址 (例如: rtmp://live-push.bilivideo.com/live-bvc/YOUR_STREAM_KEY)")
		checkOnly = flag.Bool("check", false, "仅检查连接状态，不进行推流")
	)
	flag.Parse()

	// 参数验证
	if *inputURL == "" || *outputURL == "" {
		fmt.Println("使用方法:")
		fmt.Println("  go run main.go -input <输入流地址> -output <B站推流地址>")
		fmt.Println("")
		fmt.Println("示例:")
		fmt.Println("  go run main.go -input http://example.com/live.m3u8 -output rtmp://live-push.bilivideo.com/live-bvc/YOUR_KEY")
		fmt.Println("")
		fmt.Println("选项:")
		fmt.Println("  -check  仅检查连接状态，不进行推流")
		os.Exit(1)
	}

	log.Println("🎯 Stream Pusher v1.0 - 智能流推送工具")
	log.Println("=======================================\n")

	// 创建流推送器
	pusher := NewStreamPusher(*inputURL, *outputURL)

	// 预检查阶段
	log.Println("📋 开始预检查...")
	
	// 检查输入流
	if err := pusher.CheckInputStream(); err != nil {
		log.Printf("❌ 输入流检查失败: %v", err)
		pusher.PrintErrorSuggestion(ErrorTypeInputStream, err.Error())
		if *checkOnly {
			os.Exit(1)
		}
		log.Println("⚠️  输入流有问题，但仍将尝试推流...")
	}
	
	// 检查RTMP输出
	if err := pusher.CheckRTMPOutput(); err != nil {
		log.Printf("❌ RTMP地址检查失败: %v", err)
		pusher.PrintErrorSuggestion(ErrorTypeRTMPOutput, err.Error())
		os.Exit(1)
	}
	
	if *checkOnly {
		log.Println("\n✅ 所有检查完成，连接状态良好！")
		return
	}
	
	log.Println("✅ 预检查完成，开始推流...\n")

	// 设置信号处理，优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 设置信号处理，优雅退出
	go func() {
		<-sigChan
		log.Println("收到退出信号，正在停止推流...")
		pusher.Stop()
		os.Exit(0)
	}()

	// 启动带重连的推流
	if err := pusher.StartWithRetry(); err != nil {
		log.Printf("\n❌ 推流最终失败: %v", err)
		
		// 分析错误类型并提供建议
		errorType := pusher.AnalyzeError(1, err.Error())
		pusher.PrintErrorSuggestion(errorType, err.Error())
		os.Exit(1)
	} else {
		log.Println("\n✅ 推流任务完成")
	}
}
# Stream Pusher 参数分析与问题排查文档

## 🎯 概述
本文档详细分析了Stream Pusher中可能导致"Broken pipe"错误的每个参数，并提供逐一验证的方法。

## 📊 参数对比表

| 参数类别 | 问题参数 | OBS参数 | 状态 | 可能影响 |
|---------|----------|---------|------|----------|
| 视频分辨率 | 640x480 | 1920x1080 | ❌ | 高 |
| 视频帧率 | 25fps | 30fps | ❌ | 中 |
| 视频码率 | 1000kbps | 2500kbps | ❌ | 中 |
| 码率控制 | CRF | CBR | ❌ | 高 |
| 编码预设 | ultrafast | veryfast | ❌ | 低 |
| 关键帧间隔 | 默认 | 250 | ❌ | 中 |
| 音频采样率 | 44100Hz | 48000Hz | ❌ | 高 |
| 音频码率 | 默认 | 160kbps | ❌ | 低 |
| 音频声道 | 默认 | 双声道 | ❌ | 低 |

## 🔍 逐参数验证计划

### 1. 视频分辨率验证 (优先级：高)
```bash
# 测试1：保持原分辨率
-s 640x480

# 测试2：升级到720p
-s 1280x720

# 测试3：升级到1080p
-s 1920x1080
```

**验证方法**：
1. 修改main.go中的`-s`参数
2. 运行推流测试
3. 记录连接持续时间
4. 观察是否还有"Broken pipe"错误

**预期结果**：
- 640x480：可能在10-15秒后断开
- 1280x720：可能延长到30秒
- 1920x1080：应该稳定连接

### 2. 视频帧率验证 (优先级：中)
```bash
# 测试1：保持原帧率
-r 25

# 测试2：标准帧率
-r 30

# 测试3：更高帧率
-r 60
```

**验证方法**：
1. 固定分辨率为1920x1080
2. 仅修改帧率参数
3. 测试各种帧率的稳定性

**预期结果**：
- 25fps：可能被某些服务器拒绝
- 30fps：标准直播帧率，应该稳定
- 60fps：可能对带宽要求过高

### 3. 码率控制验证 (优先级：高)
```bash
# 测试1：CRF模式（质量优先）
-crf 23

# 测试2：CBR模式（码率恒定）
-b:v 2500k -maxrate 2500k -bufsize 2500k

# 测试3：VBR模式（可变码率）
-b:v 2500k -maxrate 3000k -bufsize 2000k
```

**验证方法**：
1. 分别测试三种码率控制模式
2. 观察网络稳定性
3. 记录服务器接受度

**预期结果**：
- CRF：可变码率可能导致服务器断开
- CBR：恒定码率应该最稳定
- VBR：介于两者之间

### 4. 音频采样率验证 (优先级：高)
```bash
# 测试1：CD音质采样率
-ar 44100

# 测试2：专业音频采样率
-ar 48000

# 测试3：高音质采样率
-ar 96000
```

**验证方法**：
1. 固定视频参数
2. 仅修改音频采样率
3. 测试音频参数对连接稳定性的影响

**预期结果**：
- 44100Hz：可能不被某些直播服务器接受
- 48000Hz：专业标准，应该稳定
- 96000Hz：可能过高，浪费带宽

### 5. 关键帧间隔验证 (优先级：中)
```bash
# 测试1：短间隔
-g 60 -keyint_min 6

# 测试2：标准间隔
-g 250 -keyint_min 25

# 测试3：长间隔
-g 500 -keyint_min 50
```

**验证方法**：
1. 测试不同关键帧间隔的影响
2. 观察网络中断后的恢复能力
3. 记录带宽使用情况

**预期结果**：
- 短间隔：带宽高但恢复快
- 标准间隔：平衡性能和带宽
- 长间隔：可能影响直播稳定性

## 🧪 验证脚本生成

### 脚本1：分辨率验证
```go
// 在main.go中修改这行：
"-s", "640x480",    // 测试1
"-s", "1280x720",   // 测试2  
"-s", "1920x1080",  // 测试3
```

### 脚本2：帧率验证
```go
// 在main.go中修改这行：
"-r", "25",         // 测试1
"-r", "30",         // 测试2
"-r", "60",         // 测试3
```

### 脚本3：码率控制验证
```go
// 测试1：CRF模式
args := []string{
    // ... 其他参数
    "-crf", "23",
    // 移除 -b:v, -maxrate, -bufsize
}

// 测试2：CBR模式
args := []string{
    // ... 其他参数
    "-b:v", "2500k",
    "-maxrate", "2500k", 
    "-bufsize", "2500k",
}
```

### 脚本4：音频采样率验证
```go
// 在main.go中修改这两行：
"-i", "anullsrc=channel_layout=stereo:sample_rate=44100",  // 测试1
"-ar", "44100",                                            // 测试1

"-i", "anullsrc=channel_layout=stereo:sample_rate=48000",  // 测试2
"-ar", "48000",                                            // 测试2
```

## 📋 验证记录表

### 测试记录模板
```
测试日期：____
测试参数：____
连接时长：____
错误信息：____
CPU使用率：____
网络带宽：____
服务器响应：____
```

### 详细测试表格

| 测试编号 | 参数配置 | 连接时长 | 错误类型 | 成功率 | 备注 |
|----------|----------|----------|----------|--------|------|
| T001 | 640x480,25fps,CRF | ___秒 | Broken pipe | __% | 原始配置 |
| T002 | 1280x720,25fps,CRF | ___秒 | ___ | __% | 升级分辨率 |
| T003 | 1920x1080,25fps,CRF | ___秒 | ___ | __% | 全高清 |
| T004 | 1920x1080,30fps,CRF | ___秒 | ___ | __% | 标准帧率 |
| T005 | 1920x1080,30fps,CBR | ___秒 | ___ | __% | 恒定码率 |
| T006 | 1920x1080,30fps,CBR,48kHz | ___秒 | ___ | __% | 专业音频 |

## 🎯 验证步骤

### 步骤1：基线测试
```bash
# 使用原始"问题"配置
go run main.go -input "http://192.168.201.124/webcam/?action=stream" -output "rtmp://192.168.200.68:1935/livehime"
```
记录：连接时长、错误信息、CPU使用率

### 步骤2：单参数验证
按优先级逐一修改参数：
1. 分辨率 (640x480 → 1920x1080)
2. 码率控制 (CRF → CBR)  
3. 音频采样率 (44100Hz → 48000Hz)
4. 帧率 (25fps → 30fps)
5. 关键帧间隔 (默认 → 250)

### 步骤3：组合验证
找到关键参数后，测试参数组合的效果

### 步骤4：最终验证
使用完整OBS配置进行长时间稳定性测试

## 📊 预期验证结果

### 高影响参数 (预测)
1. **分辨率**：从640x480升级到1920x1080应该显著提升稳定性
2. **码率控制**：从CRF改为CBR应该解决大部分连接问题
3. **音频采样率**：从44100Hz改为48000Hz应该提升兼容性

### 中影响参数 (预测)
1. **帧率**：从25fps改为30fps应该提升标准化程度
2. **关键帧间隔**：设置为250应该提升网络适应性

### 低影响参数 (预测)
1. **编码预设**：从ultrafast改为veryfast影响较小
2. **音频码率**：明确设置为160kbps影响较小

## 🔧 验证工具

### 监控脚本
```bash
# 监控推流状态
while true; do
    echo "$(date): 检查推流状态..."
    ps aux | grep ffmpeg
    sleep 5
done
```

### 网络监控
```bash
# 监控网络使用
netstat -i
iftop -i en0
```

### 日志分析
```bash
# 分析FFmpeg日志
go run main.go ... 2>&1 | tee test_log.txt
```

## 📝 结论模板

完成验证后，在这里记录：

### 关键发现
1. ___参数是导致连接断开的主要原因
2. ___参数组合能实现稳定连接
3. ___参数对性能影响最大

### 最优配置
```
分辨率：___
帧率：___
码率控制：___
音频采样率：___
其他关键参数：___
```

### 性能数据
- 平均连接时长：___
- CPU使用率：___
- 内存使用：___
- 网络带宽：___

---

**验证建议**：建议按照优先级从高到低逐一验证，每个参数至少测试3次以确保结果可靠。
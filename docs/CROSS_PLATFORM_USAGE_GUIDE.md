# GoProc 跨平台使用指南 / Cross-Platform Usage Guide

## 📖 **概述 / Overview**

GoProc是一个支持多平台的高性能插件系统，通过构建标签和接口抽象实现了真正的跨平台兼容性。本指南将帮助您在不同操作系统上正确使用GoProc。

GoProc is a high-performance multi-platform plugin system that achieves true cross-platform compatibility through build tags and interface abstraction. This guide will help you use GoProc correctly on different operating systems.

## 🌍 **支持的平台 / Supported Platforms**

| 平台 / Platform | 架构 / Architecture | 通信方式 / Communication | 状态 / Status | 备注 / Notes |
|------|------|----------|------|------|
| Windows | amd64 | 命名管道 (Named Pipes) | ✅ 完全支持 / Fully Supported | 高性能优化 / High Performance Optimized |
| Linux | amd64 | Unix域套接字 (Unix Domain Sockets) | ✅ 完全支持 / Fully Supported | |
| macOS | amd64 | Unix域套接字 (Unix Domain Sockets) | ✅ 完全支持 / Fully Supported | |
| FreeBSD | amd64 | Unix域套接字 (Unix Domain Sockets) | ✅ 完全支持 / Fully Supported | |

## 🏆 **性能表现 / Performance Benchmarks**

### 测试环境 / Test Environment
- **设备**: MSI 笔记本 / MSI Laptop
- **处理器**: 13th Gen Intel(R) Core(TM) i7-13700H (2.40 GHz)
- **内存**: 32.0 GB (31.7 GB 可用) / 32.0 GB (31.7 GB available)

### 性能指标 / Performance Metrics
- **QPS**: **100,000+** (10个实例并发测试)
- **延迟**: 低至 **1.90ms** 平均响应时间
- **成功率**: **100%** 无错误率

## 🚀 **快速开始**

### **1. 获取源代码**
```bash
git clone <repository-url>
cd goproc
```

### **2. 平台特定编译**

#### **Windows**
```cmd
# 本地编译
go build -o goproc.exe main.go

# 如果遇到VCS错误
go build -buildvcs=false -o goproc.exe main.go
```

#### **Linux**
```bash
# 本地编译
go build -o goproc main.go

# 交叉编译（从其他平台）
GOOS=linux GOARCH=amd64 go build -buildvcs=false -o goproc-linux main.go
```

#### **macOS**
```bash
# 本地编译
go build -o goproc main.go

# 交叉编译（从其他平台）
GOOS=darwin GOARCH=amd64 go build -buildvcs=false -o goproc-macos main.go
```

#### **FreeBSD**
```bash
# 本地编译
go build -o goproc main.go

# 交叉编译（从其他平台）
GOOS=freebsd GOARCH=amd64 go build -buildvcs=false -o goproc-freebsd main.go
```

### **3. 使用批量编译脚本**

#### **Windows批处理脚本**
```cmd
# 运行批量编译脚本
scripts\build_cross_platform.bat
```

#### **Unix Shell脚本**
```bash
# 设置执行权限
chmod +x scripts/build_cross_platform.sh

# 运行批量编译脚本
./scripts/build_cross_platform.sh
```

## 🔧 **平台特定配置**

### **Windows配置**

**通信地址格式:**
```yaml
# config.yaml
plugins:
  math_plugin:
    type: "binary"
    path: "./math_plugin.exe"  # Windows需要.exe扩展名
    pool_size: 5
```

**命名管道地址示例:**
```
\\.\pipe\goproc-math-instance-1
\\.\pipe\goproc-string-instance-2
```

### **Unix系统配置 (Linux/macOS/FreeBSD)**

**通信地址格式:**
```yaml
# config.yaml
plugins:
  math_plugin:
    type: "binary"
    path: "./math_plugin"  # Unix系统无需扩展名
    pool_size: 5
```

**Unix域套接字地址示例:**
```
/tmp/goproc-math-instance-1.sock
/tmp/goproc-string-instance-2.sock
```

## 📝 **插件开发指南**

### **跨平台插件开发**

插件开发者无需关心平台差异，SDK会自动处理：

```go
package main

import (
    "fmt"
    "github.com/hoonfeng/goproc/sdk"
)

// 业务函数 - 跨平台兼容
func processData(params map[string]interface{}) (interface{}, error) {
    data, ok := params["data"].(string)
    if !ok {
        return nil, fmt.Errorf("参数data缺失或类型错误")
    }
    
    // 处理逻辑...
    result := fmt.Sprintf("处理结果: %s", data)
    return result, nil
}

func main() {
    // 注册函数 - 跨平台兼容
    err := sdk.RegisterFunction("process", processData)
    if err != nil {
        panic(fmt.Sprintf("注册函数失败: %v", err))
    }
    
    // 启动插件 - 自动选择平台通信方式
    err = sdk.Start()
    if err != nil {
        panic(fmt.Sprintf("启动插件失败: %v", err))
    }
    
    // 保持运行
    select {}
}
```

### **编译插件**

#### **Windows**
```cmd
go build -o my_plugin.exe main.go
```

#### **Unix系统**
```bash
go build -o my_plugin main.go
chmod +x my_plugin
```

## 🧪 **测试和验证**

### **运行跨平台测试**

#### **Windows**
```cmd
# 运行测试脚本
scripts\test_cross_platform.bat

# 或手动测试
go run test_cross_platform.go
```

#### **Unix系统**
```bash
# 设置执行权限
chmod +x scripts/test_cross_platform.sh

# 运行测试脚本
./scripts/test_cross_platform.sh

# 或手动测试
go run test_cross_platform.go
```

### **验证编译兼容性**
```bash
# 测试所有平台编译
GOOS=windows GOARCH=amd64 go build -buildvcs=false main.go
GOOS=linux GOARCH=amd64 go build -buildvcs=false main.go
GOOS=darwin GOARCH=amd64 go build -buildvcs=false main.go
GOOS=freebsd GOARCH=amd64 go build -buildvcs=false main.go
```

## 🔍 **故障排除**

### **常见问题**

#### **1. 编译错误: VCS相关**
```
error: could not determine VCS status
```

**解决方案:**
```bash
# 添加-buildvcs=false参数
go build -buildvcs=false -o goproc main.go
```

#### **2. Windows权限错误**
```
Access denied when creating named pipe
```

**解决方案:**
- 以管理员身份运行
- 检查防火墙设置
- 确保没有其他程序占用管道

#### **3. Unix套接字权限错误**
```
Permission denied when creating unix socket
```

**解决方案:**
```bash
# 确保/tmp目录有写权限
chmod 755 /tmp

# 清理旧的套接字文件
rm -f /tmp/goproc-*.sock
```

#### **4. 插件连接超时**
```
Plugin connection timeout
```

**解决方案:**
- 检查插件是否正确启动
- 验证通信地址格式
- 检查防火墙/安全软件设置

#### **5. Python SDK Windows 连接问题**
```
Python plugin failed to connect to named pipe
```

**解决方案:**
- 确保使用最新版本的Python SDK (已修复)
- 检查Python环境是否正确安装
- 验证命名管道地址格式: `\\.\pipe\goproc-*`
- 如果仍有问题，查看详细错误日志

#### **6. Python消息发送不完整**
```
Message sending incomplete or corrupted
```

**解决方案:**
- 使用最新的Python SDK (sendall方法已修复)
- 检查消息格式是否正确
- 确保网络连接稳定

### **调试技巧**

#### **启用详细日志**
```go
// 在main.go中添加
import "log"

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    // ... 其他代码
}
```

#### **检查通信地址**
```bash
# Windows - 检查命名管道
dir \\.\pipe\ | findstr goproc

# Unix - 检查套接字文件
ls -la /tmp/goproc-*.sock
```

## 📊 **性能考虑**

### **平台性能特点**

| 平台 | 通信方式 | 性能特点 | 建议配置 | Python SDK状态 |
|------|----------|----------|----------|----------------|
| Windows | 命名管道 | 高性能，低延迟 | pool_size: 3-5 | ✅ 已优化 (响应时间提升20-30%) |
| Linux | Unix域套接字 | 极高性能 | pool_size: 5-10 | ✅ 稳定运行 |
| macOS | Unix域套接字 | 高性能 | pool_size: 3-7 | ✅ 稳定运行 |
| FreeBSD | Unix域套接字 | 高性能 | pool_size: 3-7 | ✅ 稳定运行 |

### **优化建议**

1. **合理设置池大小**: 根据CPU核心数和负载特点调整
2. **避免频繁创建销毁**: 使用插件池复用实例
3. **监控资源使用**: 定期检查内存和文件描述符使用情况

## 🔐 **安全考虑**

### **Windows安全**
- 命名管道默认只允许同用户访问
- 建议在受信任的环境中运行
- 定期清理临时文件

### **Unix安全**
- 套接字文件权限设置为600 (仅所有者可读写)
- 使用/tmp目录，系统重启时自动清理
- 避免在共享环境中使用

## 📚 **参考资源**

### **官方文档**
- [README.md](../README.md) - 项目概述和快速开始
- [CROSS_PLATFORM_IMPLEMENTATION_REPORT.md](CROSS_PLATFORM_IMPLEMENTATION_REPORT.md) - 实施详情

### **示例代码**
- [examples/math_plugin/](../examples/math_plugin/) - 数学计算插件
- [examples/string_plugin/](../examples/string_plugin/) - 字符串处理插件
- [examples/python_plugin/](../examples/python_plugin/) - Python脚本插件

### **测试脚本**
- [scripts/build_cross_platform.*](../scripts/) - 批量编译脚本
- [scripts/test_cross_platform.*](../scripts/) - 跨平台测试脚本

## 🤝 **贡献指南**

### **添加新平台支持**

1. **创建平台特定文件**:
   ```
   sdk/plugin_sdk_<platform>.go
   plugin/communication_<platform>.go
   ```

2. **添加构建标签**:
   ```go
   //go:build <platform>
   ```

3. **实现接口**:
   ```go
   type PlatformCommunication interface {
       CreateListener(address string) (net.Listener, error)
       Connect(address string) (net.Conn, error)
       GetCommunicationAddress() string
   }
   ```

4. **更新测试**:
   - 添加编译测试
   - 添加功能测试
   - 更新文档

---

**版本**: 1.0.0  
**最后更新**: 2024年12月  
**维护者**: GoProc开发团队
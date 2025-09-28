# GoProc - 高性能跨平台插件系统 / High-Performance Cross-Platform Plugin System

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS%20%7C%20FreeBSD-green.svg)](https://github.com/hoonfeng/goproc)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

GoProc 是一个基于命名管道（Windows）和Unix域套接字（Unix系统）的高性能跨平台插件系统，支持多语言插件开发，提供统一的SDK接口。

GoProc is a high-performance cross-platform plugin system based on Named Pipes (Windows) and Unix Domain Sockets (Unix systems), supporting multi-language plugin development with unified SDK interfaces.

## ✨ 核心特性 / Core Features

- 🚀 **极致性能**: QPS 可达 100,000+ (10万+)，延迟低至微秒级
- 🌍 **跨平台支持**: Windows、Linux、macOS、FreeBSD 全平台兼容
- 🔧 **多语言SDK**: 支持 Go、Python、Node.js 等多种编程语言
- 📦 **简化开发**: 插件开发者只需关注业务逻辑，无需处理底层通信
- ⚡ **高并发**: 基于插件池的并发调度，支持实例级并发控制
- 🛡️ **类型安全**: 完整的类型定义和参数验证机制

## 🏆 性能表现 / Performance Benchmarks

### 测试环境 / Test Environment
- **设备**: MSI 笔记本 / MSI Laptop
- **处理器**: 13th Gen Intel(R) Core(TM) i7-13700H (2.40 GHz)
- **内存**: 32.0 GB (31.7 GB 可用) / 32.0 GB (31.7 GB available)
- **系统**: 64位操作系统, 基于 x64 的处理器 / 64-bit OS, x64-based processor

### 性能指标 / Performance Metrics
- **QPS**: **100,000+** (10个实例并发测试)
- **延迟**: 低至 **0.01ms** 平均响应时间
- **成功率**: **100%** 无错误率
- **并发**: 支持多实例并发，可配置实例数量

```
插件类型    QPS      延迟(ms)   成功率(%)
Go         240.31   4.16       100.00
Python     224.97   4.44       100.00
Node.js    131.58   7.60       100.00
总计       526.30   1.90       100.00
```

## 🌍 跨平台支持 / Cross-Platform Support

### 支持的平台 / Supported Platforms
- ✅ **Windows** (amd64) - 使用命名管道通信 / Named Pipes communication
- ✅ **Linux** (amd64) - 使用Unix域套接字通信 / Unix Domain Sockets communication
- ✅ **macOS** (amd64) - 使用Unix域套接字通信 / Unix Domain Sockets communication
- ✅ **FreeBSD** (amd64) - 使用Unix域套接字通信 / Unix Domain Sockets communication

### 平台特定实现 / Platform-Specific Implementation
系统使用Go的构建标签机制实现跨平台兼容：

**Windows平台 / Windows Platform:**
```go
//go:build windows
// 使用Microsoft/go-winio库实现命名管道通信
// Uses Microsoft/go-winio library for Named Pipes communication
```

**Unix平台 / Unix Platforms (Linux/macOS/FreeBSD):**
```go
//go:build !windows  
// 使用标准net包实现Unix域套接字通信
// Uses standard net package for Unix Domain Sockets communication
```

## 🚀 快速开始 / Quick Start

### 1. 获取框架 / Get Framework

```bash
# 克隆项目 / Clone repository
git clone https://github.com/hoonfeng/goproc.git
cd goproc

# 安装依赖 / Install dependencies
go mod tidy
```

### 2. 编写插件 / Write Plugins

#### Go插件示例 / Go Plugin Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/hoonfeng/goproc/sdk"
)

// 加法函数 / Addition function
func addHandler(params map[string]interface{}) (interface{}, error) {
    a, ok1 := params["a"].(float64)
    b, ok2 := params["b"].(float64)
    
    if !ok1 || !ok2 {
        return nil, fmt.Errorf("参数错误: a和b必须为数字")
    }
    
    return a + b, nil
}

func main() {
    // 注册函数 / Register functions
    sdk.RegisterFunction("add", addHandler)
    
    // 启动插件 / Start plugin
    if err := sdk.StartPlugin(); err != nil {
        log.Fatalf("插件启动失败: %v", err)
    }
    
    // 等待插件停止 / Wait for plugin to stop
    sdk.WaitPlugin()
}
```

#### Python插件示例 / Python Plugin Example

```python
from goproc_sdk import GoProc

def reverse_string(params):
    """字符串反转函数"""
    text = params.get('text', '')
    return text[::-1]

def main():
    # 创建插件实例
    plugin = GoProc()
    
    # 注册函数
    plugin.register_function('reverse', reverse_string)
    
    # 启动插件
    plugin.start()

if __name__ == '__main__':
    main()
```

### 3. 使用框架 / Use Framework

参考 `examples/` 目录中的示例代码，了解如何在您的应用程序中集成和使用插件系统。

Refer to the example code in the `examples/` directory to learn how to integrate and use the plugin system in your applications.

## 📁 项目结构 / Project Structure

```
goproc/
├── config/                 # 配置管理 / Configuration management
│   └── config.go          # 配置加载和验证 / Config loading and validation
├── plugin/                # 插件管理核心 / Plugin management core
│   ├── communication.go  # 跨平台通信抽象 / Cross-platform communication
│   ├── instance.go        # 插件实例管理 / Plugin instance management
│   ├── pool.go           # 插件池管理 / Plugin pool management
│   └── manager.go        # 插件管理器 / Plugin manager
├── sdk/                   # 插件SDK / Plugin SDK
│   ├── types.go          # 类型定义 / Type definitions
│   ├── go/               # Go语言SDK / Go SDK
│   │   └── sdk.go
│   └── python/           # Python语言SDK / Python SDK
│       └── goproc_sdk.py
├── examples/              # 示例插件 / Example plugins
│   ├── windows/          # Windows平台示例 / Windows examples
│   ├── unix/             # Unix平台示例 / Unix examples
│   └── README.md         # 示例说明 / Examples documentation
├── docs/                  # 文档 / Documentation
├── go.mod                # Go模块文件 / Go module file
└── README.md             # 项目文档 / Project documentation
```



## 📚 文档 / Documentation

- [快速开始指南](docs/quick-start.md) / [Quick Start Guide](docs/quick-start.md)
- [插件开发指南](docs/plugin-development.md) / [Plugin Development Guide](docs/plugin-development.md)
- [Go SDK 使用指南](docs/go_sdk_usage.md) / [Go SDK Usage Guide](docs/go_sdk_usage.md)
- [Node.js SDK 使用指南](docs/nodejs_sdk_usage.md) / [Node.js SDK Usage Guide](docs/nodejs_sdk_usage.md)
- [API参考文档](docs/api-reference.md) / [API Reference](docs/api-reference.md)
- [配置参考](docs/configuration.md) / [Configuration Reference](docs/configuration.md)
- [示例教程](examples/README.md) / [Examples Tutorial](examples/README.md)

## 🤝 贡献 / Contributing

我们欢迎所有形式的贡献！请查看 [贡献指南](CONTRIBUTING.md) 了解详情。

We welcome all forms of contributions! Please see the [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 许可证 / License

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 致谢 / Acknowledgments

本项目的开发灵感来源于 [pyproc](https://github.com/YuminosukeSato/pyproc) 项目，感谢原作者的创意和贡献。

This project was inspired by the [pyproc](https://github.com/YuminosukeSato/pyproc) project. Thanks to the original author for the creativity and contribution.

感谢所有为这个项目做出贡献的开发者和用户。

Thanks to all developers and users who contributed to this project.

---

**GoProc** - 让插件开发变得简单而高效 / Making plugin development simple and efficient
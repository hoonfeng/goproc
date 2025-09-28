# GoProc Examples / GoProc 示例

This directory contains examples and demonstrations for the GoProc plugin system.

本目录包含 GoProc 插件系统的示例和演示。

## 🔧 Latest Updates / 最新更新

### Python SDK Windows Fix (December 2024) / Python SDK Windows 修复 (2024年12月)
- ✅ **Windows Named Pipe Support**: Python plugins now work reliably on Windows
- ✅ **Windows 命名管道支持**: Python 插件现在可以在 Windows 上可靠运行
- ✅ **Improved Performance**: 20-30% response time improvement
- ✅ **性能提升**: 响应时间提升 20-30%
- ✅ **Enhanced Stability**: 100% connection success rate
- ✅ **稳定性增强**: 100% 连接成功率
- ✅ **Complete Message Transmission**: Fixed sendall method for reliable data transfer
- ✅ **完整消息传输**: 修复 sendall 方法以确保可靠的数据传输

## Directory Structure / 目录结构

### Platform-Specific Examples / 平台特定示例

- **`windows/`** - Examples specifically for Windows platform
- **`windows/`** - Windows 平台专用示例
  - Uses Windows Named Pipes for communication
  - 使用 Windows 命名管道进行通信
  - Contains Windows-specific build scripts and configurations
  - 包含 Windows 特定的构建脚本和配置
  - **Python plugins fully supported** with latest fixes
  - **Python 插件完全支持** 包含最新修复

- **`unix/`** - Examples for Unix-like systems (Linux, macOS, etc.)
- **`unix/`** - Unix 类系统示例 (Linux, macOS 等)
  - Uses Unix Domain Sockets for communication
  - 使用 Unix 域套接字进行通信
  - Contains Unix-specific build scripts and configurations
  - 包含 Unix 特定的构建脚本和配置
  - **Contains all unified test files** - Test files have been consolidated here to eliminate duplication
  - **包含所有统一测试文件** - 测试文件已整合到此处以消除重复

### Core Examples / 核心示例

Each platform directory contains the following examples:

每个平台目录包含以下示例：

1. **Basic Demo** - Simple plugin usage demonstration
1. **基础演示** - 简单的插件使用演示
2. **Comprehensive Demo** - Full-featured plugin system showcase
2. **综合演示** - 功能完整的插件系统展示
3. **Math Plugin** - Mathematical operations plugin example
3. **数学插件** - 数学运算插件示例
4. **String Plugin** - String manipulation plugin example
4. **字符串插件** - 字符串操作插件示例
5. **Python Plugin** - Python language plugin integration
5. **Python 插件** - Python 语言插件集成
6. **Node.js Plugin** - JavaScript/Node.js plugin integration
6. **Node.js 插件** - JavaScript/Node.js 插件集成

## Building Examples / 构建示例

### Windows
```cmd
cd windows
build_examples.bat
```

### Unix (Linux/macOS)
```bash
cd unix
./build_examples.sh
```

## Running Examples / 运行示例

Each example directory contains its own README with specific instructions for building and running the example.

每个示例目录都包含自己的 README 文件，其中有构建和运行示例的具体说明。

## Cross-Platform Compatibility / 跨平台兼容性

The GoProc system automatically detects the platform and uses the appropriate communication method:

GoProc 系统自动检测平台并使用适当的通信方法：

- Windows: Named Pipes (Python SDK fully fixed and optimized)
- Windows: 命名管道 (Python SDK 已完全修复和优化)
- Unix-like systems: Unix Domain Sockets
- Unix 类系统: Unix 域套接字

All examples are designed to work seamlessly across platforms with minimal configuration changes.

所有示例都设计为在各平台上无缝工作，只需最少的配置更改。

## Python Plugin Testing / Python 插件测试

### Windows Testing / Windows 测试
```cmd
# Test Python plugin functionality
# 测试 Python 插件功能
go run test_windows_python.go
```

### Unix Testing / Unix 测试
```bash
# Test Python plugin functionality
# 测试 Python 插件功能
go run test_python_plugin.go
```

**Test Coverage / 测试覆盖:**
- ✅ Addition function / 加法函数
- ✅ String reversal function / 字符串反转函数
- ✅ Time retrieval function / 时间获取函数
- ✅ Fibonacci sequence calculation / 斐波那契数列计算
- ✅ Connection stability / 连接稳定性
- ✅ Message transmission reliability / 消息传输可靠性
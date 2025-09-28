# Unix Platform Examples / Unix 平台示例

This directory contains examples specifically designed for Unix-like systems (Linux, macOS, FreeBSD).

本目录包含专为 Unix 类系统（Linux、macOS、FreeBSD）设计的示例。

## Communication Method / 通信方法
- **Unix Domain Sockets** - Used for inter-process communication between the main process and plugins
- **Unix 域套接字** - 用于主进程和插件之间的进程间通信

## Test Files / 测试文件
**Note**: All test files have been unified in this directory to eliminate duplication. These tests work on both Unix and Windows platforms due to the cross-platform nature of the GoProc system. Use `run_tests.sh` on Unix systems or `run_tests.bat` on Windows to run all tests.

**注意**: 所有测试文件已统一在此目录中以消除重复。由于 GoProc 系统的跨平台特性，这些测试在 Unix 和 Windows 平台上都能工作。在 Unix 系统上使用 `run_tests.sh`，在 Windows 上使用 `run_tests.bat` 来运行所有测试。

## Building Examples / 构建示例

### Prerequisites / 前置要求
- Go 1.19 or later / Go 1.19 或更高版本
- Python 3.7+ (for Python plugin) / Python 3.7+（用于 Python 插件）
- Node.js 14+ (for Node.js plugin) / Node.js 14+（用于 Node.js 插件）

### Build All Examples / 构建所有示例
```bash
./build_examples.sh
```

### Build Individual Examples / 构建单个示例
```bash
# Math Plugin / 数学插件
go build -o math_plugin/math_plugin math_plugin/main.go
chmod +x math_plugin/math_plugin

# String Plugin / 字符串插件
go build -o string_plugin/string_plugin string_plugin/main.go
chmod +x string_plugin/string_plugin

# Demo Applications / 演示应用程序
go build -o basic_demo/basic_demo basic_demo/main.go
go build -o comprehensive_demo/comprehensive_demo comprehensive_demo/main.go
go build -o demo_app/demo_app demo_app/main.go
go build -o performance_test/performance_test performance_test/main.go
go build -o plugin_manager/plugin_manager plugin_manager/main.go
```

## Running Examples / 运行示例

### Basic Demo / 基础演示
```bash
cd basic_demo
./basic_demo
```

### Comprehensive Demo / 综合演示
```bash
cd comprehensive_demo
./comprehensive_demo
```

### Performance Test / 性能测试
```bash
cd performance_test
./performance_test
```

## Plugin Configuration / 插件配置

Unix plugins don't require file extensions and use Unix Domain Sockets:

Unix 插件不需要文件扩展名，使用 Unix 域套接字：

```yaml
plugins:
  math_plugin:
    type: "binary"
    path: "./math_plugin"
    pool_size: 5
  string_plugin:
    type: "binary"
    path: "./string_plugin"
    pool_size: 3
```

## Cross-Platform Building / 跨平台构建

You can build for multiple Unix platforms:

您可以为多个 Unix 平台构建：

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o math_plugin/math_plugin-linux math_plugin/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o math_plugin/math_plugin-macos math_plugin/main.go

# FreeBSD
GOOS=freebsd GOARCH=amd64 go build -o math_plugin/math_plugin-freebsd math_plugin/main.go
```

## Notes / 注意事项

- Binary plugins don't need file extensions on Unix systems
- 二进制插件在 Unix 系统上不需要文件扩展名
- Unix Domain Sockets are automatically managed by the GoProc system
- Unix 域套接字由 GoProc 系统自动管理
- Make sure to set execute permissions (`chmod +x`) for compiled binaries
- 确保为编译的二进制文件设置执行权限（`chmod +x`）
- Ensure all dependencies (Python, Node.js) are installed and in PATH
- 确保所有依赖项（Python、Node.js）已安装并在 PATH 中
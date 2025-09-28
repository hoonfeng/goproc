# Windows Platform Examples / Windows 平台示例

This directory contains examples specifically designed for Windows platform.

本目录包含专为 Windows 平台设计的示例。

## Communication Method / 通信方法
- **Windows Named Pipes** - Used for inter-process communication between the main process and plugins
- **Windows 命名管道** - 用于主进程和插件之间的进程间通信

## Test Files / 测试文件
**Note**: Test files have been unified and moved to the `examples/unix/` directory to eliminate duplication. These tests work on both Unix and Windows platforms due to the cross-platform nature of the GoProc system. Use `run_tests.bat` to run all tests on Windows.

**注意**: 测试文件已统一并移至 `examples/unix/` 目录以消除重复。由于 GoProc 系统的跨平台特性，这些测试在 Unix 和 Windows 平台上都能工作。在 Windows 上使用 `run_tests.bat` 运行所有测试。

## Building Examples / 构建示例

### Prerequisites / 前置要求
- Go 1.19 or later / Go 1.19 或更高版本
- Python 3.7+ (for Python plugin) / Python 3.7+（用于 Python 插件）
- Node.js 14+ (for Node.js plugin) / Node.js 14+（用于 Node.js 插件）

### Build All Examples / 构建所有示例
```cmd
build_examples.ps1
```

### Build Individual Examples / 构建单个示例
```cmd
# Math Plugin / 数学插件
go build -o math_plugin\math_plugin.exe math_plugin\main.go

# String Plugin / 字符串插件
go build -o string_plugin\string_plugin.exe string_plugin\main.go

# Demo Applications / 演示应用程序
go build -o basic_demo\basic_demo.exe basic_demo\main.go
go build -o comprehensive_demo\comprehensive_demo.exe comprehensive_demo\main.go
go build -o demo_app\demo_app.exe demo_app\main.go
go build -o performance_test\performance_test.exe performance_test\main.go
go build -o plugin_manager\plugin_manager.exe plugin_manager\main.go
```

## Running Examples / 运行示例

### Basic Demo / 基础演示
```cmd
cd basic_demo
basic_demo.exe
```

### Comprehensive Demo / 综合演示
```cmd
cd comprehensive_demo
comprehensive_demo.exe
```

### Performance Test / 性能测试
```cmd
cd performance_test
performance_test.exe
```

## Plugin Configuration / 插件配置

Windows plugins use `.exe` extension and Named Pipes for communication:

Windows 插件使用 `.exe` 扩展名和命名管道进行通信：

```yaml
plugins:
  math_plugin:
    type: "binary"
    path: "./math_plugin.exe"
    pool_size: 5
  string_plugin:
    type: "binary"
    path: "./string_plugin.exe"
    pool_size: 3
```

## Notes / 注意事项

- All binary plugins must have `.exe` extension on Windows
- 所有二进制插件在 Windows 上必须有 `.exe` 扩展名
- Named Pipes are automatically managed by the GoProc system
- 命名管道由 GoProc 系统自动管理
- Ensure all dependencies (Python, Node.js) are installed and in PATH
- 确保所有依赖项（Python、Node.js）已安装并在 PATH 中
# GoProc Go SDK 使用指南 / GoProc Go SDK Usage Guide

## 概述 / Overview

GoProc Go SDK 是一个高性能的插件开发框架，支持跨平台通信（Windows命名管道和Unix域套接字）。SDK提供了简洁的API接口，让开发者专注于业务逻辑实现。

GoProc Go SDK is a high-performance plugin development framework that supports cross-platform communication (Windows Named Pipes and Unix Domain Sockets). The SDK provides a clean API interface, allowing developers to focus on business logic implementation.

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

## 🚀 快速开始 / Quick Start

### 1. 导入SDK / Import SDK

```go
package main

import (
    "fmt"
    "log"
    "github.com/hoonfeng/goproc/sdk"
)
```

### 2. 注册函数 / Register Functions

#### 方法一：使用全局函数 / Method 1: Using Global Functions

```go
// 加法函数 / Addition function
func addHandler(params map[string]interface{}) (interface{}, error) {
    a, ok1 := params["a"].(float64)
    b, ok2 := params["b"].(float64)
    
    if !ok1 || !ok2 {
        return nil, fmt.Errorf("参数错误: a和b必须为数字")
    }
    
    return a + b, nil
}

// 减法函数 / Subtraction function
func subtractHandler(params map[string]interface{}) (interface{}, error) {
    a, ok1 := params["a"].(float64)
    b, ok2 := params["b"].(float64)
    
    if !ok1 || !ok2 {
        return nil, fmt.Errorf("参数错误: a和b必须为数字")
    }
    
    return a - b, nil
}

func main() {
    // 注册函数 / Register functions
    if err := sdk.RegisterFunction("add", addHandler); err != nil {
        log.Fatalf("注册函数失败: %v", err)
    }
    
    if err := sdk.RegisterFunction("subtract", subtractHandler); err != nil {
        log.Fatalf("注册函数失败: %v", err)
    }
    
    // 启动插件 / Start plugin
    if err := sdk.Start(); err != nil {
        log.Fatalf("插件启动失败: %v", err)
    }
    
    // 等待插件停止 / Wait for plugin to stop
    sdk.Wait()
}
```

#### 方法二：使用SDK实例 / Method 2: Using SDK Instance

```go
func main() {
    // 创建SDK实例 / Create SDK instance
    pluginSDK := sdk.NewPluginSDK()
    
    // 注册函数 / Register functions
    pluginSDK.RegisterFunction("multiply", func(params map[string]interface{}) (interface{}, error) {
        a, ok1 := params["a"].(float64)
        b, ok2 := params["b"].(float64)
        
        if !ok1 || !ok2 {
            return nil, fmt.Errorf("参数错误: a和b必须为数字")
        }
        
        return a * b, nil
    })
    
    pluginSDK.RegisterFunction("divide", func(params map[string]interface{}) (interface{}, error) {
        a, ok1 := params["a"].(float64)
        b, ok2 := params["b"].(float64)
        
        if !ok1 || !ok2 {
            return nil, fmt.Errorf("参数错误: a和b必须为数字")
        }
        
        if b == 0 {
            return nil, fmt.Errorf("除数不能为零")
        }
        
        return a / b, nil
    })
    
    // 启动插件 / Start plugin
    if err := pluginSDK.Start(); err != nil {
        log.Fatalf("插件启动失败: %v", err)
    }
    
    // 等待插件停止 / Wait for plugin to stop
    pluginSDK.Wait()
}
```

## 📚 API 参考 / API Reference

### 核心类型 / Core Types

#### FunctionHandler

```go
type FunctionHandler func(params map[string]interface{}) (interface{}, error)
```

函数处理器类型，用于定义插件函数的签名。

Function handler type used to define the signature of plugin functions.

**参数 / Parameters:**
- `params`: 函数参数映射 / Function parameter mapping

**返回值 / Returns:**
- `interface{}`: 函数执行结果 / Function execution result
- `error`: 错误信息（如果有）/ Error information (if any)

#### PluginSDK

```go
type PluginSDK struct {
    // 私有字段 / Private fields
}
```

插件SDK主要结构体，提供插件开发的核心功能。

Main plugin SDK structure providing core functionality for plugin development.

### 全局函数 / Global Functions

#### RegisterFunction

```go
func RegisterFunction(name string, handler FunctionHandler) error
```

注册插件函数到全局SDK实例。

Register a plugin function to the global SDK instance.

**参数 / Parameters:**
- `name`: 函数名称 / Function name
- `handler`: 函数处理器 / Function handler

**返回值 / Returns:**
- `error`: 错误信息（如果有）/ Error information (if any)

#### Start

```go
func Start() error
```

启动全局SDK实例。

Start the global SDK instance.

**返回值 / Returns:**
- `error`: 错误信息（如果有）/ Error information (if any)

#### Wait

```go
func Wait()
```

等待全局SDK实例停止。

Wait for the global SDK instance to stop.

#### Stop

```go
func Stop()
```

停止全局SDK实例。

Stop the global SDK instance.

### SDK实例方法 / SDK Instance Methods

#### NewPluginSDK

```go
func NewPluginSDK() *PluginSDK
```

创建新的插件SDK实例。

Create a new plugin SDK instance.

**返回值 / Returns:**
- `*PluginSDK`: SDK实例指针 / SDK instance pointer

#### RegisterFunction (实例方法 / Instance Method)

```go
func (sdk *PluginSDK) RegisterFunction(name string, handler FunctionHandler) error
```

注册插件函数到SDK实例。

Register a plugin function to the SDK instance.

#### Start (实例方法 / Instance Method)

```go
func (sdk *PluginSDK) Start() error
```

启动SDK实例。

Start the SDK instance.

#### Wait (实例方法 / Instance Method)

```go
func (sdk *PluginSDK) Wait()
```

等待SDK实例停止。

Wait for the SDK instance to stop.

#### Stop (实例方法 / Instance Method)

```go
func (sdk *PluginSDK) Stop()
```

停止SDK实例。

Stop the SDK instance.

## 🔧 高级用法 / Advanced Usage

### 复杂数据类型处理 / Complex Data Type Handling

```go
// 处理结构体数据 / Handle struct data
func processUserHandler(params map[string]interface{}) (interface{}, error) {
    // 解析用户数据 / Parse user data
    userData, ok := params["user"].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("无效的用户数据格式")
    }
    
    name, _ := userData["name"].(string)
    age, _ := userData["age"].(float64)
    
    // 处理业务逻辑 / Process business logic
    result := map[string]interface{}{
        "processed": true,
        "message":   fmt.Sprintf("用户 %s，年龄 %.0f 岁，处理完成", name, age),
        "timestamp": time.Now().Unix(),
    }
    
    return result, nil
}

// 处理数组数据 / Handle array data
func sumArrayHandler(params map[string]interface{}) (interface{}, error) {
    numbersInterface, ok := params["numbers"].([]interface{})
    if !ok {
        return nil, fmt.Errorf("参数 numbers 必须是数组")
    }
    
    var sum float64
    for i, numInterface := range numbersInterface {
        num, ok := numInterface.(float64)
        if !ok {
            return nil, fmt.Errorf("数组第 %d 个元素不是数字", i)
        }
        sum += num
    }
    
    return sum, nil
}
```

### 错误处理最佳实践 / Error Handling Best Practices

```go
func validateAndProcessHandler(params map[string]interface{}) (interface{}, error) {
    // 参数验证 / Parameter validation
    value, exists := params["value"]
    if !exists {
        return nil, fmt.Errorf("缺少必需参数: value")
    }
    
    strValue, ok := value.(string)
    if !ok {
        return nil, fmt.Errorf("参数 value 必须是字符串类型，实际类型: %T", value)
    }
    
    if len(strValue) == 0 {
        return nil, fmt.Errorf("参数 value 不能为空")
    }
    
    // 业务逻辑处理 / Business logic processing
    if len(strValue) > 100 {
        return nil, fmt.Errorf("参数 value 长度不能超过100个字符")
    }
    
    // 返回处理结果 / Return processing result
    return map[string]interface{}{
        "original": strValue,
        "length":   len(strValue),
        "upper":    strings.ToUpper(strValue),
        "lower":    strings.ToLower(strValue),
    }, nil
}
```

### 异步处理 / Asynchronous Processing

```go
func asyncProcessHandler(params map[string]interface{}) (interface{}, error) {
    taskID, _ := params["task_id"].(string)
    if taskID == "" {
        taskID = fmt.Sprintf("task_%d", time.Now().UnixNano())
    }
    
    // 启动异步任务 / Start asynchronous task
    go func() {
        // 模拟长时间运行的任务 / Simulate long-running task
        time.Sleep(5 * time.Second)
        
        // 这里可以通过其他方式通知任务完成
        // Here you can notify task completion through other means
        log.Printf("任务 %s 完成", taskID)
    }()
    
    // 立即返回任务ID / Return task ID immediately
    return map[string]interface{}{
        "task_id": taskID,
        "status":  "started",
        "message": "任务已启动，正在后台处理",
    }, nil
}
```

## 🌍 跨平台支持 / Cross-Platform Support

SDK自动检测运行平台并使用相应的通信机制：

The SDK automatically detects the running platform and uses the appropriate communication mechanism:

- **Windows**: 使用命名管道 (Named Pipes) / Uses Named Pipes
- **Linux/macOS/FreeBSD**: 使用Unix域套接字 (Unix Domain Sockets) / Uses Unix Domain Sockets

开发者无需关心底层通信细节，SDK会自动处理平台差异。

Developers don't need to worry about underlying communication details; the SDK automatically handles platform differences.

## 🔍 调试和日志 / Debugging and Logging

### 启用调试模式 / Enable Debug Mode

```go
import (
    "log"
    "os"
)

func main() {
    // 设置日志输出 / Set log output
    log.SetOutput(os.Stdout)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    // 注册函数 / Register functions
    sdk.RegisterFunction("debug_test", func(params map[string]interface{}) (interface{}, error) {
        log.Printf("收到调用，参数: %+v", params)
        return "调试测试成功", nil
    })
    
    // 启动插件 / Start plugin
    if err := sdk.Start(); err != nil {
        log.Fatalf("插件启动失败: %v", err)
    }
    
    log.Println("插件已启动，等待调用...")
    sdk.Wait()
}
```

## 📝 完整示例 / Complete Example

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "time"
    "github.com/hoonfeng/goproc/sdk"
)

// 字符串处理插件 / String processing plugin
func main() {
    // 字符串反转 / String reverse
    sdk.RegisterFunction("reverse", func(params map[string]interface{}) (interface{}, error) {
        text, ok := params["text"].(string)
        if !ok {
            return nil, fmt.Errorf("参数 text 必须是字符串")
        }
        
        runes := []rune(text)
        for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
            runes[i], runes[j] = runes[j], runes[i]
        }
        
        return string(runes), nil
    })
    
    // 字符串转大写 / String to uppercase
    sdk.RegisterFunction("uppercase", func(params map[string]interface{}) (interface{}, error) {
        text, ok := params["text"].(string)
        if !ok {
            return nil, fmt.Errorf("参数 text 必须是字符串")
        }
        
        return strings.ToUpper(text), nil
    })
    
    // 字符串转小写 / String to lowercase
    sdk.RegisterFunction("lowercase", func(params map[string]interface{}) (interface{}, error) {
        text, ok := params["text"].(string)
        if !ok {
            return nil, fmt.Errorf("参数 text 必须是字符串")
        }
        
        return strings.ToLower(text), nil
    })
    
    // 字符串统计 / String statistics
    sdk.RegisterFunction("stats", func(params map[string]interface{}) (interface{}, error) {
        text, ok := params["text"].(string)
        if !ok {
            return nil, fmt.Errorf("参数 text 必须是字符串")
        }
        
        return map[string]interface{}{
            "length":     len(text),
            "rune_count": len([]rune(text)),
            "word_count": len(strings.Fields(text)),
            "line_count": len(strings.Split(text, "\n")),
        }, nil
    })
    
    log.Println("字符串处理插件启动中...")
    
    // 启动插件 / Start plugin
    if err := sdk.Start(); err != nil {
        log.Fatalf("插件启动失败: %v", err)
    }
    
    log.Println("字符串处理插件已启动，等待调用...")
    
    // 等待插件停止 / Wait for plugin to stop
    sdk.Wait()
    
    log.Println("字符串处理插件已停止")
}
```

## 🚨 注意事项 / Important Notes

1. **线程安全**: SDK是线程安全的，可以在多个goroutine中使用 / The SDK is thread-safe and can be used in multiple goroutines

2. **函数注册**: 必须在调用`Start()`之前注册所有函数 / All functions must be registered before calling `Start()`

3. **参数类型**: JSON反序列化后，数字类型统一为`float64` / After JSON deserialization, numeric types are unified as `float64`

4. **错误处理**: 始终检查和处理错误，提供有意义的错误信息 / Always check and handle errors, providing meaningful error messages

5. **资源清理**: 使用`defer`语句确保资源正确清理 / Use `defer` statements to ensure proper resource cleanup

## 🔗 相关链接 / Related Links

- [项目主页 / Project Homepage](https://github.com/hoonfeng/goproc)
- [Python SDK 使用指南 / Python SDK Usage Guide](python_sdk_usage.md)
- [Node.js SDK 使用指南 / Node.js SDK Usage Guide](nodejs_sdk_usage.md)
- [跨平台使用指南 / Cross-Platform Usage Guide](CROSS_PLATFORM_USAGE_GUIDE.md)
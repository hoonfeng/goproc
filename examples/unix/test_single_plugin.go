package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hoonfeng/goproc/config"
	"github.com/hoonfeng/goproc/plugin"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run test_single_plugin.go <plugin_name>")
		fmt.Println("可用插件: math, string, node, python")
		os.Exit(1)
	}

	pluginName := os.Args[1]
	
	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败: %v", err)
	}

	// 创建插件配置
	var pluginConfig *config.PluginConfig
	
	switch pluginName {
	case "math":
		pluginConfig = &config.PluginConfig{
			Type:                config.PluginTypeBinary,
			Path:                filepath.Join(currentDir, "math_plugin", "math_plugin-linux"),
			PoolSize:            1,
			MaxInstances:        1,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"add", "multiply", "subtract", "divide"},
		}
	case "string":
		pluginConfig = &config.PluginConfig{
			Type:                config.PluginTypeBinary,
			Path:                filepath.Join(currentDir, "string_plugin", "string_plugin-linux"),
			PoolSize:            1,
			MaxInstances:        1,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"toUpper", "toLower", "reverse", "trim"},
		}
	case "node":
		pluginConfig = &config.PluginConfig{
			Type:                config.PluginTypeScript,
			Interpreter:         "node",
			ScriptPath:          filepath.Join(currentDir, "node_plugin", "index.js"),
			PoolSize:            1,
			MaxInstances:        1,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"jsonStringify", "arraySum", "timestamp", "uuid"},
		}
	case "python":
		pluginConfig = &config.PluginConfig{
			Type:                config.PluginTypeScript,
			Interpreter:         "python3",
			ScriptPath:          filepath.Join(currentDir, "python_plugin", "main.py"),
			PoolSize:            1,
			MaxInstances:        1,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"datetime_utils", "text_processing", "generate_id"},
		}
	default:
		log.Fatalf("未知插件: %s", pluginName)
	}

	fmt.Printf("=== 测试 %s 插件 ===\n", pluginName)
	fmt.Printf("插件类型: %s\n", pluginConfig.Type)
	if pluginConfig.Type == config.PluginTypeBinary {
		fmt.Printf("插件路径: %s\n", pluginConfig.Path)
	} else {
		fmt.Printf("解释器: %s\n", pluginConfig.Interpreter)
		fmt.Printf("脚本路径: %s\n", pluginConfig.ScriptPath)
	}

	// 创建插件池
	pool := plugin.NewPluginPool(pluginName, pluginConfig)
	
	// 启动插件池
	fmt.Println("正在启动插件...")
	if err := pool.Start(); err != nil {
		log.Fatalf("启动插件失败: %v", err)
	}
	fmt.Println("插件启动成功!")

	// 等待插件初始化
	fmt.Println("等待插件初始化...")
	time.Sleep(2 * time.Second)

	// 测试插件功能
	fmt.Println("测试插件功能...")
	
	switch pluginName {
	case "math":
		testMathPlugin(pool)
	case "string":
		testStringPlugin(pool)
	case "node":
		testNodePlugin(pool)
	case "python":
		testPythonPlugin(pool)
	}

	// 停止插件池
	fmt.Println("停止插件...")
	pool.Stop()
	fmt.Println("测试完成!")
}

func testMathPlugin(pool *plugin.PluginPool) {
	params := map[string]interface{}{
		"a": 10,
		"b": 5,
	}
	
	result, err := pool.CallFunction("add", params)
	if err != nil {
		fmt.Printf("调用add函数失败: %v\n", err)
	} else {
		fmt.Printf("add(10, 5) = %v\n", result)
	}
}

func testStringPlugin(pool *plugin.PluginPool) {
	params := map[string]interface{}{
		"text": "hello world",
	}
	
	result, err := pool.CallFunction("toUpper", params)
	if err != nil {
		fmt.Printf("调用toUpper函数失败: %v\n", err)
	} else {
		fmt.Printf("toUpper('hello world') = %v\n", result)
	}
}

func testNodePlugin(pool *plugin.PluginPool) {
	params := map[string]interface{}{
		"data": map[string]interface{}{"test": "value"},
	}
	
	result, err := pool.CallFunction("jsonStringify", params)
	if err != nil {
		fmt.Printf("调用jsonStringify函数失败: %v\n", err)
	} else {
		fmt.Printf("jsonStringify({test: 'value'}) = %v\n", result)
	}
}

func testPythonPlugin(pool *plugin.PluginPool) {
	params := map[string]interface{}{}
	
	result, err := pool.CallFunction("datetime_utils", params)
	if err != nil {
		fmt.Printf("调用datetime_utils函数失败: %v\n", err)
	} else {
		fmt.Printf("datetime_utils() = %v\n", result)
	}
}
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

// main 主函数 - 演示插件系统的基本使用方法
// Main function - demonstrates basic usage of the plugin system
func main() {
	fmt.Println("=== 插件系统基本使用演示 ===")
	fmt.Println("=== Plugin System Basic Usage Demo ===")

	// 获取插件路径 / Get plugin paths
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败 / Failed to get current directory: %v", err)
	}
	// 如果当前在basic_demo目录，则回到examples目录 / If in basic_demo directory, go back to examples
	examplesDir := filepath.Dir(currentDir)
	if filepath.Base(currentDir) != "basic_demo" {
		// 如果不在basic_demo目录，假设在goproc根目录 / If not in basic_demo, assume in goproc root
		examplesDir = filepath.Join(currentDir, "examples")
	}

	// 创建插件配置 / Create plugin configurations
	configs := createPluginConfigs(examplesDir)

	// 创建并启动插件池 / Create and start plugin pools
	pools := make(map[string]*plugin.PluginPool)
	for name, config := range configs {
		pool := plugin.NewPluginPool(name, config)
		if err := pool.Start(); err != nil {
			log.Fatalf("启动%s插件失败 / Failed to start %s plugin: %v", name, name, err)
		}
		pools[name] = pool
		defer pool.Stop()
	}

	// 等待插件初始化 / Wait for plugin initialization
	time.Sleep(2 * time.Second)

	// 演示各插件功能 / Demonstrate plugin functions
	demonstratePlugins(pools)

	fmt.Println("\n=== 演示完成 ===")
	fmt.Println("=== Demo Completed ===")
}

// createPluginConfigs 创建插件配置
// Create plugin configurations
func createPluginConfigs(examplesDir string) map[string]*config.PluginConfig {
	return map[string]*config.PluginConfig{
		"math": {
			Type:                config.PluginTypeBinary,
			Path:                filepath.Join(examplesDir, "math_plugin", "math_plugin.exe"),
			PoolSize:            2,
			MaxInstances:        5,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"add", "multiply"},
		},
		"string": {
			Type:                config.PluginTypeBinary,
			Path:                filepath.Join(examplesDir, "string_plugin", "string_plugin.exe"),
			PoolSize:            2,
			MaxInstances:        5,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"toUpper", "reverse"},
		},
		"node": {
			Type:                config.PluginTypeScript,
			Interpreter:         "node",
			ScriptPath:          filepath.Join(examplesDir, "node_plugin", "index.js"),
			PoolSize:            2,
			MaxInstances:        5,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"jsonStringify", "arraySum"},
		},
		"python": {
			Type:                config.PluginTypeScript,
			Interpreter:         "python",
			ScriptPath:          filepath.Join(examplesDir, "python_plugin", "main.py"),
			PoolSize:            2,
			MaxInstances:        5,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"datetime_utils", "text_processing"},
		},
	}
}

// demonstratePlugins 演示各插件的功能
// Demonstrate plugin functions
func demonstratePlugins(pools map[string]*plugin.PluginPool) {
	fmt.Println("\n=== 功能演示 ===")
	fmt.Println("=== Function Demonstrations ===")

	// 数学插件演示 / Math plugin demo
	fmt.Println("\n1. 数学插件 / Math Plugin:")
	result, err := pools["math"].CallFunction("add", map[string]interface{}{"a": 15.0, "b": 25.0})
	if err != nil {
		fmt.Printf("   ❌ 加法失败 / Addition failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ 15 + 25 = %v\n", result)
	}

	result, err = pools["math"].CallFunction("multiply", map[string]interface{}{"a": 6.0, "b": 7.0})
	if err != nil {
		fmt.Printf("   ❌ 乘法失败 / Multiplication failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ 6 × 7 = %v\n", result)
	}

	// 字符串插件演示 / String plugin demo
	fmt.Println("\n2. 字符串插件 / String Plugin:")
	result, err = pools["string"].CallFunction("toUpper", map[string]interface{}{"str": "hello world"})
	if err != nil {
		fmt.Printf("   ❌ 大写转换失败 / Uppercase conversion failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ 'hello world' → '%v'\n", result)
	}

	result, err = pools["string"].CallFunction("reverse", map[string]interface{}{"str": "golang"})
	if err != nil {
		fmt.Printf("   ❌ 字符串反转失败 / String reverse failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ 'golang' → '%v'\n", result)
	}

	// Node.js插件演示 / Node.js plugin demo
	fmt.Println("\n3. Node.js插件 / Node.js Plugin:")
	testObj := map[string]interface{}{"name": "demo", "value": 42}
	result, err = pools["node"].CallFunction("jsonStringify", map[string]interface{}{"obj": testObj})
	if err != nil {
		fmt.Printf("   ❌ JSON序列化失败 / JSON stringify failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ JSON序列化 / JSON stringify: %v\n", result)
	}

	result, err = pools["node"].CallFunction("arraySum", map[string]interface{}{"array": []float64{1, 2, 3, 4, 5}})
	if err != nil {
		fmt.Printf("   ❌ 数组求和失败 / Array sum failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ [1,2,3,4,5] 求和 / sum = %v\n", result)
	}

	// Python插件演示 / Python plugin demo
	fmt.Println("\n4. Python插件 / Python Plugin:")
	result, err = pools["python"].CallFunction("datetime_utils", map[string]interface{}{"operation": "now"})
	if err != nil {
		fmt.Printf("   ❌ 时间获取失败 / Datetime failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ 当前时间 / Current time: %v\n", result)
	}

	result, err = pools["python"].CallFunction("text_processing", map[string]interface{}{
		"text":      "Hello World from Python",
		"operation": "word_count",
	})
	if err != nil {
		fmt.Printf("   ❌ 文本处理失败 / Text processing failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ 'Hello World from Python' 词数 / word count: %v\n", result)
	}
}

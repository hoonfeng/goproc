package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"goproc/config"
	"goproc/plugin"
)

func main() {
	fmt.Println("=== Python插件简化测试 ===")

	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败: %v", err)
	}
	fmt.Printf("✓ 当前工作目录: %s\n", cwd)

	// 确定Python插件路径
	pythonPluginPath := filepath.Join(cwd, "python_plugin", "main.py")
	fmt.Printf("✓ Python插件路径: %s\n", pythonPluginPath)

	// 检查插件文件是否存在
	if _, err := os.Stat(pythonPluginPath); os.IsNotExist(err) {
		log.Fatalf("❌ Python插件文件不存在: %s", pythonPluginPath)
	}
	fmt.Printf("✓ Python插件文件存在\n")

	// 创建插件配置
	pythonConfig := &config.PluginConfig{
		Type:                config.PluginTypeScript,
		Interpreter:         "python3",
		ScriptPath:          pythonPluginPath,
		PoolSize:            1,
		MaxInstances:        3,
		HealthCheckInterval: 30 * time.Second,
		Functions:           []string{"add", "subtract", "multiply", "divide"},
		Environment:         map[string]string{
			"PLUGIN_NAME": "python_plugin",
			"DEBUG": "true",
		},
	}

	fmt.Printf("✓ 插件配置创建完成\n")

	// 确保sockets目录存在
	socketsDir := filepath.Join(cwd, "sockets")
	if err := os.MkdirAll(socketsDir, 0755); err != nil {
		log.Fatalf("❌ 创建sockets目录失败: %v", err)
	}
	fmt.Printf("✓ sockets目录已准备好\n")

	// 创建单个插件实例
	fmt.Println("\n=== 创建插件实例 ===")
	instance := plugin.NewPluginInstance("python", pythonConfig, "python-test-001")
	fmt.Printf("✓ 插件实例创建成功\n")

	// 启动插件实例
	fmt.Println("\n=== 启动插件实例 ===")
	if err := instance.Start(); err != nil {
		log.Fatalf("❌ 启动插件实例失败: %v", err)
	}
	fmt.Printf("✓ 插件实例启动成功\n")

	// 等待插件初始化
	fmt.Println("\n=== 等待插件初始化 ===")
	time.Sleep(3 * time.Second)
	fmt.Printf("✓ 插件初始化等待完成\n")

	// 测试插件功能
	fmt.Println("\n=== 测试插件功能 ===")

	// 1. 测试加法函数
	fmt.Println("\n1. 测试加法函数:")
	addParams := map[string]interface{}{
		"a": 15,
		"b": 25,
	}
	fmt.Printf("  调用参数: %+v\n", addParams)
	result, err := instance.CallFunction("add", addParams)
	if err != nil {
		fmt.Printf("  ❌ 加法调用失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ 加法结果: 15 + 25 = %v\n", result)
	}

	// 2. 测试减法函数
	fmt.Println("\n2. 测试减法函数:")
	subParams := map[string]interface{}{
		"a": 50,
		"b": 20,
	}
	fmt.Printf("  调用参数: %+v\n", subParams)
	result, err = instance.CallFunction("subtract", subParams)
	if err != nil {
		fmt.Printf("  ❌ 减法调用失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ 减法结果: 50 - 20 = %v\n", result)
	}

	// 3. 测试乘法函数
	fmt.Println("\n3. 测试乘法函数:")
	mulParams := map[string]interface{}{
		"a": 6,
		"b": 7,
	}
	fmt.Printf("  调用参数: %+v\n", mulParams)
	result, err = instance.CallFunction("multiply", mulParams)
	if err != nil {
		fmt.Printf("  ❌ 乘法调用失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ 乘法结果: 6 * 7 = %v\n", result)
	}

	// 4. 测试除法函数
	fmt.Println("\n4. 测试除法函数:")
	divParams := map[string]interface{}{
		"a": 20,
		"b": 4,
	}
	fmt.Printf("  调用参数: %+v\n", divParams)
	result, err = instance.CallFunction("divide", divParams)
	if err != nil {
		fmt.Printf("  ❌ 除法调用失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ 除法结果: 20 / 4 = %v\n", result)
	}

	// 5. 测试错误处理
	fmt.Println("\n5. 测试错误处理:")
	errorParams := map[string]interface{}{
		"a": 10,
		"b": 0,
	}
	fmt.Printf("  调用参数(除零错误): %+v\n", errorParams)
	result, err = instance.CallFunction("divide", errorParams)
	if err != nil {
		fmt.Printf("  ✓ 错误处理正常: %v\n", err)
	} else {
		fmt.Printf("  ❌ 应该返回错误，但得到结果: %v\n", result)
	}

	// 停止插件实例
	fmt.Println("\n=== 停止插件实例 ===")
	if err := instance.Stop(); err != nil {
		log.Printf("❌ 停止插件实例失败: %v", err)
	} else {
		fmt.Printf("✓ 插件实例已停止\n")
	}

	fmt.Println("\n=== Python插件简化测试完成 ===")
}
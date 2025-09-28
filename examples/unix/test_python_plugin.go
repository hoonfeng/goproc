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
	fmt.Println("=== Python插件测试 ===")
	fmt.Println("=== Python Plugin Test ===")

	// 设置详细日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ 获取当前工作目录失败: %v", err)
	}
	fmt.Printf("✓ 当前工作目录: %s\n", cwd)

	// 构建Python插件路径
	pythonPluginPath := filepath.Join(cwd, "python_plugin", "main.py")
	fmt.Printf("✓ Python插件路径: %s\n", pythonPluginPath)

	// 检查插件文件是否存在
	if _, err := os.Stat(pythonPluginPath); os.IsNotExist(err) {
		log.Fatalf("❌ Python插件文件不存在: %s", pythonPluginPath)
	}
	fmt.Printf("✓ Python插件文件存在\n")

	// 检查Python是否可用
	fmt.Println("\n=== 检查Python环境 ===")
	pythonCmd := "python3"
	if _, err := os.Stat("/usr/bin/python3"); os.IsNotExist(err) {
		pythonCmd = "python"
	}
	fmt.Printf("✓ 使用Python命令: %s\n", pythonCmd)

	// 创建插件配置
	pythonConfig := &config.PluginConfig{
		Type:                config.PluginTypeScript,
		Interpreter:         pythonCmd,
		ScriptPath:          pythonPluginPath,
		PoolSize:            1,
		MaxInstances:        3,
		HealthCheckInterval: 30 * time.Second,
		Functions:           []string{"add", "subtract", "multiply", "divide", "datetime_utils", "reverse", "uppercase", "lowercase", "length", "fibonacci", "text_processing"},
		Environment:         map[string]string{
			"PLUGIN_NAME": "python_plugin",
			"DEBUG": "true",
		},
	}

	fmt.Printf("✓ 插件配置创建完成\n")
	fmt.Printf("  - 类型: %s\n", pythonConfig.Type)
	fmt.Printf("  - 解释器: %s\n", pythonConfig.Interpreter)
	fmt.Printf("  - 脚本路径: %s\n", pythonConfig.ScriptPath)
	fmt.Printf("  - 池大小: %d\n", pythonConfig.PoolSize)
	fmt.Printf("  - 最大实例数: %d\n", pythonConfig.MaxInstances)

	// 创建插件池
	fmt.Println("\n=== 创建插件池 ===")
	pythonPool := plugin.NewPluginPool("python", pythonConfig)
	if pythonPool == nil {
		log.Fatalf("❌ 创建插件池失败")
	}
	fmt.Printf("✓ 插件池创建成功\n")

	// 启动插件池
	fmt.Println("\n=== 启动插件池 ===")
	
	// 检查sockets目录是否存在，如果不存在则创建
	socketsDir := "./sockets"
	if _, err := os.Stat(socketsDir); os.IsNotExist(err) {
		fmt.Printf("创建sockets目录: %s\n", socketsDir)
		err = os.MkdirAll(socketsDir, 0755)
		if err != nil {
			fmt.Printf("❌ 创建sockets目录失败: %v\n", err)
			return
		}
	}
	fmt.Printf("✓ sockets目录已准备好\n")
	
	err = pythonPool.Start()
	if err != nil {
		fmt.Printf("❌ 启动插件池失败: %v\n", err)
		log.Fatalf("插件池启动失败")
	}
	fmt.Printf("✓ 插件池启动成功\n")

	// 等待插件初始化
	fmt.Println("\n=== 等待插件初始化 ===")
	time.Sleep(3 * time.Second)
	fmt.Printf("✓ 插件初始化等待完成\n")

	// 测试插件功能
	fmt.Println("\n=== 测试插件功能 ===")
	
	// 测试加法
	fmt.Println("\n1. 测试加法函数:")
	params := map[string]interface{}{
		"a": 15.0,
		"b": 25.0,
	}
	fmt.Printf("  调用参数: %+v\n", params)
	
	result, err := pythonPool.CallFunction("add", params)
	if err != nil {
		fmt.Printf("  ❌ 加法调用失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ 加法结果: 15 + 25 = %v\n", result)
	}

	// 测试减法
	fmt.Println("\n2. 测试减法函数:")
	params2 := map[string]interface{}{
		"a": 50.0,
		"b": 20.0,
	}
	fmt.Printf("  调用参数: %+v\n", params2)
	
	result2, err := pythonPool.CallFunction("subtract", params2)
	if err != nil {
		fmt.Printf("  ❌ 减法调用失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ 减法结果: 50 - 20 = %v\n", result2)
	}

	// 测试字符串反转
	fmt.Println("\n3. 测试字符串反转函数:")
	params3 := map[string]interface{}{
		"text": "Hello World",
	}
	fmt.Printf("  调用参数: %+v\n", params3)
	
	result3, err := pythonPool.CallFunction("reverse", params3)
	if err != nil {
		fmt.Printf("  ❌ 字符串反转调用失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ 字符串反转结果: 'Hello World' -> '%v'\n", result3)
	}

	// 测试获取当前时间
	fmt.Println("\n4. 测试获取当前时间函数:")
	params4 := map[string]interface{}{}
	fmt.Printf("  调用参数: %+v\n", params4)
	
	result4, err := pythonPool.CallFunction("datetime_utils", params4)
	if err != nil {
		fmt.Printf("  ❌ 获取时间调用失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ 当前时间: %v\n", result4)
	}

	// 测试斐波那契数列
	fmt.Println("\n5. 测试斐波那契数列函数:")
	params5 := map[string]interface{}{
		"n": 10.0,
	}
	fmt.Printf("  调用参数: %+v\n", params5)
	
	result5, err := pythonPool.CallFunction("fibonacci", params5)
	if err != nil {
		fmt.Printf("  ❌ 斐波那契调用失败: %v\n", err)
	} else {
		fmt.Printf("  ✓ 斐波那契数列(前10项): %v\n", result5)
	}

	// 测试错误处理
	fmt.Println("\n6. 测试错误处理:")
	params6 := map[string]interface{}{
		"a": 10.0,
		"b": 0.0,
	}
	fmt.Printf("  调用参数(除零错误): %+v\n", params6)
	
	result6, err := pythonPool.CallFunction("divide", params6)
	if err != nil {
		fmt.Printf("  ✓ 错误处理正常: %v\n", err)
	} else {
		fmt.Printf("  ❌ 应该返回错误，但得到结果: %v\n", result6)
	}

	// 停止插件池
	fmt.Println("\n=== 停止插件池 ===")
	pythonPool.Stop()
	fmt.Printf("✓ 插件池已停止\n")

	fmt.Println("\n=== Python插件测试完成 ===")
}
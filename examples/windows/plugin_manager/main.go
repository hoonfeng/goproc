package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"goproc/config"
	"goproc/plugin"
)

// PluginManagerDemo 插件管理器演示程序
// Plugin Manager Demo Program
type PluginManagerDemo struct {
	manager *plugin.PluginManager
}

func main() {
	fmt.Println("=== 插件管理器演示程序 ===")
	fmt.Println("=== Plugin Manager Demo ===")
	fmt.Println()

	// 获取examples目录路径 / Get examples directory path
	examplesDir, err := filepath.Abs("examples")
	if err != nil {
		log.Fatalf("获取examples目录路径失败 / Failed to get examples directory path: %v", err)
	}

	// 创建系统配置 / Create system configuration
	systemConfig := createSystemConfig(examplesDir)

	// 创建插件管理器 / Create plugin manager
	manager := plugin.NewPluginManager(systemConfig)

	demo := &PluginManagerDemo{
		manager: manager,
	}

	// 启动插件管理器 / Start plugin manager
	fmt.Println("正在启动插件管理器... / Starting plugin manager...")
	if err := demo.manager.Start(); err != nil {
		log.Fatalf("启动插件管理器失败 / Failed to start plugin manager: %v", err)
	}

	fmt.Println("插件管理器启动成功！/ Plugin manager started successfully!")
	fmt.Println()

	// 显示初始状态 / Show initial status
	demo.showAllStatus()

	// 进入交互模式 / Enter interactive mode
	demo.runInteractiveMode()

	// 停止插件管理器 / Stop plugin manager
	fmt.Println("\n正在停止插件管理器... / Stopping plugin manager...")
	demo.manager.Stop()
	fmt.Println("插件管理器已停止 / Plugin manager stopped")
}

// createSystemConfig 创建系统配置
// Create system configuration
func createSystemConfig(examplesDir string) *config.SystemConfig {
	return &config.SystemConfig{
		Plugins: map[string]config.PluginConfig{
			"math": {
				Type:                "binary",
				Path:                filepath.Join(examplesDir, "windows", "math_plugin", "math_plugin.exe"),
				PoolSize:            3,
				MaxInstances:        6,
				HealthCheckInterval: 30 * time.Second,
				Functions:           []string{"add", "subtract", "multiply", "divide"},
			},
			"string": {
				Type:                "binary",
				Path:                filepath.Join(examplesDir, "windows", "string_plugin", "string_plugin.exe"),
				PoolSize:            2,
				MaxInstances:        4,
				HealthCheckInterval: 30 * time.Second,
				Functions:           []string{"toUpper", "toLower", "reverse", "length"},
			},
		},
	}
}

// runInteractiveMode 运行交互模式
// Run interactive mode
func (demo *PluginManagerDemo) runInteractiveMode() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		demo.showMenu()
		fmt.Print("请选择操作 / Please select an option: ")

		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())
		
		switch choice {
		case "1":
			demo.showAllStatus()
		case "2":
			demo.callPluginFunction(scanner)
		case "3":
			demo.showPluginStatus(scanner)
		case "4":
			demo.restartPlugin(scanner)
		case "5":
			demo.addPlugin(scanner)
		case "6":
			demo.removePlugin(scanner)
		case "7":
			demo.testPluginPerformance(scanner)
		case "0", "q", "quit", "exit":
			fmt.Println("退出程序... / Exiting...")
			return
		default:
			fmt.Println("无效选择，请重试 / Invalid choice, please try again")
		}
		
		fmt.Println()
	}
}

// showMenu 显示菜单
// Show menu
func (demo *PluginManagerDemo) showMenu() {
	fmt.Println("=== 插件管理器操作菜单 / Plugin Manager Menu ===")
	fmt.Println("1. 显示所有插件状态 / Show all plugin status")
	fmt.Println("2. 调用插件函数 / Call plugin function")
	fmt.Println("3. 显示特定插件状态 / Show specific plugin status")
	fmt.Println("4. 重启插件 / Restart plugin")
	fmt.Println("5. 添加插件 / Add plugin")
	fmt.Println("6. 移除插件 / Remove plugin")
	fmt.Println("7. 测试插件性能 / Test plugin performance")
	fmt.Println("0. 退出 / Exit")
	fmt.Println("================================================")
}

// showAllStatus 显示所有插件状态
// Show all plugin status
func (demo *PluginManagerDemo) showAllStatus() {
	fmt.Println("=== 所有插件状态 / All Plugin Status ===")
	status := demo.manager.GetAllStatus()
	
	fmt.Printf("管理器运行状态 / Manager Running: %v\n", status["is_running"])
	fmt.Printf("插件总数 / Total Plugins: %v\n", status["total_plugins"])
	fmt.Println()

	if plugins, ok := status["plugins"].(map[string]interface{}); ok {
		for pluginName, pluginStatus := range plugins {
			fmt.Printf("插件 / Plugin: %s\n", pluginName)
			if statusMap, ok := pluginStatus.(map[string]interface{}); ok {
				for key, value := range statusMap {
					fmt.Printf("  %s: %v\n", key, value)
				}
			}
			fmt.Println()
		}
	}
}

// callPluginFunction 调用插件函数
// Call plugin function
func (demo *PluginManagerDemo) callPluginFunction(scanner *bufio.Scanner) {
	fmt.Print("请输入插件名称 / Enter plugin name: ")
	if !scanner.Scan() {
		return
	}
	pluginName := strings.TrimSpace(scanner.Text())

	fmt.Print("请输入函数名称 / Enter function name: ")
	if !scanner.Scan() {
		return
	}
	functionName := strings.TrimSpace(scanner.Text())

	// 根据插件和函数提供示例参数 / Provide example parameters based on plugin and function
	params := demo.getExampleParams(pluginName, functionName, scanner)

	fmt.Printf("正在调用 %s.%s... / Calling %s.%s...\n", pluginName, functionName, pluginName, functionName)
	
	start := time.Now()
	result, err := demo.manager.CallFunction(pluginName, functionName, params)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("调用失败 / Call failed: %v\n", err)
	} else {
		fmt.Printf("调用成功 / Call successful: %v\n", result)
		fmt.Printf("耗时 / Duration: %v\n", duration)
	}
}

// getExampleParams 获取示例参数
// Get example parameters
func (demo *PluginManagerDemo) getExampleParams(pluginName, functionName string, scanner *bufio.Scanner) map[string]interface{} {
	params := make(map[string]interface{})

	switch pluginName {
	case "math":
		switch functionName {
		case "add", "subtract", "multiply", "divide":
			fmt.Print("请输入第一个数字 / Enter first number: ")
			if scanner.Scan() {
				if a, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64); err == nil {
					params["a"] = a
				}
			}
			fmt.Print("请输入第二个数字 / Enter second number: ")
			if scanner.Scan() {
				if b, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64); err == nil {
					params["b"] = b
				}
			}
		}
	case "string":
		switch functionName {
		case "toUpper", "toLower", "reverse", "length":
			fmt.Print("请输入字符串 / Enter string: ")
			if scanner.Scan() {
				params["text"] = strings.TrimSpace(scanner.Text())
			}
		}
	}

	return params
}

// showPluginStatus 显示特定插件状态
// Show specific plugin status
func (demo *PluginManagerDemo) showPluginStatus(scanner *bufio.Scanner) {
	fmt.Print("请输入插件名称 / Enter plugin name: ")
	if !scanner.Scan() {
		return
	}
	pluginName := strings.TrimSpace(scanner.Text())

	status, err := demo.manager.GetPluginStatus(pluginName)
	if err != nil {
		fmt.Printf("获取插件状态失败 / Failed to get plugin status: %v\n", err)
		return
	}

	fmt.Printf("=== 插件 %s 状态 / Plugin %s Status ===\n", pluginName, pluginName)
	for key, value := range status {
		fmt.Printf("%s: %v\n", key, value)
	}
}

// restartPlugin 重启插件
// Restart plugin
func (demo *PluginManagerDemo) restartPlugin(scanner *bufio.Scanner) {
	fmt.Print("请输入要重启的插件名称 / Enter plugin name to restart: ")
	if !scanner.Scan() {
		return
	}
	pluginName := strings.TrimSpace(scanner.Text())

	fmt.Printf("正在重启插件 %s... / Restarting plugin %s...\n", pluginName, pluginName)
	
	if err := demo.manager.RestartPlugin(pluginName); err != nil {
		fmt.Printf("重启插件失败 / Failed to restart plugin: %v\n", err)
	} else {
		fmt.Printf("插件 %s 重启成功 / Plugin %s restarted successfully\n", pluginName, pluginName)
	}
}

// addPlugin 添加插件
// Add plugin
func (demo *PluginManagerDemo) addPlugin(scanner *bufio.Scanner) {
	fmt.Print("请输入新插件名称 / Enter new plugin name: ")
	if !scanner.Scan() {
		return
	}
	pluginName := strings.TrimSpace(scanner.Text())

	fmt.Print("请输入插件路径 / Enter plugin path: ")
	if !scanner.Scan() {
		return
	}
	pluginPath := strings.TrimSpace(scanner.Text())

	// 创建插件配置 / Create plugin configuration
	pluginConfig := config.PluginConfig{
		Type:                "binary",
		Path:                pluginPath,
		PoolSize:            2,
		MaxInstances:        4,
		HealthCheckInterval: 30 * time.Second,
		Functions:           []string{"test"},
	}

	fmt.Printf("正在添加插件 %s... / Adding plugin %s...\n", pluginName, pluginName)
	
	if err := demo.manager.AddPlugin(pluginName, pluginConfig); err != nil {
		fmt.Printf("添加插件失败 / Failed to add plugin: %v\n", err)
	} else {
		fmt.Printf("插件 %s 添加成功 / Plugin %s added successfully\n", pluginName, pluginName)
	}
}

// removePlugin 移除插件
// Remove plugin
func (demo *PluginManagerDemo) removePlugin(scanner *bufio.Scanner) {
	fmt.Print("请输入要移除的插件名称 / Enter plugin name to remove: ")
	if !scanner.Scan() {
		return
	}
	pluginName := strings.TrimSpace(scanner.Text())

	fmt.Printf("正在移除插件 %s... / Removing plugin %s...\n", pluginName, pluginName)
	
	if err := demo.manager.RemovePlugin(pluginName); err != nil {
		fmt.Printf("移除插件失败 / Failed to remove plugin: %v\n", err)
	} else {
		fmt.Printf("插件 %s 移除成功 / Plugin %s removed successfully\n", pluginName, pluginName)
	}
}

// testPluginPerformance 测试插件性能
// Test plugin performance
func (demo *PluginManagerDemo) testPluginPerformance(scanner *bufio.Scanner) {
	fmt.Print("请输入要测试的插件名称 / Enter plugin name to test: ")
	if !scanner.Scan() {
		return
	}
	pluginName := strings.TrimSpace(scanner.Text())

	fmt.Print("请输入测试次数 / Enter number of tests (default 100): ")
	testCount := 100
	if scanner.Scan() {
		if count, err := strconv.Atoi(strings.TrimSpace(scanner.Text())); err == nil && count > 0 {
			testCount = count
		}
	}

	fmt.Printf("正在测试插件 %s 性能，测试次数：%d... / Testing plugin %s performance, test count: %d...\n", 
		pluginName, testCount, pluginName, testCount)

	// 选择测试函数 / Select test function
	var functionName string
	var params map[string]interface{}

	switch pluginName {
	case "math":
		functionName = "add"
		params = map[string]interface{}{"a": 10.0, "b": 20.0}
	case "string":
		functionName = "toUpper"
		params = map[string]interface{}{"text": "hello world"}
	default:
		fmt.Printf("不支持的插件 %s / Unsupported plugin %s\n", pluginName, pluginName)
		return
	}

	// 执行性能测试 / Execute performance test
	successCount := 0
	totalDuration := time.Duration(0)
	
	start := time.Now()
	for i := 0; i < testCount; i++ {
		callStart := time.Now()
		_, err := demo.manager.CallFunction(pluginName, functionName, params)
		callDuration := time.Since(callStart)
		
		if err == nil {
			successCount++
			totalDuration += callDuration
		}
	}
	totalTime := time.Since(start)

	// 输出测试结果 / Output test results
	fmt.Println("=== 性能测试结果 / Performance Test Results ===")
	fmt.Printf("插件名称 / Plugin Name: %s\n", pluginName)
	fmt.Printf("测试函数 / Test Function: %s\n", functionName)
	fmt.Printf("总测试次数 / Total Tests: %d\n", testCount)
	fmt.Printf("成功次数 / Successful: %d\n", successCount)
	fmt.Printf("失败次数 / Failed: %d\n", testCount-successCount)
	fmt.Printf("成功率 / Success Rate: %.2f%%\n", float64(successCount)/float64(testCount)*100)
	fmt.Printf("总耗时 / Total Time: %v\n", totalTime)
	fmt.Printf("平均延迟 / Average Latency: %v\n", totalDuration/time.Duration(successCount))
	fmt.Printf("QPS: %.2f\n", float64(successCount)/totalTime.Seconds())
}
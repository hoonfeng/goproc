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
	fmt.Println("=== Node.js插件单独测试 ===")

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}
	fmt.Printf("当前工作目录: %s\n", cwd)

	// 构建Node.js插件路径
	nodePluginPath := filepath.Join(cwd, "node_plugin")
	fmt.Printf("Node.js插件路径: %s\n", nodePluginPath)

	// 检查插件是否存在
	if _, err := os.Stat(nodePluginPath); os.IsNotExist(err) {
		log.Fatalf("Node.js插件不存在: %v", err)
	}
	fmt.Println("✓ Node.js插件存在")

	// 创建Node.js插件配置
	nodeConfig := &config.PluginConfig{
		Type:         config.PluginTypeScript,
		Interpreter:  "node",
		ScriptPath:   filepath.Join(nodePluginPath, "index.js"),
		PoolSize:     3,
		MaxInstances: 5,
		Functions: []string{
			"httpGet",
			"jsonParse", 
			"jsonStringify",
			"arraySum",
			"arrayFilter",
			"timestamp",
			"uuid",
		},
	}
	fmt.Println("✓ Node.js插件配置创建完成")

	// 创建插件池
	nodePool := plugin.NewPluginPool("node_plugin", nodeConfig)
	fmt.Println("✓ Node.js插件池创建完成")

	// 启动插件池
	fmt.Println("正在启动Node.js插件池...")
	err = nodePool.Start()
	if err != nil {
		log.Fatalf("启动Node.js插件池失败: %v", err)
	}
	fmt.Println("✓ Node.js插件池启动成功")

	// 等待插件初始化
	fmt.Println("等待插件初始化...")
	time.Sleep(3 * time.Second)
	fmt.Println("✓ 插件初始化完成")

	// 测试函数
	fmt.Println("\n=== 开始测试Node.js插件函数 ===")

	// 测试jsonStringify函数
	fmt.Println("测试jsonStringify函数...")
	result, err := nodePool.CallFunction("jsonStringify", map[string]interface{}{
		"obj": map[string]interface{}{
			"name":  "test",
			"value": 123,
			"active": true,
		},
	})
	if err != nil {
		fmt.Printf("❌ jsonStringify测试失败: %v\n", err)
	} else {
		fmt.Printf("✓ jsonStringify测试成功: %v\n", result)
	}

	// 测试arraySum函数
	fmt.Println("测试arraySum函数...")
	result, err = nodePool.CallFunction("arraySum", map[string]interface{}{
		"array": []float64{1.0, 2.0, 3.0, 4.0, 5.0},
	})
	if err != nil {
		fmt.Printf("❌ arraySum测试失败: %v\n", err)
	} else {
		fmt.Printf("✓ arraySum测试成功: 1+2+3+4+5 = %v\n", result)
	}

	// 测试timestamp函数
	fmt.Println("测试timestamp函数...")
	result, err = nodePool.CallFunction("timestamp", map[string]interface{}{})
	if err != nil {
		fmt.Printf("❌ timestamp测试失败: %v\n", err)
	} else {
		fmt.Printf("✓ timestamp测试成功: %v\n", result)
	}

	// 测试uuid函数
	fmt.Println("测试uuid函数...")
	result, err = nodePool.CallFunction("uuid", map[string]interface{}{})
	if err != nil {
		fmt.Printf("❌ uuid测试失败: %v\n", err)
	} else {
		fmt.Printf("✓ uuid测试成功: %v\n", result)
	}

	// 停止插件池
	fmt.Println("\n停止Node.js插件池...")
	nodePool.Stop()
	fmt.Println("✓ Node.js插件池已停止")

	fmt.Println("\n=== Node.js插件测试完成 ===")
}
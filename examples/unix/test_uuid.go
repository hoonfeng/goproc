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
	fmt.Println("=== UUID函数单独测试 ===")
	
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}
	
	// 构建Node.js插件路径
	nodePluginPath := filepath.Join(cwd, "node_plugin")
	
	// 创建Node.js插件配置
	nodeConfig := &config.PluginConfig{
		Type:         config.PluginTypeScript,
		Interpreter:  "node",
		ScriptPath:   filepath.Join(nodePluginPath, "index.js"),
		PoolSize:     1,
		MaxInstances: 1,
		Functions: []string{"uuid"},
	}
	
	// 创建插件池
	nodePool := plugin.NewPluginPool("uuid_test", nodeConfig)
	
	// 启动插件池
	fmt.Println("正在启动插件池...")
	err = nodePool.Start()
	if err != nil {
		log.Fatalf("启动插件池失败: %v", err)
	}
	fmt.Println("✓ 插件池启动成功")
	
	// 等待插件初始化
	fmt.Println("等待插件初始化...")
	time.Sleep(3 * time.Second)
	fmt.Println("✓ 插件初始化完成")
	
	// 测试uuid函数多次
	fmt.Println("测试uuid函数...")
	for i := 0; i < 5; i++ {
		result, err := nodePool.CallFunction("uuid", map[string]interface{}{})
		if err != nil {
			fmt.Printf("❌ uuid测试失败 (第%d次): %v\n", i+1, err)
		} else {
			fmt.Printf("✓ uuid测试成功 (第%d次): %v\n", i+1, result)
		}
		time.Sleep(100 * time.Millisecond)
	}
	
	// 停止插件池
	fmt.Println("停止插件池...")
	nodePool.Stop()
	fmt.Println("✓ 插件池已停止")
}
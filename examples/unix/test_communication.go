package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("=== 通信测试 ===")
	
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录失败: %v\n", err)
		return
	}
	
	// 构建Node.js插件路径
	nodePluginPath := filepath.Join(cwd, "node_plugin")
	pluginFile := filepath.Join(nodePluginPath, "index.js")
	
	// 测试地址
	address := "\\\\.\\pipe\\test-comm-123"
	
	fmt.Printf("工作目录: %s\n", cwd)
	fmt.Printf("插件文件: %s\n", pluginFile)
	fmt.Printf("测试地址: %s\n", address)
	
	// 设置环境变量
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOPROC_PLUGIN_ADDRESS=%s", address))
	
	// 启动Node.js插件进程
	fmt.Println("\n=== 启动Node.js插件 ===")
	cmd := exec.Command("node", pluginFile, address)
	cmd.Env = env
	cmd.Dir = nodePluginPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err = cmd.Start()
	if err != nil {
		fmt.Printf("启动Node.js插件失败: %v\n", err)
		return
	}
	
	fmt.Printf("Node.js插件进程已启动，PID: %d\n", cmd.Process.Pid)
	
	// 等待插件启动
	fmt.Println("等待插件启动...")
	time.Sleep(5 * time.Second)
	
	// 检查进程状态
	if cmd.Process != nil {
		fmt.Printf("插件进程状态: 运行中\n")
	} else {
		fmt.Printf("插件进程状态: 已退出\n")
	}
	
	// 停止进程
	fmt.Println("\n=== 停止插件进程 ===")
	if cmd.Process != nil {
		err = cmd.Process.Kill()
		if err != nil {
			fmt.Printf("停止进程失败: %v\n", err)
		} else {
			fmt.Println("进程已停止")
		}
	}
	
	fmt.Println("\n=== 通信测试完成 ===")
}
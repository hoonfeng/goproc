package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"encoding/json"
)

func main() {
	// 测试连接到Node.js插件
	fmt.Println("测试连接到Node.js插件...")
	
	// 连接到命名管道
	conn, err := net.Dial("pipe", "\\\\.\\pipe\\test-fixed-123")
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()
	
	fmt.Println("成功连接到Node.js插件")
	
	// 发送注册消息
	registerMsg := map[string]interface{}{
		"type": "register",
		"functions": []string{"jsonStringify", "arraySum", "timestamp", "uuid"},
	}
	
	registerData, err := json.Marshal(registerMsg)
	if err != nil {
		log.Fatalf("JSON编码失败: %v", err)
	}
	
	// 发送注册消息
	_, err = conn.Write(append(registerData, '\n'))
	if err != nil {
		log.Fatalf("发送注册消息失败: %v", err)
	}
	
	fmt.Println("已发送注册消息")
	
	// 等待响应
	time.Sleep(1 * time.Second)
	
	// 测试调用jsonStringify函数
	callMsg := map[string]interface{}{
		"type": "call",
		"function": "jsonStringify",
		"params": map[string]interface{}{"name": "test", "value": 123},
		"id": "test-1",
	}
	
	callData, err := json.Marshal(callMsg)
	if err != nil {
		log.Fatalf("JSON编码失败: %v", err)
	}
	
	// 发送调用消息
	_, err = conn.Write(append(callData, '\n'))
	if err != nil {
		log.Fatalf("发送调用消息失败: %v", err)
	}
	
	fmt.Println("已发送jsonStringify调用消息")
	
	// 读取响应
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}
	
	fmt.Printf("收到响应: %s\n", string(buffer[:n]))
	
	fmt.Println("测试完成")
}
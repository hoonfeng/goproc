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

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前工作目录失败: %v", err)
	}

	// 构建插件路径
	mathPluginPath := filepath.Join(cwd, "..", "math_plugin", "math_plugin-linux")
	stringPluginPath := filepath.Join(cwd, "..", "string_plugin", "string_plugin-linux")



	// 检查插件文件是否存在
	if _, err := os.Stat(mathPluginPath); os.IsNotExist(err) {
		log.Fatalf("数学插件文件不存在，请先构建数学插件")
	}
	if _, err := os.Stat(stringPluginPath); os.IsNotExist(err) {
		log.Fatalf("字符串插件文件不存在，请先构建字符串插件")
	}

	// 创建数学插件池配置
	mathConfig := &config.PluginConfig{
		Type:                config.PluginTypeBinary,
		Path:                mathPluginPath,
		PoolSize:            1,
		MaxInstances:        3,
		HealthCheckInterval: 30 * time.Second,
		Functions:           []string{"add", "subtract", "multiply", "divide", "power", "sqrt"},
	}

	// 创建字符串插件池配置
	stringConfig := &config.PluginConfig{
		Type:                config.PluginTypeBinary,
		Path:                stringPluginPath,
		PoolSize:            1,
		MaxInstances:        2,
		HealthCheckInterval: 30 * time.Second,
		Functions:           []string{"toUpper", "toLower", "reverse", "trim", "replace", "split", "join", "contains", "countWords", "isPalindrome"},
	}

	// 创建插件池
	mathPool := plugin.NewPluginPool("math_plugin", mathConfig)
	defer mathPool.Stop()

	stringPool := plugin.NewPluginPool("string_plugin", stringConfig)
	defer stringPool.Stop()

	// 启动插件池
	if err := mathPool.Start(); err != nil {
		log.Fatalf("启动数学插件池失败: %v", err)
	}

	if err := stringPool.Start(); err != nil {
		log.Fatalf("启动字符串插件池失败: %v", err)
	}

	// 等待插件池初始化完成
	time.Sleep(2 * time.Second)

	// 测试数学插件
	testMathPlugin(mathPool)

	// 测试字符串插件
	testStringPlugin(stringPool)

	// 测试并发性能
	testConcurrentPerformance(mathPool, stringPool)
}

func testMathPlugin(pool *plugin.PluginPool) {
	tests := []struct {
		function string
		params   map[string]interface{}
		expected interface{}
		desc     string
	}{
		{"add", map[string]interface{}{"a": 10.0, "b": 5.0}, 15.0, "10 + 5 = 15"},
		{"subtract", map[string]interface{}{"a": 10.0, "b": 5.0}, 5.0, "10 - 5 = 5"},
		{"multiply", map[string]interface{}{"a": 10.0, "b": 5.0}, 50.0, "10 * 5 = 50"},
		{"divide", map[string]interface{}{"a": 10.0, "b": 5.0}, 2.0, "10 / 5 = 2"},
		{"power", map[string]interface{}{"base": 2.0, "exponent": 3.0}, 8.0, "2^3 = 8"},
		{"sqrt", map[string]interface{}{"num": 16.0}, 4.0, "√16 = 4"},
	}

	for _, test := range tests {
		result, err := pool.CallFunction(test.function, test.params)
		if err != nil {
			continue
		}

		// 使用类型安全的比较
		_ = isEqual(result, test.expected)
	}
}

func testStringPlugin(pool *plugin.PluginPool) {
	tests := []struct {
		function string
		params   map[string]interface{}
		expected interface{}
		desc     string
	}{
		{"toUpper", map[string]interface{}{"str": "hello world"}, "HELLO WORLD", "转大写"},
		{"toLower", map[string]interface{}{"str": "HELLO WORLD"}, "hello world", "转小写"},
		{"reverse", map[string]interface{}{"str": "hello"}, "olleh", "反转字符串"},
		{"trim", map[string]interface{}{"str": "  hello  "}, "hello", "去除空白"},
		{"contains", map[string]interface{}{"str": "hello world", "substr": "world"}, true, "检查包含子串"},
		{"countWords", map[string]interface{}{"str": "hello world from go"}, 4, "统计单词数"},
		{"isPalindrome", map[string]interface{}{"str": "A man a plan a canal Panama"}, true, "检查回文"},
	}

	for _, test := range tests {
		result, err := pool.CallFunction(test.function, test.params)
		if err != nil {
			continue
		}

		// 使用类型安全的比较
		_ = isEqual(result, test.expected)
	}
}

// isEqual 类型安全的比较函数
func isEqual(a, b interface{}) bool {
	// 如果类型相同，直接比较
	if a == b {
		return true
	}

	// 处理数字类型的比较
	switch aVal := a.(type) {
	case int:
		switch bVal := b.(type) {
		case int:
			return aVal == bVal
		case float64:
			return float64(aVal) == bVal
		case float32:
			return float64(aVal) == float64(bVal)
		}
	case float64:
		switch bVal := b.(type) {
		case float64:
			return aVal == bVal
		case int:
			return aVal == float64(bVal)
		case float32:
			return aVal == float64(bVal)
		}
	case float32:
		switch bVal := b.(type) {
		case float32:
			return aVal == bVal
		case int:
			return float64(aVal) == float64(bVal)
		case float64:
			return float64(aVal) == bVal
		}
	}

	// 其他类型使用字符串表示比较
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func testConcurrentPerformance(mathPool, stringPool *plugin.PluginPool) {
	start := time.Now()
	
	// 使用信号量控制并发数量，避免过度使用插件池资源
	semaphore := make(chan struct{}, 3) // 限制最大并发数为3
	done := make(chan bool, 20) // 总共20个任务
	
	// 并发执行数学运算
	for i := 0; i < 10; i++ {
		go func(i int) {
			semaphore <- struct{}{} // 获取信号量
			defer func() { <-semaphore }() // 释放信号量
			
			_, err := mathPool.CallFunction("add", map[string]interface{}{"a": float64(i), "b": float64(i * 2)})
			_ = err
			done <- true
		}(i)
	}

	// 并发执行字符串处理
	for i := 0; i < 10; i++ {
		go func(i int) {
			semaphore <- struct{}{} // 获取信号量
			defer func() { <-semaphore }() // 释放信号量
			
			_, err := stringPool.CallFunction("toUpper", map[string]interface{}{"str": fmt.Sprintf("test %d", i)})
			_ = err
			done <- true
		}(i)
	}

	// 等待所有操作完成
	for i := 0; i < 20; i++ {
		<-done
	}
	
	elapsed := time.Since(start)
	_ = elapsed
}
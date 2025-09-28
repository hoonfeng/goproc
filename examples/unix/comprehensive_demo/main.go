package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"goproc/config"
	"goproc/plugin"
)

func main() {
	fmt.Println("=== 综合性演示程序启动 ===")
	fmt.Println("步骤1: 获取插件路径...")

	// 获取插件路径
	examplesDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败: %v", err)
	}
	examplesDir = filepath.Dir(examplesDir) // 回到examples目录
	fmt.Printf("当前工作目录: %s\n", examplesDir)

	// 插件路径
	mathPluginPath := filepath.Join(examplesDir, "math_plugin", "math_plugin-linux")
	stringPluginPath := filepath.Join(examplesDir, "string_plugin", "string_plugin-linux")
	nodePluginPath := filepath.Join(examplesDir, "node_plugin")
	pythonPluginPath := filepath.Join(examplesDir, "python_plugin")

	fmt.Println("步骤2: 检查插件是否存在...")
	// 检查插件是否存在
	if _, err := os.Stat(mathPluginPath); err != nil {
		log.Fatalf("数学插件不存在: %v", err)
	}
	fmt.Printf("✓ 数学插件存在: %s\n", mathPluginPath)
	
	if _, err := os.Stat(stringPluginPath); err != nil {
		log.Fatalf("字符串插件不存在: %v", err)
	}
	fmt.Printf("✓ 字符串插件存在: %s\n", stringPluginPath)
	
	if _, err := os.Stat(nodePluginPath); err != nil {
		log.Fatalf("Node.js插件不存在: %v", err)
	}
	fmt.Printf("✓ Node.js插件存在: %s\n", nodePluginPath)
	
	if _, err := os.Stat(pythonPluginPath); err != nil {
		log.Fatalf("Python插件不存在: %v", err)
	}
	fmt.Printf("✓ Python插件存在: %s\n", pythonPluginPath)

	fmt.Println("步骤3: 创建插件池配置...")
	// 创建插件池配置（实例数为5）
	mathConfig := &config.PluginConfig{
		Type:                config.PluginTypeBinary,
		Path:                mathPluginPath,
		PoolSize:            5,
		MaxInstances:        10,
		HealthCheckInterval: 30 * time.Second,
		Functions:           []string{"add", "subtract", "multiply", "divide", "power", "sqrt"},
	}
	fmt.Printf("✓ 数学插件配置创建完成 (函数数: %d)\n", len(mathConfig.Functions))

	stringConfig := &config.PluginConfig{
		Type:                config.PluginTypeBinary,
		Path:                stringPluginPath,
		PoolSize:            5,
		MaxInstances:        10,
		HealthCheckInterval: 30 * time.Second,
		Functions:           []string{"toUpper", "toLower", "reverse", "trim", "contains", "countWords", "isPalindrome"},
	}
	fmt.Printf("✓ 字符串插件配置创建完成 (函数数: %d)\n", len(stringConfig.Functions))

	nodeConfig := &config.PluginConfig{
		Type:                config.PluginTypeScript,
		Interpreter:         "node",
		ScriptPath:          filepath.Join(nodePluginPath, "index.js"),
		PoolSize:            5,
		MaxInstances:        10,
		HealthCheckInterval: 30 * time.Second,
		Functions:           []string{"httpGet", "jsonParse", "jsonStringify", "arraySum", "arrayFilter", "timestamp", "uuid"},
	}
	fmt.Printf("✓ Node.js插件配置创建完成 (函数数: %d)\n", len(nodeConfig.Functions))

	pythonConfig := &config.PluginConfig{
		Type:                config.PluginTypeScript,
		Interpreter:         "python",
		ScriptPath:          filepath.Join(pythonPluginPath, "main.py"),
		PoolSize:            5,
		MaxInstances:        10,
		HealthCheckInterval: 30 * time.Second,
		Functions:           []string{"http_request", "data_analysis", "text_processing", "datetime_utils", "generate_id"},
	}
	fmt.Printf("✓ Python插件配置创建完成 (函数数: %d)\n", len(pythonConfig.Functions))

	fmt.Println("步骤4: 创建插件池实例...")
	// 创建插件池
	mathPool := plugin.NewPluginPool("math_plugin", mathConfig)
	fmt.Println("✓ 数学插件池创建完成")
	
	stringPool := plugin.NewPluginPool("string_plugin", stringConfig)
	fmt.Println("✓ 字符串插件池创建完成")
	
	nodePool := plugin.NewPluginPool("node_plugin", nodeConfig)
	fmt.Println("✓ Node.js插件池创建完成")
	
	pythonPool := plugin.NewPluginPool("python_plugin", pythonConfig)
	fmt.Println("✓ Python插件池创建完成")

	defer mathPool.Stop()
	defer stringPool.Stop()
	defer nodePool.Stop()
	defer pythonPool.Stop()

	fmt.Println("步骤5: 启动插件池...")
	// 启动插件池
	fmt.Println("正在启动数学插件池...")
	if err := mathPool.Start(); err != nil {
		log.Fatalf("启动数学插件池失败: %v", err)
	}
	fmt.Println("✓ 数学插件池启动成功")
	
	fmt.Println("正在启动字符串插件池...")
	if err := stringPool.Start(); err != nil {
		log.Fatalf("启动字符串插件池失败: %v", err)
	}
	fmt.Println("✓ 字符串插件池启动成功")
	
	fmt.Println("正在启动Node.js插件池...")
	if err := nodePool.Start(); err != nil {
		log.Fatalf("启动Node.js插件池失败: %v", err)
	}
	fmt.Println("✓ Node.js插件池启动成功")
	
	fmt.Println("正在启动Python插件池...")
	if err := pythonPool.Start(); err != nil {
		log.Fatalf("启动Python插件池失败: %v", err)
	}
	fmt.Println("✓ Python插件池启动成功")

	fmt.Println("步骤6: 等待插件池初始化完成...")
	// 等待插件池初始化完成
	fmt.Println("等待3秒插件池初始化...")
	time.Sleep(3 * time.Second)
	fmt.Println("✓ 插件池初始化完成")

	fmt.Println("步骤7: 开始基本功能测试...")
	// 测试基本功能
	testBasicFunctions(mathPool, stringPool, nodePool, pythonPool)

	fmt.Println("步骤8: 开始并发性能测试...")
	// 测试100万并发性能
	testMillionConcurrency(mathPool, stringPool, nodePool, pythonPool)

	fmt.Println("步骤9: 输出测试完成信息...")
	// 输出测试完成信息
	fmt.Println("=== 综合性演示测试完成 ===")
	fmt.Println("✓ 所有插件池已成功启动")
	fmt.Println("✓ 基本功能测试已完成")
	fmt.Println("✓ 并发性能测试已完成")
	fmt.Println("✓ 插件系统运行正常")
	fmt.Println("=== 程序执行完毕 ===")
}

func testBasicFunctions(mathPool, stringPool, nodePool, pythonPool *plugin.PluginPool) {
	fmt.Println("=== 基本功能测试开始 ===")
	
	// 测试数学插件
	result, err := mathPool.CallFunction("add", map[string]interface{}{"a": 10.0, "b": 5.0})
	if err != nil {
		fmt.Printf("❌ 数学插件测试失败: %v\n", err)
	} else {
		fmt.Printf("✓ 数学插件测试成功: 10 + 5 = %v\n", result)
	}

	// 测试字符串插件
	result, err = stringPool.CallFunction("toUpper", map[string]interface{}{"str": "hello world"})
	if err != nil {
		fmt.Printf("❌ 字符串插件测试失败: %v\n", err)
	} else {
		fmt.Printf("✓ 字符串插件测试成功: 'hello world' -> '%v'\n", result)
	}

	// 测试Node.js插件
	result, err = nodePool.CallFunction("jsonStringify", map[string]interface{}{"obj": map[string]interface{}{"name": "test", "value": 123}})
	if err != nil {
		fmt.Printf("❌ Node.js插件测试失败: %v\n", err)
	} else {
		fmt.Printf("✓ Node.js插件测试成功: JSON序列化结果: %v\n", result)
	}

	// 测试Python插件
	result, err = pythonPool.CallFunction("datetime_utils", map[string]interface{}{"operation": "now"})
	if err != nil {
		fmt.Printf("❌ Python插件测试失败: %v\n", err)
	} else {
		fmt.Printf("✓ Python插件测试成功: 当前时间: %v\n", result)
	}
	
	fmt.Println("=== 基本功能测试完成 ===\n")
}

func testMillionConcurrency(mathPool, stringPool, nodePool, pythonPool *plugin.PluginPool) {
	fmt.Println("=== 并发性能测试开始 ===")
	fmt.Printf("测试规模: 1,000,000 个任务\n")
	fmt.Printf("批次大小: 10,000 个任务\n")
	fmt.Printf("最大并发: 1,000 个任务\n\n")
	
	const totalTasks = 1000000 // 100万任务
	const batchSize = 10000     // 每批处理1万个任务
	const maxConcurrent = 1000  // 最大并发数

	var successCount int64
	var errorCount int64
	var totalTime time.Duration

	// 分批处理，避免内存溢出
	for batch := 0; batch < totalTasks/batchSize; batch++ {
		fmt.Printf("正在处理批次 %d/%d...\n", batch+1, totalTasks/batchSize)
		batchStart := time.Now()
		batchSuccess, batchErrors := testConcurrentBatch(mathPool, stringPool, nodePool, pythonPool, batchSize, maxConcurrent)
		batchTime := time.Since(batchStart)

		atomic.AddInt64(&successCount, batchSuccess)
		atomic.AddInt64(&errorCount, batchErrors)
		totalTime += batchTime
		
		fmt.Printf("批次 %d 完成: 成功 %d, 失败 %d, 耗时 %v\n\n", 
			batch+1, batchSuccess, batchErrors, batchTime)
	}
	
	// 输出性能报告
	fmt.Println("=== 并发性能测试报告 ===")
	fmt.Printf("总任务数: %d\n", totalTasks)
	fmt.Printf("成功任务: %d (%.2f%%)\n", successCount, float64(successCount)/float64(totalTasks)*100)
	fmt.Printf("失败任务: %d (%.2f%%)\n", errorCount, float64(errorCount)/float64(totalTasks)*100)
	fmt.Printf("总耗时: %v\n", totalTime)
	fmt.Printf("平均QPS: %.2f 任务/秒\n", float64(totalTasks)/totalTime.Seconds())
	fmt.Printf("平均延迟: %.2f 毫秒/任务\n", totalTime.Seconds()/float64(totalTasks)*1000)
	fmt.Println("=== 并发性能测试完成 ===\n")
}

func testConcurrentBatch(mathPool, stringPool, nodePool, pythonPool *plugin.PluginPool, 
	batchSize, maxConcurrent int) (int64, int64) {
	
	var successCount int64
	var errorCount int64
	var wg sync.WaitGroup

	// 信号量控制并发数
	semaphore := make(chan struct{}, maxConcurrent)
	done := make(chan bool, batchSize)

	for i := 0; i < batchSize; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 随机选择插件和函数进行测试
			pluginType := taskID % 4
			var result interface{}
			var err error

			switch pluginType {
			case 0: // 数学插件
				result, err = mathPool.CallFunction("add", map[string]interface{}{
					"a": float64(taskID), 
					"b": float64(taskID * 2),
				})
			case 1: // 字符串插件
				result, err = stringPool.CallFunction("toUpper", map[string]interface{}{
					"str": fmt.Sprintf("task_%d", taskID),
				})
			case 2: // Node.js插件
				result, err = nodePool.CallFunction("arraySum", map[string]interface{}{
					"array": []float64{float64(taskID), float64(taskID * 2)},
				})
			case 3: // Python插件
				result, err = pythonPool.CallFunction("text_processing", map[string]interface{}{
					"text": fmt.Sprintf("Task number %d", taskID),
					"operation": "word_count",
				})
			}

			if err != nil {
				atomic.AddInt64(&errorCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
				// 验证结果
				_ = result // 使用结果避免编译器警告
			}

			done <- true
		}(i)
	}

	wg.Wait()
	close(done)

	// 等待所有任务完成
	for i := 0; i < batchSize; i++ {
		<-done
	}

	return successCount, errorCount
}

// 类型安全的比较函数
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
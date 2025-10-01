package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hoonfeng/goproc/config"
	"github.com/hoonfeng/goproc/plugin"
)

// TestConfig 测试配置结构体
// Test configuration structure
type TestConfig struct {
	TotalTasks    int // 总任务数 / Total number of tasks
	BatchSize     int // 批次大小 / Batch size
	MaxConcurrent int // 最大并发数 / Maximum concurrent tasks
}

// PerformanceResult 性能测试结果结构体
// Performance test result structure
type PerformanceResult struct {
	PluginName   string        // 插件名称 / Plugin name
	TotalTasks   int64         // 总任务数 / Total tasks
	SuccessCount int64         // 成功任务数 / Successful tasks
	ErrorCount   int64         // 失败任务数 / Failed tasks
	TotalTime    time.Duration // 总耗时 / Total time
	QPS          float64       // 每秒查询数 / Queries per second
	AvgLatency   float64       // 平均延迟(毫秒) / Average latency (ms)
	SuccessRate  float64       // 成功率 / Success rate
}

func main() {
	fmt.Println("=== 插件系统并发性能测试 ===")
	fmt.Println("=== Plugin System Concurrency Performance Test ===")

	// 测试配置 / Test configuration
	testConfig := &TestConfig{
		TotalTasks:    1000000, // 总任务数 / Total number of tasks
		BatchSize:     100000,  // 批次大小 / Batch size
		MaxConcurrent: 10,      // 最大并发数，匹配实例数 / Maximum concurrent tasks, match instance count
	}

	fmt.Printf("测试配置 / Test Configuration:\n")
	fmt.Printf("- 总任务数 / Total Tasks: %d\n", testConfig.TotalTasks)
	fmt.Printf("- 批次大小 / Batch Size: %d\n", testConfig.BatchSize)
	fmt.Printf("- 最大并发 / Max Concurrent: %d\n", testConfig.MaxConcurrent)
	fmt.Println()

	// 获取插件路径 / Get plugin paths
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败 / Failed to get current directory: %v", err)
	}
	// 如果当前在performance_test目录，则回到examples目录 / If in performance_test directory, go back to examples
	examplesDir := filepath.Dir(currentDir)
	if filepath.Base(currentDir) != "performance_test" {
		// 如果不在performance_test目录，假设在goproc根目录 / If not in performance_test, assume in goproc root
		examplesDir = filepath.Join(currentDir, "examples")
	}

	// 创建插件配置 / Create plugin configurations
	configs := createPluginConfigs(examplesDir)

	// 创建并启动插件池 / Create and start plugin pools
	pools := make(map[string]*plugin.PluginPool)
	for name, config := range configs {
		fmt.Printf("正在启动 %s 插件... / Starting %s plugin...\n", name, name)
		pool := plugin.NewPluginPool(name, config)
		if err := pool.Start(); err != nil {
			log.Fatalf("启动%s插件失败 / Failed to start %s plugin: %v", name, name, err)
		}
		fmt.Printf("%s 插件启动成功 / %s plugin started successfully\n", name, name)
		pools[name] = pool
		defer pool.Stop()
	}

	// 等待插件初始化 / Wait for plugin initialization
	fmt.Println("等待插件初始化... / Waiting for plugin initialization...")
	time.Sleep(3 * time.Second)
	fmt.Println("开始性能测试... / Starting performance test...")
	fmt.Println()

	// 执行性能测试 / Execute performance tests
	results := runPerformanceTests(pools, testConfig)

	// 输出最终对比报告 / Output final comparison report
	fmt.Println("\n" + strings.Repeat("*", 80))
	fmt.Println("=== 所有插件测试完成，生成对比报告 ===")
	fmt.Println("=== All Plugin Tests Completed, Generating Comparison Report ===")
	fmt.Println(strings.Repeat("*", 80))
	printPerformanceReport(results, testConfig)

	// 等待所有插件池停止 / Wait for all plugin pools to stop
	fmt.Println("等待所有插件池停止... / Waiting for all plugin pools to stop...")
	for _, pool := range pools {
		pool.Stop()
	}
	fmt.Println("所有插件池已停止 / All plugin pools stopped")
	fmt.Println("=== 插件系统并发性能测试完成 ===")
	fmt.Println("=== Plugin System Concurrency Performance Test Completed ===")
}

// createPluginConfigs 创建插件配置
// Create plugin configurations
func createPluginConfigs(examplesDir string) map[string]*config.PluginConfig {
	return map[string]*config.PluginConfig{
		"math": {
			Type:                config.PluginTypeBinary,
			Path:                filepath.Join(examplesDir, "math_plugin", "math_plugin.exe"),
			PoolSize:            10, // 减少池大小以加快启动 / Reduce pool size for faster startup
			MaxInstances:        10,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"add", "multiply", "subtract", "divide"},
		},
		"string": {
			Type:                config.PluginTypeBinary,
			Path:                filepath.Join(examplesDir, "string_plugin", "string_plugin.exe"),
			PoolSize:            10,
			MaxInstances:        10,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"toUpper", "toLower", "reverse", "trim"},
		},
		"node": {
			Type:                config.PluginTypeScript,
			Interpreter:         "node",
			ScriptPath:          filepath.Join(examplesDir, "node_plugin", "index.js"),
			PoolSize:            10,
			MaxInstances:        10,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"jsonStringify", "arraySum", "timestamp", "uuid"},
		},
		"python": {
			Type:                config.PluginTypeScript,
			Interpreter:         "python",
			ScriptPath:          filepath.Join(examplesDir, "python_plugin", "main.py"),
			PoolSize:            2,
			MaxInstances:        5,
			HealthCheckInterval: 30 * time.Second,
			Functions:           []string{"datetime_utils", "text_processing", "generate_id"},
		},
	}
}

// runPerformanceTests 运行性能测试 - 串行测试每个插件
// Run performance tests - test each plugin sequentially
func runPerformanceTests(pools map[string]*plugin.PluginPool, testConfig *TestConfig) []*PerformanceResult {
	var results []*PerformanceResult

	// 串行测试每个插件 / Test each plugin sequentially
	pluginNames := []string{"math", "string", "node", "python"} // 固定顺序 / Fixed order

	for _, pluginName := range pluginNames {
		pool, exists := pools[pluginName]
		if !exists {
			fmt.Printf("警告：插件 %s 不存在，跳过测试 / Warning: Plugin %s not found, skipping\n", pluginName, pluginName)
			continue
		}

		fmt.Printf("\n开始测试 %s 插件... / Starting test for %s plugin...\n", pluginName, pluginName)
		fmt.Printf("任务数量：%d / Task count: %d\n", testConfig.TotalTasks, testConfig.TotalTasks)

		// 测试单个插件 / Test single plugin
		result := testPluginPerformance(pluginName, pool, testConfig)
		results = append(results, result)

		// 立即输出该插件的性能报告 / Immediately output performance report for this plugin
		printSinglePluginReport(result, testConfig)

		fmt.Printf("完成测试 %s 插件 / Completed test for %s plugin\n\n", pluginName, pluginName)
	}

	return results
}

// testPluginPerformance 测试单个插件的性能
// Test performance of a single plugin
func testPluginPerformance(pluginName string, pool *plugin.PluginPool, testConfig *TestConfig) *PerformanceResult {
	var successCount int64
	var errorCount int64
	startTime := time.Now()

	// 分批处理任务 / Process tasks in batches
	numBatches := testConfig.TotalTasks / testConfig.BatchSize
	fmt.Printf("开始分批测试，共 %d 批，每批 %d 任务 / Starting batch testing, %d batches, %d tasks per batch\n",
		numBatches, testConfig.BatchSize, numBatches, testConfig.BatchSize)

	for batch := 0; batch < numBatches; batch++ {
		fmt.Printf("执行第 %d/%d 批... / Executing batch %d/%d...\n", batch+1, numBatches, batch+1, numBatches)
		batchSuccess, batchErrors := runBatch(pluginName, pool, testConfig.BatchSize, testConfig.MaxConcurrent)
		atomic.AddInt64(&successCount, batchSuccess)
		atomic.AddInt64(&errorCount, batchErrors)

		// 显示进度 / Show progress
		currentTotal := atomic.LoadInt64(&successCount) + atomic.LoadInt64(&errorCount)
		fmt.Printf("批次完成，当前进度：%d/%d (成功:%d, 失败:%d) / Batch completed, progress: %d/%d (success:%d, failed:%d)\n",
			currentTotal, testConfig.TotalTasks, atomic.LoadInt64(&successCount), atomic.LoadInt64(&errorCount))
	}

	totalTime := time.Since(startTime)
	totalTasks := int64(testConfig.TotalTasks)

	return &PerformanceResult{
		PluginName:   pluginName,
		TotalTasks:   totalTasks,
		SuccessCount: successCount,
		ErrorCount:   errorCount,
		TotalTime:    totalTime,
		QPS:          float64(totalTasks) / totalTime.Seconds(),
		AvgLatency:   totalTime.Seconds() / float64(totalTasks) * 1000, // 转换为毫秒 / Convert to milliseconds
		SuccessRate:  float64(successCount) / float64(totalTasks) * 100,
	}
}

// runBatch 运行一批任务
// Run a batch of tasks
func runBatch(pluginName string, pool *plugin.PluginPool, batchSize, maxConcurrent int) (int64, int64) {
	var successCount int64
	var errorCount int64
	var wg sync.WaitGroup

	fmt.Printf("开始执行批次：插件=%s, 批次大小=%d, 最大并发=%d / Starting batch: plugin=%s, batch_size=%d, max_concurrent=%d\n",
		pluginName, batchSize, maxConcurrent, pluginName, batchSize, maxConcurrent)

	// 信号量控制并发数 / Semaphore to control concurrency
	semaphore := make(chan struct{}, maxConcurrent)

	for i := 0; i < batchSize; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 根据插件类型选择测试函数 / Select test function based on plugin type
			var err error
			switch pluginName {
			case "math":
				_, err = pool.CallFunction("add", map[string]interface{}{
					"a": float64(taskID),
					"b": float64(taskID * 2),
				})
			case "string":
				_, err = pool.CallFunction("toUpper", map[string]interface{}{
					"str": fmt.Sprintf("task_%d", taskID),
				})
			case "node":
				_, err = pool.CallFunction("arraySum", map[string]interface{}{
					"array": []float64{float64(taskID), float64(taskID * 2)},
				})
			case "python":
				_, err = pool.CallFunction("text_processing", map[string]interface{}{
					"text":      fmt.Sprintf("Task number %d", taskID),
					"operation": "word_count",
				})
			}

			if err != nil {
				atomic.AddInt64(&errorCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()
	return successCount, errorCount
}

// printSinglePluginReport 打印单个插件的性能报告
// Print performance report for a single plugin
func printSinglePluginReport(result *PerformanceResult, testConfig *TestConfig) {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("=== %s 插件性能报告 / %s Plugin Performance Report ===\n", result.PluginName, result.PluginName)
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("总任务数 / Total Tasks: %d\n", result.TotalTasks)
	fmt.Printf("成功任务 / Successful Tasks: %d\n", result.SuccessCount)
	fmt.Printf("失败任务 / Failed Tasks: %d\n", result.ErrorCount)
	fmt.Printf("总耗时 / Total Time: %.2f 秒 / seconds\n", result.TotalTime.Seconds())
	fmt.Printf("QPS (每秒查询数 / Queries Per Second): %.2f\n", result.QPS)
	fmt.Printf("平均延迟 / Average Latency: %.2f 毫秒 / ms\n", result.AvgLatency)
	fmt.Printf("成功率 / Success Rate: %.2f%%\n", result.SuccessRate)

	// 性能等级评估 / Performance level assessment
	var performanceLevel string
	if result.QPS >= 10000 {
		performanceLevel = "优秀 / Excellent"
	} else if result.QPS >= 5000 {
		performanceLevel = "良好 / Good"
	} else if result.QPS >= 1000 {
		performanceLevel = "一般 / Average"
	} else {
		performanceLevel = "需要优化 / Needs Optimization"
	}
	fmt.Printf("性能等级 / Performance Level: %s\n", performanceLevel)

	fmt.Println(strings.Repeat("=", 60))
}

// printPerformanceReport 打印最终对比报告
// Print final comparison report
func printPerformanceReport(results []*PerformanceResult, testConfig *TestConfig) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("=== 插件性能对比报告 / Plugin Performance Comparison Report ===")
	fmt.Println(strings.Repeat("=", 80))

	// 表头 / Table header
	fmt.Printf("%-12s %-10s %-10s %-10s %-12s %-12s %-12s %-12s\n",
		"插件名称", "总任务", "成功", "失败", "耗时(秒)", "QPS", "延迟(ms)", "成功率(%)")
	fmt.Printf("%-12s %-10s %-10s %-10s %-12s %-12s %-12s %-12s\n",
		"Plugin", "Total", "Success", "Failed", "Time(s)", "QPS", "Latency", "Success%")
	fmt.Println(strings.Repeat("-", 80))

	var totalTasks int64
	var totalSuccess int64
	var totalErrors int64
	var totalTime time.Duration

	// 输出每个插件的结果 / Output results for each plugin
	for _, result := range results {
		fmt.Printf("%-12s %-10d %-10d %-10d %-12.2f %-12.2f %-12.2f %-12.2f\n",
			result.PluginName,
			result.TotalTasks,
			result.SuccessCount,
			result.ErrorCount,
			result.TotalTime.Seconds(),
			result.QPS,
			result.AvgLatency,
			result.SuccessRate)

		totalTasks += result.TotalTasks
		totalSuccess += result.SuccessCount
		totalErrors += result.ErrorCount
		if result.TotalTime > totalTime {
			totalTime = result.TotalTime // 使用最长时间 / Use the longest time
		}
	}

	fmt.Println(strings.Repeat("-", 80))

	// 总计 / Summary
	overallQPS := float64(totalTasks) / totalTime.Seconds()
	overallLatency := totalTime.Seconds() / float64(totalTasks) * 1000
	overallSuccessRate := float64(totalSuccess) / float64(totalTasks) * 100

	fmt.Printf("%-12s %-10d %-10d %-10d %-12.2f %-12.2f %-12.2f %-12.2f\n",
		"总计/Total",
		totalTasks,
		totalSuccess,
		totalErrors,
		totalTime.Seconds(),
		overallQPS,
		overallLatency,
		overallSuccessRate)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("测试总结 / Test Summary:\n")
	fmt.Printf("- 测试插件数量 / Plugins Tested: %d\n", len(results))
	fmt.Printf("- 总任务数量 / Total Tasks: %d\n", totalTasks)
	fmt.Printf("- 总成功率 / Overall Success Rate: %.2f%%\n", overallSuccessRate)
	fmt.Printf("- 总体QPS / Overall QPS: %.2f\n", overallQPS)
	fmt.Printf("- 平均延迟 / Average Latency: %.2f ms\n", overallLatency)
	fmt.Println(strings.Repeat("=", 80))
}

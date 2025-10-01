package plugin

import (
	"fmt"
	"sync"
	"time"

	"github.com/hoonfeng/goproc/config"

	"github.com/google/uuid"
)

// PluginPool 插件池
type PluginPool struct {
	PluginName   string
	Config       *config.PluginConfig
	Instances    map[string]*PluginInstance
	Available    chan *PluginInstance
	Mutex        sync.RWMutex
	IsRunning    bool
	MaxInstances int

	// 简化的等待队列
	waitQueue chan *PluginInstance // 直接存储等待的实例，而不是等待者通道
}

// NewPluginPool 创建新的插件池
func NewPluginPool(pluginName string, config *config.PluginConfig) *PluginPool {
	// 可用队列容量设置为最大实例数，避免队列中有重复实例
	queueCapacity := config.MaxInstances
	if queueCapacity < 10 {
		queueCapacity = 10 // 最小容量为10
	}

	// 简化的等待队列，容量为最大实例数
	waitQueueCapacity := config.MaxInstances
	if waitQueueCapacity < 10 {
		waitQueueCapacity = 10
	}

	pool := &PluginPool{
		PluginName:   pluginName,
		Config:       config,
		Instances:    make(map[string]*PluginInstance),
		Available:    make(chan *PluginInstance, queueCapacity),
		waitQueue:    make(chan *PluginInstance, waitQueueCapacity),
		IsRunning:    false,
		MaxInstances: config.MaxInstances,
	}

	return pool
}

// Start 启动插件池
func (pp *PluginPool) Start() error {
	if pp.IsRunning {
		return fmt.Errorf("插件池 %s 已经在运行", pp.PluginName)
	}

	// 创建初始实例（同步执行，不需要锁）
	successCount := 0
	for i := 0; i < pp.Config.PoolSize; i++ {
		_, err := pp.createNewInstance()
		if err != nil {
			// 记录错误但不中断启动过程
			continue
		}
		successCount++
	}

	// 如果没有任何实例创建成功，返回错误
	if successCount == 0 {
		errMsg := fmt.Errorf("插件池 %s 启动失败：无法创建任何插件实例", pp.PluginName)
		return errMsg
	}

	pp.IsRunning = true

	return nil
}

// GetInstance 获取可用插件实例
func (pp *PluginPool) GetInstance() (*PluginInstance, error) {
	if !pp.IsRunning {
		return nil, fmt.Errorf("插件池 %s 未运行", pp.PluginName)
	}

	//fmt.Printf("[Pool] 正在获取实例，当前实例数: %d, 最大实例数: %d\n", len(pp.Instances), pp.MaxInstances)

	// 首先尝试从可用队列获取实例（非阻塞）
	select {
	case instance := <-pp.Available:
		//fmt.Printf("[Pool] 从可用队列获取实例: %s\n", instance.ID)
		// 简化健康检查：只在必要时检查，避免频繁检查导致性能问题
		// 如果实例最近被使用过，假设它是健康的
		return instance, nil
	default:
		//fmt.Printf("[Pool] 可用队列为空\n")
		// 队列为空，继续后续逻辑
	}

	// 检查当前实例数量
	pp.Mutex.RLock()
	currentCount := len(pp.Instances)
	pp.Mutex.RUnlock()

	if currentCount < pp.MaxInstances {
		//fmt.Printf("[Pool] 创建新实例（当前: %d, 最大: %d）\n", currentCount, pp.MaxInstances)
		// 未达到最大实例数，创建新实例
		return pp.createNewInstance()
	} else {
		//fmt.Printf("[Pool] 已达到最大实例数，等待可用实例\n")
		// 已达到最大实例数，使用更智能的等待机制
		return pp.waitForAvailableInstance()
	}
}

// createNewInstance 创建新实例
func (pp *PluginPool) createNewInstance() (*PluginInstance, error) {
	pp.Mutex.RLock()
	currentCount := len(pp.Instances)
	pp.Mutex.RUnlock()

	if currentCount >= pp.MaxInstances {
		return nil, fmt.Errorf("已达到最大实例数 %d", pp.MaxInstances)
	}

	// 使用UUID生成唯一实例ID，移除连字符以缩短长度
	instanceUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("生成实例ID失败: %w", err)
	}

	// 移除UUID中的连字符，缩短名称长度
	uuidStr := instanceUUID.String()
	uuidStr = uuidStr[:8] + uuidStr[9:13] + uuidStr[14:18] + uuidStr[19:23] + uuidStr[24:]
	instanceID := fmt.Sprintf("%s%s", pp.PluginName, uuidStr)

	instance := NewPluginInstance(pp.PluginName, pp.Config, instanceID)

	// 先启动实例，如果失败则不添加到映射中
	if err := instance.Start(); err != nil {
		return nil, fmt.Errorf("启动实例 %s 失败: %w", instanceID, err)
	}

	// 启动成功后再添加到实例映射
	pp.Mutex.Lock()
	pp.Instances[instanceID] = instance
	pp.Mutex.Unlock()

	// 将实例放入可用队列
	pp.Available <- instance

	return instance, nil
}

// ReturnInstance 归还插件实例
func (pp *PluginPool) ReturnInstance(instance *PluginInstance) {
	pp.Mutex.RLock()
	if !pp.IsRunning {
		pp.Mutex.RUnlock()
		return
	}

	// 检查实例是否仍然存在
	_, exists := pp.Instances[instance.ID]
	pp.Mutex.RUnlock()

	if !exists {
		return
	}

	// 简化健康检查：只在实例调用失败时进行健康检查
	// 正常归还时假设实例是健康的

	// 优先处理等待队列中的请求
	select {
	case pp.waitQueue <- instance:
		// 成功放入等待队列
	default:
		// 等待队列已满，将实例放回可用队列
		select {
		case pp.Available <- instance:
			// 成功归还
		default:
			// 队列已满，不归还（理论上不应该发生）
		}
	}
}

// restartInstance 重启插件实例
func (pp *PluginPool) restartInstance(instance *PluginInstance) error {
	// 停止实例
	if err := instance.Stop(); err != nil {
		return fmt.Errorf("停止实例失败: %w", err)
	}

	// 重新启动实例
	if err := instance.Start(); err != nil {
		return fmt.Errorf("启动实例失败: %w", err)
	}

	return nil
}

// waitForAvailableInstance 等待可用实例（简化版本）
func (pp *PluginPool) waitForAvailableInstance() (*PluginInstance, error) {
	// 设置超时时间（5秒）
	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	// 使用纯channel阻塞等待，移除default分支避免CPU忙等待
	// Use pure channel blocking wait, remove default branch to avoid CPU busy waiting
	select {
	case instance := <-pp.Available:
		return instance, nil
	case instance := <-pp.waitQueue:
		return instance, nil
	case <-timeout.C:
		return nil, fmt.Errorf("等待可用实例超时")
	}
}

// removeInstance 移除插件实例
func (pp *PluginPool) removeInstance(instanceID string) {
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	if instance, exists := pp.Instances[instanceID]; exists {
		// 停止实例
		instance.Stop()

		// 从实例映射中移除
		delete(pp.Instances, instanceID)
	}
}

// CallFunction 调用插件函数
func (pp *PluginPool) CallFunction(functionName string, params map[string]interface{}) (interface{}, error) {
	//fmt.Printf("[Pool] 调用函数: %s, 参数: %v\n", functionName, params)

	// 获取实例
	instance, err := pp.GetInstance()
	if err != nil {
		//fmt.Printf("[Pool] 获取实例失败: %v\n", err)
		return nil, err
	}

	// 检查实例是否为nil
	if instance == nil {
		//fmt.Printf("[Pool] 获取到的插件实例为nil\n")
		return nil, fmt.Errorf("获取到的插件实例为nil")
	}

	//fmt.Printf("[Pool] 获取实例成功: %s\n", instance.ID)

	// 确保实例被归还
	defer pp.ReturnInstance(instance)

	// 调用函数
	//fmt.Printf("[Pool] 在实例 %s 上调用函数: %s\n", instance.ID, functionName)
	result, err := instance.CallFunction(functionName, params)
	if err != nil {
		//fmt.Printf("[Pool] 函数调用失败: %v\n", err)
		return nil, err
	}

	//fmt.Printf("[Pool] 函数调用成功: %v\n", result)
	return result, nil
}

// GetStatus 获取插件池状态
func (pp *PluginPool) GetStatus() map[string]interface{} {
	pp.Mutex.RLock()
	defer pp.Mutex.RUnlock()

	instancesStatus := make(map[string]interface{})
	for id, instance := range pp.Instances {
		instancesStatus[id] = instance.GetStatus()
	}

	// 计算可用实例数：通道中的实例数，但不能超过总实例数
	// 由于通道可能有重复实例，我们需要确保可用实例数不超过总实例数
	availableCount := len(pp.Available)
	totalInstances := len(pp.Instances)

	// 确保可用实例数不超过总实例数
	if availableCount > totalInstances {
		availableCount = totalInstances
	}

	return map[string]interface{}{
		"plugin_name":     pp.PluginName,
		"is_running":      pp.IsRunning,
		"total_instances": totalInstances,
		"max_instances":   pp.MaxInstances,
		"available_count": availableCount,
		"instances":       instancesStatus,
	}
}

// Stop 停止插件池
func (pp *PluginPool) Stop() {
	// 检查是否已经在运行状态
	pp.Mutex.RLock()
	if !pp.IsRunning {
		pp.Mutex.RUnlock()
		return
	}
	pp.Mutex.RUnlock()

	// 设置停止标志（原子操作，不需要锁）
	pp.Mutex.Lock()
	pp.IsRunning = false
	pp.Mutex.Unlock()

	// 关闭等待队列，防止新的等待请求
	close(pp.waitQueue)

	// 停止所有实例（不需要锁保护，因为IsRunning=false会阻止新操作）

	instances := make([]*PluginInstance, 0, len(pp.Instances))
	for _, instance := range pp.Instances {
		instances = append(instances, instance)
	}

	for _, instance := range instances {
		if err := instance.Stop(); err != nil {
			// 记录错误但不中断停止过程
		}
	}

	// 清空实例映射
	pp.Mutex.Lock()
	pp.Instances = make(map[string]*PluginInstance)
	pp.Mutex.Unlock()

	// 关闭可用通道
	close(pp.Available)
}

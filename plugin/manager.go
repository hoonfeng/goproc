package plugin

import (
	"fmt"
	"sync"

	"goproc/config"
)

// PluginManager 插件管理器
type PluginManager struct {
	Config   *config.SystemConfig
	Pools    map[string]*PluginPool
	Mutex    sync.RWMutex
	IsRunning bool
}

// NewPluginManager 创建新的插件管理器
func NewPluginManager(config *config.SystemConfig) *PluginManager {
	return &PluginManager{
		Config:    config,
		Pools:     make(map[string]*PluginPool),
		IsRunning: false,
	}
}

// Start 启动插件管理器
func (pm *PluginManager) Start() error {
	pm.Mutex.Lock()
	defer pm.Mutex.Unlock()
	
	if pm.IsRunning {
		return fmt.Errorf("插件管理器已经在运行")
	}
	
	// 验证配置
	if err := config.ValidateConfig(pm.Config); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}
	
	// 创建并启动所有插件池
	for pluginName, pluginConfig := range pm.Config.Plugins {
		pool := NewPluginPool(pluginName, &pluginConfig)
		
		if err := pool.Start(); err != nil {
			continue
		}
		
		pm.Pools[pluginName] = pool
	}
	
	if len(pm.Pools) == 0 {
		return fmt.Errorf("没有成功启动任何插件池")
	}
	
	pm.IsRunning = true
	
	return nil
}

// CallFunction 调用插件函数
func (pm *PluginManager) CallFunction(pluginName string, functionName string, params map[string]interface{}) (interface{}, error) {
	pm.Mutex.RLock()
	pool, exists := pm.Pools[pluginName]
	pm.Mutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("插件 %s 不存在", pluginName)
	}
	
	if !pm.IsRunning {
		return nil, fmt.Errorf("插件管理器未运行")
	}
	
	// 调用函数
	result, err := pool.CallFunction(functionName, params)
	if err != nil {
		return nil, fmt.Errorf("调用函数 %s 失败: %w", functionName, err)
	}
	
	return result, nil
}

// GetPluginStatus 获取插件状态
func (pm *PluginManager) GetPluginStatus(pluginName string) (map[string]interface{}, error) {
	pm.Mutex.RLock()
	pool, exists := pm.Pools[pluginName]
	pm.Mutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("插件 %s 不存在", pluginName)
	}
	
	return pool.GetStatus(), nil
}

// GetAllStatus 获取所有插件状态
func (pm *PluginManager) GetAllStatus() map[string]interface{} {
	pm.Mutex.RLock()
	defer pm.Mutex.RUnlock()
	
	status := make(map[string]interface{})
	
	for pluginName, pool := range pm.Pools {
		status[pluginName] = pool.GetStatus()
	}
	
	return map[string]interface{}{
		"is_running": pm.IsRunning,
		"total_plugins": len(pm.Pools),
		"plugins": status,
	}
}

// Stop 停止插件管理器
func (pm *PluginManager) Stop() {
	pm.Mutex.Lock()
	defer pm.Mutex.Unlock()
	
	if !pm.IsRunning {
		return
	}
	
	// 停止所有插件池
	for _, pool := range pm.Pools {
		pool.Stop()
	}
	
	// 清空插件池映射
	pm.Pools = make(map[string]*PluginPool)
	
	pm.IsRunning = false
}

// RestartPlugin 重启插件
func (pm *PluginManager) RestartPlugin(pluginName string) error {
	pm.Mutex.Lock()
	defer pm.Mutex.Unlock()
	
	if !pm.IsRunning {
		return fmt.Errorf("插件管理器未运行")
	}
	
	pool, exists := pm.Pools[pluginName]
	if !exists {
		return fmt.Errorf("插件 %s 不存在", pluginName)
	}
	
	// 停止插件池
	pool.Stop()
	
	// 重新创建插件池
	pluginConfig, exists := pm.Config.Plugins[pluginName]
	if !exists {
		return fmt.Errorf("插件 %s 的配置不存在", pluginName)
	}
	
	newPool := NewPluginPool(pluginName, &pluginConfig)
	if err := newPool.Start(); err != nil {
		return fmt.Errorf("重启插件池 %s 失败: %w", pluginName, err)
	}
	
	// 更新插件池
	pm.Pools[pluginName] = newPool
	
	return nil
}

// AddPlugin 动态添加插件
func (pm *PluginManager) AddPlugin(pluginName string, pluginConfig config.PluginConfig) error {
	pm.Mutex.Lock()
	defer pm.Mutex.Unlock()
	
	if !pm.IsRunning {
		return fmt.Errorf("插件管理器未运行")
	}
	
	if _, exists := pm.Pools[pluginName]; exists {
		return fmt.Errorf("插件 %s 已存在", pluginName)
	}
	
	// 更新配置
	pm.Config.Plugins[pluginName] = pluginConfig
	
	// 创建并启动插件池
	pool := NewPluginPool(pluginName, &pluginConfig)
	if err := pool.Start(); err != nil {
		delete(pm.Config.Plugins, pluginName)
		return fmt.Errorf("启动插件池 %s 失败: %w", pluginName, err)
	}
	
	pm.Pools[pluginName] = pool
	
	return nil
}

// RemovePlugin 动态移除插件
func (pm *PluginManager) RemovePlugin(pluginName string) error {
	pm.Mutex.Lock()
	defer pm.Mutex.Unlock()
	
	if !pm.IsRunning {
		return fmt.Errorf("插件管理器未运行")
	}
	
	pool, exists := pm.Pools[pluginName]
	if !exists {
		return fmt.Errorf("插件 %s 不存在", pluginName)
	}
	
	// 停止插件池
	pool.Stop()
	
	// 从配置和池中移除
	delete(pm.Config.Plugins, pluginName)
	delete(pm.Pools, pluginName)
	
	return nil
}
package plugin

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"goproc/config"
	"goproc/sdk"
)

// PluginInstance 插件实例
type PluginInstance struct {
	ID            string
	PluginName    string
	Config        *config.PluginConfig
	Process       *exec.Cmd
	Conn          net.Conn
	Address       string
	IsRunning     bool
	IsConnected   bool
	Functions     []string
	LastUsed      time.Time
	Mutex         sync.RWMutex
	ConnMutex     sync.Mutex // 连接级别的互斥锁，确保同一时间只有一个操作使用连接
	Communication CommunicationChannel
}

// NewPluginInstance 创建新的插件实例
func NewPluginInstance(pluginName string, config *config.PluginConfig, instanceID string) *PluginInstance {
	return &PluginInstance{
		ID:            instanceID,
		PluginName:    pluginName,
		Config:        config,
		IsRunning:     false,
		IsConnected:   false,
		Functions:     make([]string, 0),
		LastUsed:      time.Now(),
		Communication: NewCommunicationChannel(),
	}
}

// Start 启动插件实例
func (pi *PluginInstance) Start() error {
	pi.Mutex.Lock()

	if pi.IsRunning {
		pi.Mutex.Unlock()
		return fmt.Errorf("插件实例 %s 已经在运行", pi.ID)
	}

	// 生成通信地址
	pi.Address = pi.Communication.GenerateAddress(pi.PluginName, pi.ID)

	// 启动插件进程
	if err := pi.startProcess(); err != nil {
		pi.Mutex.Unlock()
		return fmt.Errorf("启动插件进程失败: %w", err)
	}

	// 检查进程是否成功启动
	if err := pi.waitForProcessReady(); err != nil {
		pi.Process.Process.Kill()
		pi.Mutex.Unlock()
		return fmt.Errorf("等待进程启动失败: %w", err)
	}

	// 连接到插件进程创建的监听器
	if err := pi.connectToPlugin(); err != nil {
		pi.Process.Process.Kill()
		pi.Mutex.Unlock()
		return fmt.Errorf("连接到插件进程失败: %w", err)
	}

	// 在等待注册之前释放锁，避免死锁
	pi.Mutex.Unlock()

	// 等待插件注册
	if err := pi.waitForRegistration(); err != nil {
		pi.Conn.Close()
		pi.Process.Process.Kill()
		return fmt.Errorf("等待插件注册失败: %w", err)
	}

	// 注册完成后重新获取锁来更新状态
	pi.Mutex.Lock()
	pi.IsRunning = true
	pi.IsConnected = true
	pi.Mutex.Unlock()

	return nil
}

// startProcess 启动插件进程
func (pi *PluginInstance) startProcess() error {
	command, args := pi.Config.GetPluginCommand()

	// 添加通信地址参数
	args = append(args, pi.Address)

	pi.Process = exec.Command(command, args...)

	// 设置环境变量
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOPROC_PLUGIN_ADDRESS=%s", pi.Address))

	// 添加配置中的环境变量
	for key, value := range pi.Config.Environment {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	pi.Process.Env = env

	// 设置标准输出和错误输出
	pi.Process.Stdout = os.Stdout
	pi.Process.Stderr = os.Stderr

	// 设置独立的工作目录，避免权限冲突
	if pi.Process.Dir == "" {
		pi.Process.Dir = filepath.Dir(command)
	}

	// 启动进程
	if err := pi.Process.Start(); err != nil {
		return fmt.Errorf("启动进程失败: %w", err)
	}

	return nil
}

// waitForProcessReady 等待进程启动就绪
// Wait for process to be ready
func (pi *PluginInstance) waitForProcessReady() error {
	// 设置最大等待时间为500ms，比原来的1秒快很多
	// Set maximum wait time to 500ms, much faster than original 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 检查进程状态的间隔时间
	// Interval for checking process status
	checkInterval := 10 * time.Millisecond
	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("等待进程启动超时")
		default:
			// 检查进程是否还在运行
			// Check if process is still running
			if pi.Process.Process != nil {
				// 检查进程是否已退出
				// Check if process has exited
				if processState := pi.Process.ProcessState; processState != nil && processState.Exited() {
					return fmt.Errorf("进程已退出")
				}
				
				// 给进程一个最小的启动时间（100ms），然后就认为可以尝试连接了
				// Give process a minimum startup time (100ms), then consider it ready for connection attempts
				if time.Since(startTime) >= 100*time.Millisecond {
					return nil
				}
				
				// 短暂等待后继续检查
				// Wait briefly before checking again
				time.Sleep(checkInterval)
			} else {
				return fmt.Errorf("进程对象为空")
			}
		}
	}
}

// connectToPlugin 连接到插件进程
// Connect to plugin process
func (pi *PluginInstance) connectToPlugin() error {
	// 设置连接超时为10秒，比原来的30秒更合理
	// Set connection timeout to 10 seconds, more reasonable than original 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 尝试连接插件进程
	// Attempt to connect to plugin process
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("连接到插件进程超时")
		default:
			// 尝试连接
			// Attempt connection
			conn, err := pi.Communication.Dial(pi.Address)
			if err == nil {
				pi.Conn = conn
				return nil
			}

			// 连接失败，等待后重试（从500ms减少到100ms）
			// Connection failed, wait before retry (reduced from 500ms to 100ms)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// waitForRegistration 等待插件注册
// Wait for plugin registration
func (pi *PluginInstance) waitForRegistration() error {
	protocol := NewMessageProtocol(pi.Conn)

	// 设置总超时（10秒，比原来的30秒更合理）
	// Set total timeout (10 seconds, more reasonable than original 30 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		// 检查是否超时
		select {
		case <-ctx.Done():
			return fmt.Errorf("等待插件注册超时")
		default:
			// 继续等待
		}

		// 设置单次读取超时（1秒）
		pi.Conn.SetReadDeadline(time.Now().Add(1 * time.Second))

		data, err := protocol.ReceiveMessage()
		if err != nil {
			// 如果是超时错误，继续等待
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			return err
		}

		msg, err := sdk.DecodeMessage(data)
		if err != nil {
			continue // 继续等待有效消息
		}

		// 检查是否为注册消息
		if msg.Type == sdk.MessageTypeRegister {
			// 提取注册的函数列表
			functionsParam := msg.Params["functions"]
			if functionsParam == nil {
				continue
			}

			// 处理不同类型的functions字段
			var functionNames []string

			switch functions := functionsParam.(type) {
			case []interface{}:
				// 处理[]interface{}类型
				functionNames = make([]string, len(functions))
				for i, fn := range functions {
					if name, ok := fn.(string); ok {
						functionNames[i] = name
					}
				}
			case []string:
				// 处理[]string类型
				functionNames = functions
			default:
				continue
			}

			if len(functionNames) > 0 {
				pi.RegisterFunctions(functionNames)

				// 发送注册确认消息
				ackMsg := &sdk.Message{
					Type: sdk.MessageTypeRegisterAck,
				}
				if err := pi.sendMessage(ackMsg); err != nil {
					return err
				}
				return nil
			}
		}
		// 继续等待注册消息
	}
}

// CallFunction 调用插件函数
func (pi *PluginInstance) CallFunction(functionName string, params map[string]interface{}) (interface{}, error) {
	// 优化锁操作：一次性检查所有前置条件，减少锁的获取和释放次数
	// Optimize lock operations: check all preconditions at once, reduce lock acquisition/release frequency
	pi.Mutex.RLock()
	isConnected := pi.IsConnected && pi.Conn != nil
	hasFunc := pi.hasFunction(functionName)
	pi.Mutex.RUnlock()

	if !isConnected {
		return nil, fmt.Errorf("插件实例 %s 未连接", pi.ID)
	}

	if !hasFunc {
		return nil, fmt.Errorf("插件实例 %s 不支持函数 %s", pi.ID, functionName)
	}

	// 获取连接级别的互斥锁，确保同一时间只有一个操作使用连接
	pi.ConnMutex.Lock()
	defer pi.ConnMutex.Unlock()

	// 生成消息ID
	messageID := fmt.Sprintf("call-%d", time.Now().UnixNano())

	// 构造调用消息
	callMsg := &sdk.Message{
		Type:     sdk.MessageTypeCall,
		ID:       messageID,
		Function: functionName,
		Params:   params,
	}

	// 发送消息
	if err := pi.sendMessage(callMsg); err != nil {
		return nil, fmt.Errorf("发送调用消息失败: %w", err)
	}

	// 接收响应（设置30秒超时）
	// Receive response with 30-second timeout
	pi.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	defer pi.Conn.SetReadDeadline(time.Time{}) // 清除超时设置

	protocol := NewMessageProtocol(pi.Conn)

	// 使用高效的阻塞读取，避免轮询式等待
	// Use efficient blocking read, avoid polling-style waiting
	for {
		// 直接阻塞读取消息，依赖连接级别的超时
		// Directly block read message, rely on connection-level timeout
		data, err := protocol.ReceiveMessage()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				//fmt.Printf("[Instance] 等待响应超时，消息ID: %s\n", messageID)
				return nil, fmt.Errorf("等待响应超时")
			}
			//fmt.Printf("[Instance] 接收消息错误: %v\n", err)
			return nil, fmt.Errorf("接收响应失败: %w", err)
		}

		msg, err := sdk.DecodeMessage(data)
		if err != nil {
			//fmt.Printf("[Instance] 解码消息错误: %v\n", err)
			continue // 解码错误继续等待
		}

		//fmt.Printf("[Instance] 收到消息: ID=%s, Type=%s\n", msg.ID, msg.Type)

		// 检查消息ID是否匹配
		if msg.ID == messageID {
			pi.LastUsed = time.Now()

			if msg.Type == sdk.MessageTypeResult {
				//fmt.Printf("[Instance] 收到结果响应: %v\n", msg.Result)
				return msg.Result, nil
			} else if msg.Type == sdk.MessageTypeError {
				//fmt.Printf("[Instance] 收到错误响应: %s\n", msg.Error)
				return nil, fmt.Errorf("插件返回错误: %s", msg.Error)
			} else {
				//fmt.Printf("[Instance] 收到未知的响应类型: %s\n", msg.Type)
				return nil, fmt.Errorf("收到未知的响应类型: %s", msg.Type)
			}
		}

		// 处理心跳消息
		if msg.Type == sdk.MessageTypePing {
			// 发送pong响应
			pongMsg := &sdk.Message{
				Type: sdk.MessageTypePong,
			}
			pongData, _ := sdk.EncodeMessage(pongMsg)
			protocol.SendMessage(pongData)
			continue
		}
	}
}

// sendMessage 发送消息
func (pi *PluginInstance) sendMessage(msg *sdk.Message) error {
	data, err := sdk.EncodeMessage(msg)
	if err != nil {
		return err
	}

	// 使用消息协议发送
	protocol := NewMessageProtocol(pi.Conn)
	return protocol.SendMessage(data)
}

// waitForResponse 等待响应（保留旧方法，但不再使用）
func (pi *PluginInstance) waitForResponse(messageID string) (interface{}, error) {
	protocol := NewMessageProtocol(pi.Conn)

	// 设置读取超时
	pi.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	for {
		data, err := protocol.ReceiveMessage()
		if err != nil {
			return nil, err
		}

		msg, err := sdk.DecodeMessage(data)
		if err != nil {
			return nil, err
		}

		// 检查消息ID是否匹配
		if msg.ID == messageID {
			switch msg.Type {
			case sdk.MessageTypeResult:
				return msg.Result, nil
			case sdk.MessageTypeError:
				return nil, fmt.Errorf("插件返回错误: %s", msg.Error)
			default:
				return nil, fmt.Errorf("收到未知的响应类型: %s", msg.Type)
			}
		}
		// 不匹配的消息，继续等待目标消息
	}
}

// hasFunction 检查是否支持指定函数
func (pi *PluginInstance) hasFunction(functionName string) bool {
	for _, fn := range pi.Functions {
		if fn == functionName {
			return true
		}
	}
	return false
}

// RegisterFunctions 注册函数列表
func (pi *PluginInstance) RegisterFunctions(functions []string) {
	pi.Mutex.Lock()
	defer pi.Mutex.Unlock()

	pi.Functions = functions
}

// HealthCheck 健康检查
func (pi *PluginInstance) HealthCheck() bool {
	pi.Mutex.RLock()
	defer pi.Mutex.RUnlock()

	if !pi.IsConnected || pi.Conn == nil {
		return false
	}

	// 发送ping消息检查连接
	pingMsg := &sdk.Message{
		Type: sdk.MessageTypePing,
		ID:   fmt.Sprintf("healthcheck-%d", time.Now().UnixNano()),
	}

	if err := pi.sendMessage(pingMsg); err != nil {
		return false
	}

	// 等待pong响应
	_, err := pi.waitForResponse(pingMsg.ID)
	return err == nil
}

// Stop 停止插件实例
func (pi *PluginInstance) Stop() error {
	pi.Mutex.Lock()
	defer pi.Mutex.Unlock()

	if !pi.IsRunning {
		return nil
	}

	// 1. 先发送停止信号给插件进程（优雅关闭）
	if pi.IsConnected && pi.Conn != nil {
		// 尝试发送停止消息
		stopMsg := &sdk.Message{
			Type: sdk.MessageTypeStop,
			ID:   fmt.Sprintf("stop-%d", time.Now().UnixNano()),
		}

		pi.sendMessage(stopMsg) // 移除未使用的err变量

		// 等待插件进程优雅退出（最多等待2秒）
		done := make(chan bool, 1)
		go func() {
			if pi.Process != nil && pi.Process.Process != nil {
				pi.Process.Wait() // 移除未使用的err变量
				done <- true
			} else {
				done <- false
			}
		}()

		select {
		case <-done:
			pi.IsRunning = false
			pi.IsConnected = false

			// 清理通信资源
			pi.Communication.Cleanup(pi.Address)

			return nil
		case <-time.After(2 * time.Second):
			// 优雅退出超时，继续强制终止
		}
	}

	// 2. 如果优雅关闭失败，强制终止进程
	if pi.Process != nil && pi.Process.Process != nil {
		// 先关闭连接
		if pi.Conn != nil {
			pi.Conn.Close()
		}

		// 终止进程
		pi.Process.Process.Kill()

		// 等待进程结束
		done := make(chan error, 1)
		go func() {
			done <- pi.Process.Wait()
		}()

		select {
		case <-time.After(5 * time.Second):
			// 强制终止
		case <-done:
			// 进程结束
		}
	}

	// 3. 清理通信资源
	pi.Communication.Cleanup(pi.Address)

	pi.IsRunning = false
	pi.IsConnected = false

	return nil
}

// GetStatus 获取实例状态
func (pi *PluginInstance) GetStatus() map[string]interface{} {
	pi.Mutex.RLock()
	defer pi.Mutex.RUnlock()

	return map[string]interface{}{
		"id":           pi.ID,
		"plugin_name":  pi.PluginName,
		"is_running":   pi.IsRunning,
		"is_connected": pi.IsConnected,
		"functions":    pi.Functions,
		"last_used":    pi.LastUsed.Format(time.RFC3339),
		"address":      pi.Address,
	}
}

package sdk

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// PlatformCommunication 平台特定通信接口
// PlatformCommunication Platform-specific communication interface
type PlatformCommunication interface {
	// CreateListener 创建监听器 / Create listener
	CreateListener(address string) (net.Listener, error)
	// Connect 连接到指定地址 / Connect to specified address
	Connect(address string) (net.Conn, error)
	// GetCommunicationAddress 获取通信地址 / Get communication address
	GetCommunicationAddress() string
}

// PluginSDK 插件SDK
// PluginSDK Plugin SDK
type PluginSDK struct {
	functions map[string]FunctionHandler // 注册的函数 / Registered functions
	conn      net.Conn                   // 连接对象 / Connection object
	listener  net.Listener               // 监听器对象 / Listener object
	isRunning bool                       // 运行状态 / Running status
	platform  PlatformCommunication     // 平台特定通信实现 / Platform-specific communication implementation
}

// NewPluginSDK 创建新的插件SDK
// NewPluginSDK Create new plugin SDK
func NewPluginSDK() *PluginSDK {
	return &PluginSDK{
		functions: make(map[string]FunctionHandler),
		isRunning: false,
		platform:  newPlatformCommunication(), // 使用平台特定实现 / Use platform-specific implementation
	}
}

// RegisterFunction 注册函数
func (sdk *PluginSDK) RegisterFunction(name string, handler FunctionHandler) error {
	if sdk.isRunning {
		return fmt.Errorf("插件已启动，无法注册新函数")
	}

	if _, exists := sdk.functions[name]; exists {
		return fmt.Errorf("函数 %s 已注册", name)
	}

	sdk.functions[name] = handler
	return nil
}

// Start 启动插件SDK
func (sdk *PluginSDK) Start() error {
	if sdk.isRunning {
		return fmt.Errorf("插件已启动")
	}

	// 获取通信地址
	address := sdk.getCommunicationAddress()
	if address == "" {
		return fmt.Errorf("无法获取通信地址")
	}

	// 创建监听器并等待连接
	listener, conn, err := sdk.createListenerAndWait(address)
	if err != nil {
		return fmt.Errorf("创建监听器失败: %w", err)
	}
	sdk.listener = listener
	sdk.conn = conn
	sdk.isRunning = true

	// 发送注册消息并等待确认
	if err := sdk.sendRegisterMessageAndWait(); err != nil {
		return fmt.Errorf("注册失败: %w", err)
	}

	// 启动消息处理循环
	go sdk.messageLoop()

	return nil
}

// getCommunicationAddress 获取通信地址
// getCommunicationAddress Get communication address
func (sdk *PluginSDK) getCommunicationAddress() string {
	// 从环境变量获取通信地址 / Get communication address from environment variable
	if address := os.Getenv("GOPROC_PLUGIN_ADDRESS"); address != "" {
		return address
	}

	// 从命令行参数获取 / Get from command line arguments
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	// 使用平台特定的默认地址 / Use platform-specific default address
	return sdk.platform.GetCommunicationAddress()
}

// createListenerAndWait 创建监听器并等待连接
// createListenerAndWait Create listener and wait for connection
func (sdk *PluginSDK) createListenerAndWait(address string) (net.Listener, net.Conn, error) {
	// 使用平台特定的监听器创建方法 / Use platform-specific listener creation method
	listener, err := sdk.platform.CreateListener(address)
	if err != nil {
		return nil, nil, fmt.Errorf("创建监听器失败: %w", err)
	}

	// 设置连接超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	connChan := make(chan net.Conn, 1)
	errChan := make(chan error, 1)

	// 异步接受连接
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- err
			return
		}
		connChan <- conn
	}()

	// 等待连接或超时
	select {
	case conn := <-connChan:
		return listener, conn, nil
	case err := <-errChan:
		listener.Close()
		return nil, nil, fmt.Errorf("接受连接失败: %w", err)
	case <-ctx.Done():
		listener.Close()
		return nil, nil, fmt.Errorf("等待连接超时")
	}
}

// connect 建立连接
// connect Establish connection
func (sdk *PluginSDK) connect(address string) (net.Conn, error) {
	// 尝试连接，最多重试5次 / Try to connect, retry up to 5 times
	var conn net.Conn
	var err error

	for i := 0; i < 5; i++ {
		// 使用平台特定的连接方法 / Use platform-specific connection method
		conn, err = sdk.platform.Connect(address)

		if err == nil {
			return conn, nil
		}

		time.Sleep(time.Second * time.Duration(i+1))
	}

	return nil, err
}

// sendRegisterMessage 发送注册消息
func (sdk *PluginSDK) sendRegisterMessage() error {
	functions := make([]string, 0, len(sdk.functions))
	for name := range sdk.functions {
		functions = append(functions, name)
	}

	msg := &Message{
		Type: MessageTypeRegister,
		Params: map[string]interface{}{
			"functions": functions,
		},
	}

	return sdk.sendMessage(msg)
}

// sendRegisterMessageAndWait 发送注册消息并等待确认
func (sdk *PluginSDK) sendRegisterMessageAndWait() error {
	// 发送注册消息
	if err := sdk.sendRegisterMessage(); err != nil {
		return fmt.Errorf("发送注册消息失败: %w", err)
	}

	// 移除超时设置，使用阻塞模式
	sdk.conn.SetReadDeadline(time.Time{})

	// 等待注册确认响应
	buffer := make([]byte, 4096)
	var messageBuffer []byte

	for {
		n, err := sdk.conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("连接已关闭")
			}
			return fmt.Errorf("读取注册确认失败: %w", err)
		}

		messageBuffer = append(messageBuffer, buffer[:n]...)

		// 处理完整的消息
		for len(messageBuffer) >= 4 {
			length := int(messageBuffer[0])<<24 | int(messageBuffer[1])<<16 |
				int(messageBuffer[2])<<8 | int(messageBuffer[3])

			if len(messageBuffer) < 4+length {
				break // 消息不完整，继续读取
			}

			messageData := messageBuffer[4 : 4+length]
			messageBuffer = messageBuffer[4+length:]

			// 解码消息
			msg, err := DecodeMessage(messageData)
			if err != nil {
				return fmt.Errorf("解码注册确认消息失败: %w", err)
			}

			// 检查是否为注册确认消息
			if msg.Type == MessageTypeRegisterAck {
				return nil
			}
			// 如果是其他消息类型，继续等待注册确认
		}
	}
}

// sendMessage 发送消息
func (sdk *PluginSDK) sendMessage(msg *Message) error {
	data, err := EncodeMessage(msg)
	if err != nil {
		return err
	}

	// 添加长度前缀
	length := len(data)
	header := []byte{
		byte(length >> 24),
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	}

	_, err = sdk.conn.Write(append(header, data...))
	return err
}

// messageLoop 消息处理循环
func (sdk *PluginSDK) messageLoop() {
	defer sdk.conn.Close()

	buffer := make([]byte, 4096)
	var messageBuffer []byte

	for sdk.isRunning {
		n, err := sdk.conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		messageBuffer = append(messageBuffer, buffer[:n]...)

		// 处理完整的消息
		for len(messageBuffer) >= 4 {
			length := int(messageBuffer[0])<<24 | int(messageBuffer[1])<<16 |
				int(messageBuffer[2])<<8 | int(messageBuffer[3])

			if len(messageBuffer) < 4+length {
				break // 消息不完整，继续读取
			}

			messageData := messageBuffer[4 : 4+length]
			messageBuffer = messageBuffer[4+length:]

			sdk.handleMessage(messageData)
		}
	}

	sdk.isRunning = false
}

// handleMessage 处理消息
func (sdk *PluginSDK) handleMessage(data []byte) {
	msg, err := DecodeMessage(data)
	if err != nil {
		return
	}

	switch msg.Type {
	case MessageTypeCall:
		sdk.handleCallMessage(msg)
	case MessageTypePing:
		sdk.handlePingMessage(msg)
	case MessageTypePong:
		sdk.handlePongMessage(msg)
	case MessageTypeStop:
		sdk.handleStopMessage(msg)
	}
}

// handleCallMessage 处理调用消息
func (sdk *PluginSDK) handleCallMessage(msg *Message) {
	handler, exists := sdk.functions[msg.Function]
	if !exists {
		sdk.sendErrorMessage(msg.ID, fmt.Sprintf("函数 %s 不存在", msg.Function))
		return
	}

	// 异步处理函数调用
	go func() {
		result, err := handler(msg.Params)
		if err != nil {
			sdk.sendErrorMessage(msg.ID, err.Error())
			return
		}

		sdk.sendResultMessage(msg.ID, result)
	}()
}

// handlePingMessage 处理心跳消息
func (sdk *PluginSDK) handlePingMessage(msg *Message) {
	pongMsg := &Message{
		Type: MessageTypePong,
		ID:   msg.ID,
	}
	sdk.sendMessage(pongMsg)
}

// handlePongMessage 处理心跳响应消息
func (sdk *PluginSDK) handlePongMessage(msg *Message) {
}

// handleStopMessage 处理停止消息
func (sdk *PluginSDK) handleStopMessage(msg *Message) {
	sdk.isRunning = false
}

// sendResultMessage 发送结果消息
func (sdk *PluginSDK) sendResultMessage(id string, result interface{}) {
	msg := &Message{
		Type:   MessageTypeResult,
		ID:     id,
		Result: result,
	}
	sdk.sendMessage(msg)
}

// sendErrorMessage 发送错误消息
func (sdk *PluginSDK) sendErrorMessage(id string, errorMsg string) {
	msg := &Message{
		Type:  MessageTypeError,
		ID:    id,
		Error: errorMsg,
	}
	sdk.sendMessage(msg)
}

// Wait 等待插件运行直到停止
func (sdk *PluginSDK) Wait() {
	for sdk.isRunning {
		time.Sleep(100 * time.Millisecond)
	}
}

// Stop 停止插件SDK
func (sdk *PluginSDK) Stop() {
	sdk.isRunning = false
	if sdk.conn != nil {
		sdk.conn.Close()
	}
}

// 全局SDK实例
var globalSDK = NewPluginSDK()

// RegisterFunction 全局注册函数
func RegisterFunction(name string, handler FunctionHandler) error {
	return globalSDK.RegisterFunction(name, handler)
}

// Start 全局启动函数
func Start() error {
	return globalSDK.Start()
}

// Wait 全局等待函数
func Wait() {
	globalSDK.Wait()
}

// Stop 全局停止函数
func Stop() {
	globalSDK.Stop()
}

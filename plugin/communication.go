package plugin

import (
	"fmt"
	"io"
	"net"
	"time"
)

// CommunicationType 通信类型
// CommunicationType Communication type
type CommunicationType string

const (
	CommunicationTypePipe CommunicationType = "pipe" // 命名管道（Windows） / Named pipe (Windows)
	CommunicationTypeUnix CommunicationType = "unix" // Unix域套接字（非Windows） / Unix domain socket (non-Windows)
)

// CommunicationChannel 通信通道接口
// CommunicationChannel Communication channel interface
type CommunicationChannel interface {
	// 连接相关 / Connection related
	Dial(address string) (net.Conn, error)
	Listen(address string) (net.Listener, error)
	
	// 地址生成 / Address generation
	GenerateAddress(pluginName string, instanceID string) string
	
	// 清理 / Cleanup
	Cleanup(address string) error
}

// NewCommunicationChannel 创建通信通道
// NewCommunicationChannel Create communication channel
func NewCommunicationChannel() CommunicationChannel {
	return newPlatformCommunicationChannel()
}

// MessageProtocol 消息协议处理
// MessageProtocol Message protocol handling
type MessageProtocol struct {
	conn net.Conn
}

// NewMessageProtocol 创建消息协议
// NewMessageProtocol Create message protocol
func NewMessageProtocol(conn net.Conn) *MessageProtocol {
	return &MessageProtocol{conn: conn}
}

// SendMessage 发送消息
// SendMessage Send message
func (mp *MessageProtocol) SendMessage(data []byte) error {
	// 添加长度前缀 / Add length prefix
	length := len(data)
	header := []byte{
		byte(length >> 24),
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	}
	
	// 发送头部和数据 / Send header and data
	fullData := append(header, data...)
	
	_, err := mp.conn.Write(fullData)
	if err != nil {
		return err
	}
	
	return nil
}

// ReceiveMessage 接收消息
// ReceiveMessage Receive message
func (mp *MessageProtocol) ReceiveMessage() ([]byte, error) {
	// 读取消息头部（4字节长度） / Read message header (4-byte length)
	header := make([]byte, 4)
	_, err := io.ReadFull(mp.conn, header)
	if err != nil {
		return nil, err
	}
	
	length := int(header[0])<<24 | int(header[1])<<16 | 
		int(header[2])<<8 | int(header[3])
	
	// 读取消息体 / Read message body
	data := make([]byte, length)
	_, err = io.ReadFull(mp.conn, data)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

// HeartbeatChecker 心跳检查器
// HeartbeatChecker Heartbeat checker
type HeartbeatChecker struct {
	conn     net.Conn
	timeout  time.Duration
	interval time.Duration
}

// NewHeartbeatChecker 创建心跳检查器
// NewHeartbeatChecker Create heartbeat checker
func NewHeartbeatChecker(conn net.Conn, timeout time.Duration, interval time.Duration) *HeartbeatChecker {
	return &HeartbeatChecker{
		conn:     conn,
		timeout:  timeout,
		interval: interval,
	}
}

// Start 开始心跳检查
// Start Start heartbeat checking
func (h *HeartbeatChecker) Start() <-chan error {
	errChan := make(chan error, 1)
	
	go func() {
		ticker := time.NewTicker(h.interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				if err := h.sendPing(); err != nil {
					errChan <- err
					return
				}
				
				if err := h.waitForPong(); err != nil {
					errChan <- err
					return
				}
			}
		}
	}()
	
	return errChan
}

// sendPing 发送ping消息
// sendPing Send ping message
func (h *HeartbeatChecker) sendPing() error {
	// 使用消息协议发送ping消息 / Use message protocol to send ping message
	protocol := NewMessageProtocol(h.conn)
	pingMsg := []byte("ping")
	
	// 设置写超时 / Set write timeout
	h.conn.SetWriteDeadline(time.Now().Add(h.timeout))
	
	if err := protocol.SendMessage(pingMsg); err != nil {
		return fmt.Errorf("发送ping失败: %w", err)
	}
	
	return nil
}

// waitForPong 等待pong响应
// waitForPong Wait for pong response
func (h *HeartbeatChecker) waitForPong() error {
	// 使用消息协议接收pong消息 / Use message protocol to receive pong message
	protocol := NewMessageProtocol(h.conn)
	
	// 设置读超时 / Set read timeout
	h.conn.SetReadDeadline(time.Now().Add(h.timeout))
	
	data, err := protocol.ReceiveMessage()
	if err != nil {
		return fmt.Errorf("等待pong失败: %w", err)
	}
	
	if string(data) != "pong" {
		return fmt.Errorf("收到无效的pong响应: %s", string(data))
	}
	
	return nil
}
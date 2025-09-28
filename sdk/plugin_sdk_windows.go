//go:build windows
// +build windows

package sdk

import (
	"fmt"
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

// WindowsCommunication Windows平台通信实现
// WindowsCommunication Windows platform communication implementation
type WindowsCommunication struct{}

// newPlatformCommunication 创建平台特定的通信实现
// newPlatformCommunication Create platform-specific communication implementation
func newPlatformCommunication() PlatformCommunication {
	return &WindowsCommunication{}
}

// CreateListener 创建Windows命名管道监听器
// CreateListener Create Windows named pipe listener
func (w *WindowsCommunication) CreateListener(address string) (net.Listener, error) {
	// 使用winio.ListenPipe创建Windows命名管道监听器
	// Use winio.ListenPipe to create Windows named pipe listener
	listener, err := winio.ListenPipe(address, nil)
	if err != nil {
		return nil, fmt.Errorf("创建Windows命名管道监听器失败: %w", err)
	}
	return listener, nil
}

// Connect 连接到Windows命名管道
// Connect Connect to Windows named pipe
func (w *WindowsCommunication) Connect(address string) (net.Conn, error) {
	// 使用winio.DialPipe连接Windows命名管道
	// Use winio.DialPipe to connect to Windows named pipe
	conn, err := winio.DialPipe(address, nil)
	if err != nil {
		return nil, fmt.Errorf("连接Windows命名管道失败: %w", err)
	}
	return conn, nil
}

// GetCommunicationAddress 获取Windows平台默认通信地址
// GetCommunicationAddress Get Windows platform default communication address
func (w *WindowsCommunication) GetCommunicationAddress() string {
	// 生成Windows命名管道地址
	// Generate Windows named pipe address
	processID := os.Getpid()
	return fmt.Sprintf("\\\\.\\pipe\\goproc_plugin_%d", processID)
}
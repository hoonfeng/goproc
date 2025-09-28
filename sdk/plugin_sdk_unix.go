//go:build !windows
// +build !windows

package sdk

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// UnixCommunication Unix平台通信实现
// UnixCommunication Unix platform communication implementation
type UnixCommunication struct{}

// newPlatformCommunication 创建平台特定的通信实现
// newPlatformCommunication Create platform-specific communication implementation
func newPlatformCommunication() PlatformCommunication {
	return &UnixCommunication{}
}

// CreateListener 创建Unix域套接字监听器
// CreateListener Create Unix domain socket listener
func (u *UnixCommunication) CreateListener(address string) (net.Listener, error) {
	// 确保套接字文件不存在 / Ensure socket file does not exist
	os.Remove(address)
	
	// 创建目录 / Create directory
	dir := filepath.Dir(address)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}
	
	// 使用net.Listen创建Unix域套接字监听器
	// Use net.Listen to create Unix domain socket listener
	listener, err := net.Listen("unix", address)
	if err != nil {
		return nil, fmt.Errorf("创建Unix域套接字监听器失败: %w", err)
	}
	
	return listener, nil
}

// Connect 连接到Unix域套接字
// Connect Connect to Unix domain socket
func (u *UnixCommunication) Connect(address string) (net.Conn, error) {
	// 使用net.Dial连接Unix域套接字
	// Use net.Dial to connect to Unix domain socket
	conn, err := net.Dial("unix", address)
	if err != nil {
		return nil, fmt.Errorf("连接Unix域套接字失败: %w", err)
	}
	return conn, nil
}

// GetCommunicationAddress 获取Unix平台默认通信地址
// GetCommunicationAddress Get Unix platform default communication address
func (u *UnixCommunication) GetCommunicationAddress() string {
	// 生成Unix域套接字地址
	// Generate Unix domain socket address
	processID := os.Getpid()
	tmpDir := os.TempDir()
	return filepath.Join(tmpDir, fmt.Sprintf("goproc_plugin_%d.sock", processID))
}
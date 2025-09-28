//go:build !windows
// +build !windows

package plugin

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// UnixCommunication Unix域套接字通信（非Windows）
// UnixCommunication Unix domain socket communication (non-Windows)
type UnixCommunication struct{}

// newPlatformCommunicationChannel 创建平台特定的通信通道
// newPlatformCommunicationChannel Create platform-specific communication channel
func newPlatformCommunicationChannel() CommunicationChannel {
	return &UnixCommunication{}
}

// Dial 建立连接（Unix套接字）
// Dial Establish connection (Unix socket)
func (u *UnixCommunication) Dial(address string) (net.Conn, error) {
	return net.Dial("unix", address)
}

// Listen 监听连接（Unix套接字）
// Listen Listen for connections (Unix socket)
func (u *UnixCommunication) Listen(address string) (net.Listener, error) {
	// 确保套接字文件不存在
	// Ensure socket file does not exist
	os.Remove(address)
	
	// 创建目录
	// Create directory
	dir := filepath.Dir(address)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	
	return net.Listen("unix", address)
}

// GenerateAddress 生成Unix套接字地址
// GenerateAddress Generate Unix socket address
func (u *UnixCommunication) GenerateAddress(pluginName string, instanceID string) string {
	// 使用/tmp/目录作为套接字文件存放位置
	// Use /tmp/ directory as socket file location
	socketsDir := "/tmp/goproc_sockets"
	
	// Unix域套接字格式：/tmp/goproc_sockets/pluginname-instanceid.sock
	// Unix domain socket format: /tmp/goproc_sockets/pluginname-instanceid.sock
	return filepath.Join(socketsDir, fmt.Sprintf("%s-%s.sock", pluginName, instanceID))
}

// Cleanup 清理Unix套接字
// Cleanup Clean up Unix socket
func (u *UnixCommunication) Cleanup(address string) error {
	return os.Remove(address)
}
//go:build windows
// +build windows

package plugin

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

// PipeCommunication 命名管道通信（Windows）
// PipeCommunication Named pipe communication (Windows)
type PipeCommunication struct{}

// newPlatformCommunicationChannel 创建平台特定的通信通道
// newPlatformCommunicationChannel Create platform-specific communication channel
func newPlatformCommunicationChannel() CommunicationChannel {
	return &PipeCommunication{}
}

// Dial 建立连接（管道）
// Dial Establish connection (pipe)
func (p *PipeCommunication) Dial(address string) (net.Conn, error) {
	// 使用go-winio库连接Windows命名管道
	// Use go-winio library to connect to Windows named pipe
	return winio.DialPipe(address, nil)
}

// Listen 监听连接（管道）
// Listen Listen for connections (pipe)
func (p *PipeCommunication) Listen(address string) (net.Listener, error) {
	// 使用go-winio库创建Windows命名管道监听器
	// Use go-winio library to create Windows named pipe listener
	return winio.ListenPipe(address, nil)
}

// GenerateAddress 生成管道地址
// GenerateAddress Generate pipe address
func (p *PipeCommunication) GenerateAddress(pluginName string, instanceID string) string {
	// Windows命名管道格式：\\.\pipe\instanceid
	// 注意：instanceID已经包含了插件名称，不需要重复拼接
	// Windows named pipe format: \\.\pipe\instanceid
	// Note: instanceID already contains the plugin name, no need to concatenate again
	return fmt.Sprintf("\\\\.\\pipe\\%s", instanceID)
}

// Cleanup 清理管道
// Cleanup Clean up pipe
func (p *PipeCommunication) Cleanup(address string) error {
	// Windows命名管道不需要手动清理文件
	// Windows named pipes do not require manual file cleanup
	return nil
}
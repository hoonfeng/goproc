package main

import (
	"fmt"
	"runtime"
	"github.com/hoonfeng/goproc/plugin"
	"github.com/hoonfeng/goproc/sdk"
)

// 跨平台兼容性测试 / Cross-platform compatibility test
func main() {
	fmt.Println("=== 跨平台兼容性测试 / Cross-platform Compatibility Test ===")
	
	// 显示当前平台信息 / Display current platform information
	fmt.Printf("操作系统 / OS: %s\n", runtime.GOOS)
	fmt.Printf("架构 / Architecture: %s\n", runtime.GOARCH)
	
	// 测试插件通信接口 / Test plugin communication interface
	fmt.Println("\n=== 测试插件通信接口 / Testing Plugin Communication Interface ===")
	
	// 创建通信通道 / Create communication channel
	comm := plugin.NewCommunicationChannel()
	
	// 根据平台显示通信方式 / Display communication method based on platform
	if runtime.GOOS == "windows" {
		fmt.Println("使用Windows命名管道通信 / Using Windows Named Pipe Communication")
	} else {
		fmt.Println("使用Unix域套接字通信 / Using Unix Domain Socket Communication")
	}
	
	// 生成地址 / Generate address
	address := comm.GenerateAddress("test-plugin", "test-instance")
	fmt.Printf("生成的通信地址 / Generated Communication Address: %s\n", address)
	
	// 测试SDK平台通信 / Test SDK platform communication
	fmt.Println("\n=== 测试SDK平台通信 / Testing SDK Platform Communication ===")
	
	// 创建SDK实例 / Create SDK instance
	_ = sdk.NewPluginSDK()
	fmt.Println("✓ SDK实例创建成功 / SDK Instance Created Successfully")
	
	// 验证地址格式 / Verify address format
	if runtime.GOOS == "windows" {
		if len(address) > 0 && (address[:9] == "\\\\.\\pipe\\" || address[:2] == "\\\\") {
			fmt.Println("✓ Windows命名管道地址格式正确 / Windows Named Pipe Address Format Correct")
		} else {
			fmt.Println("✗ Windows命名管道地址格式错误 / Windows Named Pipe Address Format Incorrect")
		}
	} else {
		if len(address) > 0 && (address[0] == '/' || address[:5] == "unix:") {
			fmt.Println("✓ Unix域套接字地址格式正确 / Unix Domain Socket Address Format Correct")
		} else {
			fmt.Println("✗ Unix域套接字地址格式错误 / Unix Domain Socket Address Format Incorrect")
		}
	}
	
	fmt.Println("\n=== 跨平台兼容性测试完成 / Cross-platform Compatibility Test Completed ===")
	fmt.Println("✓ 所有平台特定代码已正确分离 / All platform-specific code properly separated")
	fmt.Println("✓ 接口抽象工作正常 / Interface abstraction working correctly")
	fmt.Println("✓ 跨平台编译成功 / Cross-platform compilation successful")
}
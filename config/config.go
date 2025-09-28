package config

import (
	"fmt"
	"time"
)

// PluginType 插件类型
type PluginType string

const (
	PluginTypeBinary PluginType = "binary" // 二进制插件
	PluginTypeScript PluginType = "script" // 脚本插件
)

// PluginConfig 插件配置
type PluginConfig struct {
	Type                PluginType        `yaml:"type"`                  // 插件类型
	Path                string            `yaml:"path"`                  // 插件可执行文件路径（二进制插件）
	Interpreter         string            `yaml:"interpreter"`           // 解释器（脚本插件）
	ScriptPath          string            `yaml:"script_path"`           // 脚本路径（脚本插件）
	PoolSize            int               `yaml:"pool_size"`             // 初始池大小
	MaxInstances        int               `yaml:"max_instances"`         // 最大实例数
	HealthCheckInterval time.Duration     `yaml:"health_check_interval"` // 健康检查间隔
	Args                []string          `yaml:"args"`                  // 启动参数
	Functions           []string          `yaml:"functions"`             // 插件提供的函数列表
	Environment         map[string]string `yaml:"environment"`           // 环境变量
}

// SystemConfig 系统配置 / System Configuration
type SystemConfig struct {
	Version     string                  `yaml:"version"`     // 配置文件版本 / Configuration version
	Plugins     map[string]PluginConfig `yaml:"plugins"`     // 插件配置 / Plugin configurations
	System      SystemSettings          `yaml:"system"`      // 系统设置 / System settings
	Platform    PlatformConfig          `yaml:"platform"`    // 平台特定配置 / Platform-specific configuration
	Development DevelopmentConfig       `yaml:"development"` // 开发配置 / Development configuration
}

// SystemSettings 系统设置 / System Settings
type SystemSettings struct {
	LogLevel           string `yaml:"log_level"`            // 日志级别 / Log level
	LogFile            string `yaml:"log_file"`             // 日志文件路径 / Log file path
	MaxConcurrentCalls int    `yaml:"max_concurrent_calls"` // 最大并发调用数 / Max concurrent calls
	CallTimeout        string `yaml:"call_timeout"`         // 调用超时时间 / Call timeout
	EnableMetrics      bool   `yaml:"enable_metrics"`       // 启用性能指标收集 / Enable metrics collection
	MetricsPort        int    `yaml:"metrics_port"`         // 指标服务端口 / Metrics service port
	EnableAuth         bool   `yaml:"enable_auth"`          // 启用认证 / Enable authentication
	AuthToken          string `yaml:"auth_token"`           // 认证令牌 / Authentication token
}

// PlatformConfig 平台特定配置 / Platform-specific Configuration
type PlatformConfig struct {
	Windows WindowsConfig `yaml:"windows"` // Windows配置 / Windows configuration
	Unix    UnixConfig    `yaml:"unix"`    // Unix配置 / Unix configuration
}

// WindowsConfig Windows平台配置 / Windows Platform Configuration
type WindowsConfig struct {
	NamedPipePrefix string `yaml:"named_pipe_prefix"` // 命名管道前缀 / Named pipe prefix
	PipeBufferSize  int    `yaml:"pipe_buffer_size"`  // 管道缓冲区大小 / Pipe buffer size
}

// UnixConfig Unix平台配置 / Unix Platform Configuration
type UnixConfig struct {
	SocketDir         string `yaml:"socket_dir"`         // 套接字目录 / Socket directory
	SocketPermissions string `yaml:"socket_permissions"` // 套接字权限 / Socket permissions
}

// DevelopmentConfig 开发配置 / Development Configuration
type DevelopmentConfig struct {
	DebugMode bool `yaml:"debug_mode"` // 调试模式 / Debug mode
	HotReload bool `yaml:"hot_reload"` // 热重载 / Hot reload
	Profiling bool `yaml:"profiling"`  // 性能分析 / Profiling
}

// LoadConfig 加载配置文件

// ValidateConfig 验证配置
func ValidateConfig(config *SystemConfig) error {
	if len(config.Plugins) == 0 {
		return fmt.Errorf("至少需要配置一个插件")
	}

	for name, pluginConfig := range config.Plugins {
		switch pluginConfig.Type {
		case PluginTypeBinary:
			if pluginConfig.Path == "" {
				return fmt.Errorf("二进制插件 %s 的路径不能为空", name)
			}
		case PluginTypeScript:
			if pluginConfig.Interpreter == "" {
				return fmt.Errorf("脚本插件 %s 的解释器不能为空", name)
			}
			if pluginConfig.ScriptPath == "" {
				return fmt.Errorf("脚本插件 %s 的脚本路径不能为空", name)
			}
		default:
			return fmt.Errorf("插件 %s 的类型 %s 不支持", name, pluginConfig.Type)
		}

		if pluginConfig.PoolSize <= 0 {
			pluginConfig.PoolSize = 3
		}
		if pluginConfig.MaxInstances <= 0 {
			pluginConfig.MaxInstances = 10
		}
		if pluginConfig.HealthCheckInterval <= 0 {
			pluginConfig.HealthCheckInterval = 30 * time.Second
		}

		if len(pluginConfig.Functions) == 0 {
			return fmt.Errorf("插件 %s 必须至少提供一个函数", name)
		}
	}

	return nil
}

// GetPluginCommand 获取插件启动命令
func (p *PluginConfig) GetPluginCommand() (string, []string) {
	switch p.Type {
	case PluginTypeBinary:
		return p.Path, p.Args
	case PluginTypeScript:
		args := []string{p.ScriptPath}
		args = append(args, p.Args...)
		return p.Interpreter, args
	default:
		return "", nil
	}
}

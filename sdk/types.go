package sdk

import (
	"encoding/json"
)

// FunctionHandler 函数处理器类型
type FunctionHandler func(params map[string]interface{}) (interface{}, error)

// MessageType 消息类型
type MessageType string

const (
	MessageTypeRegister     MessageType = "register"     // 注册消息
	MessageTypeRegisterAck  MessageType = "register_ack" // 注册确认消息
	MessageTypeCall         MessageType = "call"         // 调用消息
	MessageTypeResult       MessageType = "result"       // 结果消息
	MessageTypeError        MessageType = "error"        // 错误消息
	MessageTypePing         MessageType = "ping"         // 心跳消息
	MessageTypePong         MessageType = "pong"         // 心跳响应
	MessageTypeStop         MessageType = "stop"         // 停止消息
)

// Message 消息结构
type Message struct {
	Type      MessageType       `json:"type"`               // 消息类型
	ID        string            `json:"id,omitempty"`       // 消息ID
	Function  string            `json:"function,omitempty"` // 函数名
	Params    map[string]interface{} `json:"params,omitempty"` // 参数
	Result    interface{}       `json:"result,omitempty"`   // 结果
	Error     string            `json:"error,omitempty"`    // 错误信息
}

// RegisterMessage 注册消息
type RegisterMessage struct {
	Functions []string `json:"functions"` // 注册的函数列表
}

// CallMessage 调用消息
type CallMessage struct {
	Function string                 `json:"function"` // 函数名
	Params   map[string]interface{} `json:"params"`   // 参数
}

// ResultMessage 结果消息
type ResultMessage struct {
	Result interface{} `json:"result"` // 结果
}

// ErrorMessage 错误消息
type ErrorMessage struct {
	Error string `json:"error"` // 错误信息
}

// EncodeMessage 编码消息
func EncodeMessage(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}

// DecodeMessage 解码消息
func DecodeMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// PluginInfo 插件信息
type PluginInfo struct {
	Name      string   `json:"name"`      // 插件名称
	Version   string   `json:"version"`   // 插件版本
	Functions []string `json:"functions"` // 支持的函数
}
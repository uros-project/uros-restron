package actor

import (
	"encoding/json"
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	// FunctionCall 函数调用消息
	FunctionCall MessageType = "function_call"
	// FunctionResponse 函数响应消息
	FunctionResponse MessageType = "function_response"
	// StatusQuery 状态查询消息
	StatusQuery MessageType = "status_query"
	// StatusResponse 状态响应消息
	StatusResponse MessageType = "status_response"
	// Error 错误消息
	Error MessageType = "error"
	// Heartbeat 心跳消息
	Heartbeat MessageType = "heartbeat"
)

// Message Actor 消息结构
type Message struct {
	ID            string                 `json:"id"`             // 消息ID
	Type          MessageType            `json:"type"`           // 消息类型
	From          string                 `json:"from"`           // 发送者
	To            string                 `json:"to"`             // 接收者
	Function      string                 `json:"function"`       // 要调用的函数名
	Payload       map[string]interface{} `json:"payload"`        // 消息载荷
	Timestamp     time.Time              `json:"timestamp"`      // 时间戳
	CorrelationID string                 `json:"correlation_id"` // 关联ID，用于请求-响应匹配
}

// NewMessage 创建新消息
func NewMessage(msgType MessageType, from, to string) *Message {
	return &Message{
		ID:        generateMessageID(),
		Type:      msgType,
		From:      from,
		To:        to,
		Timestamp: time.Now(),
	}
}

// NewFunctionCallMessage 创建函数调用消息
func NewFunctionCallMessage(from, to, function string, params map[string]interface{}) *Message {
	msg := NewMessage(FunctionCall, from, to)
	msg.Function = function
	msg.Payload = map[string]interface{}{
		"function": function,
		"params":   params,
	}
	return msg
}

// NewFunctionResponseMessage 创建函数响应消息
func NewFunctionResponseMessage(from, to string, success bool, result map[string]interface{}, err string) *Message {
	msg := NewMessage(FunctionResponse, from, to)
	msg.Payload = map[string]interface{}{
		"success": success,
		"result":  result,
		"error":   err,
	}
	return msg
}

// NewStatusMessage 创建状态消息
func NewStatusMessage(from, to, status string, details map[string]interface{}) *Message {
	msg := NewMessage(StatusResponse, from, to)
	msg.Payload = map[string]interface{}{
		"status":    status,
		"details":   details,
		"timestamp": time.Now(),
	}
	return msg
}

// ToJSON 将消息转换为JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON 从JSON创建消息
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// SetCorrelationID 设置关联ID
func (m *Message) SetCorrelationID(id string) {
	m.CorrelationID = id
}

// generateMessageID 生成消息ID
func generateMessageID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

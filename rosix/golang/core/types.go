package core

import (
	"time"
)

// ResourceDescriptor 资源描述符 - 类似POSIX的文件描述符
type ResourceDescriptor int64

// ResourcePath 资源路径 - 类似文件系统路径
type ResourcePath string

// Resource 资源接口 - 所有资源的抽象
type Resource interface {
	// ID 返回资源的唯一标识
	ID() string

	// Path 返回资源路径
	Path() ResourcePath

	// Type 返回资源类型
	Type() ResourceType

	// Attributes 返回资源的静态属性
	Attributes() map[string]interface{}

	// Features 返回资源的动态属性
	Features() map[string]interface{}

	// Behaviors 返回资源支持的行为列表
	Behaviors() []string

	// Metadata 返回资源元数据
	Metadata() ResourceMetadata
}

// ResourceType 资源类型
type ResourceType string

const (
	// TypeDevice 设备类型资源
	TypeDevice ResourceType = "device"
	// TypeObject 对象类型资源
	TypeObject ResourceType = "object"
	// TypePerson 人员类型资源
	TypePerson ResourceType = "person"
	// TypeService 服务类型资源
	TypeService ResourceType = "service"
	// TypeActor Actor类型资源
	TypeActor ResourceType = "actor"
)

// ResourceMetadata 资源元数据
type ResourceMetadata struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Tags        []string               `json:"tags"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Owner       string                 `json:"owner"`
	Extra       map[string]interface{} `json:"extra"`
}

// ResourceState 资源状态
type ResourceState string

const (
	// StateActive 活跃状态
	StateActive ResourceState = "active"
	// StateIdle 空闲状态
	StateIdle ResourceState = "idle"
	// StateBusy 繁忙状态
	StateBusy ResourceState = "busy"
	// StateError 错误状态
	StateError ResourceState = "error"
	// StateOffline 离线状态
	StateOffline ResourceState = "offline"
)

// BehaviorDefinition 行为定义
type BehaviorDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Parameters  map[string]ParameterDefinition `json:"parameters"`
	Returns     map[string]ParameterDefinition `json:"returns"`
	Async       bool                           `json:"async"`
}

// ParameterDefinition 参数定义
type ParameterDefinition struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default"`
}

// ResourceHandle 资源句柄 - 打开的资源实例
type ResourceHandle struct {
	RD         ResourceDescriptor
	Resource   Resource
	Mode       OpenMode
	Context    *Context
	OpenedAt   time.Time
	LastAccess time.Time
}

// OpenMode 打开模式
type OpenMode int

const (
	// ModeRead 只读模式
	ModeRead OpenMode = 1 << iota
	// ModeWrite 写入模式
	ModeWrite
	// ModeInvoke 调用模式（执行行为）
	ModeInvoke
	// ModeWatch 监听模式
	ModeWatch
)

// Context 执行上下文
type Context struct {
	ID        string
	UserID    string
	SessionID string
	Metadata  map[string]interface{}
	Deadline  time.Time
	Cancel    chan struct{}
}

// Query 资源查询
type Query struct {
	Type       ResourceType           `json:"type,omitempty"`
	Category   string                 `json:"category,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	Features   map[string]interface{} `json:"features,omitempty"`
	Limit      int                    `json:"limit,omitempty"`
	Offset     int                    `json:"offset,omitempty"`
}

// Relation 资源关系
type Relation struct {
	Type        string                 `json:"type"`
	Source      ResourceDescriptor     `json:"source"`
	Target      ResourceDescriptor     `json:"target"`
	Metadata    map[string]interface{} `json:"metadata"`
	Bidirection bool                   `json:"bidirection"`
}

// Pipe 资源管道 - 用于资源间数据流
type Pipe struct {
	Source ResourceDescriptor
	Target ResourceDescriptor
	Filter func(interface{}) interface{}
}

// Event 资源事件
type Event struct {
	Type      EventType              `json:"type"`
	Resource  ResourceDescriptor     `json:"resource"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// EventType 事件类型
type EventType string

const (
	// EventStateChange 状态变化事件
	EventStateChange EventType = "state_change"
	// EventFeatureUpdate 特征更新事件
	EventFeatureUpdate EventType = "feature_update"
	// EventBehaviorInvoked 行为调用事件
	EventBehaviorInvoked EventType = "behavior_invoked"
	// EventError 错误事件
	EventError EventType = "error"
)

// WatchCallback 监听回调函数
type WatchCallback func(event Event) error

// Error 错误类型
type Error struct {
	Code    ErrorCode
	Message string
	Details map[string]interface{}
}

func (e *Error) Error() string {
	return e.Message
}

// ErrorCode 错误码
type ErrorCode int

const (
	// ErrNotFound 资源未找到
	ErrNotFound ErrorCode = 404
	// ErrPermissionDenied 权限拒绝
	ErrPermissionDenied ErrorCode = 403
	// ErrInvalidParameter 无效参数
	ErrInvalidParameter ErrorCode = 400
	// ErrResourceBusy 资源繁忙
	ErrResourceBusy ErrorCode = 409
	// ErrNotImplemented 未实现
	ErrNotImplemented ErrorCode = 501
	// ErrInternalError 内部错误
	ErrInternalError ErrorCode = 500
)

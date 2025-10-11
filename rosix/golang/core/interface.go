package core

// ROSIX 主接口 - 类似POSIX的系统调用接口
type ROSIX interface {
	// ==================== 资源操作原语 ====================

	// Open 打开资源，返回资源描述符
	// 类似 POSIX: open()
	Open(path ResourcePath, mode OpenMode, ctx *Context) (ResourceDescriptor, error)

	// Close 关闭资源
	// 类似 POSIX: close()
	Close(rd ResourceDescriptor) error

	// Read 读取资源属性或特征
	// 类似 POSIX: read()
	Read(rd ResourceDescriptor, key string) (interface{}, error)

	// Write 写入资源属性或特征
	// 类似 POSIX: write()
	Write(rd ResourceDescriptor, key string, value interface{}) error

	// Invoke 调用资源行为
	// 类似 POSIX: ioctl()
	Invoke(rd ResourceDescriptor, behavior string, params map[string]interface{}) (interface{}, error)

	// ==================== 资源发现和查询 ====================

	// Find 查找资源
	Find(query Query) ([]Resource, error)

	// List 列出指定路径下的资源
	// 类似 POSIX: readdir()
	List(path ResourcePath) ([]Resource, error)

	// Stat 获取资源信息
	// 类似 POSIX: stat()
	Stat(rd ResourceDescriptor) (Resource, error)

	// ==================== 资源监听 ====================

	// Watch 监听资源变化
	// 类似 Linux: inotify
	Watch(rd ResourceDescriptor, events []EventType, callback WatchCallback) error

	// Unwatch 取消监听
	Unwatch(rd ResourceDescriptor) error

	// ==================== 资源关系 ====================

	// Link 建立资源关系
	// 类似 POSIX: link()
	Link(source, target ResourceDescriptor, relationType string, metadata map[string]interface{}) error

	// Unlink 解除资源关系
	// 类似 POSIX: unlink()
	Unlink(source, target ResourceDescriptor) error

	// GetRelations 获取资源的关系
	GetRelations(rd ResourceDescriptor) ([]Relation, error)

	// ==================== 资源协同 ====================

	// Pipe 创建资源间的数据管道
	// 类似 POSIX: pipe()
	Pipe(source, target ResourceDescriptor, filter func(interface{}) interface{}) (*Pipe, error)

	// Fork 复制资源（创建资源实例）
	// 类似 POSIX: fork()
	Fork(rd ResourceDescriptor, params map[string]interface{}) (ResourceDescriptor, error)

	// ==================== 会话和上下文 ====================

	// CreateContext 创建执行上下文
	CreateContext(userID, sessionID string, metadata map[string]interface{}) (*Context, error)

	// DestroyContext 销毁执行上下文
	DestroyContext(ctx *Context) error

	// ==================== 批量操作 ====================

	// Batch 批量执行操作
	Batch(operations []Operation) ([]Result, error)

	// Transaction 事务性执行操作
	Transaction(operations []Operation) ([]Result, error)
}

// ResourceRegistry 资源注册表接口
type ResourceRegistry interface {
	// Register 注册资源
	Register(resource Resource) error

	// Unregister 注销资源
	Unregister(id string) error

	// Get 获取资源
	Get(id string) (Resource, error)

	// GetByPath 通过路径获取资源
	GetByPath(path ResourcePath) (Resource, error)

	// Query 查询资源
	Query(query Query) ([]Resource, error)

	// Watch 监听资源注册变化
	Watch(callback func(resource Resource, action string)) error
}

// BehaviorExecutor 行为执行器接口
type BehaviorExecutor interface {
	// Execute 执行行为
	Execute(resource Resource, behavior string, params map[string]interface{}) (interface{}, error)

	// ExecuteAsync 异步执行行为
	ExecuteAsync(resource Resource, behavior string, params map[string]interface{}) (chan interface{}, chan error)

	// GetDefinition 获取行为定义
	GetDefinition(resource Resource, behavior string) (*BehaviorDefinition, error)

	// ListBehaviors 列出资源的所有行为
	ListBehaviors(resource Resource) ([]BehaviorDefinition, error)
}

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(event Event) error

	// Subscribe 订阅事件
	Subscribe(eventType EventType, callback func(Event) error) (string, error)

	// Unsubscribe 取消订阅
	Unsubscribe(subscriptionID string) error

	// PublishAsync 异步发布事件
	PublishAsync(event Event) error
}

// Operation 操作定义（用于批量和事务操作）
type Operation struct {
	Type   OperationType          `json:"type"`
	Target ResourceDescriptor     `json:"target"`
	Data   map[string]interface{} `json:"data"`
}

// OperationType 操作类型
type OperationType string

const (
	OpRead   OperationType = "read"
	OpWrite  OperationType = "write"
	OpInvoke OperationType = "invoke"
	OpLink   OperationType = "link"
	OpUnlink OperationType = "unlink"
)

// Result 操作结果
type Result struct {
	Success bool                   `json:"success"`
	Data    interface{}            `json:"data"`
	Error   *Error                 `json:"error,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

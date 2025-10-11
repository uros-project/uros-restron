package syscall

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"uros-restron/internal/actor"
	"uros-restron/internal/models"
	"uros-restron/rosix/core"
	"uros-restron/rosix/resource"
)

// System ROSIX系统调用实现
type System struct {
	registry        core.ResourceRegistry
	eventBus        core.EventBus
	actorManager    *actor.ActorManager
	thingService    *models.ThingService
	behaviorService *models.BehaviorService

	// 资源句柄管理
	nextRD    int64
	handles   map[core.ResourceDescriptor]*core.ResourceHandle
	handlesMu sync.RWMutex

	// 监听管理
	watchers   map[core.ResourceDescriptor]*watcher
	watchersMu sync.RWMutex
}

type watcher struct {
	events   []core.EventType
	callback core.WatchCallback
	cancel   chan struct{}
}

// NewSystem 创建ROSIX系统实例
func NewSystem(
	actorManager *actor.ActorManager,
	thingService *models.ThingService,
	behaviorService *models.BehaviorService,
) *System {
	sys := &System{
		registry:        resource.NewRegistry(),
		actorManager:    actorManager,
		thingService:    thingService,
		behaviorService: behaviorService,
		nextRD:          1000, // 从1000开始分配资源描述符
		handles:         make(map[core.ResourceDescriptor]*core.ResourceHandle),
		watchers:        make(map[core.ResourceDescriptor]*watcher),
	}

	// 初始化时同步资源
	sys.syncResources()

	return sys
}

// syncResources 同步现有资源到注册表
func (s *System) syncResources() {
	// 同步Things
	if things, err := s.thingService.GetAllThings(); err == nil {
		for _, thing := range things {
			adapter := resource.NewThingAdapter(&thing)
			s.registry.Register(adapter)
		}
	}

	// 同步Actors
	actors := s.actorManager.ListActors()
	for _, act := range actors {
		status := act.GetStatus()
		actorID := fmt.Sprintf("%v", status["id"])
		name := fmt.Sprintf("%v", status["name"])
		statusStr := fmt.Sprintf("%v", status["status"])

		var functions []string
		if funcs, ok := status["functions"].([]string); ok {
			functions = funcs
		}

		adapter := resource.NewActorAdapter(actorID, name, statusStr, functions)
		s.registry.Register(adapter)
	}
}

// ==================== 资源操作原语 ====================

// Open 打开资源
func (s *System) Open(path core.ResourcePath, mode core.OpenMode, ctx *core.Context) (core.ResourceDescriptor, error) {
	// 从注册表查找资源
	res, err := s.registry.GetByPath(path)
	if err != nil {
		return 0, err
	}

	// 分配资源描述符
	rd := core.ResourceDescriptor(atomic.AddInt64(&s.nextRD, 1))

	// 创建资源句柄
	handle := &core.ResourceHandle{
		RD:         rd,
		Resource:   res,
		Mode:       mode,
		Context:    ctx,
		OpenedAt:   time.Now(),
		LastAccess: time.Now(),
	}

	// 保存句柄
	s.handlesMu.Lock()
	s.handles[rd] = handle
	s.handlesMu.Unlock()

	return rd, nil
}

// Close 关闭资源
func (s *System) Close(rd core.ResourceDescriptor) error {
	s.handlesMu.Lock()
	defer s.handlesMu.Unlock()

	if _, exists := s.handles[rd]; !exists {
		return &core.Error{
			Code:    core.ErrNotFound,
			Message: "invalid resource descriptor",
		}
	}

	// 取消监听
	s.Unwatch(rd)

	// 删除句柄
	delete(s.handles, rd)

	return nil
}

// Read 读取资源属性或特征
func (s *System) Read(rd core.ResourceDescriptor, key string) (interface{}, error) {
	handle, err := s.getHandle(rd)
	if err != nil {
		return nil, err
	}

	// 更新访问时间
	handle.LastAccess = time.Now()

	res := handle.Resource

	// 先从特征中查找
	features := res.Features()
	if value, ok := features[key]; ok {
		return value, nil
	}

	// 再从属性中查找
	attributes := res.Attributes()
	if value, ok := attributes[key]; ok {
		return value, nil
	}

	return nil, &core.Error{
		Code:    core.ErrNotFound,
		Message: fmt.Sprintf("key %s not found", key),
	}
}

// Write 写入资源属性或特征
func (s *System) Write(rd core.ResourceDescriptor, key string, value interface{}) error {
	handle, err := s.getHandle(rd)
	if err != nil {
		return err
	}

	// 检查权限
	if handle.Mode&core.ModeWrite == 0 {
		return &core.Error{
			Code:    core.ErrPermissionDenied,
			Message: "resource not opened for writing",
		}
	}

	// 更新访问时间
	handle.LastAccess = time.Now()

	// 根据资源类型执行写入
	res := handle.Resource

	// 如果是Thing，更新其状态
	if thingAdapter, ok := res.(*resource.ThingAdapter); ok {
		thing := thingAdapter.GetThing()
		if thing.Status == nil {
			thing.Status = make(map[string]interface{})
		}
		thing.Status[key] = value

		// 更新到数据库
		return s.thingService.UpdateThing(thing)
	}

	return &core.Error{
		Code:    core.ErrNotImplemented,
		Message: "write not supported for this resource type",
	}
}

// Invoke 调用资源行为
func (s *System) Invoke(rd core.ResourceDescriptor, behavior string, params map[string]interface{}) (interface{}, error) {
	handle, err := s.getHandle(rd)
	if err != nil {
		return nil, err
	}

	// 检查权限
	if handle.Mode&core.ModeInvoke == 0 {
		return nil, &core.Error{
			Code:    core.ErrPermissionDenied,
			Message: "resource not opened for invocation",
		}
	}

	// 更新访问时间
	handle.LastAccess = time.Now()

	res := handle.Resource

	// 根据资源类型执行调用
	switch res.Type() {
	case core.TypeActor:
		// 调用Actor函数
		return s.actorManager.CallFunction(res.ID(), behavior, params)

	case core.TypeDevice, core.TypeObject:
		// 通过Thing的关联Actor执行
		if thingAdapter, ok := res.(*resource.ThingAdapter); ok {
			thing := thingAdapter.GetThing()
			// 查找关联的Actor并调用
			for _, beh := range thing.Behaviors {
				if _, exists := beh.Functions[behavior]; exists {
					return s.actorManager.CallFunction(beh.ID, behavior, params)
				}
			}
		}
		return nil, &core.Error{
			Code:    core.ErrNotFound,
			Message: fmt.Sprintf("behavior %s not found", behavior),
		}

	default:
		return nil, &core.Error{
			Code:    core.ErrNotImplemented,
			Message: "invoke not supported for this resource type",
		}
	}
}

// ==================== 资源发现和查询 ====================

// Find 查找资源
func (s *System) Find(query core.Query) ([]core.Resource, error) {
	return s.registry.Query(query)
}

// List 列出指定路径下的资源
func (s *System) List(path core.ResourcePath) ([]core.Resource, error) {
	// 简化实现：列出所有资源
	return s.registry.Query(core.Query{}), nil
}

// Stat 获取资源信息
func (s *System) Stat(rd core.ResourceDescriptor) (core.Resource, error) {
	handle, err := s.getHandle(rd)
	if err != nil {
		return nil, err
	}
	return handle.Resource, nil
}

// ==================== 资源监听 ====================

// Watch 监听资源变化
func (s *System) Watch(rd core.ResourceDescriptor, events []core.EventType, callback core.WatchCallback) error {
	handle, err := s.getHandle(rd)
	if err != nil {
		return err
	}

	if handle.Mode&core.ModeWatch == 0 {
		return &core.Error{
			Code:    core.ErrPermissionDenied,
			Message: "resource not opened for watching",
		}
	}

	w := &watcher{
		events:   events,
		callback: callback,
		cancel:   make(chan struct{}),
	}

	s.watchersMu.Lock()
	s.watchers[rd] = w
	s.watchersMu.Unlock()

	return nil
}

// Unwatch 取消监听
func (s *System) Unwatch(rd core.ResourceDescriptor) error {
	s.watchersMu.Lock()
	defer s.watchersMu.Unlock()

	if w, exists := s.watchers[rd]; exists {
		close(w.cancel)
		delete(s.watchers, rd)
	}

	return nil
}

// ==================== 辅助方法 ====================

// getHandle 获取资源句柄
func (s *System) getHandle(rd core.ResourceDescriptor) (*core.ResourceHandle, error) {
	s.handlesMu.RLock()
	defer s.handlesMu.RUnlock()

	handle, exists := s.handles[rd]
	if !exists {
		return nil, &core.Error{
			Code:    core.ErrNotFound,
			Message: "invalid resource descriptor",
		}
	}

	return handle, nil
}

// ==================== 未实现的方法（待扩展）====================

func (s *System) Link(source, target core.ResourceDescriptor, relationType string, metadata map[string]interface{}) error {
	return &core.Error{Code: core.ErrNotImplemented, Message: "Link not implemented yet"}
}

func (s *System) Unlink(source, target core.ResourceDescriptor) error {
	return &core.Error{Code: core.ErrNotImplemented, Message: "Unlink not implemented yet"}
}

func (s *System) GetRelations(rd core.ResourceDescriptor) ([]core.Relation, error) {
	return nil, &core.Error{Code: core.ErrNotImplemented, Message: "GetRelations not implemented yet"}
}

func (s *System) Pipe(source, target core.ResourceDescriptor, filter func(interface{}) interface{}) (*core.Pipe, error) {
	return nil, &core.Error{Code: core.ErrNotImplemented, Message: "Pipe not implemented yet"}
}

func (s *System) Fork(rd core.ResourceDescriptor, params map[string]interface{}) (core.ResourceDescriptor, error) {
	return 0, &core.Error{Code: core.ErrNotImplemented, Message: "Fork not implemented yet"}
}

func (s *System) CreateContext(userID, sessionID string, metadata map[string]interface{}) (*core.Context, error) {
	return &core.Context{
		ID:        fmt.Sprintf("ctx_%d", time.Now().UnixNano()),
		UserID:    userID,
		SessionID: sessionID,
		Metadata:  metadata,
		Cancel:    make(chan struct{}),
	}, nil
}

func (s *System) DestroyContext(ctx *core.Context) error {
	close(ctx.Cancel)
	return nil
}

func (s *System) Batch(operations []core.Operation) ([]core.Result, error) {
	return nil, &core.Error{Code: core.ErrNotImplemented, Message: "Batch not implemented yet"}
}

func (s *System) Transaction(operations []core.Operation) ([]core.Result, error) {
	return nil, &core.Error{Code: core.ErrNotImplemented, Message: "Transaction not implemented yet"}
}

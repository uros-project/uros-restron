package actor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"uros-restron/internal/models"
)

// ActorManager Actor管理器
type ActorManager struct {
	actors          map[string]Actor
	mu              sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
	behaviorService *models.BehaviorService
}

// NewActorManager 创建Actor管理器
func NewActorManager(behaviorService *models.BehaviorService) *ActorManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ActorManager{
		actors:          make(map[string]Actor),
		ctx:             ctx,
		cancel:          cancel,
		behaviorService: behaviorService,
	}
}

// CreateActorFromBehavior 从Behavior创建Actor
func (am *ActorManager) CreateActorFromBehavior(behaviorID string) (Actor, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// 检查Actor是否已存在
	if actor, exists := am.actors[behaviorID]; exists {
		return actor, nil
	}

	// 获取Behavior
	behavior, err := am.behaviorService.GetBehavior(behaviorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get behavior %s: %v", behaviorID, err)
	}

	// 创建BehaviorActor
	actor := NewBehaviorActor(behavior)

	// 启动Actor
	if err := actor.Start(am.ctx); err != nil {
		return nil, fmt.Errorf("failed to start actor for behavior %s: %v", behaviorID, err)
	}

	// 注册Actor
	am.actors[behaviorID] = actor

	return actor, nil
}

// CreateActorFromBehaviorData 从Behavior数据直接创建Actor
func (am *ActorManager) CreateActorFromBehaviorData(behavior *models.Behavior) (Actor, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// 检查Actor是否已存在
	if actor, exists := am.actors[behavior.ID]; exists {
		return actor, nil
	}

	// 创建BehaviorActor
	actor := NewBehaviorActor(behavior)

	// 启动Actor
	if err := actor.Start(am.ctx); err != nil {
		return nil, fmt.Errorf("failed to start actor for behavior %s: %v", behavior.ID, err)
	}

	// 注册Actor
	am.actors[behavior.ID] = actor

	return actor, nil
}

// GetActor 获取Actor
func (am *ActorManager) GetActor(actorID string) (Actor, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	actor, exists := am.actors[actorID]
	if !exists {
		return nil, fmt.Errorf("actor %s not found", actorID)
	}

	return actor, nil
}

// ListActors 列出所有Actor
func (am *ActorManager) ListActors() []Actor {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var actors []Actor
	for _, actor := range am.actors {
		actors = append(actors, actor)
	}

	return actors
}

// StopActor 停止Actor
func (am *ActorManager) StopActor(actorID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	actor, exists := am.actors[actorID]
	if !exists {
		return fmt.Errorf("actor %s not found", actorID)
	}

	// 停止Actor
	if err := actor.Stop(); err != nil {
		return fmt.Errorf("failed to stop actor %s: %v", actorID, err)
	}

	// 从注册表中移除
	delete(am.actors, actorID)

	return nil
}

// StopAllActors 停止所有Actor
func (am *ActorManager) StopAllActors() error {
	am.mu.Lock()
	defer am.mu.Unlock()

	var errors []error
	for actorID, actor := range am.actors {
		if err := actor.Stop(); err != nil {
			errors = append(errors, fmt.Errorf("failed to stop actor %s: %v", actorID, err))
		}
	}

	// 清空注册表
	am.actors = make(map[string]Actor)

	if len(errors) > 0 {
		return fmt.Errorf("errors stopping actors: %v", errors)
	}

	return nil
}

// SendMessage 发送消息到Actor
func (am *ActorManager) SendMessage(actorID string, msg *Message) error {
	actor, err := am.GetActor(actorID)
	if err != nil {
		return err
	}

	return actor.Send(msg)
}

// CallFunction 调用Actor的函数
func (am *ActorManager) CallFunction(actorID, functionName string, params map[string]interface{}) (map[string]interface{}, error) {
	actor, err := am.GetActor(actorID)
	if err != nil {
		return nil, err
	}

	// 创建函数调用消息
	msg := NewFunctionCallMessage("manager", actorID, functionName, params)
	msg.SetCorrelationID(fmt.Sprintf("%d", time.Now().UnixNano()))

	// 发送消息
	if err := actor.Send(msg); err != nil {
		return nil, fmt.Errorf("failed to send message to actor %s: %v", actorID, err)
	}

	// 这里应该等待响应，但为了简化，我们直接调用函数
	if behaviorActor, ok := actor.(*BehaviorActor); ok {
		return behaviorActor.CallFunction(functionName, params)
	}

	return nil, fmt.Errorf("actor %s is not a BehaviorActor", actorID)
}

// GetActorStatus 获取Actor状态
func (am *ActorManager) GetActorStatus(actorID string) (ActorState, error) {
	actor, err := am.GetActor(actorID)
	if err != nil {
		return "", err
	}

	return actor.State(), nil
}

// GetActorInfo 获取Actor信息
func (am *ActorManager) GetActorInfo(actorID string) (map[string]interface{}, error) {
	actor, err := am.GetActor(actorID)
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"id":    actor.ID(),
		"state": actor.State(),
	}

	// 如果是BehaviorActor，添加更多信息
	if behaviorActor, ok := actor.(*BehaviorActor); ok {
		behavior := behaviorActor.GetBehavior()
		info["behavior_id"] = behavior.ID
		info["behavior_name"] = behavior.Name
		info["behavior_type"] = behavior.Type
		info["behavior_category"] = behavior.Category
		info["available_functions"] = behaviorActor.GetAvailableFunctions()
	}

	return info, nil
}

// StartHeartbeat 启动心跳监控
func (am *ActorManager) StartHeartbeat(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				am.sendHeartbeatToAllActors()
			case <-am.ctx.Done():
				return
			}
		}
	}()
}

// sendHeartbeatToAllActors 向所有Actor发送心跳
func (am *ActorManager) sendHeartbeatToAllActors() {
	am.mu.RLock()
	actors := make([]Actor, 0, len(am.actors))
	for _, actor := range am.actors {
		actors = append(actors, actor)
	}
	am.mu.RUnlock()

	for _, actor := range actors {
		heartbeatMsg := NewMessage(Heartbeat, "manager", actor.ID())
		heartbeatMsg.SetCorrelationID(fmt.Sprintf("heartbeat_%d", time.Now().UnixNano()))

		if err := actor.Send(heartbeatMsg); err != nil {
			fmt.Printf("Failed to send heartbeat to actor %s: %v\n", actor.ID(), err)
		}
	}
}

// Shutdown 关闭Actor管理器
func (am *ActorManager) Shutdown() error {
	// 停止所有Actor
	if err := am.StopAllActors(); err != nil {
		return err
	}

	// 取消上下文
	am.cancel()

	return nil
}

// GetActorCount 获取Actor数量
func (am *ActorManager) GetActorCount() int {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return len(am.actors)
}

// GetActorsByCategory 根据分类获取Actor
func (am *ActorManager) GetActorsByCategory(category string) []Actor {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var actors []Actor
	for _, actor := range am.actors {
		if behaviorActor, ok := actor.(*BehaviorActor); ok {
			if behaviorActor.GetBehavior().Category == category {
				actors = append(actors, actor)
			}
		}
	}

	return actors
}

// GetActorsByType 根据类型获取Actor
func (am *ActorManager) GetActorsByType(behaviorType string) []Actor {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var actors []Actor
	for _, actor := range am.actors {
		if behaviorActor, ok := actor.(*BehaviorActor); ok {
			if behaviorActor.GetBehavior().Type == models.BehaviorType(behaviorType) {
				actors = append(actors, actor)
			}
		}
	}

	return actors
}

// GetAllActorStatuses 获取所有Actor状态
func (am *ActorManager) GetAllActorStatuses() map[string]interface{} {
	am.mu.RLock()
	defer am.mu.RUnlock()

	statuses := make(map[string]interface{})
	for actorID, actor := range am.actors {
		statuses[actorID] = actor.GetStatus()
	}

	return statuses
}

// BroadcastMessage 广播消息到所有Actor
func (am *ActorManager) BroadcastMessage(msg *Message) {
	am.mu.RLock()
	actors := make([]Actor, 0, len(am.actors))
	for _, actor := range am.actors {
		actors = append(actors, actor)
	}
	am.mu.RUnlock()

	for _, actor := range actors {
		if err := actor.Send(msg); err != nil {
			fmt.Printf("Failed to broadcast message to actor %s: %v\n", actor.ID(), err)
		}
	}
}

// HealthCheck 健康检查
func (am *ActorManager) HealthCheck() map[string]interface{} {
	am.mu.RLock()
	defer am.mu.RUnlock()

	healthy := 0
	total := len(am.actors)

	for _, actor := range am.actors {
		if actor.State() == ActorStateRunning {
			healthy++
		}
	}

	return map[string]interface{}{
		"status":    "healthy",
		"total":     total,
		"healthy":   healthy,
		"unhealthy": total - healthy,
	}
}

// RegisterBehaviorsFromService 从服务注册所有行为为Actor
func (am *ActorManager) RegisterBehaviorsFromService(behaviorService *models.BehaviorService) error {
	// 获取所有行为
	behaviors, err := behaviorService.GetAllBehaviors()
	if err != nil {
		return fmt.Errorf("failed to get behaviors: %v", err)
	}

	// 为每个行为创建Actor
	for _, behavior := range behaviors {
		_, err := am.CreateActorFromBehaviorData(&behavior)
		if err != nil {
			fmt.Printf("Failed to create actor for behavior %s: %v\n", behavior.ID, err)
		}
	}

	return nil
}

// Start 启动Actor管理器
func (am *ActorManager) Start() error {
	// 启动心跳监控
	am.StartHeartbeat(30 * time.Second)
	return nil
}

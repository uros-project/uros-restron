package actor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"uros-restron/internal/models"
)

// FunctionDefinition 函数定义
type FunctionDefinition struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	InputParams  map[string]interface{} `json:"input_params"`
	OutputParams map[string]interface{} `json:"output_params"`
}

// BehaviorActor Behavior Actor 实现
type BehaviorActor struct {
	id          string
	name        string
	behavior    *models.Behavior
	Functions   map[string]FunctionHandler `json:"-"`
	MessageChan chan *Message              `json:"-"`
	Context     context.Context            `json:"-"`
	Cancel      context.CancelFunc         `json:"-"`
	Status      string                     `json:"status"`
	LastActive  time.Time                  `json:"last_active"`
	mu          sync.RWMutex               `json:"-"`
}

// FunctionHandler 函数处理器接口
type FunctionHandler interface {
	Execute(params map[string]interface{}) (map[string]interface{}, error)
	GetDefinition() FunctionDefinition
}

// NewBehaviorActor 创建新的 Behavior Actor
func NewBehaviorActor(behavior *models.Behavior) *BehaviorActor {
	ctx, cancel := context.WithCancel(context.Background())

	actor := &BehaviorActor{
		id:          behavior.ID,
		name:        behavior.Name,
		behavior:    behavior,
		Functions:   make(map[string]FunctionHandler),
		MessageChan: make(chan *Message, 100),
		Context:     ctx,
		Cancel:      cancel,
		Status:      "initializing",
		LastActive:  time.Now(),
	}

	// 注册函数处理器
	actor.registerFunctionHandlers()

	return actor
}

// ID 返回Actor ID
func (ba *BehaviorActor) ID() string {
	return ba.id
}

// State 返回Actor状态
func (ba *BehaviorActor) State() ActorState {
	ba.mu.RLock()
	defer ba.mu.RUnlock()

	switch ba.Status {
	case "running":
		return ActorStateRunning
	case "stopped":
		return ActorStateStopped
	case "error":
		return ActorStateError
	default:
		return ActorStateIdle
	}
}

// Send 发送消息到Actor
func (ba *BehaviorActor) Send(msg *Message) error {
	select {
	case ba.MessageChan <- msg:
		return nil
	case <-ba.Context.Done():
		return fmt.Errorf("actor %s is stopped", ba.id)
	default:
		return fmt.Errorf("actor %s message channel is full", ba.id)
	}
}

// SetMessageHandler 设置消息处理器（BehaviorActor不需要外部处理器）
func (ba *BehaviorActor) SetMessageHandler(handler MessageHandler) {
	// BehaviorActor使用内置的消息处理逻辑
}

// registerFunctionHandlers 注册函数处理器
func (ba *BehaviorActor) registerFunctionHandlers() {
	if ba.behavior.Functions == nil {
		return
	}

	// 遍历行为中的所有函数
	for funcName, funcData := range ba.behavior.Functions {
		// funcData 现在是 models.Function 类型，不需要类型断言

		// 创建函数处理器
		handler := &DefaultFunctionHandler{
			Name:           funcData.Name,
			Description:    funcData.Description,
			InputParams:    convertParametersToMap(funcData.InputParams),
			OutputParams:   convertParametersToMap(funcData.OutputParams),
			Implementation: convertImplementationToMap(funcData.Implementation),
		}

		ba.Functions[funcName] = handler
	}
}

// convertParametersToMap 将 Parameter 映射转换为 map[string]interface{}
func convertParametersToMap(params map[string]models.Parameter) map[string]interface{} {
	result := make(map[string]interface{})
	for name, param := range params {
		result[name] = map[string]interface{}{
			"type":        param.Type,
			"description": param.Description,
			"required":    param.Required,
			"default":     param.Default,
			"min":         param.Min,
			"max":         param.Max,
			"enum":        param.Enum,
		}
	}
	return result
}

// convertImplementationToMap 将 FunctionImplementation 转换为 map[string]interface{}
func convertImplementationToMap(impl models.FunctionImplementation) map[string]interface{} {
	steps := make([]map[string]interface{}, len(impl.Steps))
	for i, step := range impl.Steps {
		steps[i] = map[string]interface{}{
			"step":        step.Step,
			"action":      step.Action,
			"description": step.Description,
			"condition":   step.Condition,
		}
	}
	return map[string]interface{}{
		"steps": steps,
	}
}

// Start 启动 Actor
func (ba *BehaviorActor) Start(ctx context.Context) error {
	ba.mu.Lock()
	ba.Status = "running"
	ba.mu.Unlock()

	log.Printf("Behavior Actor %s (%s) started", ba.name, ba.id)

	go ba.messageLoop()
	return nil
}

// Stop 停止 Actor
func (ba *BehaviorActor) Stop() error {
	ba.mu.Lock()
	ba.Status = "stopped"
	ba.mu.Unlock()

	ba.Cancel()
	close(ba.MessageChan)

	log.Printf("Behavior Actor %s (%s) stopped", ba.name, ba.id)
	return nil
}

// SendMessage 发送消息到 Actor
func (ba *BehaviorActor) SendMessage(msg *Message) error {
	select {
	case ba.MessageChan <- msg:
		return nil
	case <-ba.Context.Done():
		return fmt.Errorf("actor %s is stopped", ba.id)
	default:
		return fmt.Errorf("actor %s message channel is full", ba.id)
	}
}

// messageLoop 消息处理循环
func (ba *BehaviorActor) messageLoop() {
	for {
		select {
		case msg := <-ba.MessageChan:
			ba.handleMessage(msg)
		case <-ba.Context.Done():
			return
		}
	}
}

// handleMessage 处理消息
func (ba *BehaviorActor) handleMessage(msg *Message) {
	ba.mu.Lock()
	ba.LastActive = time.Now()
	ba.mu.Unlock()

	log.Printf("Actor %s received message: %s", ba.id, msg.Type)

	switch msg.Type {
	case FunctionCall:
		ba.handleFunctionCall(msg)
	case StatusQuery:
		ba.handleStatusQuery(msg)
	case Heartbeat:
		ba.handleHeartbeat(msg)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// handleFunctionCall 处理函数调用
func (ba *BehaviorActor) handleFunctionCall(msg *Message) {
	funcName := msg.Function
	if funcName == "" {
		// 从 payload 中获取函数名
		if payload, ok := msg.Payload["function"].(string); ok {
			funcName = payload
		}
	}

	// 获取函数参数
	var params map[string]interface{}
	if payload, ok := msg.Payload["params"].(map[string]interface{}); ok {
		params = payload
	} else {
		params = make(map[string]interface{})
	}

	// 查找函数处理器
	handler, exists := ba.Functions[funcName]
	if !exists {
		response := NewFunctionResponseMessage(ba.id, msg.From, false, nil, fmt.Sprintf("function %s not found", funcName))
		response.CorrelationID = msg.CorrelationID
		// 这里应该发送响应，但为了简化，我们只记录日志
		log.Printf("Function %s not found in actor %s", funcName, ba.id)
		return
	}

	// 执行函数
	result, err := handler.Execute(params)

	// 创建响应消息
	success := err == nil
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	response := NewFunctionResponseMessage(ba.id, msg.From, success, result, errorMsg)
	response.CorrelationID = msg.CorrelationID

	log.Printf("Actor %s executed function %s: success=%v", ba.id, funcName, success)
}

// handleStatusQuery 处理状态查询
func (ba *BehaviorActor) handleStatusQuery(msg *Message) {
	ba.mu.RLock()
	status := ba.Status
	lastActive := ba.LastActive
	ba.mu.RUnlock()

	details := map[string]interface{}{
		"id":          ba.id,
		"name":        ba.name,
		"status":      status,
		"last_active": lastActive,
		"functions":   ba.getAvailableFunctions(),
	}

	response := NewStatusMessage(ba.id, msg.From, status, details)
	response.CorrelationID = msg.CorrelationID

	log.Printf("Actor %s status: %s", ba.id, status)
}

// handleHeartbeat 处理心跳消息
func (ba *BehaviorActor) handleHeartbeat(msg *Message) {
	ba.mu.Lock()
	ba.LastActive = time.Now()
	ba.mu.Unlock()

	// 可以在这里添加心跳响应逻辑
	// 目前只是更新最后活跃时间
}

// getAvailableFunctions 获取可用函数列表
func (ba *BehaviorActor) getAvailableFunctions() []string {
	var functions []string
	for funcName := range ba.Functions {
		functions = append(functions, funcName)
	}
	return functions
}

// GetStatus 获取 Actor 状态
func (ba *BehaviorActor) GetStatus() map[string]interface{} {
	ba.mu.RLock()
	defer ba.mu.RUnlock()

	return map[string]interface{}{
		"id":          ba.id,
		"name":        ba.name,
		"status":      ba.Status,
		"last_active": ba.LastActive,
		"functions":   ba.getAvailableFunctions(),
	}
}

// GetBehavior 获取行为信息
func (ba *BehaviorActor) GetBehavior() *models.Behavior {
	return ba.behavior
}

// GetAvailableFunctions 获取可用函数列表（公开方法）
func (ba *BehaviorActor) GetAvailableFunctions() []string {
	return ba.getAvailableFunctions()
}

// CallFunction 调用函数
func (ba *BehaviorActor) CallFunction(functionName string, params map[string]interface{}) (map[string]interface{}, error) {
	handler, exists := ba.Functions[functionName]
	if !exists {
		return nil, fmt.Errorf("function %s not found", functionName)
	}

	return handler.Execute(params)
}

// GetFunctionInfo 获取函数信息
func (ba *BehaviorActor) GetFunctionInfo(functionName string) (FunctionDefinition, error) {
	handler, exists := ba.Functions[functionName]
	if !exists {
		return FunctionDefinition{}, fmt.Errorf("function %s not found", functionName)
	}

	return handler.GetDefinition(), nil
}

// 辅助函数
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getMap(data map[string]interface{}, key string) map[string]interface{} {
	if val, ok := data[key].(map[string]interface{}); ok {
		return val
	}
	return make(map[string]interface{})
}

package actor

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ActorState Actor状态
type ActorState string

const (
	// ActorStateIdle 空闲状态
	ActorStateIdle ActorState = "idle"
	// ActorStateRunning 运行状态
	ActorStateRunning ActorState = "running"
	// ActorStateStopped 停止状态
	ActorStateStopped ActorState = "stopped"
	// ActorStateError 错误状态
	ActorStateError ActorState = "error"
)

// Actor 定义Actor接口
type Actor interface {
	// ID 返回Actor的唯一标识
	ID() string

	// Start 启动Actor
	Start(ctx context.Context) error

	// Stop 停止Actor
	Stop() error

	// Send 发送消息到Actor
	Send(msg *Message) error

	// State 返回Actor当前状态
	State() ActorState

	// SetMessageHandler 设置消息处理器
	SetMessageHandler(handler MessageHandler)

	// GetStatus 获取Actor状态信息
	GetStatus() map[string]interface{}
}

// MessageHandler 消息处理器接口
type MessageHandler interface {
	// HandleMessage 处理消息
	HandleMessage(ctx context.Context, msg *Message) (*Message, error)

	// CanHandle 检查是否能处理指定类型的消息
	CanHandle(msgType MessageType) bool
}

// BaseActor 基础Actor实现
type BaseActor struct {
	id             string
	state          ActorState
	messageChan    chan *Message
	messageHandler MessageHandler
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.RWMutex
	wg             sync.WaitGroup
}

// NewBaseActor 创建基础Actor
func NewBaseActor(id string) *BaseActor {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseActor{
		id:          id,
		state:       ActorStateIdle,
		messageChan: make(chan *Message, 100), // 缓冲100个消息
		ctx:         ctx,
		cancel:      cancel,
	}
}

// ID 返回Actor ID
func (a *BaseActor) ID() string {
	return a.id
}

// Start 启动Actor
func (a *BaseActor) Start(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.state != ActorStateIdle {
		return fmt.Errorf("actor %s is not in idle state", a.id)
	}

	a.state = ActorStateRunning
	a.wg.Add(1)

	go a.messageLoop()

	return nil
}

// Stop 停止Actor
func (a *BaseActor) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.state == ActorStateStopped {
		return nil
	}

	a.state = ActorStateStopped
	a.cancel()

	// 等待消息循环结束
	a.wg.Wait()

	// 关闭消息通道
	close(a.messageChan)

	return nil
}

// Send 发送消息到Actor
func (a *BaseActor) Send(msg *Message) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.state != ActorStateRunning {
		return fmt.Errorf("actor %s is not running", a.id)
	}

	select {
	case a.messageChan <- msg:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout sending message to actor %s", a.id)
	}
}

// State 返回Actor状态
func (a *BaseActor) State() ActorState {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state
}

// SetMessageHandler 设置消息处理器
func (a *BaseActor) SetMessageHandler(handler MessageHandler) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.messageHandler = handler
}

// GetStatus 获取Actor状态信息
func (a *BaseActor) GetStatus() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return map[string]interface{}{
		"id":    a.id,
		"state": a.state,
	}
}

// messageLoop 消息处理循环
func (a *BaseActor) messageLoop() {
	defer a.wg.Done()

	for {
		select {
		case msg, ok := <-a.messageChan:
			if !ok {
				return // 通道已关闭
			}
			a.handleMessage(msg)

		case <-a.ctx.Done():
			return // 上下文取消
		}
	}
}

// handleMessage 处理单个消息
func (a *BaseActor) handleMessage(msg *Message) {
	if a.messageHandler == nil {
		// 如果没有设置消息处理器，记录错误
		fmt.Printf("Actor %s received message but no handler set\n", a.id)
		return
	}

	// 检查是否能处理此消息类型
	if !a.messageHandler.CanHandle(msg.Type) {
		fmt.Printf("Actor %s cannot handle message type %s\n", a.id, msg.Type)
		return
	}

	// 处理消息
	response, err := a.messageHandler.HandleMessage(a.ctx, msg)
	if err != nil {
		fmt.Printf("Error handling message in actor %s: %v\n", a.id, err)
		// 发送错误响应
		errorMsg := NewMessage(Error, a.id, msg.From)
		errorMsg.CorrelationID = msg.CorrelationID
		errorMsg.Payload["error"] = err.Error()
		// 这里可以添加错误响应发送逻辑
		return
	}

	// 如果有响应消息，发送响应
	if response != nil {
		// 这里可以添加响应发送逻辑
		fmt.Printf("Actor %s processed message and generated response\n", a.id)
	}
}

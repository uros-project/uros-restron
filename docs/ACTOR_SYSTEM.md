# UROS Actor系统

## 概述

UROS Actor系统是一个基于Go的Actor模型实现，允许将Behavior编译为可执行的Actor，接受外部消息并驱动不同的function执行。

## 核心特性

- **Actor模型**: 每个Behavior可以编译为一个独立的Actor
- **消息驱动**: 通过消息传递机制驱动函数执行
- **异步处理**: 支持并发消息处理
- **函数执行**: 根据消息类型执行不同的function
- **状态管理**: 完整的Actor生命周期管理
- **类型安全**: 强类型的消息和参数验证

## 架构组件

### 1. 消息系统 (Message System)

```go
// 消息类型
type MessageType string

const (
    FunctionCall     MessageType = "function_call"
    FunctionResponse MessageType = "function_response"
    StatusUpdate     MessageType = "status_update"
    Heartbeat        MessageType = "heartbeat"
    Error            MessageType = "error"
)

// 消息结构
type Message struct {
    ID            string                 `json:"id"`
    Type          MessageType           `json:"type"`
    From          string                 `json:"from"`
    To            string                 `json:"to"`
    Function      string                 `json:"function"`
    Payload       map[string]interface{} `json:"payload"`
    Timestamp     time.Time              `json:"timestamp"`
    CorrelationID string                `json:"correlation_id"`
}
```

### 2. Actor接口 (Actor Interface)

```go
type Actor interface {
    ID() string
    Start(ctx context.Context) error
    Stop() error
    Send(msg *Message) error
    State() ActorState
    SetMessageHandler(handler MessageHandler)
}
```

### 3. BehaviorActor实现

```go
type BehaviorActor struct {
    *BaseActor
    behavior *models.Behavior
    executor *FunctionExecutor
}
```

### 4. 函数执行器 (Function Executor)

```go
type FunctionExecutor struct {
    behavior  *models.Behavior
    functions map[string]interface{}
}
```

## 使用方法

### 1. 创建BehaviorActor

```go
// 从Behavior创建Actor
behaviorActor := actor.NewBehaviorActor(behavior)

// 启动Actor
ctx := context.Background()
if err := behaviorActor.Start(ctx); err != nil {
    log.Fatalf("启动Actor失败: %v", err)
}
```

### 2. 调用函数

```go
// 直接调用函数
params := map[string]interface{}{
    "air_quality":    150.0,
    "target_quality": 50.0,
}

result, err := behaviorActor.CallFunction("purify_air", params)
if err != nil {
    log.Printf("函数调用失败: %v", err)
} else {
    fmt.Printf("函数调用成功: %+v\n", result)
}
```

### 3. 发送消息

```go
// 创建函数调用消息
msg := actor.NewFunctionCallMessage("sender", "actor-id", "function-name", params)
msg.SetCorrelationID("unique-id")

// 发送消息
if err := behaviorActor.Send(msg); err != nil {
    log.Printf("发送消息失败: %v", err)
}
```

### 4. 使用Actor管理器

```go
// 创建Actor管理器
actorManager := actor.NewActorManager(behaviorService)

// 从Behavior创建Actor
actor, err := actorManager.CreateActorFromBehavior("behavior-id")
if err != nil {
    log.Printf("创建Actor失败: %v", err)
}

// 调用Actor函数
result, err := actorManager.CallFunction("actor-id", "function-name", params)
if err != nil {
    log.Printf("调用函数失败: %v", err)
}
```

## API接口

### HTTP API

系统提供了完整的HTTP API接口来管理Actor：

#### 创建Actor
```http
POST /api/v1/actors
Content-Type: application/json

{
    "behavior_id": "purifier-001"
}
```

#### 获取Actor信息
```http
GET /api/v1/actors/info?actor_id=purifier-001
```

#### 列出所有Actor
```http
GET /api/v1/actors
```

#### 调用Actor函数
```http
POST /api/v1/actors/call?actor_id=purifier-001
Content-Type: application/json

{
    "function": "purify_air",
    "params": {
        "air_quality": 150.0,
        "target_quality": 50.0
    }
}
```

#### 获取Actor状态
```http
GET /api/v1/actors/status?actor_id=purifier-001
```

#### 停止Actor
```http
DELETE /api/v1/actors?actor_id=purifier-001
```

## 函数定义格式

Behavior中的函数定义需要遵循以下格式：

```json
{
    "function_name": {
        "name": "函数名称",
        "description": "函数描述",
        "input_params": {
            "param_name": {
                "type": "string|number|boolean|object",
                "description": "参数描述",
                "required": true,
                "min": 0,
                "max": 100,
                "enum": ["value1", "value2"]
            }
        },
        "output_params": {
            "result_name": {
                "type": "string|number|boolean|object",
                "description": "结果描述"
            }
        },
        "implementation": {
            "steps": [
                {
                    "step": 1,
                    "action": "action_name",
                    "description": "步骤描述",
                    "condition": "condition_expression"
                }
            ]
        }
    }
}
```

## 支持的动作类型

系统内置支持以下动作类型：

- `check_air_quality`: 检查空气质量
- `start_fan`: 启动风扇
- `activate_filter`: 激活过滤器
- `monitor_progress`: 监控进度
- `read_filter_sensor`: 读取过滤器传感器
- `calculate_usage`: 计算使用量
- `determine_status`: 确定状态
- `validate_speed`: 验证速度
- `set_motor_speed`: 设置电机速度
- `confirm_speed`: 确认速度
- `parse_input`: 解析输入
- `understand_intent`: 理解意图
- `generate_response`: 生成响应
- `format_output`: 格式化输出
- `validate_credentials`: 验证凭据
- `check_permissions`: 检查权限
- `generate_session`: 生成会话
- `load_user_profile`: 加载用户档案
- `extract_preferences`: 提取偏好
- `format_preferences`: 格式化偏好

## 示例代码

完整的示例代码请参考 `examples/actor_example.go` 文件。

## 注意事项

1. **并发安全**: Actor系统是并发安全的，支持多个Actor同时运行
2. **消息缓冲**: 每个Actor有100个消息的缓冲队列
3. **超时处理**: 消息发送有5秒超时限制
4. **错误处理**: 所有操作都有完整的错误处理机制
5. **资源管理**: 需要正确调用Stop()方法来释放资源

## 扩展开发

要添加新的动作类型，需要在 `FunctionExecutor` 中添加对应的执行方法：

```go
func (fe *FunctionExecutor) executeNewAction(params map[string]interface{}) (map[string]interface{}, error) {
    // 实现新的动作逻辑
    return map[string]interface{}{
        "result": "success",
        "timestamp": time.Now().Unix(),
    }, nil
}
```

然后在 `executeStep` 方法中添加对应的case：

```go
case "new_action":
    return fe.executeNewAction(params)
```

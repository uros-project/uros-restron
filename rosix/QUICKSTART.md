# ROSIX 快速开始指南

## 什么是ROSIX？

ROSIX（Resource Operating System Interface eXtension）是一个面向资源的编程层，为资源管理系统提供类似POSIX的标准化接口。它将底层的Things、Actors、Behaviors等概念统一抽象为"资源"，并提供统一的操作接口。

## 核心概念

### 1. 资源 (Resource)
系统中可管理的实体，包括：
- 设备（传感器、执行器）
- 对象（容器、物品）
- Actor（行为实例）
- 服务（通知、计算）

### 2. 资源描述符 (Resource Descriptor)
类似文件描述符，用于引用打开的资源。通过RD进行资源操作。

### 3. 资源路径 (Resource Path)
类似文件路径的资源标识：
```
/actors/{id}
/things/{type}/{id}
/devices/{category}/{id}
```

### 4. 系统调用 (System Calls)
统一的资源操作接口：
- `Open()` - 打开资源
- `Close()` - 关闭资源  
- `Read()` - 读取属性
- `Write()` - 写入属性
- `Invoke()` - 调用行为
- `Find()` - 查找资源

## 30秒快速体验

### 1. 通过HTTP API

```bash
# 查找资源
curl -X POST http://localhost:8080/api/v1/rosix/resources/find \
  -H "Content-Type: application/json" \
  -d '{"type": "actor", "limit": 5}'

# 调用资源
curl -X POST http://localhost:8080/api/v1/rosix/resources/invoke \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/actors/your-actor-id",
    "behavior": "purify_air",
    "params": {"mode": "auto", "intensity": 3}
  }'

# AI驱动调用
curl -X POST http://localhost:8080/api/v1/rosix/ai/invoke \
  -H "Content-Type: application/json" \
  -d '{"prompt": "打开空气净化器"}'
```

### 2. 通过Go代码

```go
package main

import (
    "log"
    "uros-restron/rosix/core"
    "uros-restron/rosix/syscall"
)

func main() {
    // 创建ROSIX系统
    rosix := syscall.NewSystem(actorManager, thingService, behaviorService)
    
    // 创建上下文
    ctx, _ := rosix.CreateContext("user_001", "session_123", nil)
    defer rosix.DestroyContext(ctx)
    
    // 查找资源
    resources, _ := rosix.Find(core.Query{
        Type:     core.TypeActor,
        Category: "purifier",
        Limit:    1,
    })
    
    // 打开资源
    rd, _ := rosix.Open(resources[0].Path(), core.ModeInvoke, ctx)
    defer rosix.Close(rd)
    
    // 调用行为
    result, _ := rosix.Invoke(rd, "purify_air", map[string]interface{}{
        "mode": "auto",
        "intensity": 3,
    })
    
    log.Printf("结果: %v", result)
}
```

### 3. 通过CLI工具

```bash
# 查找资源
rosix find --type actor --category purifier

# 读取属性
rosix read /actors/abc123 status

# 调用行为
rosix invoke /actors/abc123 purify_air --params '{"mode":"auto"}'

# AI调用
rosix ai invoke "打开空气净化器"

# AI编排
rosix ai orchestrate "进入睡眠模式"
```

## 5分钟深入了解

### 完整的资源操作流程

```go
// 1. 初始化系统
rosix := syscall.NewSystem(actorManager, thingService, behaviorService)

// 2. 创建执行上下文
ctx, err := rosix.CreateContext("user_001", "session_123", map[string]interface{}{
    "location": "客厅",
    "device": "mobile",
})
defer rosix.DestroyContext(ctx)

// 3. 查找资源
resources, err := rosix.Find(core.Query{
    Type:     core.TypeActor,
    Category: "purifier",
    Tags:     []string{"indoor"},
    Limit:    5,
})

if len(resources) == 0 {
    log.Fatal("未找到资源")
}

fmt.Printf("找到 %d 个资源\n", len(resources))

// 4. 打开资源（支持多种模式）
rd, err := rosix.Open(
    resources[0].Path(),
    core.ModeRead | core.ModeWrite | core.ModeInvoke,
    ctx,
)
defer rosix.Close(rd)

// 5. 读取资源信息
status, _ := rosix.Read(rd, "status")
name, _ := rosix.Read(rd, "name")
fmt.Printf("资源: %s, 状态: %v\n", name, status)

// 6. 写入资源属性
err = rosix.Write(rd, "mode", "auto")

// 7. 调用资源行为
result, err := rosix.Invoke(rd, "purify_air", map[string]interface{}{
    "mode":      "auto",
    "intensity": 3,
    "duration":  3600,
})

fmt.Printf("调用结果: %v\n", result)
```

### AI驱动的资源管理

```go
// 创建AI编排器
orchestrator := ai.NewSimpleOrchestrator(rosix)

// 场景1：自然语言调用
result, err := orchestrator.Invoke("打开客厅的空气净化器", ctx)
fmt.Printf("AI识别意图: %s\n", result.Intent)
fmt.Printf("执行结果: %s\n", result.Message)

// 场景2：多资源编排
plan, err := orchestrator.Orchestrate("进入睡眠模式", ctx)
fmt.Printf("生成计划，共%d个步骤:\n", len(plan.Steps))
for _, step := range plan.Steps {
    fmt.Printf("  %d. %s\n", step.Order, step.Description)
}

// 场景3：信息查询
answer, err := orchestrator.Query("客厅的温度是多少？", ctx)
fmt.Printf("回答: %s\n", answer.Answer)

// 场景4：获取建议
suggestion, err := orchestrator.Suggest("如何改善空气质量？", ctx)
for _, item := range suggestion.Suggestions {
    fmt.Printf("建议: %s\n", item.Title)
}
```

## 常见使用场景

### 场景1：设备监控

```go
// 打开传感器
sensor, _ := rosix.Open("/devices/sensor/temp_001", core.ModeRead, ctx)
defer rosix.Close(sensor)

// 定期读取数据
ticker := time.NewTicker(5 * time.Second)
for range ticker.C {
    temp, _ := rosix.Read(sensor, "temperature")
    humidity, _ := rosix.Read(sensor, "humidity")
    fmt.Printf("温度: %.1f°C, 湿度: %.1f%%\n", temp, humidity)
}
```

### 场景2：设备控制

```go
// 打开空调
ac, _ := rosix.Open("/devices/ac/living_room", core.ModeInvoke, ctx)
defer rosix.Close(ac)

// 设置温度
rosix.Invoke(ac, "set_temperature", map[string]interface{}{
    "temperature": 26,
    "mode": "cool",
})

// 设置风速
rosix.Invoke(ac, "set_fan_speed", map[string]interface{}{
    "speed": "medium",
})
```

### 场景3：批量操作

```go
// 查找所有灯光设备
lights, _ := rosix.Find(core.Query{
    Category: "light",
})

// 批量关闭
for _, light := range lights {
    rd, _ := rosix.Open(light.Path(), core.ModeInvoke, ctx)
    rosix.Invoke(rd, "turn_off", nil)
    rosix.Close(rd)
}
```

### 场景4：资源监听

```go
// 打开资源（监听模式）
rd, _ := rosix.Open("/devices/sensor/motion_001", core.ModeWatch, ctx)
defer rosix.Close(rd)

// 设置监听回调
callback := func(event core.Event) error {
    fmt.Printf("检测到事件: %s\n", event.Type)
    fmt.Printf("数据: %v\n", event.Data)
    
    // 根据事件触发动作
    if event.Type == core.EventStateChange {
        // 执行相应操作
    }
    return nil
}

// 开始监听
rosix.Watch(rd, []core.EventType{
    core.EventStateChange,
    core.EventFeatureUpdate,
}, callback)
```

## API端点速查

### 资源操作
- `POST /api/v1/rosix/resources/find` - 查找资源
- `POST /api/v1/rosix/resources/read` - 读取资源
- `POST /api/v1/rosix/resources/write` - 写入资源
- `POST /api/v1/rosix/resources/invoke` - 调用资源

### AI接口
- `POST /api/v1/rosix/ai/invoke` - AI调用
- `POST /api/v1/rosix/ai/orchestrate` - AI编排
- `POST /api/v1/rosix/ai/query` - AI查询
- `POST /api/v1/rosix/ai/suggest` - AI建议

### 系统信息
- `GET /api/v1/rosix/info` - 系统信息

## 与POSIX的对比

| POSIX | ROSIX | 说明 |
|-------|-------|------|
| `open()` | `Open()` | 打开文件/资源 |
| `close()` | `Close()` | 关闭文件/资源 |
| `read()` | `Read()` | 读取数据/属性 |
| `write()` | `Write()` | 写入数据/属性 |
| `ioctl()` | `Invoke()` | 设备控制/行为调用 |
| `stat()` | `Stat()` | 获取信息 |
| `readdir()` | `List()` | 列出内容 |
| 文件描述符 | 资源描述符 | 资源引用 |
| 文件路径 | 资源路径 | 资源标识 |

## 下一步

- 查看 [架构设计](ARCHITECTURE.md) 了解详细设计
- 查看 [集成指南](INTEGRATION.md) 了解如何集成到现有系统
- 查看 [API文档](API.md) 了解完整的API说明
- 查看 [示例代码](examples/) 学习更多用法

## 常见问题

**Q: ROSIX与直接调用Actor有什么区别？**
A: ROSIX提供了统一的抽象层，让你可以用相同的方式操作不同类型的资源，而不需要关心底层实现。同时提供了AI驱动、资源发现、监听等高级功能。

**Q: 性能开销如何？**
A: ROSIX层的开销很小，主要是资源查找和句柄管理。实际的行为执行直接调用底层Actor，没有额外开销。

**Q: 如何扩展新的资源类型？**
A: 只需实现Resource接口并创建适配器，然后注册到Registry即可。参见[扩展指南](ARCHITECTURE.md#扩展性设计)。

**Q: AI功能需要外部模型吗？**
A: 当前提供的SimpleOrchestrator是基于规则的简单实现。你可以集成真实的AI模型（如GPT、Claude）来增强功能。

**Q: 支持分布式部署吗？**
A: 当前版本是单机版本，分布式支持在规划中。未来将支持分布式资源注册表和跨节点资源访问。


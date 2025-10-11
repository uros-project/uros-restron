# ROSIX 集成指南

## 概述

本文档说明如何将ROSIX编程层集成到现有的UROS Restron系统中。

## 集成步骤

### 1. 在main.go中初始化ROSIX系统

```go
package main

import (
	"log"
	"uros-restron/internal/actor"
	"uros-restron/internal/api"
	"uros-restron/internal/config"
	"uros-restron/internal/database"
	"uros-restron/internal/models"
	"uros-restron/internal/utils"
	"uros-restron/rosix/ai"
	rosixapi "uros-restron/rosix/api"
	"uros-restron/rosix/syscall"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.InitDB(cfg.Database.DSN)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 运行数据库迁移
	migrationUtils := utils.NewMigrationUtils(db)
	if err := migrationUtils.RunMigrations(&models.Thing{}, &models.ThingType{}, &models.Relationship{}, &models.Behavior{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 创建索引和种子数据
	if err := migrationUtils.CreateIndexes(); err != nil {
		log.Fatal("Failed to create indexes:", err)
	}
	if err := migrationUtils.SeedData(); err != nil {
		log.Fatal("Failed to seed data:", err)
	}

	// 初始化服务
	thingService := models.NewThingService(db)
	thingTypeService := models.NewThingTypeService(db)
	relationshipService := models.NewRelationshipService(db)
	behaviorService := models.NewBehaviorService(db)
	actorManager := actor.NewActorManager(behaviorService)
	hub := api.NewHub()

	// 启动 Actor 管理器
	actorManager.Start()

	// 注册行为到 Actor 系统
	if err := actorManager.RegisterBehaviorsFromService(behaviorService); err != nil {
		log.Printf("Warning: Failed to register behaviors: %v", err)
	}

	// 填充预定义行为
	if err := behaviorService.SeedPredefinedBehaviors(); err != nil {
		log.Printf("Warning: Failed to seed behaviors: %v", err)
	}

	// ========== 初始化ROSIX系统 ==========
	
	// 创建ROSIX系统实例
	rosixSystem := syscall.NewSystem(actorManager, thingService, behaviorService)
	log.Println("ROSIX system initialized")

	// 创建AI编排器
	orchestrator := ai.NewSimpleOrchestrator(rosixSystem)
	log.Println("AI orchestrator initialized")

	// 创建ROSIX API处理器
	rosixHandler := rosixapi.NewROSIXHandler(rosixSystem, orchestrator)

	// =====================================

	// 启动 WebSocket 服务
	go hub.Run()

	// 启动 HTTP 服务器
	server := api.NewServer(cfg, thingService, thingTypeService, relationshipService, behaviorService, actorManager, hub)
	
	// 在服务器中添加ROSIX路由
	// 注意：需要修改 api.Server 以支持添加额外的路由
	// server.AddROSIXRoutes(rosixHandler)

	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
```

### 2. 修改 internal/api/server.go 以支持ROSIX路由

在 `server.go` 的 `setupRoutes()` 方法中添加：

```go
// 在 setupRoutes() 方法中添加
import rosixapi "uros-restron/rosix/api"

// 在 API 路由组中添加
func (s *Server) setupRoutes() {
	// ... 现有代码 ...

	// API 路由组
	api := s.router.Group("/api/v1")
	{
		// ... 现有路由 ...

		// ROSIX 路由（如果已初始化）
		if s.rosixHandler != nil {
			rosixapi.SetupROSIXRoutes(api, s.rosixHandler)
		}
	}
}
```

### 3. 更新 Server 结构体

```go
type Server struct {
	config              *config.Config
	thingService        *models.ThingService
	thingTypeService    *models.ThingTypeService
	relationshipService *models.RelationshipService
	behaviorService     *models.BehaviorService
	actorManager        *actor.ActorManager
	hub                 *Hub
	router              *gin.Engine
	
	// 添加ROSIX支持
	rosixHandler        *rosixapi.ROSIXHandler
}

// 添加设置ROSIX处理器的方法
func (s *Server) SetROSIXHandler(handler *rosixapi.ROSIXHandler) {
	s.rosixHandler = handler
}
```

## API端点

集成完成后，以下ROSIX端点将可用：

### 系统信息
- `GET /api/v1/rosix/info` - 获取ROSIX系统信息

### 资源操作
- `POST /api/v1/rosix/resources/find` - 查找资源
- `POST /api/v1/rosix/resources/read` - 读取资源属性
- `POST /api/v1/rosix/resources/write` - 写入资源属性
- `POST /api/v1/rosix/resources/invoke` - 调用资源行为

### AI接口
- `POST /api/v1/rosix/ai/invoke` - AI驱动的资源调用
- `POST /api/v1/rosix/ai/orchestrate` - AI编排多资源协同
- `POST /api/v1/rosix/ai/query` - AI查询资源信息
- `POST /api/v1/rosix/ai/suggest` - AI提供建议

## 使用示例

### 1. 查找资源

```bash
curl -X POST http://localhost:8080/api/v1/rosix/resources/find \
  -H "Content-Type: application/json" \
  -d '{
    "type": "actor",
    "category": "purifier",
    "limit": 5
  }'
```

### 2. 调用资源行为

```bash
curl -X POST http://localhost:8080/api/v1/rosix/resources/invoke \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/actors/d6c3adb7-2071-47a1-8abd-abfe2987ca6e",
    "behavior": "purify_air",
    "params": {
      "mode": "auto",
      "intensity": 3
    },
    "user_id": "user_001",
    "session": "session_123"
  }'
```

### 3. AI驱动调用

```bash
curl -X POST http://localhost:8080/api/v1/rosix/ai/invoke \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "打开空气净化器",
    "user_id": "user_001",
    "session": "session_123"
  }'
```

### 4. AI编排

```bash
curl -X POST http://localhost:8080/api/v1/rosix/ai/orchestrate \
  -H "Content-Type: application/json" \
  -d '{
    "goal": "进入睡眠模式",
    "user_id": "user_001",
    "session": "session_123"
  }'
```

### 5. AI查询

```bash
curl -X POST http://localhost:8080/api/v1/rosix/ai/query \
  -H "Content-Type: application/json" \
  -d '{
    "question": "客厅的温度是多少？",
    "user_id": "user_001",
    "session": "session_123"
  }'
```

## 编程接口使用

### Go代码示例

```go
package main

import (
	"log"
	"uros-restron/rosix/core"
	"uros-restron/rosix/syscall"
	"uros-restron/rosix/ai"
)

func main() {
	// 假设已经有了 actorManager, thingService, behaviorService
	
	// 创建ROSIX系统
	rosix := syscall.NewSystem(actorManager, thingService, behaviorService)
	
	// 创建上下文
	ctx, _ := rosix.CreateContext("user_001", "session_123", nil)
	defer rosix.DestroyContext(ctx)
	
	// 查找资源
	resources, err := rosix.Find(core.Query{
		Type:     core.TypeActor,
		Category: "purifier",
		Limit:    1,
	})
	
	if err != nil || len(resources) == 0 {
		log.Fatal("No resources found")
	}
	
	// 打开资源
	rd, err := rosix.Open(resources[0].Path(), 
		core.ModeRead|core.ModeInvoke, ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rosix.Close(rd)
	
	// 读取资源状态
	status, _ := rosix.Read(rd, "status")
	log.Printf("Resource status: %v", status)
	
	// 调用资源行为
	result, err := rosix.Invoke(rd, "purify_air", map[string]interface{}{
		"mode":      "auto",
		"intensity": 3,
	})
	
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Invoke result: %v", result)
	
	// 使用AI编排器
	orchestrator := ai.NewSimpleOrchestrator(rosix)
	aiResult, err := orchestrator.Invoke("启动空气净化", ctx)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("AI result: %s", aiResult.Message)
}
```

## 优势

1. **统一接口**: 所有资源使用相同的操作接口
2. **简化开发**: 开发者无需关心底层实现细节
3. **AI原生**: 内置AI驱动的资源管理能力
4. **可扩展**: 易于添加新的资源类型和行为
5. **标准化**: 类似POSIX的设计理念，易于理解和使用

## 后续扩展

1. 完善资源关系管理（Link/Unlink）
2. 实现资源管道（Pipe）
3. 添加资源监听机制（Watch）
4. 支持批量和事务操作
5. 集成真实的AI模型（如接入GPT/Claude等）
6. 添加资源权限控制
7. 实现资源缓存机制
8. 支持分布式资源管理


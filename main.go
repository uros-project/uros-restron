package main

import (
	"log"
	"uros-restron/internal/actor"
	"uros-restron/internal/api"
	"uros-restron/internal/config"
	"uros-restron/internal/database"
	"uros-restron/internal/models"
	"uros-restron/internal/utils"
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

	// 创建索引
	if err := migrationUtils.CreateIndexes(); err != nil {
		log.Fatal("Failed to create indexes:", err)
	}

	// 种子数据
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

	// 启动 WebSocket 服务
	go hub.Run()

	// 启动 HTTP 服务器
	server := api.NewServer(cfg, thingService, thingTypeService, relationshipService, behaviorService, actorManager, hub)

	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

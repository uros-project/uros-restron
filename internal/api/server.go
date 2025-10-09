package api

import (
	"net/http"
	"uros-restron/internal/actor"
	"uros-restron/internal/config"
	"uros-restron/internal/models"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config              *config.Config
	thingService        *models.ThingService
	thingTypeService    *models.ThingTypeService
	relationshipService *models.RelationshipService
	behaviorService     *models.BehaviorService
	actorManager        *actor.ActorManager
	hub                 *Hub
	router              *gin.Engine
}

func NewServer(cfg *config.Config, thingService *models.ThingService, thingTypeService *models.ThingTypeService, relationshipService *models.RelationshipService, behaviorService *models.BehaviorService, actorManager *actor.ActorManager, hub *Hub) *Server {
	server := &Server{
		config:              cfg,
		thingService:        thingService,
		thingTypeService:    thingTypeService,
		relationshipService: relationshipService,
		behaviorService:     behaviorService,
		actorManager:        actorManager,
		hub:                 hub,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	gin.SetMode(gin.DebugMode)
	s.router = gin.New()
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.router.Use(corsMiddleware())

	// 静态文件服务
	s.router.Static("/static", "./web/static")
	s.router.LoadHTMLGlob("./web/templates/*")

	// 主页面
	s.router.GET("/", s.indexPage)
	s.router.GET("/types", s.typesPage)
	s.router.GET("/things", s.thingsPage)
	s.router.GET("/behaviors", s.behaviorsPage)
	s.router.GET("/relationships", s.relationshipsPage)
	s.router.GET("/actors", s.actorsPage)
	s.router.GET("/graph", s.graphPage)

	// API 路由组
	api := s.router.Group("/api/v1")
	{
		// 数字孪生相关路由 - 使用独立的处理器
		thingHandler := NewThingHandler(s.thingService, s.relationshipService, s.behaviorService, s.hub)
		SetupThingRoutes(api, thingHandler)

		// 事物类型相关路由 - 使用独立的处理器
		thingTypeHandler := NewThingTypeHandler(s.thingTypeService, s.thingService, s.hub)
		SetupThingTypeRoutes(api, thingTypeHandler)

		// 关系管理相关路由 - 使用独立的处理器
		relationshipHandler := NewRelationshipHandler(s.relationshipService, s.hub)
		SetupRelationshipRoutes(api, relationshipHandler)

		// 行为管理相关路由 - 使用独立的处理器
		behaviorHandler := NewBehaviorHandler(s.behaviorService, s.thingTypeService, s.thingService, s.hub)
		SetupBehaviorRoutes(api, behaviorHandler)

		// Actor系统相关路由 - 使用独立的处理器
		actorHandler := NewActorHandler(s.actorManager, s.hub)
		SetupActorRoutes(api, actorHandler)

		// WebSocket 路由
		api.GET("/ws", s.handleWebSocket)

	}

	// 健康检查
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func (s *Server) Start() error {
	addr := s.config.Server.Host + ":" + s.config.Server.Port
	return s.router.Run(addr)
}

// CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

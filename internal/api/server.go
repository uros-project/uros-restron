package api

import (
	"net/http"
	"uros-restron/internal/config"
	"uros-restron/internal/models"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config              *config.Config
	thingService        *models.ThingService
	thingTypeService    *models.ThingTypeService
	relationshipService *models.RelationshipService
	hub                 *Hub
	router              *gin.Engine
}

func NewServer(cfg *config.Config, thingService *models.ThingService, thingTypeService *models.ThingTypeService, relationshipService *models.RelationshipService, hub *Hub) *Server {
	server := &Server{
		config:              cfg,
		thingService:        thingService,
		thingTypeService:    thingTypeService,
		relationshipService: relationshipService,
		hub:                 hub,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	gin.SetMode(gin.ReleaseMode)
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
	s.router.GET("/graph", s.graphPage)

	// API 路由组
	api := s.router.Group("/api/v1")
	{
		// 数字孪生相关路由
		api.GET("/things", s.listThings)
		api.POST("/things", s.createThing)
		api.GET("/things/:id", s.getThing)
		api.PUT("/things/:id", s.updateThing)
		api.DELETE("/things/:id", s.deleteThing)

		// 属性相关路由
		api.PUT("/things/:id/properties/:name", s.updateProperty)
		api.PUT("/things/:id/status", s.updateStatus)

		// 事物类型相关路由
		api.GET("/thing-types", s.listThingTypes)
		api.POST("/thing-types", s.createThingType)
		api.GET("/thing-types/:id", s.getThingType)
		api.PUT("/thing-types/:id", s.updateThingType)
		api.DELETE("/thing-types/:id", s.deleteThingType)

		// 根据类型创建事物
		api.POST("/thing-types/:id/things", s.createThingFromType)

		// 关系管理相关路由
		api.GET("/relationships", s.listRelationships)
		api.POST("/relationships", s.createRelationship)
		api.GET("/relationships/:id", s.getRelationship)
		api.PUT("/relationships/:id", s.updateRelationship)
		api.DELETE("/relationships/:id", s.deleteRelationship)
		api.GET("/relationships/types", s.getRelationshipTypes)

		// 事物关系路由
		api.GET("/things/:id/relationships", s.getThingRelationships)

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

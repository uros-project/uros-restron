package api

import (
	"github.com/gin-gonic/gin"
)

// SetupActorRoutes 设置Actor系统相关的路由
func SetupActorRoutes(router *gin.RouterGroup, handler *ActorHandler) {
	// Actor 管理
	router.GET("/actors", handler.ListActors)
	router.GET("/actors/:id", handler.GetActor)

	// Actor 函数调用
	router.POST("/actors/:id/functions/:function", handler.CallActorFunction)
	router.POST("/actors/:id/messages", handler.SendMessageToActor)
	router.GET("/actors/:id/functions", handler.GetActorFunctions)
	router.GET("/actors/health", handler.HealthCheck)
}

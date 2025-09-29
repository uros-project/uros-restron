package api

import (
	"github.com/gin-gonic/gin"
)

// SetupThingRoutes 设置数字孪生相关的路由
func SetupThingRoutes(router *gin.RouterGroup, handler *ThingHandler) {
	// 数字孪生 CRUD 操作
	router.GET("/things", handler.ListThings)
	router.POST("/things", handler.CreateThing)
	router.GET("/things/:id", handler.GetThing)
	router.PUT("/things/:id", handler.UpdateThing)
	router.DELETE("/things/:id", handler.DeleteThing)

	// 状态更新
	router.PUT("/things/:id/status", handler.UpdateStatus)

	// 事物关系路由
	router.GET("/things/:id/relationships", handler.GetThingRelationships)

	// 事物行为路由
	router.GET("/things/:id/behaviors", handler.GetThingBehaviors)
	router.POST("/things/:id/behaviors", handler.AssignBehaviorToThing)
	router.DELETE("/things/:id/behaviors/:behaviorId", handler.RemoveBehaviorFromThing)
}

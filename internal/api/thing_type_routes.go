package api

import (
	"github.com/gin-gonic/gin"
)

// SetupThingTypeRoutes 设置事物类型相关的路由
func SetupThingTypeRoutes(router *gin.RouterGroup, handler *ThingTypeHandler) {
	// 事物类型 CRUD 操作
	router.GET("/thing-types", handler.ListThingTypes)
	router.POST("/thing-types", handler.CreateThingType)
	router.GET("/thing-types/:id", handler.GetThingType)
	router.PUT("/thing-types/:id", handler.UpdateThingType)
	router.DELETE("/thing-types/:id", handler.DeleteThingType)

	// 根据类型创建事物实例
	router.POST("/thing-types/:id/things", handler.CreateThingFromType)
}

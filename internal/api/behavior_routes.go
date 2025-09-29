package api

import (
	"github.com/gin-gonic/gin"
)

// SetupBehaviorRoutes 设置行为相关的路由
func SetupBehaviorRoutes(router *gin.RouterGroup, handler *BehaviorHandler) {
	// 行为 CRUD 操作
	router.GET("/behaviors", handler.ListBehaviors)
	router.POST("/behaviors", handler.CreateBehavior)
	router.GET("/behaviors/:id", handler.GetBehavior)
	router.PUT("/behaviors/:id", handler.UpdateBehavior)
	router.DELETE("/behaviors/:id", handler.DeleteBehavior)

	// 行为分类和预定义
	router.GET("/behaviors/category/:category", handler.GetBehaviorsByCategory)
	router.GET("/behaviors/predefined", handler.GetPredefinedBehaviors)
	router.POST("/behaviors/seed", handler.SeedBehaviors)

	// 事物类型行为路由
	router.GET("/thing-types/:id/behaviors", handler.GetThingTypeBehaviors)
	router.POST("/thing-types/:id/behaviors", handler.AssignBehaviorToThingType)
	router.DELETE("/thing-types/:id/behaviors/:behaviorId", handler.RemoveBehaviorFromThingType)
}

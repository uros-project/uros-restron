package api

import (
	"github.com/gin-gonic/gin"
)

// SetupROSIXRoutes 设置ROSIX路由
func SetupROSIXRoutes(router *gin.RouterGroup, handler *ROSIXHandler) {
	// 系统信息
	router.GET("/rosix/info", handler.GetSystemInfo)

	// 资源操作
	router.POST("/rosix/resources/find", handler.FindResources)
	router.POST("/rosix/resources/read", handler.ReadResource)
	router.POST("/rosix/resources/write", handler.WriteResource)
	router.POST("/rosix/resources/invoke", handler.InvokeResource)

	// AI接口
	router.POST("/rosix/ai/invoke", handler.AIInvoke)
	router.POST("/rosix/ai/orchestrate", handler.AIOrchestrate)
	router.POST("/rosix/ai/query", handler.AIQuery)
	router.POST("/rosix/ai/suggest", handler.AISuggest)
}

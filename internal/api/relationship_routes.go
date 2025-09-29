package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRelationshipRoutes 设置关系相关的路由
func SetupRelationshipRoutes(router *gin.RouterGroup, handler *RelationshipHandler) {
	// 关系 CRUD 操作
	router.GET("/relationships", handler.ListRelationships)
	router.POST("/relationships", handler.CreateRelationship)
	router.GET("/relationships/:id", handler.GetRelationship)
	router.PUT("/relationships/:id", handler.UpdateRelationship)
	router.DELETE("/relationships/:id", handler.DeleteRelationship)

	// 关系类型
	router.GET("/relationships/types", handler.GetRelationshipTypes)
}

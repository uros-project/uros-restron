package api

import (
	"net/http"
	"strconv"
	"uros-restron/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 关系管理相关处理器

// listRelationships 获取关系列表
func (s *Server) listRelationships(c *gin.Context) {
	sourceID := c.Query("source_id")
	targetID := c.Query("target_id")
	relationshipType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	relationships, err := s.relationshipService.ListRelationships(sourceID, targetID, models.RelationshipType(relationshipType), limit, offset)
	if err != nil {
		logrus.Error("Failed to list relationships:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list relationships"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  relationships,
		"count": len(relationships),
	})
}

// createRelationship 创建关系
func (s *Server) createRelationship(c *gin.Context) {
	var relationship models.Relationship
	if err := c.ShouldBindJSON(&relationship); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.relationshipService.CreateRelationship(&relationship); err != nil {
		logrus.Error("Failed to create relationship:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create relationship"})
		return
	}

	c.JSON(http.StatusCreated, relationship)
}

// getRelationship 获取单个关系
func (s *Server) getRelationship(c *gin.Context) {
	id := c.Param("id")
	relationship, err := s.relationshipService.GetRelationship(id)
	if err != nil {
		logrus.Error("Failed to get relationship:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relationship"})
		return
	}

	if relationship == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relationship not found"})
		return
	}

	c.JSON(http.StatusOK, relationship)
}

// updateRelationship 更新关系
func (s *Server) updateRelationship(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.relationshipService.UpdateRelationship(id, updates); err != nil {
		logrus.Error("Failed to update relationship:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update relationship"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Relationship updated successfully"})
}

// deleteRelationship 删除关系
func (s *Server) deleteRelationship(c *gin.Context) {
	id := c.Param("id")
	if err := s.relationshipService.DeleteRelationship(id); err != nil {
		logrus.Error("Failed to delete relationship:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete relationship"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Relationship deleted successfully"})
}

// getThingRelationships 获取事物的所有关系
func (s *Server) getThingRelationships(c *gin.Context) {
	thingID := c.Param("id")
	relationships, err := s.relationshipService.GetThingRelationships(thingID)
	if err != nil {
		logrus.Error("Failed to get thing relationships:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get thing relationships"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  relationships,
		"count": len(relationships),
	})
}

// getRelationshipTypes 获取关系类型列表
func (s *Server) getRelationshipTypes(c *gin.Context) {
	types := s.relationshipService.GetRelationshipTypes()
	c.JSON(http.StatusOK, gin.H{
		"data":  types,
		"count": len(types),
	})
}

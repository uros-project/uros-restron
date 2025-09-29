package api

import (
	"net/http"
	"strconv"
	"uros-restron/internal/models"
	"uros-restron/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RelationshipHandler 关系处理器
type RelationshipHandler struct {
	relationshipService *models.RelationshipService
	hub                 *Hub
}

// NewRelationshipHandler 创建新的关系处理器
func NewRelationshipHandler(relationshipService *models.RelationshipService, hub *Hub) *RelationshipHandler {
	return &RelationshipHandler{
		relationshipService: relationshipService,
		hub:                 hub,
	}
}

// ListRelationships 获取关系列表
func (h *RelationshipHandler) ListRelationships(c *gin.Context) {
	relationshipType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid offset parameter")
		return
	}

	relationships, err := h.relationshipService.ListRelationships(relationshipType, "", "", limit, offset)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list relationships")
		return
	}

	utils.RespondWithData(c, gin.H{
		"data":  relationships,
		"count": len(relationships),
	})
}

// CreateRelationship 创建关系
func (h *RelationshipHandler) CreateRelationship(c *gin.Context) {
	var relationship models.Relationship
	if err := c.ShouldBindJSON(&relationship); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.relationshipService.CreateRelationship(&relationship); err != nil {
		logrus.Error("Failed to create relationship:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create relationship")
		return
	}

	utils.RespondWithDataStatus(c, relationship, http.StatusCreated)
}

// GetRelationship 获取单个关系
func (h *RelationshipHandler) GetRelationship(c *gin.Context) {
	id := c.Param("id")

	relationship, err := h.relationshipService.GetRelationship(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Relationship not found")
		return
	}

	utils.RespondWithData(c, relationship)
}

// UpdateRelationship 更新关系
func (h *RelationshipHandler) UpdateRelationship(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.relationshipService.UpdateRelationship(id, updates); err != nil {
		logrus.Error("Failed to update relationship:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update relationship")
		return
	}

	// 获取更新后的数据
	relationship, err := h.relationshipService.GetRelationship(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get updated relationship")
		return
	}

	utils.RespondWithData(c, relationship)
}

// DeleteRelationship 删除关系
func (h *RelationshipHandler) DeleteRelationship(c *gin.Context) {
	id := c.Param("id")

	if err := h.relationshipService.DeleteRelationship(id); err != nil {
		logrus.Error("Failed to delete relationship:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete relationship")
		return
	}

	utils.RespondWithData(c, gin.H{"message": "Relationship deleted successfully"})
}

// GetRelationshipTypes 获取关系类型
func (h *RelationshipHandler) GetRelationshipTypes(c *gin.Context) {
	types := h.relationshipService.GetRelationshipTypes()
	utils.RespondWithData(c, types)
}
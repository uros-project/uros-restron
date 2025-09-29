package api

import (
	"net/http"
	"strconv"
	"uros-restron/internal/models"
	"uros-restron/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ThingHandler 数字孪生处理器
type ThingHandler struct {
	thingService        *models.ThingService
	relationshipService *models.RelationshipService
	behaviorService     *models.BehaviorService
	hub                 *Hub
}

// NewThingHandler 创建新的数字孪生处理器
func NewThingHandler(thingService *models.ThingService, relationshipService *models.RelationshipService, behaviorService *models.BehaviorService, hub *Hub) *ThingHandler {
	return &ThingHandler{
		thingService:        thingService,
		relationshipService: relationshipService,
		behaviorService:     behaviorService,
		hub:                 hub,
	}
}

// ListThings 获取数字孪生列表
func (h *ThingHandler) ListThings(c *gin.Context) {
	thingType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "10")
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

	things, err := h.thingService.ListThings(thingType, limit, offset)
	if err != nil {
		utils.HandleError(c, err, "Failed to list things")
		return
	}

	utils.SuccessResponseWithCount(c, things, len(things))
}

// CreateThing 创建数字孪生
func (h *ThingHandler) CreateThing(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// 手动构建 Thing
	thing := &models.Thing{
		Name:        request["name"].(string),
		Type:        request["type"].(string),
		Description: request["description"].(string),
	}

	// 处理 Attributes
	if attrs, ok := request["attributes"].(map[string]interface{}); ok {
		thing.Attributes = attrs
	}

	// 处理 Features
	if feats, ok := request["features"].(map[string]interface{}); ok {
		thing.Features = feats
	}

	// 处理 BehaviorID
	if behaviorID, ok := request["behaviorId"].(string); ok {
		thing.BehaviorID = behaviorID
	}

	if err := h.thingService.CreateThing(thing); err != nil {
		logrus.Error("Failed to create thing:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create thing")
		return
	}

	// 广播新事物创建事件
	h.hub.Broadcast("thing_created", thing)

	utils.RespondWithDataStatus(c, thing, http.StatusCreated)
}

// GetThing 获取单个数字孪生
func (h *ThingHandler) GetThing(c *gin.Context) {
	id := c.Param("id")

	thing, err := h.thingService.GetThing(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Thing not found")
		return
	}

	utils.RespondWithData(c, thing)
}

// UpdateThing 更新数字孪生
func (h *ThingHandler) UpdateThing(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// 处理 BehaviorID 字段名映射
	if behaviorID, ok := updates["behaviorId"].(string); ok {
		updates["behavior_id"] = behaviorID
		delete(updates, "behaviorId")
	}

	if err := h.thingService.UpdateThing(id, updates); err != nil {
		logrus.Error("Failed to update thing:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update thing")
		return
	}

	// 获取更新后的数据
	thing, err := h.thingService.GetThing(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get updated thing")
		return
	}

	// 广播更新事件
	h.hub.Broadcast("thing_updated", thing)

	utils.RespondWithData(c, thing)
}

// DeleteThing 删除数字孪生
func (h *ThingHandler) DeleteThing(c *gin.Context) {
	id := c.Param("id")

	if err := h.thingService.DeleteThing(id); err != nil {
		logrus.Error("Failed to delete thing:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete thing")
		return
	}

	// 广播删除事件
	h.hub.Broadcast("thing_deleted", gin.H{"id": id})

	utils.RespondWithData(c, gin.H{"message": "Thing deleted successfully"})
}

// UpdateStatus 更新状态
func (h *ThingHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var status map[string]interface{}
	if err := c.ShouldBindJSON(&status); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.thingService.UpdateStatus(id, status); err != nil {
		logrus.Error("Failed to update status:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update status")
		return
	}

	// 获取更新后的数据
	thing, err := h.thingService.GetThing(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get updated thing")
		return
	}

	// 广播状态更新事件
	h.hub.Broadcast("status_updated", gin.H{
		"thingId": id,
		"status":  status,
		"thing":   thing,
	})

	utils.RespondWithData(c, gin.H{"message": "Status updated successfully"})
}

// GetThingRelationships 获取事物的关系
func (h *ThingHandler) GetThingRelationships(c *gin.Context) {
	thingID := c.Param("id")

	relationships, err := h.relationshipService.GetThingRelationships(thingID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get thing relationships")
		return
	}

	utils.RespondWithData(c, relationships)
}

// GetThingBehaviors 获取事物的行为
func (h *ThingHandler) GetThingBehaviors(c *gin.Context) {
	thingID := c.Param("id")

	behavior, err := h.thingService.GetThingBehavior(thingID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Thing behavior not found")
		return
	}

	utils.RespondWithData(c, behavior)
}

// AssignBehaviorToThing 为事物分配行为
func (h *ThingHandler) AssignBehaviorToThing(c *gin.Context) {
	thingID := c.Param("id")

	var request struct {
		BehaviorID string `json:"behaviorId"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.thingService.SetBehavior(thingID, request.BehaviorID); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to assign behavior to thing")
		return
	}

	utils.RespondWithData(c, gin.H{"message": "Behavior assigned successfully"})
}

// RemoveBehaviorFromThing 从事物中移除行为
func (h *ThingHandler) RemoveBehaviorFromThing(c *gin.Context) {
	thingID := c.Param("id")

	if err := h.thingService.RemoveBehavior(thingID); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to remove behavior from thing")
		return
	}

	utils.RespondWithData(c, gin.H{"message": "Behavior removed successfully"})
}
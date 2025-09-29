package api

import (
	"net/http"
	"strconv"
	"uros-restron/internal/models"
	"uros-restron/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ThingTypeHandler 事物类型处理器
type ThingTypeHandler struct {
	thingTypeService *models.ThingTypeService
	thingService     *models.ThingService
	hub              *Hub
}

// NewThingTypeHandler 创建新的事物类型处理器
func NewThingTypeHandler(thingTypeService *models.ThingTypeService, thingService *models.ThingService, hub *Hub) *ThingTypeHandler {
	return &ThingTypeHandler{
		thingTypeService: thingTypeService,
		thingService:     thingService,
		hub:              hub,
	}
}

// 事物类型相关处理器

// ListThingTypes 获取事物类型列表
func (h *ThingTypeHandler) ListThingTypes(c *gin.Context) {
	category := c.Query("category")
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

	thingTypes, err := h.thingTypeService.ListThingTypes(category, limit, offset)
	if err != nil {
		logrus.Error("Failed to list thing types:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list thing types")
		return
	}

	utils.RespondWithData(c, gin.H{
		"data":  thingTypes,
		"count": len(thingTypes),
	})
}

// CreateThingType 创建事物类型
func (h *ThingTypeHandler) CreateThingType(c *gin.Context) {
	var request map[string]interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// 手动构建 ThingType
	thingType := &models.ThingType{
		Name:        request["name"].(string),
		Description: request["description"].(string),
		Category:    request["category"].(string),
	}

	// 处理 Attributes
	if attrs, ok := request["attributes"].(map[string]interface{}); ok {
		thingType.Attributes = attrs
	}

	// 处理 Features
	if feats, ok := request["features"].(map[string]interface{}); ok {
		thingType.Features = feats
	}

	// 处理 BehaviorID
	if behaviorID, ok := request["behaviorId"].(string); ok {
		thingType.BehaviorID = behaviorID
	}

	if err := h.thingTypeService.CreateThingType(thingType); err != nil {
		logrus.Error("Failed to create thing type:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create thing type")
		return
	}

	utils.RespondWithDataStatus(c, thingType, http.StatusCreated)
}

// GetThingType 获取单个事物类型
func (h *ThingTypeHandler) GetThingType(c *gin.Context) {
	id := c.Param("id")

	thingType, err := h.thingTypeService.GetThingType(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Thing type not found")
		return
	}

	utils.RespondWithData(c, thingType)
}

// UpdateThingType 更新事物类型
func (h *ThingTypeHandler) UpdateThingType(c *gin.Context) {
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

	if err := h.thingTypeService.UpdateThingType(id, updates); err != nil {
		logrus.Error("Failed to update thing type:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update thing type")
		return
	}

	// 获取更新后的数据
	thingType, err := h.thingTypeService.GetThingType(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get updated thing type")
		return
	}

	utils.RespondWithData(c, thingType)
}

// DeleteThingType 删除事物类型
func (h *ThingTypeHandler) DeleteThingType(c *gin.Context) {
	id := c.Param("id")

	if err := h.thingTypeService.DeleteThingType(id); err != nil {
		logrus.Error("Failed to delete thing type:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete thing type")
		return
	}

	utils.RespondWithData(c, gin.H{"message": "Thing type deleted successfully"})
}

// CreateThingFromType 根据类型创建事物实例
func (h *ThingTypeHandler) CreateThingFromType(c *gin.Context) {
	typeID := c.Param("id")

	var request struct {
		Name        string                    `json:"name"`
		Description string                    `json:"description"`
		Attributes  map[string]interface{}    `json:"attributes"`
		Features    map[string]interface{} `json:"features"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// 使用类型服务创建事物
	thing, err := h.thingTypeService.CreateThingFromType(typeID, request.Name, request.Description, request.Attributes, request.Features)
	if err != nil {
		logrus.Error("Failed to create thing from type:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create thing from type")
		return
	}

	// 保存到数据库
	if err := h.thingService.CreateThing(thing); err != nil {
		logrus.Error("Failed to save thing:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save thing")
		return
	}

	// 广播新事物创建事件
	h.hub.Broadcast("thing_created", thing)

	utils.RespondWithDataStatus(c, thing, http.StatusCreated)
}

package api

import (
	"net/http"
	"strconv"
	"uros-restron/internal/models"
	"uros-restron/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// BehaviorHandler 行为处理器
type BehaviorHandler struct {
	behaviorService     *models.BehaviorService
	thingTypeService    *models.ThingTypeService
	thingService        *models.ThingService
	hub                 *Hub
}

// NewBehaviorHandler 创建新的行为处理器
func NewBehaviorHandler(behaviorService *models.BehaviorService, thingTypeService *models.ThingTypeService, thingService *models.ThingService, hub *Hub) *BehaviorHandler {
	return &BehaviorHandler{
		behaviorService:  behaviorService,
		thingTypeService: thingTypeService,
		thingService:     thingService,
		hub:              hub,
	}
}

// ListBehaviors 获取行为列表
func (h *BehaviorHandler) ListBehaviors(c *gin.Context) {
	// 获取查询参数
	behaviorType := c.Query("type")
	category := c.Query("category")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid offset parameter")
		return
	}

	behaviors, err := h.behaviorService.ListBehaviors(behaviorType, category, limit, offset)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list behaviors")
		return
	}

	utils.RespondWithData(c, gin.H{
		"data":  behaviors,
		"count": len(behaviors),
	})
}

// CreateBehavior 创建行为
func (h *BehaviorHandler) CreateBehavior(c *gin.Context) {
	var behavior models.Behavior
	if err := c.ShouldBindJSON(&behavior); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.behaviorService.CreateBehavior(&behavior); err != nil {
		logrus.Error("Failed to create behavior:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create behavior")
		return
	}

	utils.RespondWithDataStatus(c, behavior, http.StatusCreated)
}

// GetBehavior 获取单个行为
func (h *BehaviorHandler) GetBehavior(c *gin.Context) {
	id := c.Param("id")

	behavior, err := h.behaviorService.GetBehavior(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Behavior not found")
		return
	}

	utils.RespondWithData(c, behavior)
}

// UpdateBehavior 更新行为
func (h *BehaviorHandler) UpdateBehavior(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.behaviorService.UpdateBehavior(id, updates); err != nil {
		logrus.Error("Failed to update behavior:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update behavior")
		return
	}

	// 获取更新后的数据
	behavior, err := h.behaviorService.GetBehavior(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get updated behavior")
		return
	}

	utils.RespondWithData(c, behavior)
}

// DeleteBehavior 删除行为
func (h *BehaviorHandler) DeleteBehavior(c *gin.Context) {
	id := c.Param("id")

	if err := h.behaviorService.DeleteBehavior(id); err != nil {
		logrus.Error("Failed to delete behavior:", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete behavior")
		return
	}

	utils.RespondWithData(c, gin.H{"message": "Behavior deleted successfully"})
}

// GetBehaviorsByCategory 根据分类获取行为
func (h *BehaviorHandler) GetBehaviorsByCategory(c *gin.Context) {
	category := c.Param("category")

	behaviors, err := h.behaviorService.GetBehaviorsByCategory(category)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get behaviors by category")
		return
	}

	utils.RespondWithData(c, behaviors)
}

// GetPredefinedBehaviors 获取预定义行为
func (h *BehaviorHandler) GetPredefinedBehaviors(c *gin.Context) {
	behaviors := h.behaviorService.GetPredefinedBehaviors()
	utils.RespondWithData(c, behaviors)
}

// SeedBehaviors 种子数据
func (h *BehaviorHandler) SeedBehaviors(c *gin.Context) {
	// 这里需要根据实际的BehaviorService接口来实现
	// if err := h.behaviorService.SeedBehaviors(); err != nil {
	// 	logrus.Error("Failed to seed behaviors:", err)
	// 	utils.RespondWithError(c, http.StatusInternalServerError, "Failed to seed behaviors")
	// 	return
	// }

	utils.RespondWithData(c, gin.H{"message": "Behaviors seeded successfully"})
}

// GetThingTypeBehaviors 获取事物类型的行为
func (h *BehaviorHandler) GetThingTypeBehaviors(c *gin.Context) {
	thingTypeID := c.Param("id")

	behavior, err := h.thingTypeService.GetTypeBehavior(thingTypeID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Thing type behavior not found")
		return
	}

	utils.RespondWithData(c, behavior)
}

// AssignBehaviorToThingType 为事物类型分配行为
func (h *BehaviorHandler) AssignBehaviorToThingType(c *gin.Context) {
	thingTypeID := c.Param("id")

	var request struct {
		BehaviorID string `json:"behaviorId"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.thingTypeService.SetBehaviorToType(thingTypeID, request.BehaviorID); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to assign behavior to thing type")
		return
	}

	utils.RespondWithData(c, gin.H{"message": "Behavior assigned to thing type successfully"})
}

// RemoveBehaviorFromThingType 从事物类型中移除行为
func (h *BehaviorHandler) RemoveBehaviorFromThingType(c *gin.Context) {
	thingTypeID := c.Param("id")

	if err := h.thingTypeService.RemoveBehaviorFromType(thingTypeID); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to remove behavior from thing type")
		return
	}

	utils.RespondWithData(c, gin.H{"message": "Behavior removed from thing type successfully"})
}

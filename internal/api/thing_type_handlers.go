package api

import (
	"net/http"
	"strconv"
	"uros-restron/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 事物类型相关处理器

// listThingTypes 获取事物类型列表
func (s *Server) listThingTypes(c *gin.Context) {
	category := c.Query("category")
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

	thingTypes, err := s.thingTypeService.ListThingTypes(category, limit, offset)
	if err != nil {
		logrus.Error("Failed to list thing types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list thing types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  thingTypes,
		"count": len(thingTypes),
	})
}

// createThingType 创建事物类型
func (s *Server) createThingType(c *gin.Context) {
	var request map[string]interface{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	features := make(map[string]models.ThingTypeFeature)
	if feats, ok := request["features"].(map[string]interface{}); ok {
		for name, featureData := range feats {
			if featureMap, ok := featureData.(map[string]interface{}); ok {
				// 直接使用 featureMap 作为 properties，而不是嵌套在 properties 字段中
				features[name] = models.ThingTypeFeature{
					Properties: featureMap,
				}
			}
		}
	}
	thingType.Features = features

	if err := s.thingTypeService.CreateThingType(thingType); err != nil {
		logrus.Error("Failed to create thing type:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thing type"})
		return
	}

	c.JSON(http.StatusCreated, thingType)
}

// getThingType 获取单个事物类型
func (s *Server) getThingType(c *gin.Context) {
	id := c.Param("id")

	thingType, err := s.thingTypeService.GetThingType(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thing type not found"})
		return
	}

	c.JSON(http.StatusOK, thingType)
}

// updateThingType 更新事物类型
func (s *Server) updateThingType(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.thingTypeService.UpdateThingType(id, updates); err != nil {
		logrus.Error("Failed to update thing type:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update thing type"})
		return
	}

	// 获取更新后的数据
	thingType, err := s.thingTypeService.GetThingType(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated thing type"})
		return
	}

	c.JSON(http.StatusOK, thingType)
}

// deleteThingType 删除事物类型
func (s *Server) deleteThingType(c *gin.Context) {
	id := c.Param("id")

	if err := s.thingTypeService.DeleteThingType(id); err != nil {
		logrus.Error("Failed to delete thing type:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete thing type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thing type deleted successfully"})
}

// createThingFromType 根据类型创建事物实例
func (s *Server) createThingFromType(c *gin.Context) {
	typeID := c.Param("id")

	var request struct {
		Name        string                    `json:"name"`
		Description string                    `json:"description"`
		Attributes  map[string]interface{}    `json:"attributes"`
		Features    map[string]models.Feature `json:"features"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用类型服务创建事物
	thing, err := s.thingTypeService.CreateThingFromType(typeID, request.Name, request.Description, request.Attributes, request.Features)
	if err != nil {
		logrus.Error("Failed to create thing from type:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thing from type"})
		return
	}

	// 保存到数据库
	if err := s.thingService.CreateThing(thing); err != nil {
		logrus.Error("Failed to save thing:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save thing"})
		return
	}

	// 广播新事物创建事件
	s.hub.Broadcast("thing_created", thing)

	c.JSON(http.StatusCreated, thing)
}

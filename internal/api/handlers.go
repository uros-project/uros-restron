package api

import (
	"net/http"
	"strconv"
	"uros-restron/internal/models"
	"uros-restron/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// listThings 获取数字孪生列表
func (s *Server) listThings(c *gin.Context) {
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

	things, err := s.thingService.ListThings(thingType, limit, offset)
	if err != nil {
		utils.HandleError(c, err, "Failed to list things")
		return
	}

	utils.SuccessResponseWithCount(c, things, len(things))
}

// createThing 创建数字孪生
func (s *Server) createThing(c *gin.Context) {
	var thing models.Thing
	if err := c.ShouldBindJSON(&thing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.thingService.CreateThing(&thing); err != nil {
		logrus.Error("Failed to create thing:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thing"})
		return
	}

	// 广播新事物创建事件
	s.hub.Broadcast("thing_created", thing)

	c.JSON(http.StatusCreated, thing)
}

// getThing 获取单个数字孪生
func (s *Server) getThing(c *gin.Context) {
	id := c.Param("id")

	thing, err := s.thingService.GetThing(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thing not found"})
		return
	}

	c.JSON(http.StatusOK, thing)
}

// updateThing 更新数字孪生
func (s *Server) updateThing(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.thingService.UpdateThing(id, updates); err != nil {
		logrus.Error("Failed to update thing:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update thing"})
		return
	}

	// 获取更新后的数据
	thing, err := s.thingService.GetThing(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated thing"})
		return
	}

	// 广播更新事件
	s.hub.Broadcast("thing_updated", thing)

	c.JSON(http.StatusOK, thing)
}

// deleteThing 删除数字孪生
func (s *Server) deleteThing(c *gin.Context) {
	id := c.Param("id")

	if err := s.thingService.DeleteThing(id); err != nil {
		logrus.Error("Failed to delete thing:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete thing"})
		return
	}

	// 广播删除事件
	s.hub.Broadcast("thing_deleted", gin.H{"id": id})

	c.JSON(http.StatusOK, gin.H{"message": "Thing deleted successfully"})
}

// updateProperty 更新属性
func (s *Server) updateProperty(c *gin.Context) {
	id := c.Param("id")
	propertyName := c.Param("name")

	var request struct {
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.thingService.UpdateProperty(id, propertyName, request.Value); err != nil {
		logrus.Error("Failed to update property:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update property"})
		return
	}

	// 获取更新后的数据
	thing, err := s.thingService.GetThing(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated thing"})
		return
	}

	// 广播属性更新事件
	s.hub.Broadcast("property_updated", gin.H{
		"thingId":  id,
		"property": propertyName,
		"value":    request.Value,
		"thing":    thing,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Property updated successfully"})
}

// updateStatus 更新状态
func (s *Server) updateStatus(c *gin.Context) {
	id := c.Param("id")

	var status map[string]interface{}
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.thingService.UpdateStatus(id, status); err != nil {
		logrus.Error("Failed to update status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	// 获取更新后的数据
	thing, err := s.thingService.GetThing(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated thing"})
		return
	}

	// 广播状态更新事件
	s.hub.Broadcast("status_updated", gin.H{
		"thingId": id,
		"status":  status,
		"thing":   thing,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

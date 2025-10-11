package api

import (
	"net/http"
	"uros-restron/rosix/ai"
	"uros-restron/rosix/core"

	"github.com/gin-gonic/gin"
)

// ROSIXHandler ROSIX API处理器
type ROSIXHandler struct {
	rosix        core.ROSIX
	orchestrator ai.AIOrchestrator
}

// NewROSIXHandler 创建ROSIX处理器
func NewROSIXHandler(rosix core.ROSIX, orchestrator ai.AIOrchestrator) *ROSIXHandler {
	return &ROSIXHandler{
		rosix:        rosix,
		orchestrator: orchestrator,
	}
}

// ==================== 资源操作接口 ====================

// FindResources 查找资源
func (h *ROSIXHandler) FindResources(c *gin.Context) {
	var query core.Query
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resources, err := h.rosix.Find(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	var result []map[string]interface{}
	for _, res := range resources {
		result = append(result, map[string]interface{}{
			"id":         res.ID(),
			"path":       res.Path(),
			"type":       res.Type(),
			"attributes": res.Attributes(),
			"features":   res.Features(),
			"behaviors":  res.Behaviors(),
			"metadata":   res.Metadata(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"resources": result,
			"count":     len(result),
		},
	})
}

// InvokeResource 调用资源行为
func (h *ROSIXHandler) InvokeResource(c *gin.Context) {
	var request struct {
		Path     string                 `json:"path" binding:"required"`
		Behavior string                 `json:"behavior" binding:"required"`
		Params   map[string]interface{} `json:"params"`
		UserID   string                 `json:"user_id"`
		Session  string                 `json:"session"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建上下文
	ctx, err := h.rosix.CreateContext(request.UserID, request.Session, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.DestroyContext(ctx)

	// 打开资源
	rd, err := h.rosix.Open(core.ResourcePath(request.Path), core.ModeInvoke, ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.Close(rd)

	// 调用行为
	result, err := h.rosix.Invoke(rd, request.Behavior, request.Params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"path":     request.Path,
			"behavior": request.Behavior,
			"result":   result,
		},
	})
}

// ReadResource 读取资源属性
func (h *ROSIXHandler) ReadResource(c *gin.Context) {
	var request struct {
		Path    string `json:"path" binding:"required"`
		Key     string `json:"key" binding:"required"`
		UserID  string `json:"user_id"`
		Session string `json:"session"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建上下文
	ctx, err := h.rosix.CreateContext(request.UserID, request.Session, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.DestroyContext(ctx)

	// 打开资源
	rd, err := h.rosix.Open(core.ResourcePath(request.Path), core.ModeRead, ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.Close(rd)

	// 读取属性
	value, err := h.rosix.Read(rd, request.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"path":  request.Path,
			"key":   request.Key,
			"value": value,
		},
	})
}

// WriteResource 写入资源属性
func (h *ROSIXHandler) WriteResource(c *gin.Context) {
	var request struct {
		Path    string      `json:"path" binding:"required"`
		Key     string      `json:"key" binding:"required"`
		Value   interface{} `json:"value" binding:"required"`
		UserID  string      `json:"user_id"`
		Session string      `json:"session"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建上下文
	ctx, err := h.rosix.CreateContext(request.UserID, request.Session, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.DestroyContext(ctx)

	// 打开资源
	rd, err := h.rosix.Open(core.ResourcePath(request.Path), core.ModeWrite, ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.Close(rd)

	// 写入属性
	err = h.rosix.Write(rd, request.Key, request.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"path":  request.Path,
			"key":   request.Key,
			"value": request.Value,
		},
	})
}

// ==================== AI接口 ====================

// AIInvoke AI调用接口
func (h *ROSIXHandler) AIInvoke(c *gin.Context) {
	var request struct {
		Prompt  string `json:"prompt" binding:"required"`
		UserID  string `json:"user_id"`
		Session string `json:"session"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建上下文
	ctx, err := h.rosix.CreateContext(request.UserID, request.Session, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.DestroyContext(ctx)

	// AI调用
	result, err := h.orchestrator.Invoke(request.Prompt, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// AIOrchestrate AI编排接口
func (h *ROSIXHandler) AIOrchestrate(c *gin.Context) {
	var request struct {
		Goal    string `json:"goal" binding:"required"`
		UserID  string `json:"user_id"`
		Session string `json:"session"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建上下文
	ctx, err := h.rosix.CreateContext(request.UserID, request.Session, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.DestroyContext(ctx)

	// AI编排
	plan, err := h.orchestrator.Orchestrate(request.Goal, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plan,
	})
}

// AIQuery AI查询接口
func (h *ROSIXHandler) AIQuery(c *gin.Context) {
	var request struct {
		Question string `json:"question" binding:"required"`
		UserID   string `json:"user_id"`
		Session  string `json:"session"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建上下文
	ctx, err := h.rosix.CreateContext(request.UserID, request.Session, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.DestroyContext(ctx)

	// AI查询
	result, err := h.orchestrator.Query(request.Question, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// AISuggest AI建议接口
func (h *ROSIXHandler) AISuggest(c *gin.Context) {
	var request struct {
		Question string `json:"question" binding:"required"`
		UserID   string `json:"user_id"`
		Session  string `json:"session"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建上下文
	ctx, err := h.rosix.CreateContext(request.UserID, request.Session, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer h.rosix.DestroyContext(ctx)

	// AI建议
	suggestion, err := h.orchestrator.Suggest(request.Question, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    suggestion,
	})
}

// ==================== 系统信息接口 ====================

// GetSystemInfo 获取系统信息
func (h *ROSIXHandler) GetSystemInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"name":        "ROSIX",
			"version":     "1.0.0",
			"description": "Resource Operating System Interface eXtension",
			"features": []string{
				"统一资源访问",
				"标准化操作接口",
				"资源生命周期管理",
				"AI驱动的资源编排",
			},
		},
	})
}

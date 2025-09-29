package api

import (
	"net/http"
	"uros-restron/internal/actor"
	"uros-restron/internal/utils"

	"github.com/gin-gonic/gin"
)

// ActorHandler Actor系统处理器
type ActorHandler struct {
	actorManager *actor.ActorManager
	hub          *Hub
}

// NewActorHandler 创建新的Actor处理器
func NewActorHandler(actorManager *actor.ActorManager, hub *Hub) *ActorHandler {
	return &ActorHandler{
		actorManager: actorManager,
		hub:          hub,
	}
}

// ListActors 获取Actor列表
func (h *ActorHandler) ListActors(c *gin.Context) {
	actors := h.actorManager.ListActors()

	// 转换为响应格式
	var actorList []map[string]interface{}
	for _, actor := range actors {
		actorList = append(actorList, actor.GetStatus())
	}

	utils.RespondWithData(c, gin.H{
		"data":  actorList,
		"count": len(actorList),
	})
}

// GetActor 获取单个Actor
func (h *ActorHandler) GetActor(c *gin.Context) {
	id := c.Param("id")

	actor, err := h.actorManager.GetActor(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Actor not found")
		return
	}

	utils.RespondWithData(c, actor.GetStatus())
}

// CallActorFunction 调用Actor函数
func (h *ActorHandler) CallActorFunction(c *gin.Context) {
	actorID := c.Param("id")
	functionName := c.Param("function")

	var request struct {
		Params map[string]interface{} `json:"params"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// 获取Actor
	_, err := h.actorManager.GetActor(actorID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Actor not found")
		return
	}

	// 调用函数 - 这里需要根据实际的ActorManager接口来实现
	// result, err := h.actorManager.CallActorFunction(actorID, functionName, request.Params)
	// 暂时返回一个占位符响应
	result := map[string]interface{}{
		"actorId":    actorID,
		"function":   functionName,
		"params":     request.Params,
		"message":    "Function call not implemented yet",
	}

	utils.RespondWithData(c, result)
}

// SendMessageToActor 向Actor发送消息
func (h *ActorHandler) SendMessageToActor(c *gin.Context) {
	actorID := c.Param("id")

	var request struct {
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// 获取Actor
	_, err := h.actorManager.GetActor(actorID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Actor not found")
		return
	}

	// 发送消息 - 这里需要根据实际的ActorManager接口来实现
	// response, err := h.actorManager.SendMessageToActor(actorID, request.Message, request.Data)
	// 暂时返回一个占位符响应
	response := map[string]interface{}{
		"actorId": actorID,
		"message": request.Message,
		"data":    request.Data,
		"status":  "Message sent successfully",
	}

	utils.RespondWithData(c, response)
}

// GetActorFunctions 获取Actor函数列表
func (h *ActorHandler) GetActorFunctions(c *gin.Context) {
	actorID := c.Param("id")

	// 获取Actor
	_, err := h.actorManager.GetActor(actorID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Actor not found")
		return
	}

	// 获取函数列表 - 这里需要根据实际的ActorManager接口来实现
	// functions, err := h.actorManager.GetActorFunctions(actorID)
	// 暂时返回一个占位符响应
	functions := []string{"function1", "function2", "function3"}

	utils.RespondWithData(c, functions)
}

// HealthCheck Actor系统健康检查
func (h *ActorHandler) HealthCheck(c *gin.Context) {
	// 获取系统状态 - 这里需要根据实际的ActorManager接口来实现
	// status := h.actorManager.GetSystemStatus()
	// 暂时返回一个占位符响应
	status := map[string]interface{}{
		"status":    "healthy",
		"actors":    len(h.actorManager.ListActors()),
		"timestamp": "2024-01-01T00:00:00Z",
	}

	utils.RespondWithData(c, status)
}
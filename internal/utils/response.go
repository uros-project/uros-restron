package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Response 通用响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Count   int         `json:"count,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessResponseWithCount 带计数的成功响应
func SuccessResponseWithCount(c *gin.Context, data interface{}, count int) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Count:   count,
	})
}

// CreatedResponse 创建成功响应
func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// APIErrorResponse API 错误响应
func APIErrorResponse(c *gin.Context, err *APIError) {
	c.JSON(err.Code, Response{
		Success: false,
		Error:   err.Message,
	})
}

// HandleError 处理错误并返回响应
func HandleError(c *gin.Context, err error, defaultMessage string) {
	if apiErr, ok := IsAPIError(err); ok {
		APIErrorResponse(c, apiErr)
		return
	}

	logrus.Error("Internal error:", err)
	ErrorResponse(c, http.StatusInternalServerError, defaultMessage)
}

// ValidationErrorResponse 验证错误响应
func ValidationErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// NotFoundResponse 未找到响应
func NotFoundResponse(c *gin.Context, resource string) {
	ErrorResponse(c, http.StatusNotFound, resource+" not found")
}

package utils

import (
	"errors"
	"net/http"
)

// APIError 表示 API 错误
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error 实现 error 接口
func (e *APIError) Error() string {
	return e.Message
}

// 预定义的错误
var (
	ErrNotFound     = &APIError{Code: http.StatusNotFound, Message: "资源未找到"}
	ErrBadRequest   = &APIError{Code: http.StatusBadRequest, Message: "请求参数错误"}
	ErrUnauthorized = &APIError{Code: http.StatusUnauthorized, Message: "未授权访问"}
	ErrForbidden    = &APIError{Code: http.StatusForbidden, Message: "禁止访问"}
	ErrConflict     = &APIError{Code: http.StatusConflict, Message: "资源冲突"}
	ErrInternal     = &APIError{Code: http.StatusInternalServerError, Message: "内部服务器错误"}
)

// NewAPIError 创建新的 API 错误
func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewAPIErrorWithDetails 创建带详情的 API 错误
func NewAPIErrorWithDetails(code int, message, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// IsAPIError 检查是否为 API 错误
func IsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

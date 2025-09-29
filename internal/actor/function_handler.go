package actor

import (
	"fmt"
	"log"
)

// DefaultFunctionHandler 默认函数处理器
type DefaultFunctionHandler struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	InputParams    map[string]interface{} `json:"input_params"`
	OutputParams   map[string]interface{} `json:"output_params"`
	Implementation map[string]interface{} `json:"implementation"`
}

// Execute 执行函数
func (h *DefaultFunctionHandler) Execute(params map[string]interface{}) (map[string]interface{}, error) {
	log.Printf("Executing function: %s", h.Name)
	log.Printf("Input parameters: %+v", params)

	// 这里应该根据Implementation中的逻辑来执行函数
	// 为了简化，我们返回一个模拟的结果
	result := map[string]interface{}{
		"function": h.Name,
		"status":   "executed",
		"params":   params,
	}

	return result, nil
}

// GetDefinition 获取函数定义
func (h *DefaultFunctionHandler) GetDefinition() FunctionDefinition {
	return FunctionDefinition{
		Name:         h.Name,
		Description:  h.Description,
		InputParams:  h.InputParams,
		OutputParams: h.OutputParams,
	}
}

// DefaultFunctionExecutor 默认函数执行器
type DefaultFunctionExecutor struct {
	functions map[string]FunctionHandler
}

// NewDefaultFunctionExecutor 创建默认函数执行器
func NewDefaultFunctionExecutor() *DefaultFunctionExecutor {
	return &DefaultFunctionExecutor{
		functions: make(map[string]FunctionHandler),
	}
}

// RegisterFunction 注册函数
func (e *DefaultFunctionExecutor) RegisterFunction(name string, handler FunctionHandler) {
	e.functions[name] = handler
}

// ExecuteFunction 执行函数
func (e *DefaultFunctionExecutor) ExecuteFunction(functionName string, params map[string]interface{}) (map[string]interface{}, error) {
	handler, exists := e.functions[functionName]
	if !exists {
		return nil, fmt.Errorf("function %s not found", functionName)
	}

	return handler.Execute(params)
}

// GetAvailableFunctions 获取可用函数列表
func (e *DefaultFunctionExecutor) GetAvailableFunctions() []string {
	var functions []string
	for name := range e.functions {
		functions = append(functions, name)
	}
	return functions
}

// GetFunctionInfo 获取函数信息
func (e *DefaultFunctionExecutor) GetFunctionInfo(functionName string) (FunctionDefinition, error) {
	handler, exists := e.functions[functionName]
	if !exists {
		return FunctionDefinition{}, fmt.Errorf("function %s not found", functionName)
	}

	return handler.GetDefinition(), nil
}

package actor

import (
	"context"
	"fmt"
	"time"

	"uros-restron/internal/models"
)

// FunctionExecutor 函数执行器
type FunctionExecutor struct {
	behavior  *models.Behavior
	functions map[string]interface{}
}

// NewFunctionExecutor 创建函数执行器
func NewFunctionExecutor(behavior *models.Behavior) *FunctionExecutor {
	// 将 models.Function 转换为 map[string]interface{} 以保持兼容性
	functions := make(map[string]interface{})
	for name, function := range behavior.Functions {
		functions[name] = map[string]interface{}{
			"name":           function.Name,
			"description":    function.Description,
			"input_params":   function.InputParams,
			"output_params":  function.OutputParams,
			"implementation": function.Implementation,
		}
	}
	
	return &FunctionExecutor{
		behavior:  behavior,
		functions: functions,
	}
}

// HasFunction 检查是否有指定函数
func (fe *FunctionExecutor) HasFunction(functionName string) bool {
	_, exists := fe.functions[functionName]
	return exists
}

// GetAvailableFunctions 获取可用函数列表
func (fe *FunctionExecutor) GetAvailableFunctions() []string {
	var functions []string
	for name := range fe.functions {
		functions = append(functions, name)
	}
	return functions
}

// GetFunctionInfo 获取函数信息
func (fe *FunctionExecutor) GetFunctionInfo(functionName string) (map[string]interface{}, error) {
	function, exists := fe.functions[functionName]
	if !exists {
		return nil, fmt.Errorf("function %s not found", functionName)
	}

	functionMap, ok := function.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid function format for %s", functionName)
	}

	return functionMap, nil
}

// ExecuteFunction 执行函数
func (fe *FunctionExecutor) ExecuteFunction(ctx context.Context, functionName string, params map[string]interface{}) (map[string]interface{}, error) {
	// 获取函数定义
	functionInfo, err := fe.GetFunctionInfo(functionName)
	if err != nil {
		return nil, err
	}

	// 验证输入参数
	if err := fe.validateInputParams(functionInfo, params); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %v", err)
	}

	// 执行函数实现
	result, err := fe.executeFunctionImplementation(ctx, functionInfo, params)
	if err != nil {
		return nil, fmt.Errorf("function execution failed: %v", err)
	}

	// 验证输出参数
	if err := fe.validateOutputParams(functionInfo, result); err != nil {
		return nil, fmt.Errorf("output validation failed: %v", err)
	}

	return result, nil
}

// validateInputParams 验证输入参数
func (fe *FunctionExecutor) validateInputParams(functionInfo map[string]interface{}, params map[string]interface{}) error {
	inputParams, exists := functionInfo["input_params"]
	if !exists {
		return nil // 没有输入参数要求
	}

	inputParamsMap, ok := inputParams.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid input_params format")
	}

	for paramName, paramDef := range inputParamsMap {
		paramDefMap, ok := paramDef.(map[string]interface{})
		if !ok {
			continue
		}

		// 检查必需参数
		if required, exists := paramDefMap["required"]; exists {
			if requiredBool, ok := required.(bool); ok && requiredBool {
				if _, exists := params[paramName]; !exists {
					return fmt.Errorf("required parameter %s is missing", paramName)
				}
			}
		}

		// 验证参数值
		if paramValue, exists := params[paramName]; exists {
			if err := fe.validateParamValue(paramName, paramValue, paramDefMap); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateParamValue 验证参数值
func (fe *FunctionExecutor) validateParamValue(paramName string, value interface{}, paramDef map[string]interface{}) error {
	paramType, exists := paramDef["type"]
	if !exists {
		return nil // 没有类型要求
	}

	typeStr, ok := paramType.(string)
	if !ok {
		return nil
	}

	switch typeStr {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("parameter %s must be a string", paramName)
		}
		// 检查字符串长度
		if maxLength, exists := paramDef["max_length"]; exists {
			if maxLen, ok := maxLength.(float64); ok {
				if len(value.(string)) > int(maxLen) {
					return fmt.Errorf("parameter %s exceeds max length %d", paramName, int(maxLen))
				}
			}
		}

	case "number":
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("parameter %s must be a number", paramName)
		}
		// 检查数值范围
		if min, exists := paramDef["min"]; exists {
			if minVal, ok := min.(float64); ok {
				if value.(float64) < minVal {
					return fmt.Errorf("parameter %s must be >= %f", paramName, minVal)
				}
			}
		}
		if max, exists := paramDef["max"]; exists {
			if maxVal, ok := max.(float64); ok {
				if value.(float64) > maxVal {
					return fmt.Errorf("parameter %s must be <= %f", paramName, maxVal)
				}
			}
		}

	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("parameter %s must be a boolean", paramName)
		}

	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("parameter %s must be an object", paramName)
		}
	}

	// 检查枚举值
	if enum, exists := paramDef["enum"]; exists {
		if enumList, ok := enum.([]interface{}); ok {
			valid := false
			for _, enumValue := range enumList {
				if enumValue == value {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("parameter %s must be one of %v", paramName, enumList)
			}
		}
	}

	return nil
}

// validateOutputParams 验证输出参数
func (fe *FunctionExecutor) validateOutputParams(functionInfo map[string]interface{}, result map[string]interface{}) error {
	outputParams, exists := functionInfo["output_params"]
	if !exists {
		return nil // 没有输出参数要求
	}

	outputParamsMap, ok := outputParams.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid output_params format")
	}

	// 检查必需的输出参数
	for paramName, paramDef := range outputParamsMap {
		paramDefMap, ok := paramDef.(map[string]interface{})
		if !ok {
			continue
		}

		// 检查参数是否存在
		if _, exists := result[paramName]; !exists {
			// 如果参数是必需的，返回错误
			if required, exists := paramDefMap["required"]; exists {
				if requiredBool, ok := required.(bool); ok && requiredBool {
					return fmt.Errorf("required output parameter %s is missing", paramName)
				}
			}
		}
	}

	return nil
}

// executeFunctionImplementation 执行函数实现
func (fe *FunctionExecutor) executeFunctionImplementation(ctx context.Context, functionInfo map[string]interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	implementation, exists := functionInfo["implementation"]
	if !exists {
		return nil, fmt.Errorf("function implementation not found")
	}

	implMap, ok := implementation.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid implementation format")
	}

	// 执行实现步骤
	steps, exists := implMap["steps"]
	if !exists {
		return nil, fmt.Errorf("implementation steps not found")
	}

	stepsList, ok := steps.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid steps format")
	}

	// 执行每个步骤
	result := make(map[string]interface{})
	for _, stepInterface := range stepsList {
		step, ok := stepInterface.(map[string]interface{})
		if !ok {
			continue
		}

		stepResult, err := fe.executeStep(ctx, step, params, result)
		if err != nil {
			return nil, fmt.Errorf("step execution failed: %v", err)
		}

		// 合并步骤结果
		for k, v := range stepResult {
			result[k] = v
		}
	}

	return result, nil
}

// executeStep 执行单个步骤
func (fe *FunctionExecutor) executeStep(ctx context.Context, step map[string]interface{}, params map[string]interface{}, currentResult map[string]interface{}) (map[string]interface{}, error) {
	action, exists := step["action"]
	if !exists {
		return nil, fmt.Errorf("step action not found")
	}

	actionStr, ok := action.(string)
	if !ok {
		return nil, fmt.Errorf("invalid action format")
	}

	// 检查条件
	if condition, exists := step["condition"]; exists {
		if conditionStr, ok := condition.(string); ok {
			if !fe.evaluateCondition(conditionStr, params, currentResult) {
				return nil, nil // 条件不满足，跳过此步骤
			}
		}
	}

	// 执行动作
	switch actionStr {
	case "check_air_quality":
		return fe.executeCheckAirQuality(params)
	case "start_fan":
		return fe.executeStartFan(params)
	case "activate_filter":
		return fe.executeActivateFilter(params)
	case "monitor_progress":
		return fe.executeMonitorProgress(params)
	case "read_filter_sensor":
		return fe.executeReadFilterSensor(params)
	case "calculate_usage":
		return fe.executeCalculateUsage(params)
	case "determine_status":
		return fe.executeDetermineStatus(params)
	case "validate_speed":
		return fe.executeValidateSpeed(params)
	case "set_motor_speed":
		return fe.executeSetMotorSpeed(params)
	case "confirm_speed":
		return fe.executeConfirmSpeed(params)
	case "parse_input":
		return fe.executeParseInput(params)
	case "understand_intent":
		return fe.executeUnderstandIntent(params)
	case "generate_response":
		return fe.executeGenerateResponse(params)
	case "format_output":
		return fe.executeFormatOutput(params)
	case "validate_credentials":
		return fe.executeValidateCredentials(params)
	case "check_permissions":
		return fe.executeCheckPermissions(params)
	case "generate_session":
		return fe.executeGenerateSession(params)
	case "load_user_profile":
		return fe.executeLoadUserProfile(params)
	case "extract_preferences":
		return fe.executeExtractPreferences(params)
	case "format_preferences":
		return fe.executeFormatPreferences(params)
	default:
		// 对于未知动作，返回模拟结果
		return fe.executeGenericAction(actionStr, params)
	}
}

// evaluateCondition 评估条件
func (fe *FunctionExecutor) evaluateCondition(condition string, params map[string]interface{}, currentResult map[string]interface{}) bool {
	// 简单的条件评估实现
	// 这里可以实现更复杂的条件评估逻辑
	return true // 默认返回true，实际实现中需要解析条件表达式
}

// 以下是各种动作的具体实现
func (fe *FunctionExecutor) executeCheckAirQuality(params map[string]interface{}) (map[string]interface{}, error) {
	airQuality, _ := params["air_quality"].(float64)
	return map[string]interface{}{
		"current_quality": airQuality,
		"timestamp":       time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeStartFan(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"fan_started": true,
		"timestamp":   time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeActivateFilter(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"filter_activated": true,
		"timestamp":        time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeMonitorProgress(params map[string]interface{}) (map[string]interface{}, error) {
	// 模拟净化进度
	progress := 0.8 // 80%完成
	return map[string]interface{}{
		"progress":  progress,
		"timestamp": time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeReadFilterSensor(params map[string]interface{}) (map[string]interface{}, error) {
	filterID, _ := params["filter_id"].(string)
	return map[string]interface{}{
		"filter_id":    filterID,
		"sensor_value": 75.5,
		"timestamp":    time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeCalculateUsage(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"usage_hours": 1200,
		"timestamp":   time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeDetermineStatus(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"status":         "good",
		"remaining_life": 0.8,
		"timestamp":      time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeValidateSpeed(params map[string]interface{}) (map[string]interface{}, error) {
	speed, _ := params["speed"].(string)
	validSpeeds := []string{"low", "medium", "high", "auto"}
	valid := false
	for _, v := range validSpeeds {
		if v == speed {
			valid = true
			break
		}
	}
	return map[string]interface{}{
		"valid":     valid,
		"timestamp": time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeSetMotorSpeed(params map[string]interface{}) (map[string]interface{}, error) {
	speed, _ := params["speed"].(string)
	return map[string]interface{}{
		"motor_speed": speed,
		"timestamp":   time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeConfirmSpeed(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"confirmed": true,
		"timestamp": time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeParseInput(params map[string]interface{}) (map[string]interface{}, error) {
	userInput, _ := params["user_input"].(string)
	return map[string]interface{}{
		"parsed_input": userInput,
		"timestamp":    time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeUnderstandIntent(params map[string]interface{}) (map[string]interface{}, error) {
	intent := "general_query"
	confidence := 0.85
	return map[string]interface{}{
		"intent":     intent,
		"confidence": confidence,
		"timestamp":  time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeGenerateResponse(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"response":  "处理完成",
		"timestamp": time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeFormatOutput(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"formatted": true,
		"timestamp": time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeValidateCredentials(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"valid":     true,
		"timestamp": time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeCheckPermissions(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"permissions": []string{"read", "write"},
		"timestamp":   time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeGenerateSession(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"session_token": "token_" + fmt.Sprintf("%d", time.Now().Unix()),
		"timestamp":     time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeLoadUserProfile(params map[string]interface{}) (map[string]interface{}, error) {
	userID, _ := params["user_id"].(string)
	return map[string]interface{}{
		"user_id":   userID,
		"profile":   map[string]interface{}{"name": "用户", "email": "user@example.com"},
		"timestamp": time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeExtractPreferences(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"preferences": map[string]interface{}{
			"language": "zh-CN",
			"theme":    "light",
		},
		"timestamp": time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeFormatPreferences(params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"formatted": true,
		"timestamp": time.Now().Unix(),
	}, nil
}

func (fe *FunctionExecutor) executeGenericAction(action string, params map[string]interface{}) (map[string]interface{}, error) {
	// 对于未知动作，返回通用结果
	return map[string]interface{}{
		"action":    action,
		"executed":  true,
		"timestamp": time.Now().Unix(),
	}, nil
}

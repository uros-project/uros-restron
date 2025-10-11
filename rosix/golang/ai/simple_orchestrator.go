package ai

import (
	"fmt"
	"strings"
	"time"
	"uros-restron/rosix/core"
)

// SimpleOrchestrator 简单的AI编排器实现（演示用）
type SimpleOrchestrator struct {
	rosix          core.ROSIX
	intentPatterns map[string]*IntentPattern
}

// NewSimpleOrchestrator 创建简单AI编排器
func NewSimpleOrchestrator(rosix core.ROSIX) *SimpleOrchestrator {
	so := &SimpleOrchestrator{
		rosix:          rosix,
		intentPatterns: make(map[string]*IntentPattern),
	}

	// 注册默认意图模式
	so.registerDefaultIntents()

	return so
}

// registerDefaultIntents 注册默认意图模式
func (so *SimpleOrchestrator) registerDefaultIntents() {
	intents := []*IntentPattern{
		{
			Name:        "read_temperature",
			Patterns:    []string{"温度", "多少度", "temperature"},
			Entities:    []string{"location", "sensor"},
			Description: "读取温度",
		},
		{
			Name:        "set_temperature",
			Patterns:    []string{"调节温度", "设置温度", "调到", "set temperature"},
			Entities:    []string{"location", "value"},
			Description: "设置温度",
		},
		{
			Name:        "air_purify",
			Patterns:    []string{"净化", "空气", "purify", "pm2.5"},
			Entities:    []string{"location", "intensity"},
			Description: "空气净化",
		},
		{
			Name:        "find_item",
			Patterns:    []string{"找", "查找", "在哪", "find"},
			Entities:    []string{"item_name"},
			Description: "查找物品",
		},
	}

	for _, intent := range intents {
		so.intentPatterns[intent.Name] = intent
	}
}

// Invoke 通过自然语言调用资源
func (so *SimpleOrchestrator) Invoke(prompt string, context *core.Context) (*InvokeResult, error) {
	// 简单的意图识别（基于关键词匹配）
	intent := so.recognizeIntent(prompt)

	result := &InvokeResult{
		Intent:    intent.Name,
		Resources: []string{},
		Actions:   []Action{},
		Result:    make(map[string]interface{}),
	}

	// 根据意图执行操作
	switch intent.Name {
	case "read_temperature":
		return so.handleReadTemperature(prompt, context)

	case "set_temperature":
		return so.handleSetTemperature(prompt, context)

	case "air_purify":
		return so.handleAirPurify(prompt, context)

	case "find_item":
		return so.handleFindItem(prompt, context)

	default:
		result.Success = false
		result.Message = "无法识别意图"
	}

	return result, nil
}

// recognizeIntent 识别意图
func (so *SimpleOrchestrator) recognizeIntent(text string) *Intent {
	text = strings.ToLower(text)

	for name, pattern := range so.intentPatterns {
		for _, keyword := range pattern.Patterns {
			if strings.Contains(text, keyword) {
				return &Intent{
					Name:       name,
					Confidence: 0.8,
					Entities:   make(map[string]interface{}),
					Slots:      make(map[string]string),
				}
			}
		}
	}

	return &Intent{
		Name:       "unknown",
		Confidence: 0,
		Entities:   make(map[string]interface{}),
		Slots:      make(map[string]string),
	}
}

// handleReadTemperature 处理读取温度意图
func (so *SimpleOrchestrator) handleReadTemperature(prompt string, ctx *core.Context) (*InvokeResult, error) {
	// 查找环境监测类型的资源
	resources, err := so.rosix.Find(core.Query{
		Category: "environment",
		Limit:    1,
	})

	if err != nil || len(resources) == 0 {
		return &InvokeResult{
			Success: false,
			Intent:  "read_temperature",
			Message: "未找到环境监测资源",
		}, nil
	}

	// 打开资源
	rd, err := so.rosix.Open(resources[0].Path(), core.ModeRead|core.ModeInvoke, ctx)
	if err != nil {
		return &InvokeResult{
			Success: false,
			Intent:  "read_temperature",
			Message: fmt.Sprintf("打开资源失败: %v", err),
		}, nil
	}
	defer so.rosix.Close(rd)

	// 调用读取环境数据的行为
	result, err := so.rosix.Invoke(rd, "read_environment", map[string]interface{}{
		"metrics": []string{"temperature"},
	})

	if err != nil {
		return &InvokeResult{
			Success: false,
			Intent:  "read_temperature",
			Message: fmt.Sprintf("调用失败: %v", err),
		}, nil
	}

	return &InvokeResult{
		Success:   true,
		Intent:    "read_temperature",
		Resources: []string{resources[0].ID()},
		Actions: []Action{
			{
				Type:     ActionInvoke,
				Behavior: "read_environment",
				Parameters: map[string]interface{}{
					"metrics": []string{"temperature"},
				},
			},
		},
		Result:  map[string]interface{}{"data": result},
		Message: "成功读取温度数据",
	}, nil
}

// handleSetTemperature 处理设置温度意图
func (so *SimpleOrchestrator) handleSetTemperature(prompt string, ctx *core.Context) (*InvokeResult, error) {
	// 这里可以扩展更复杂的参数提取逻辑
	return &InvokeResult{
		Success: false,
		Intent:  "set_temperature",
		Message: "设置温度功能正在开发中",
	}, nil
}

// handleAirPurify 处理空气净化意图
func (so *SimpleOrchestrator) handleAirPurify(prompt string, ctx *core.Context) (*InvokeResult, error) {
	// 查找空气净化器资源
	resources, err := so.rosix.Find(core.Query{
		Category: "purifier",
		Limit:    1,
	})

	if err != nil || len(resources) == 0 {
		return &InvokeResult{
			Success: false,
			Intent:  "air_purify",
			Message: "未找到空气净化器资源",
		}, nil
	}

	// 打开资源
	rd, err := so.rosix.Open(resources[0].Path(), core.ModeInvoke, ctx)
	if err != nil {
		return &InvokeResult{
			Success: false,
			Intent:  "air_purify",
			Message: fmt.Sprintf("打开资源失败: %v", err),
		}, nil
	}
	defer so.rosix.Close(rd)

	// 调用净化空气行为
	result, err := so.rosix.Invoke(rd, "purify_air", map[string]interface{}{
		"mode":      "auto",
		"intensity": 3,
	})

	if err != nil {
		return &InvokeResult{
			Success: false,
			Intent:  "air_purify",
			Message: fmt.Sprintf("调用失败: %v", err),
		}, nil
	}

	return &InvokeResult{
		Success:   true,
		Intent:    "air_purify",
		Resources: []string{resources[0].ID()},
		Actions: []Action{
			{
				Type:     ActionInvoke,
				Behavior: "purify_air",
				Parameters: map[string]interface{}{
					"mode":      "auto",
					"intensity": 3,
				},
			},
		},
		Result:  map[string]interface{}{"data": result},
		Message: "成功启动空气净化",
	}, nil
}

// handleFindItem 处理查找物品意图
func (so *SimpleOrchestrator) handleFindItem(prompt string, ctx *core.Context) (*InvokeResult, error) {
	return &InvokeResult{
		Success: false,
		Intent:  "find_item",
		Message: "查找物品功能正在开发中",
	}, nil
}

// Orchestrate 编排多资源协同
func (so *SimpleOrchestrator) Orchestrate(goal string, context *core.Context) (*Plan, error) {
	plan := &Plan{
		ID:            fmt.Sprintf("plan_%d", time.Now().UnixNano()),
		Goal:          goal,
		Steps:         []PlanStep{},
		Resources:     []string{},
		EstimatedTime: 10,
		Metadata:      make(map[string]interface{}),
	}

	// 简单示例：根据目标生成计划
	if strings.Contains(goal, "睡眠模式") || strings.Contains(goal, "晚上") {
		plan.Steps = []PlanStep{
			{
				Order:       1,
				Description: "关闭所有灯光",
				Action: Action{
					Type:       ActionInvoke,
					Behavior:   "turn_off",
					Parameters: map[string]interface{}{},
				},
			},
			{
				Order:       2,
				Description: "调低空调温度到24度",
				Action: Action{
					Type:       ActionInvoke,
					Behavior:   "set_temperature",
					Parameters: map[string]interface{}{"value": 24},
				},
				DependsOn: []int{1},
			},
			{
				Order:       3,
				Description: "启动空气净化器",
				Action: Action{
					Type:       ActionInvoke,
					Behavior:   "purify_air",
					Parameters: map[string]interface{}{"mode": "silent"},
				},
				DependsOn: []int{2},
			},
		}
	}

	return plan, nil
}

// Query 查询资源信息
func (so *SimpleOrchestrator) Query(question string, context *core.Context) (*QueryResult, error) {
	return &QueryResult{
		Answer:     "查询功能正在开发中",
		Data:       make(map[string]interface{}),
		Resources:  []string{},
		Confidence: 0.5,
	}, nil
}

// Suggest 获取建议
func (so *SimpleOrchestrator) Suggest(question string, context *core.Context) (*Suggestion, error) {
	return &Suggestion{
		Question:    question,
		Suggestions: []SuggestionItem{},
		Context:     make(map[string]interface{}),
	}, nil
}

// Learn 学习用户行为模式
func (so *SimpleOrchestrator) Learn(behavior *UserBehavior) error {
	// 实现学习逻辑
	return nil
}

// Predict 预测资源状态
func (so *SimpleOrchestrator) Predict(resourceID string, duration int) (*Prediction, error) {
	return &Prediction{
		ResourceID:  resourceID,
		Duration:    duration,
		Predictions: make(map[string]interface{}),
		Confidence:  0.5,
	}, nil
}

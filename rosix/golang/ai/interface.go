package ai

import (
	"uros-restron/rosix/core"
)

// AIOrchestrator AI编排器接口
type AIOrchestrator interface {
	// Invoke 通过自然语言调用资源
	// 示例: "将客厅温度调节到25度"
	Invoke(prompt string, context *core.Context) (*InvokeResult, error)

	// Orchestrate 编排多资源协同
	// 示例: "晚上8点，关闭所有灯光，调低空调温度"
	Orchestrate(goal string, context *core.Context) (*Plan, error)

	// Query 查询资源信息
	// 示例: "客厅的温度是多少？"
	Query(question string, context *core.Context) (*QueryResult, error)

	// Suggest 获取建议
	// 示例: "如何降低室内PM2.5？"
	Suggest(question string, context *core.Context) (*Suggestion, error)

	// Learn 学习用户行为模式
	Learn(behavior *UserBehavior) error

	// Predict 预测资源状态
	Predict(resourceID string, duration int) (*Prediction, error)
}

// InvokeResult AI调用结果
type InvokeResult struct {
	Success   bool                   `json:"success"`
	Intent    string                 `json:"intent"`    // 识别的意图
	Resources []string               `json:"resources"` // 涉及的资源
	Actions   []Action               `json:"actions"`   // 执行的动作
	Result    map[string]interface{} `json:"result"`    // 执行结果
	Message   string                 `json:"message"`   // 反馈消息
}

// Plan 执行计划
type Plan struct {
	ID            string                 `json:"id"`
	Goal          string                 `json:"goal"`
	Steps         []PlanStep             `json:"steps"`
	Resources     []string               `json:"resources"`
	EstimatedTime int                    `json:"estimated_time"` // 预计执行时间(秒)
	Metadata      map[string]interface{} `json:"metadata"`
}

// PlanStep 计划步骤
type PlanStep struct {
	Order       int    `json:"order"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      Action `json:"action"`
	DependsOn   []int  `json:"depends_on"` // 依赖的步骤序号
	Condition   string `json:"condition,omitempty"`
}

// Action 动作定义
type Action struct {
	Type       ActionType             `json:"type"`
	Behavior   string                 `json:"behavior,omitempty"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ActionType 动作类型
type ActionType string

const (
	ActionRead   ActionType = "read"
	ActionWrite  ActionType = "write"
	ActionInvoke ActionType = "invoke"
	ActionWait   ActionType = "wait"
)

// QueryResult 查询结果
type QueryResult struct {
	Answer     string                 `json:"answer"`
	Data       map[string]interface{} `json:"data"`
	Resources  []string               `json:"resources"`
	Confidence float64                `json:"confidence"` // 置信度 0-1
}

// Suggestion 建议
type Suggestion struct {
	Question    string                 `json:"question"`
	Suggestions []SuggestionItem       `json:"suggestions"`
	Context     map[string]interface{} `json:"context"`
}

// SuggestionItem 建议项
type SuggestionItem struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    int      `json:"priority"`
	Confidence  float64  `json:"confidence"`
	Actions     []Action `json:"actions"`
}

// UserBehavior 用户行为
type UserBehavior struct {
	UserID    string                 `json:"user_id"`
	Action    string                 `json:"action"`
	Resources []string               `json:"resources"`
	Timestamp int64                  `json:"timestamp"`
	Context   map[string]interface{} `json:"context"`
}

// Prediction 预测结果
type Prediction struct {
	ResourceID  string                 `json:"resource_id"`
	Duration    int                    `json:"duration"` // 预测时长(秒)
	Predictions map[string]interface{} `json:"predictions"`
	Confidence  float64                `json:"confidence"`
}

// IntentRecognizer 意图识别器
type IntentRecognizer interface {
	// Recognize 识别意图
	Recognize(text string) (*Intent, error)

	// RegisterIntent 注册意图模式
	RegisterIntent(intent *IntentPattern) error
}

// Intent 意图
type Intent struct {
	Name       string                 `json:"name"`
	Confidence float64                `json:"confidence"`
	Entities   map[string]interface{} `json:"entities"`
	Slots      map[string]string      `json:"slots"`
}

// IntentPattern 意图模式
type IntentPattern struct {
	Name        string   `json:"name"`
	Patterns    []string `json:"patterns"`
	Entities    []string `json:"entities"`
	Handler     string   `json:"handler"`
	Description string   `json:"description"`
}

// ResourceSelector AI资源选择器
type ResourceSelector interface {
	// Select 根据条件选择资源
	Select(criteria string, context *core.Context) ([]core.Resource, error)

	// Rank 对资源进行排序
	Rank(resources []core.Resource, criteria string) ([]core.Resource, error)

	// Filter 过滤资源
	Filter(resources []core.Resource, condition string) ([]core.Resource, error)
}

package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BehaviorType 行为类型
type BehaviorType string

// Behavior 定义事物的行为模式
type Behavior struct {
	ID          string       `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name"`
	Type        BehaviorType `json:"type"`
	Description string       `json:"description"`
	Category    string       `json:"category"` // 行为分类：device, person, object


	// 函数定义 - 包含多个函数及其实现
	Functions     map[string]Function `json:"functions" gorm:"-"`
	FunctionsJSON string              `json:"-" gorm:"column:functions;type:text"`

	// 行为参数
	Parameters     map[string]interface{} `json:"parameters" gorm:"-"` // 行为参数
	ParametersJSON string                 `json:"-" gorm:"column:parameters;type:text"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Function 定义行为中的函数
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputParams map[string]Parameter   `json:"input_params"`
	OutputParams map[string]Parameter `json:"output_params"`
	Implementation FunctionImplementation `json:"implementation"`
}

// Parameter 定义函数参数
type Parameter struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Min         *float64    `json:"min,omitempty"`
	Max         *float64    `json:"max,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
}

// FunctionImplementation 定义函数实现
type FunctionImplementation struct {
	Steps []ImplementationStep `json:"steps"`
}

// ImplementationStep 定义实现步骤
type ImplementationStep struct {
	Step        int    `json:"step"`
	Action      string `json:"action"`
	Description string `json:"description"`
	Condition   string `json:"condition,omitempty"`
}

// BeforeCreate GORM hook for serializing data before creation
func (b *Behavior) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New().String()
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	return b.serializeData()
}

// BeforeUpdate GORM hook for serializing data before update
func (b *Behavior) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return b.serializeData()
}

// AfterFind GORM hook for deserializing data after retrieval
func (b *Behavior) AfterFind(tx *gorm.DB) error {
	return b.deserializeData()
}

// serializeData 序列化 JSON 字段
func (b *Behavior) serializeData() error {
	if b.Functions != nil {
		if data, err := json.Marshal(b.Functions); err != nil {
			return err
		} else {
			b.FunctionsJSON = string(data)
		}
	}
	if b.Parameters != nil {
		if data, err := json.Marshal(b.Parameters); err != nil {
			return err
		} else {
			b.ParametersJSON = string(data)
		}
	}
	return nil
}

// deserializeData 反序列化 JSON 字段
func (b *Behavior) deserializeData() error {
	if b.FunctionsJSON != "" {
		if err := json.Unmarshal([]byte(b.FunctionsJSON), &b.Functions); err != nil {
			return err
		}
	}
	if b.ParametersJSON != "" {
		if err := json.Unmarshal([]byte(b.ParametersJSON), &b.Parameters); err != nil {
			return err
		}
	}
	return nil
}

// BehaviorService 提供行为相关的业务逻辑
type BehaviorService struct {
	db *gorm.DB
}

// NewBehaviorService 创建新的 BehaviorService
func NewBehaviorService(db *gorm.DB) *BehaviorService {
	return &BehaviorService{db: db}
}

// CreateBehavior 创建新的行为
func (s *BehaviorService) CreateBehavior(behavior *Behavior) error {
	return s.db.Create(behavior).Error
}

// GetBehavior 根据ID获取行为
func (s *BehaviorService) GetBehavior(id string) (*Behavior, error) {
	var behavior Behavior
	if err := s.db.First(&behavior, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &behavior, nil
}

// ListBehaviors 获取行为列表
func (s *BehaviorService) ListBehaviors(behaviorType, category string, limit, offset int) ([]Behavior, error) {
	var behaviors []Behavior
	query := s.db
	if behaviorType != "" {
		query = query.Where("type = ?", behaviorType)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	return behaviors, query.Limit(limit).Offset(offset).Find(&behaviors).Error
}

// UpdateBehavior 更新行为
func (s *BehaviorService) UpdateBehavior(id string, updates map[string]interface{}) error {
	return s.db.Model(&Behavior{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteBehavior 删除行为
func (s *BehaviorService) DeleteBehavior(id string) error {
	return s.db.Delete(&Behavior{}, "id = ?", id).Error
}

// GetBehaviorsByCategory 根据分类获取行为
func (s *BehaviorService) GetBehaviorsByCategory(category string) ([]Behavior, error) {
	var behaviors []Behavior
	return behaviors, s.db.Where("category = ?", category).Find(&behaviors).Error
}

// GetPredefinedBehaviors 获取预定义行为
func (s *BehaviorService) GetPredefinedBehaviors() []Behavior {
	if behaviors, err := LoadBehaviorsFromPath("./behaviors"); err != nil {
		return []Behavior{}
	} else {
		return behaviors
	}
}

// GetAllBehaviors 获取所有行为
func (s *BehaviorService) GetAllBehaviors() ([]Behavior, error) {
	var behaviors []Behavior
	return behaviors, s.db.Find(&behaviors).Error
}

// SeedPredefinedBehaviors 填充预定义行为到数据库
func (s *BehaviorService) SeedPredefinedBehaviors() error {
	for _, behavior := range s.GetPredefinedBehaviors() {
		var existing Behavior
		if err := s.db.Where("id = ?", behavior.ID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&behavior).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

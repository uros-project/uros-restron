package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


// ThingType 表示事物类型定义 - 符合 Ditto 标准
type ThingType struct {
	ID             string                      `json:"id" gorm:"primaryKey"`
	Name           string                      `json:"name"`
	Description    string                      `json:"description"`
	Category       string                      `json:"category"`                                         // person, machine, object
	Attributes     map[string]interface{}      `json:"attributes" gorm:"-"`                              // 属性模式定义，不存储到数据库
	AttributesJSON string                      `json:"-" gorm:"column:attributes;type:text"`             // 数据库存储的 JSON 字符串
	Features       map[string]interface{}      `json:"features" gorm:"-"`                                // 功能模式定义，不存储到数据库
	FeaturesJSON   string                      `json:"-" gorm:"column:features;type:text"`               // 数据库存储的 JSON 字符串
	BehaviorID     string                      `json:"behaviorId"`                                          // 关联的行为ID
	Behavior       *Behavior                   `json:"behavior,omitempty" gorm:"foreignKey:BehaviorID"` // 关联的行为
	CreatedAt      time.Time                   `json:"createdAt"`
	UpdatedAt      time.Time                   `json:"updatedAt"`
}

// BeforeCreate GORM 钩子，在创建前序列化 Attributes 和 Features
func (t *ThingType) BeforeCreate(tx *gorm.DB) error {
	return t.serializeData()
}

// BeforeUpdate GORM 钩子，在更新前序列化 Attributes 和 Features
func (t *ThingType) BeforeUpdate(tx *gorm.DB) error {
	return t.serializeData()
}

// AfterFind GORM 钩子，在查询后反序列化 Attributes 和 Features
func (t *ThingType) AfterFind(tx *gorm.DB) error {
	return t.deserializeData()
}

// serializeData 序列化 Attributes 和 Features 为 JSON 字符串
func (t *ThingType) serializeData() error {
	if t.Attributes != nil {
		data, err := json.Marshal(t.Attributes)
		if err != nil {
			return err
		}
		t.AttributesJSON = string(data)
	}

	if t.Features != nil {
		data, err := json.Marshal(t.Features)
		if err != nil {
			return err
		}
		t.FeaturesJSON = string(data)
	}
	return nil
}

// deserializeData 反序列化 JSON 字符串为 Attributes 和 Features
func (t *ThingType) deserializeData() error {
	if t.AttributesJSON != "" {
		err := json.Unmarshal([]byte(t.AttributesJSON), &t.Attributes)
		if err != nil {
			return err
		}
	}

	if t.FeaturesJSON != "" {
		err := json.Unmarshal([]byte(t.FeaturesJSON), &t.Features)
		if err != nil {
			return err
		}
	}
	return nil
}


// ThingTypeService 提供事物类型相关的业务逻辑
type ThingTypeService struct {
	db *gorm.DB
}

// NewThingTypeService 创建新的 ThingTypeService
func NewThingTypeService(db *gorm.DB) *ThingTypeService {
	return &ThingTypeService{db: db}
}

// CreateThingType 创建新的事物类型
func (s *ThingTypeService) CreateThingType(thingType *ThingType) error {
	if thingType.ID == "" {
		thingType.ID = uuid.New().String()
	}
	thingType.CreatedAt = time.Now()
	thingType.UpdatedAt = time.Now()

	// 手动序列化 Attributes 和 Features
	if thingType.Attributes != nil {
		data, err := json.Marshal(thingType.Attributes)
		if err != nil {
			return err
		}
		thingType.AttributesJSON = string(data)
	}

	if thingType.Features != nil {
		data, err := json.Marshal(thingType.Features)
		if err != nil {
			return err
		}
		thingType.FeaturesJSON = string(data)
	}

	return s.db.Create(thingType).Error
}

// GetThingType 根据ID获取事物类型
func (s *ThingTypeService) GetThingType(id string) (*ThingType, error) {
	var thingType ThingType
	err := s.db.First(&thingType, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// 反序列化 Attributes 和 Features
	if thingType.AttributesJSON != "" {
		err = json.Unmarshal([]byte(thingType.AttributesJSON), &thingType.Attributes)
		if err != nil {
			return nil, err
		}
	}

	if thingType.FeaturesJSON != "" {
		err = json.Unmarshal([]byte(thingType.FeaturesJSON), &thingType.Features)
		if err != nil {
			return nil, err
		}
	}

	return &thingType, nil
}

// ListThingTypes 获取所有事物类型
func (s *ThingTypeService) ListThingTypes(category string, limit, offset int) ([]ThingType, error) {
	var thingTypes []ThingType
	query := s.db

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Limit(limit).Offset(offset).Find(&thingTypes).Error
	if err != nil {
		return nil, err
	}

	// 反序列化每个 ThingType 的 Attributes 和 Features
	for i := range thingTypes {
		if thingTypes[i].AttributesJSON != "" {
			err = json.Unmarshal([]byte(thingTypes[i].AttributesJSON), &thingTypes[i].Attributes)
			if err != nil {
				return nil, err
			}
		}
		if thingTypes[i].FeaturesJSON != "" {
			err = json.Unmarshal([]byte(thingTypes[i].FeaturesJSON), &thingTypes[i].Features)
			if err != nil {
				return nil, err
			}
		}
	}

	return thingTypes, nil
}

// UpdateThingType 更新事物类型
func (s *ThingTypeService) UpdateThingType(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return s.db.Model(&ThingType{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteThingType 删除事物类型
func (s *ThingTypeService) DeleteThingType(id string) error {
	return s.db.Delete(&ThingType{}, "id = ?", id).Error
}

// CreateThingFromType 根据类型创建事物实例
func (s *ThingTypeService) CreateThingFromType(thingTypeID string, name, description string, attributes map[string]interface{}, features map[string]interface{}) (*Thing, error) {
	// 获取类型定义
	thingType, err := s.GetThingType(thingTypeID)
	if err != nil {
		return nil, err
	}

	// 创建事物实例
	thing := &Thing{
		ID:          uuid.New().String(),
		Name:        name,
		Type:        thingType.Category,
		Description: description,
		Attributes:  make(map[string]interface{}),
		Features:    make(map[string]interface{}),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 设置 Attributes
	if attributes != nil {
		thing.Attributes = attributes
	}

	// 设置 Features
	if features != nil {
		thing.Features = features
	} else {
		// 根据类型功能创建功能
		for featureName := range thingType.Features {
			thing.Features[featureName] = make(map[string]interface{})
		}
	}

	// 序列化 Attributes 和 Features
	if thing.Attributes != nil {
		data, err := json.Marshal(thing.Attributes)
		if err != nil {
			return nil, err
		}
		thing.AttributesJSON = string(data)
	}

	if thing.Features != nil {
		data, err := json.Marshal(thing.Features)
		if err != nil {
			return nil, err
		}
		thing.FeaturesJSON = string(data)
	}

	return thing, nil
}

// AssignDefaultBehavior 为 ThingType 分配默认行为
func (s *ThingTypeService) AssignDefaultBehavior(thingTypeID string) error {
	var thingType ThingType
	if err := s.db.First(&thingType, "id = ?", thingTypeID).Error; err != nil {
		return err
	}

	// 根据分类获取对应的行为
	var behavior Behavior
	if err := s.db.Where("category = ?", thingType.Category).First(&behavior).Error; err != nil {
		return err
	}

	// 分配行为
	return s.db.Model(&thingType).Update("behavior_id", behavior.ID).Error
}

// SetBehaviorToType 为 ThingType 设置行为
func (s *ThingTypeService) SetBehaviorToType(thingTypeID, behaviorID string) error {
	return s.db.Model(&ThingType{}).Where("id = ?", thingTypeID).Update("behavior_id", behaviorID).Error
}

// RemoveBehaviorFromType 从 ThingType 移除行为
func (s *ThingTypeService) RemoveBehaviorFromType(thingTypeID string) error {
	return s.db.Model(&ThingType{}).Where("id = ?", thingTypeID).Update("behavior_id", nil).Error
}

// GetTypeBehavior 获取 ThingType 的行为
func (s *ThingTypeService) GetTypeBehavior(thingTypeID string) (*Behavior, error) {
	var thingType ThingType
	if err := s.db.Preload("Behavior").First(&thingType, "id = ?", thingTypeID).Error; err != nil {
		return nil, err
	}
	return thingType.Behavior, nil
}


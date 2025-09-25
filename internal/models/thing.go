package models

import (
	"encoding/json"
	"time"

	"uros-restron/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Thing 表示数字孪生实体 - 符合 Ditto 标准
type Thing struct {
	ID             string                 `json:"id" gorm:"primaryKey"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"` // person, machine, object
	Description    string                 `json:"description"`
	Attributes     map[string]interface{} `json:"attributes" gorm:"-"`                  // 静态元数据，不存储到数据库
	AttributesJSON string                 `json:"-" gorm:"column:attributes;type:text"` // 存储为 JSON 字符串
	Features       map[string]Feature     `json:"features" gorm:"-"`                    // 动态功能，不存储到数据库
	FeaturesJSON   string                 `json:"-" gorm:"column:features;type:text"`   // 存储为 JSON 字符串
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
}

// Property 表示事物的属性
type Property struct {
	ID      string      `json:"id" gorm:"primaryKey"`
	ThingID string      `json:"thingId"`
	Name    string      `json:"name"`
	Value   interface{} `json:"value" gorm:"type:text"`
	Type    string      `json:"type"` // string, number, boolean, object
}

// Feature 表示事物的功能 - 符合 Ditto 标准
type Feature struct {
	Properties map[string]interface{} `json:"properties"` // 功能的状态和配置
}

// ThingService 提供数字孪生相关的业务逻辑
type ThingService struct {
	db        *gorm.DB
	jsonUtils *utils.JSONUtils
}

// NewThingService 创建新的 ThingService
func NewThingService(db *gorm.DB) *ThingService {
	return &ThingService{
		db:        db,
		jsonUtils: &utils.JSONUtils{},
	}
}

// CreateThing 创建新的数字孪生
func (s *ThingService) CreateThing(thing *Thing) error {
	if thing.ID == "" {
		thing.ID = uuid.New().String()
	}
	thing.CreatedAt = time.Now()
	thing.UpdatedAt = time.Now()

	// 序列化 Attributes
	attributesJSON, err := s.jsonUtils.SerializeMap(thing.Attributes)
	if err != nil {
		return err
	}
	thing.AttributesJSON = attributesJSON

	// 序列化 Features
	featuresMap := make(map[string]interface{})
	for k, v := range thing.Features {
		featuresMap[k] = v
	}
	featuresJSON, err := s.jsonUtils.SerializeFeatures(featuresMap)
	if err != nil {
		return err
	}
	thing.FeaturesJSON = featuresJSON

	return s.db.Create(thing).Error
}

// GetThing 根据ID获取数字孪生
func (s *ThingService) GetThing(id string) (*Thing, error) {
	var thing Thing
	err := s.db.First(&thing, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// 反序列化 Attributes
	attributes, err := s.jsonUtils.DeserializeMap(thing.AttributesJSON)
	if err != nil {
		return nil, err
	}
	thing.Attributes = attributes

	// 反序列化 Features
	featuresMap, err := s.jsonUtils.DeserializeFeatures(thing.FeaturesJSON)
	if err != nil {
		return nil, err
	}
	thing.Features = make(map[string]Feature)
	for k, v := range featuresMap {
		if featureMap, ok := v.(map[string]interface{}); ok {
			thing.Features[k] = Feature{Properties: featureMap}
		}
	}

	return &thing, nil
}

// ListThings 获取所有数字孪生
func (s *ThingService) ListThings(thingType string, limit, offset int) ([]Thing, error) {
	var things []Thing
	query := s.db

	if thingType != "" {
		query = query.Where("type = ?", thingType)
	}

	err := query.Limit(limit).Offset(offset).Find(&things).Error
	if err != nil {
		return nil, err
	}

	// 反序列化每个 Thing 的 Attributes 和 Features
	for i := range things {
		if things[i].AttributesJSON != "" {
			err = json.Unmarshal([]byte(things[i].AttributesJSON), &things[i].Attributes)
			if err != nil {
				return nil, err
			}
		}
		if things[i].FeaturesJSON != "" {
			err = json.Unmarshal([]byte(things[i].FeaturesJSON), &things[i].Features)
			if err != nil {
				return nil, err
			}
		}
	}

	return things, nil
}

// UpdateThing 更新数字孪生
func (s *ThingService) UpdateThing(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	// 处理 Attributes 和 Features 的序列化
	if attributes, ok := updates["attributes"]; ok {
		if attrs, ok := attributes.(map[string]interface{}); ok {
			data, err := json.Marshal(attrs)
			if err != nil {
				return err
			}
			updates["attributes"] = string(data)
		}
	}

	if features, ok := updates["features"]; ok {
		if feats, ok := features.(map[string]Feature); ok {
			data, err := json.Marshal(feats)
			if err != nil {
				return err
			}
			updates["features"] = string(data)
		}
	}

	return s.db.Model(&Thing{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteThing 删除数字孪生
func (s *ThingService) DeleteThing(id string) error {
	return s.db.Delete(&Thing{}, "id = ?", id).Error
}

// UpdateProperty 更新属性
func (s *ThingService) UpdateProperty(thingID, propertyName string, value interface{}) error {
	// 先尝试更新现有属性
	result := s.db.Model(&Property{}).Where("thing_id = ? AND name = ?", thingID, propertyName).Update("value", value)

	// 如果属性不存在，创建新属性
	if result.RowsAffected == 0 {
		property := &Property{
			ID:      uuid.New().String(),
			ThingID: thingID,
			Name:    propertyName,
			Value:   value,
		}
		return s.db.Create(property).Error
	}

	// 更新事物的更新时间
	return s.db.Model(&Thing{}).Where("id = ?", thingID).Update("updated_at", time.Now()).Error
}

// UpdateStatus 更新状态
func (s *ThingService) UpdateStatus(thingID string, status map[string]interface{}) error {
	statusJSON, err := json.Marshal(status)
	if err != nil {
		return err
	}

	return s.db.Model(&Thing{}).Where("id = ?", thingID).Updates(map[string]interface{}{
		"status":     string(statusJSON),
		"updated_at": time.Now(),
	}).Error
}

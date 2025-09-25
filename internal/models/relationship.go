package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RelationshipType 关系类型枚举
type RelationshipType string

const (
	// 强关联关系
	RelationshipTypeContains RelationshipType = "contains" // 包含关系
	RelationshipTypeComposes RelationshipType = "composes" // 组合关系
	RelationshipTypeOwns     RelationshipType = "owns"     // 拥有关系

	// 弱关联关系
	RelationshipTypeRelatesTo    RelationshipType = "relates_to"   // 关联关系
	RelationshipTypeDependsOn    RelationshipType = "depends_on"   // 依赖关系
	RelationshipTypeInfluences   RelationshipType = "influences"   // 影响关系
	RelationshipTypeCollaborates RelationshipType = "collaborates" // 协作关系
)

// Relationship 关系模型
type Relationship struct {
	ID             string                 `json:"id" gorm:"primaryKey"`
	SourceID       string                 `json:"sourceId" gorm:"not null;index"`       // 源事物ID
	TargetID       string                 `json:"targetId" gorm:"not null;index"`       // 目标事物ID
	Type           RelationshipType       `json:"type" gorm:"not null"`                 // 关系类型
	Name           string                 `json:"name"`                                 // 关系名称
	Description    string                 `json:"description"`                          // 关系描述
	Properties     map[string]interface{} `json:"properties" gorm:"-"`                  // 关系属性，不存储到数据库
	PropertiesJSON string                 `json:"-" gorm:"column:properties;type:text"` // 存储为 JSON 字符串
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`

	// 关联对象
	Source *Thing `json:"source,omitempty" gorm:"foreignKey:SourceID"`
	Target *Thing `json:"target,omitempty" gorm:"foreignKey:TargetID"`
}

// RelationshipService 关系服务
type RelationshipService struct {
	db *gorm.DB
}

// NewRelationshipService 创建关系服务
func NewRelationshipService(db *gorm.DB) *RelationshipService {
	return &RelationshipService{db: db}
}

// CreateRelationship 创建关系
func (s *RelationshipService) CreateRelationship(relationship *Relationship) error {
	if relationship.ID == "" {
		relationship.ID = uuid.New().String()
	}
	relationship.CreatedAt = time.Now()
	relationship.UpdatedAt = time.Now()

	// 序列化 Properties
	if relationship.Properties != nil {
		data, err := json.Marshal(relationship.Properties)
		if err != nil {
			return err
		}
		relationship.PropertiesJSON = string(data)
	}

	return s.db.Create(relationship).Error
}

// GetRelationship 获取单个关系
func (s *RelationshipService) GetRelationship(id string) (*Relationship, error) {
	var relationship Relationship
	err := s.db.Preload("Source").Preload("Target").First(&relationship, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// 反序列化 Properties
	if relationship.PropertiesJSON != "" {
		err = json.Unmarshal([]byte(relationship.PropertiesJSON), &relationship.Properties)
		if err != nil {
			return nil, err
		}
	}

	return &relationship, nil
}

// ListRelationships 获取关系列表
func (s *RelationshipService) ListRelationships(sourceID, targetID string, relationshipType RelationshipType, limit, offset int) ([]Relationship, error) {
	var relationships []Relationship
	query := s.db.Preload("Source").Preload("Target")

	if sourceID != "" {
		query = query.Where("source_id = ?", sourceID)
	}
	if targetID != "" {
		query = query.Where("target_id = ?", targetID)
	}
	if relationshipType != "" {
		query = query.Where("type = ?", relationshipType)
	}

	err := query.Limit(limit).Offset(offset).Find(&relationships).Error
	if err != nil {
		return nil, err
	}

	// 反序列化每个关系的 Properties
	for i := range relationships {
		if relationships[i].PropertiesJSON != "" {
			err = json.Unmarshal([]byte(relationships[i].PropertiesJSON), &relationships[i].Properties)
			if err != nil {
				return nil, err
			}
		}
	}

	return relationships, nil
}

// UpdateRelationship 更新关系
func (s *RelationshipService) UpdateRelationship(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	// 处理 Properties 的序列化
	if properties, ok := updates["properties"]; ok {
		if props, ok := properties.(map[string]interface{}); ok {
			data, err := json.Marshal(props)
			if err != nil {
				return err
			}
			updates["properties"] = string(data)
		}
	}

	return s.db.Model(&Relationship{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteRelationship 删除关系
func (s *RelationshipService) DeleteRelationship(id string) error {
	return s.db.Delete(&Relationship{}, "id = ?", id).Error
}

// GetThingRelationships 获取事物的所有关系
func (s *RelationshipService) GetThingRelationships(thingID string) ([]Relationship, error) {
	var relationships []Relationship
	err := s.db.Preload("Source").Preload("Target").
		Where("source_id = ? OR target_id = ?", thingID, thingID).
		Find(&relationships).Error
	if err != nil {
		return nil, err
	}

	// 反序列化每个关系的 Properties
	for i := range relationships {
		if relationships[i].PropertiesJSON != "" {
			err = json.Unmarshal([]byte(relationships[i].PropertiesJSON), &relationships[i].Properties)
			if err != nil {
				return nil, err
			}
		}
	}

	return relationships, nil
}

// GetRelationshipTypes 获取所有关系类型
func (s *RelationshipService) GetRelationshipTypes() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":        string(RelationshipTypeContains),
			"name":        "包含关系",
			"description": "表示一个事物包含另一个事物",
			"strength":    "strong",
		},
		{
			"type":        string(RelationshipTypeComposes),
			"name":        "组合关系",
			"description": "表示一个事物由另一个事物组成",
			"strength":    "strong",
		},
		{
			"type":        string(RelationshipTypeOwns),
			"name":        "拥有关系",
			"description": "表示一个事物拥有另一个事物",
			"strength":    "strong",
		},
		{
			"type":        string(RelationshipTypeRelatesTo),
			"name":        "关联关系",
			"description": "表示两个事物之间存在一般关联",
			"strength":    "weak",
		},
		{
			"type":        string(RelationshipTypeDependsOn),
			"name":        "依赖关系",
			"description": "表示一个事物依赖于另一个事物",
			"strength":    "weak",
		},
		{
			"type":        string(RelationshipTypeInfluences),
			"name":        "影响关系",
			"description": "表示一个事物影响另一个事物",
			"strength":    "weak",
		},
		{
			"type":        string(RelationshipTypeCollaborates),
			"name":        "协作关系",
			"description": "表示两个事物之间存在协作关系",
			"strength":    "weak",
		},
	}
}

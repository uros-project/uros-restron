package utils

import (
	"regexp"
	"strings"
)

// ValidationUtils 提供验证功能
type ValidationUtils struct{}

// ValidateID 验证 ID 格式
func (v *ValidationUtils) ValidateID(id string) error {
	if id == "" {
		return NewAPIError(400, "ID cannot be empty")
	}

	// UUID 格式验证
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(id) {
		return NewAPIError(400, "Invalid ID format")
	}

	return nil
}

// ValidateName 验证名称
func (v *ValidationUtils) ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return NewAPIError(400, "Name cannot be empty")
	}

	if len(name) > 100 {
		return NewAPIError(400, "Name too long")
	}

	return nil
}

// ValidateThingType 验证事物类型
func (v *ValidationUtils) ValidateThingType(thingType string) error {
	validTypes := []string{"person", "machine", "object"}

	for _, validType := range validTypes {
		if thingType == validType {
			return nil
		}
	}

	return NewAPIError(400, "Invalid thing type")
}

// ValidateRelationshipType 验证关系类型
func (v *ValidationUtils) ValidateRelationshipType(relType string) error {
	validTypes := []string{
		"contains", "composes", "owns",
		"relates_to", "depends_on", "influences", "collaborates",
	}

	for _, validType := range validTypes {
		if relType == validType {
			return nil
		}
	}

	return NewAPIError(400, "Invalid relationship type")
}

// ValidatePagination 验证分页参数
func (v *ValidationUtils) ValidatePagination(limit, offset int) error {
	if limit < 1 || limit > 100 {
		return NewAPIError(400, "Limit must be between 1 and 100")
	}

	if offset < 0 {
		return NewAPIError(400, "Offset must be non-negative")
	}

	return nil
}

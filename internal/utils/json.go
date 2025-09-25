package utils

import (
	"encoding/json"
)

// JSONUtils 提供 JSON 序列化和反序列化的通用方法
type JSONUtils struct{}

// SerializeMap 将 map[string]interface{} 序列化为 JSON 字符串
func (j *JSONUtils) SerializeMap(data map[string]interface{}) (string, error) {
	if data == nil {
		return "", nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// DeserializeMap 将 JSON 字符串反序列化为 map[string]interface{}
func (j *JSONUtils) DeserializeMap(jsonStr string) (map[string]interface{}, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SerializeFeatures 将 map[string]Feature 序列化为 JSON 字符串
func (j *JSONUtils) SerializeFeatures(features map[string]interface{}) (string, error) {
	if features == nil {
		return "", nil
	}
	bytes, err := json.Marshal(features)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// DeserializeFeatures 将 JSON 字符串反序列化为 map[string]interface{}
func (j *JSONUtils) DeserializeFeatures(jsonStr string) (map[string]interface{}, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var features map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &features)
	if err != nil {
		return nil, err
	}
	return features, nil
}

// SerializeRelationships 将关系属性序列化为 JSON 字符串
func (j *JSONUtils) SerializeRelationships(properties map[string]interface{}) (string, error) {
	if properties == nil {
		return "", nil
	}
	bytes, err := json.Marshal(properties)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// DeserializeRelationships 将 JSON 字符串反序列化为关系属性
func (j *JSONUtils) DeserializeRelationships(jsonStr string) (map[string]interface{}, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var properties map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &properties)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

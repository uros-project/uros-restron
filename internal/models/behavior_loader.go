package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// LoadBehaviorsFromPath 从指定路径加载所有行为定义
func LoadBehaviorsFromPath(behaviorsPath string) ([]Behavior, error) {
	var allBehaviors []Behavior

	// 扫描行为目录
	categories, err := scanCategories(behaviorsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan categories: %v", err)
	}

	// 加载每个分类的行为
	for _, category := range categories {
		behaviors, err := loadBehaviorsFromCategory(behaviorsPath, category)
		if err != nil {
			return nil, fmt.Errorf("failed to load behaviors from category %s: %v", category, err)
		}
		allBehaviors = append(allBehaviors, behaviors...)
	}

	return allBehaviors, nil
}

// scanCategories 扫描行为分类目录
func scanCategories(behaviorsPath string) ([]string, error) {
	var categories []string

	entries, err := ioutil.ReadDir(behaviorsPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			categories = append(categories, entry.Name())
		}
	}

	return categories, nil
}

// loadBehaviorsFromCategory 从指定分类加载行为
func loadBehaviorsFromCategory(behaviorsPath, category string) ([]Behavior, error) {
	var behaviors []Behavior

	categoryPath := filepath.Join(behaviorsPath, category)
	entries, err := ioutil.ReadDir(categoryPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			behavior, err := loadBehaviorFromFile(filepath.Join(categoryPath, entry.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to load behavior from file %s: %v", entry.Name(), err)
			}
			behaviors = append(behaviors, behavior)
		}
	}

	return behaviors, nil
}

// loadBehaviorFromFile 从文件加载单个行为
func loadBehaviorFromFile(filePath string) (Behavior, error) {
	var behavior Behavior

	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return behavior, err
	}

	// 解析 JSON
	if err := json.Unmarshal(data, &behavior); err != nil {
		return behavior, err
	}

	return behavior, nil
}

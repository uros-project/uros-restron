package resource

import (
	"fmt"
	"sync"
	"uros-restron/rosix/core"
)

// Registry 资源注册表实现
type Registry struct {
	resources map[string]core.Resource       // ID -> Resource
	pathIndex map[core.ResourcePath]string   // Path -> ID
	typeIndex map[core.ResourceType][]string // Type -> []ID
	mu        sync.RWMutex
	watchers  []func(resource core.Resource, action string)
}

// NewRegistry 创建资源注册表
func NewRegistry() *Registry {
	return &Registry{
		resources: make(map[string]core.Resource),
		pathIndex: make(map[core.ResourcePath]string),
		typeIndex: make(map[core.ResourceType][]string),
		watchers:  make([]func(resource core.Resource, action string), 0),
	}
}

// Register 注册资源
func (r *Registry) Register(resource core.Resource) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := resource.ID()
	path := resource.Path()
	resourceType := resource.Type()

	// 检查是否已存在
	if _, exists := r.resources[id]; exists {
		return &core.Error{
			Code:    core.ErrInvalidParameter,
			Message: fmt.Sprintf("resource %s already registered", id),
		}
	}

	// 注册资源
	r.resources[id] = resource
	r.pathIndex[path] = id

	// 更新类型索引
	if _, ok := r.typeIndex[resourceType]; !ok {
		r.typeIndex[resourceType] = make([]string, 0)
	}
	r.typeIndex[resourceType] = append(r.typeIndex[resourceType], id)

	// 通知监听器
	r.notifyWatchers(resource, "register")

	return nil
}

// Unregister 注销资源
func (r *Registry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	resource, exists := r.resources[id]
	if !exists {
		return &core.Error{
			Code:    core.ErrNotFound,
			Message: fmt.Sprintf("resource %s not found", id),
		}
	}

	// 删除索引
	delete(r.pathIndex, resource.Path())
	delete(r.resources, id)

	// 更新类型索引
	resourceType := resource.Type()
	if ids, ok := r.typeIndex[resourceType]; ok {
		newIds := make([]string, 0)
		for _, rid := range ids {
			if rid != id {
				newIds = append(newIds, rid)
			}
		}
		r.typeIndex[resourceType] = newIds
	}

	// 通知监听器
	r.notifyWatchers(resource, "unregister")

	return nil
}

// Get 获取资源
func (r *Registry) Get(id string) (core.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resource, exists := r.resources[id]
	if !exists {
		return nil, &core.Error{
			Code:    core.ErrNotFound,
			Message: fmt.Sprintf("resource %s not found", id),
		}
	}

	return resource, nil
}

// GetByPath 通过路径获取资源
func (r *Registry) GetByPath(path core.ResourcePath) (core.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.pathIndex[path]
	if !exists {
		return nil, &core.Error{
			Code:    core.ErrNotFound,
			Message: fmt.Sprintf("resource at path %s not found", path),
		}
	}

	return r.resources[id], nil
}

// Query 查询资源
func (r *Registry) Query(query core.Query) ([]core.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []core.Resource
	var candidates []core.Resource

	// 如果指定了类型，从类型索引开始
	if query.Type != "" {
		if ids, ok := r.typeIndex[query.Type]; ok {
			for _, id := range ids {
				if resource, exists := r.resources[id]; exists {
					candidates = append(candidates, resource)
				}
			}
		}
	} else {
		// 否则遍历所有资源
		for _, resource := range r.resources {
			candidates = append(candidates, resource)
		}
	}

	// 过滤候选资源
	for _, resource := range candidates {
		if r.matchQuery(resource, query) {
			results = append(results, resource)

			// 应用限制
			if query.Limit > 0 && len(results) >= query.Limit {
				break
			}
		}
	}

	return results, nil
}

// matchQuery 检查资源是否匹配查询条件
func (r *Registry) matchQuery(resource core.Resource, query core.Query) bool {
	metadata := resource.Metadata()

	// 检查分类
	if query.Category != "" && metadata.Category != query.Category {
		return false
	}

	// 检查标签
	if len(query.Tags) > 0 {
		hasTag := false
		for _, queryTag := range query.Tags {
			for _, resourceTag := range metadata.Tags {
				if queryTag == resourceTag {
					hasTag = true
					break
				}
			}
			if hasTag {
				break
			}
		}
		if !hasTag {
			return false
		}
	}

	// 检查属性
	if len(query.Attributes) > 0 {
		resourceAttrs := resource.Attributes()
		for key, value := range query.Attributes {
			if resourceAttrs[key] != value {
				return false
			}
		}
	}

	// 检查特征
	if len(query.Features) > 0 {
		resourceFeatures := resource.Features()
		for key, value := range query.Features {
			if resourceFeatures[key] != value {
				return false
			}
		}
	}

	return true
}

// Watch 监听资源注册变化
func (r *Registry) Watch(callback func(resource core.Resource, action string)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.watchers = append(r.watchers, callback)
	return nil
}

// notifyWatchers 通知所有监听器
func (r *Registry) notifyWatchers(resource core.Resource, action string) {
	for _, watcher := range r.watchers {
		go watcher(resource, action)
	}
}

// List 列出所有资源
func (r *Registry) List() []core.Resource {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resources := make([]core.Resource, 0, len(r.resources))
	for _, resource := range r.resources {
		resources = append(resources, resource)
	}
	return resources
}

// Count 返回资源数量
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.resources)
}

// CountByType 返回指定类型的资源数量
func (r *Registry) CountByType(resourceType core.ResourceType) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if ids, ok := r.typeIndex[resourceType]; ok {
		return len(ids)
	}
	return 0
}

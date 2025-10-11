package resource

import (
	"fmt"
	"uros-restron/internal/models"
	"uros-restron/rosix/core"
)

// ThingAdapter 将内部Thing模型适配为ROSIX Resource接口
type ThingAdapter struct {
	thing *models.Thing
}

// NewThingAdapter 创建Thing适配器
func NewThingAdapter(thing *models.Thing) *ThingAdapter {
	return &ThingAdapter{
		thing: thing,
	}
}

// ID 返回资源ID
func (ta *ThingAdapter) ID() string {
	return ta.thing.ID
}

// Path 返回资源路径
func (ta *ThingAdapter) Path() core.ResourcePath {
	// 构建资源路径: /things/{type}/{id}
	return core.ResourcePath(fmt.Sprintf("/things/%s/%s", ta.thing.Type, ta.thing.ID))
}

// Type 返回资源类型
func (ta *ThingAdapter) Type() core.ResourceType {
	return core.TypeDevice // 默认为设备类型，可根据thing.Type进行映射
}

// Attributes 返回资源的静态属性
func (ta *ThingAdapter) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"id":          ta.thing.ID,
		"name":        ta.thing.Name,
		"type":        ta.thing.Type,
		"type_id":     ta.thing.TypeID,
		"description": ta.thing.Description,
		"metadata":    ta.thing.Metadata,
	}
}

// Features 返回资源的动态属性
func (ta *ThingAdapter) Features() map[string]interface{} {
	features := make(map[string]interface{})

	// 从Status中提取特征
	if ta.thing.Status != nil {
		for key, value := range ta.thing.Status {
			features[key] = value
		}
	}

	return features
}

// Behaviors 返回资源支持的行为列表
func (ta *ThingAdapter) Behaviors() []string {
	behaviors := []string{}

	// 从关联的Behavior中提取行为名称
	for _, behavior := range ta.thing.Behaviors {
		if behavior.Functions != nil {
			for funcName := range behavior.Functions {
				behaviors = append(behaviors, funcName)
			}
		}
	}

	return behaviors
}

// Metadata 返回资源元数据
func (ta *ThingAdapter) Metadata() core.ResourceMetadata {
	return core.ResourceMetadata{
		Name:        ta.thing.Name,
		Description: ta.thing.Description,
		Category:    ta.thing.Type,
		Tags:        []string{ta.thing.Type},
		CreatedAt:   ta.thing.CreatedAt,
		UpdatedAt:   ta.thing.UpdatedAt,
		Owner:       "",
		Extra:       ta.thing.Metadata,
	}
}

// GetThing 获取内部Thing对象
func (ta *ThingAdapter) GetThing() *models.Thing {
	return ta.thing
}

// ActorAdapter 将Actor适配为ROSIX Resource接口
type ActorAdapter struct {
	actorID   string
	name      string
	status    string
	functions []string
}

// NewActorAdapter 创建Actor适配器
func NewActorAdapter(actorID, name, status string, functions []string) *ActorAdapter {
	return &ActorAdapter{
		actorID:   actorID,
		name:      name,
		status:    status,
		functions: functions,
	}
}

// ID 返回资源ID
func (aa *ActorAdapter) ID() string {
	return aa.actorID
}

// Path 返回资源路径
func (aa *ActorAdapter) Path() core.ResourcePath {
	return core.ResourcePath(fmt.Sprintf("/actors/%s", aa.actorID))
}

// Type 返回资源类型
func (aa *ActorAdapter) Type() core.ResourceType {
	return core.TypeActor
}

// Attributes 返回资源的静态属性
func (aa *ActorAdapter) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"id":   aa.actorID,
		"name": aa.name,
	}
}

// Features 返回资源的动态属性
func (aa *ActorAdapter) Features() map[string]interface{} {
	return map[string]interface{}{
		"status": aa.status,
	}
}

// Behaviors 返回资源支持的行为列表
func (aa *ActorAdapter) Behaviors() []string {
	return aa.functions
}

// Metadata 返回资源元数据
func (aa *ActorAdapter) Metadata() core.ResourceMetadata {
	return core.ResourceMetadata{
		Name:        aa.name,
		Description: "Actor resource",
		Category:    "actor",
		Tags:        []string{"actor", "behavior"},
		Extra:       map[string]interface{}{},
	}
}

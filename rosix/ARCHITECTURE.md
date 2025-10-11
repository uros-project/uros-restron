# ROSIX 架构设计

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                     应用层 (Applications)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Web应用     │  │  CLI工具     │  │  AI Agent    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │ HTTP/RPC
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              ROSIX 编程层 (Programming Layer)                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  ROSIX 系统调用 (System Calls)                       │   │
│  │  Open/Close/Read/Write/Invoke/Find/Watch...         │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  AI编排器    │  │  资源注册表  │  │  事件总线    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │ 适配层
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              资源层 (Resource Layer)                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Things      │  │  Actors      │  │  Behaviors   │      │
│  │  (设备对象)  │  │  (行为实例)  │  │  (行为定义)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              基础设施层 (Infrastructure)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  数据库      │  │  消息队列    │  │  日志系统    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## 核心组件

### 1. Core（核心）
- **types.go**: 核心数据类型定义
  - ResourceDescriptor: 资源描述符
  - Resource: 资源接口
  - Context: 执行上下文
  - Query: 查询条件
  - Event: 事件定义

- **interface.go**: 核心接口定义
  - ROSIX: 主系统调用接口
  - ResourceRegistry: 资源注册表接口
  - BehaviorExecutor: 行为执行器接口
  - EventBus: 事件总线接口

### 2. Resource（资源）
- **adapter.go**: 资源适配器
  - ThingAdapter: Thing到Resource的适配
  - ActorAdapter: Actor到Resource的适配
  
- **registry.go**: 资源注册表实现
  - 资源注册/注销
  - 资源查询
  - 资源索引（ID、Path、Type）

### 3. Syscall（系统调用）
- **rosix.go**: ROSIX系统调用实现
  - 资源操作：Open/Close/Read/Write/Invoke
  - 资源发现：Find/List/Stat
  - 资源监听：Watch/Unwatch
  - 资源关系：Link/Unlink/GetRelations
  - 上下文管理：CreateContext/DestroyContext

### 4. AI（AI协同）
- **interface.go**: AI接口定义
  - AIOrchestrator: AI编排器接口
  - IntentRecognizer: 意图识别接口
  - ResourceSelector: 资源选择器接口

- **simple_orchestrator.go**: 简单AI编排器实现
  - 自然语言调用
  - 多资源编排
  - 意图识别
  - 资源查询

### 5. API（HTTP接口）
- **handlers.go**: HTTP请求处理器
  - 资源操作接口
  - AI调用接口
  - 系统信息接口

- **routes.go**: 路由定义
  - RESTful API路由映射

## 数据流

### 场景1：直接资源调用

```
用户代码
  ↓
rosix.Open(path)
  ↓
Registry.GetByPath(path) → 查找资源
  ↓
创建ResourceHandle
  ↓
返回ResourceDescriptor
  ↓
rosix.Invoke(rd, behavior, params)
  ↓
ActorManager.CallFunction()
  ↓
Actor.Execute()
  ↓
返回结果
```

### 场景2：AI驱动调用

```
用户自然语言输入
  ↓
AIOrchestrator.Invoke(prompt)
  ↓
IntentRecognizer.Recognize() → 识别意图
  ↓
ResourceSelector.Select() → 选择资源
  ↓
生成执行计划
  ↓
rosix.Open() → 打开资源
  ↓
rosix.Invoke() → 调用行为
  ↓
返回执行结果
```

## 资源生命周期

```
┌──────────┐
│  创建     │  Thing/Actor创建
└────┬─────┘
     │
     ▼
┌──────────┐
│  注册     │  Registry.Register()
└────┬─────┘
     │
     ▼
┌──────────┐
│  发现     │  rosix.Find()
└────┬─────┘
     │
     ▼
┌──────────┐
│  打开     │  rosix.Open() → 获得RD
└────┬─────┘
     │
     ▼
┌──────────┐
│  使用     │  Read/Write/Invoke
└────┬─────┘
     │
     ▼
┌──────────┐
│  关闭     │  rosix.Close()
└────┬─────┘
     │
     ▼
┌──────────┐
│  注销     │  Registry.Unregister()
└────┬─────┘
     │
     ▼
┌──────────┐
│  销毁     │  Thing/Actor销毁
└──────────┘
```

## 资源路径设计

类似文件系统路径，ROSIX使用层级路径标识资源：

```
/things/{type}/{id}         - Thing类型资源
/actors/{id}                - Actor类型资源
/devices/{category}/{id}    - 设备类型资源
/objects/{category}/{id}    - 对象类型资源
/persons/{id}               - 人员类型资源
/services/{name}            - 服务类型资源
```

示例：
```
/things/purifier/device_001
/actors/bc62b42f-b7d3-4216-a3da-34408e74a7e8
/devices/sensor/temp_sensor_01
/objects/container/storage_box_01
/persons/user_001
/services/notification
```

## 资源描述符管理

```
ResourceDescriptor (RD) 类似文件描述符：
- 从1000开始分配
- 每个打开的资源对应一个RD
- RD映射到ResourceHandle
- Handle包含：
  - Resource引用
  - 打开模式（Read/Write/Invoke/Watch）
  - 执行上下文
  - 访问时间戳
```

## AI编排流程

```
1. 意图识别
   输入: "打开客厅的空气净化器"
   输出: Intent{name: "air_purify", entities: {"location": "客厅"}}

2. 资源选择
   查询: Query{Category: "purifier", Location: "客厅"}
   输出: [Resource1, Resource2, ...]

3. 生成计划
   Plan{
     Steps: [
       {Resource: purifier_001, Action: Invoke{purify_air}}
     ]
   }

4. 执行计划
   - Open(resource.Path())
   - Invoke(rd, "purify_air", params)
   - Close(rd)

5. 返回结果
   InvokeResult{
     Success: true,
     Message: "空气净化器已启动"
   }
```

## 扩展性设计

### 1. 新增资源类型
```go
// 1. 定义资源适配器
type NewResourceAdapter struct {
    data *NewResourceType
}

func (a *NewResourceAdapter) ID() string { ... }
func (a *NewResourceAdapter) Path() core.ResourcePath { ... }
// 实现Resource接口的其他方法

// 2. 注册到Registry
registry.Register(newAdapter)
```

### 2. 新增系统调用
```go
// 在core/interface.go中添加接口方法
type ROSIX interface {
    // ... 现有方法
    NewOperation(params) error
}

// 在syscall/rosix.go中实现
func (s *System) NewOperation(params) error {
    // 实现逻辑
}
```

### 3. 新增AI意图
```go
// 在AI编排器中注册新意图
orchestrator.RegisterIntent(&IntentPattern{
    Name:     "new_intent",
    Patterns: []string{"关键词1", "关键词2"},
    Entities: []string{"entity1", "entity2"},
    Handler:  "handleNewIntent",
})
```

## 性能优化

### 1. 资源索引
- ID索引：O(1)查找
- Path索引：O(1)查找
- Type索引：O(1)分类查找

### 2. 并发控制
- Registry使用RWMutex
- Handle映射使用RWMutex
- 支持并发读取

### 3. 资源缓存
- 缓存频繁访问的资源
- LRU淘汰策略
- 可配置缓存大小

### 4. 事件聚合
- 批量发送事件
- 事件过滤和去重
- 异步事件处理

## 安全性设计

### 1. 权限控制
- 基于Context的权限验证
- 资源级别的访问控制
- 操作级别的权限检查

### 2. 资源隔离
- 不同用户的资源隔离
- Session级别的资源管理
- 租户级别的多租户支持

### 3. 审计日志
- 记录所有资源操作
- 记录AI决策过程
- 支持操作回溯

## 监控指标

### 1. 资源指标
- 注册资源数量
- 打开句柄数量
- 资源类型分布

### 2. 性能指标
- 操作响应时间
- 并发请求数
- 错误率

### 3. AI指标
- 意图识别准确率
- 资源选择准确率
- 执行成功率

## 未来扩展

### 1. 分布式支持
- 分布式资源注册表
- 远程资源访问
- 跨节点资源协同

### 2. 高级AI功能
- 深度学习模型集成
- 用户行为预测
- 自动化运维

### 3. 资源编排语言
- 声明式资源编排
- 可视化编排工具
- 编排模板库

### 4. 生态系统
- 资源插件系统
- 第三方资源接入
- 资源市场


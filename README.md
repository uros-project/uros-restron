# Uros Restron - 简化版数字孪生平台

一个基于 Go 的简化版 Eclipse Ditto 实现，用于管理人、机、物等资源的数字化镜像。

## 特性

- 🚀 **高性能**: 基于 Go 构建，支持高并发
- 🔄 **实时通信**: WebSocket 支持实时数据同步
- 📊 **RESTful API**: 完整的 REST API 接口
- 💾 **数据持久化**: 基于 SQLite 的数据存储
- 🎯 **类型支持**: 支持人(person)、机(machine)、物(object)等类型
- 📡 **实时广播**: 支持属性、状态变更的实时通知

## 快速开始

### 安装依赖

```bash
go mod tidy
```

### 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 文档

### 数字孪生管理

#### 创建数字孪生
```bash
POST /api/v1/things
Content-Type: application/json

{
  "name": "智能传感器001",
  "type": "machine",
  "description": "温度湿度传感器",
  "properties": [
    {
      "name": "temperature",
      "value": 25.5,
      "type": "number"
    }
  ],
  "status": {
    "online": true,
    "battery": 85
  }
}
```

#### 获取数字孪生列表
```bash
GET /api/v1/things?type=machine&limit=10&offset=0
```

#### 获取单个数字孪生
```bash
GET /api/v1/things/{id}
```

#### 更新数字孪生
```bash
PUT /api/v1/things/{id}
Content-Type: application/json

{
  "name": "更新后的名称",
  "description": "更新后的描述"
}
```

#### 删除数字孪生
```bash
DELETE /api/v1/things/{id}
```

### 属性管理

#### 更新属性
```bash
PUT /api/v1/things/{id}/properties/{propertyName}
Content-Type: application/json

{
  "value": 26.8
}
```

#### 更新状态
```bash
PUT /api/v1/things/{id}/status
Content-Type: application/json

{
  "online": true,
  "battery": 90,
  "lastSeen": "2024-01-01T12:00:00Z"
}
```

### WebSocket 实时通信

连接到 WebSocket 端点：
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws');

// 订阅特定事物的更新
ws.send(JSON.stringify({
  type: 'subscribe',
  data: 'thing-id-here'
}));

// 监听消息
ws.onmessage = function(event) {
  const message = JSON.parse(event.data);
  console.log('收到消息:', message);
};
```

## 支持的消息类型

- `thing_created`: 新数字孪生创建
- `thing_updated`: 数字孪生更新
- `thing_deleted`: 数字孪生删除
- `property_updated`: 属性更新
- `status_updated`: 状态更新

## 项目结构

```
uros-restron/
├── main.go                 # 应用入口
├── go.mod                  # Go 模块文件
├── internal/
│   ├── api/               # API 层
│   │   ├── server.go      # HTTP 服务器
│   │   ├── handlers.go    # API 处理器
│   │   └── websocket.go   # WebSocket 处理
│   ├── config/            # 配置管理
│   ├── database/          # 数据库连接
│   └── models/            # 数据模型
└── README.md              # 项目文档
```

## 环境变量

- `PORT`: 服务端口 (默认: 8080)
- `HOST`: 服务主机 (默认: localhost)
- `DATABASE_DSN`: 数据库连接字符串 (默认: things.db)

## 示例使用场景

### 1. 智能设备监控
```bash
# 创建设备数字孪生
curl -X POST http://localhost:8080/api/v1/things \
  -H "Content-Type: application/json" \
  -d '{
    "name": "智能空调001",
    "type": "machine",
    "description": "客厅空调",
    "properties": [
      {"name": "temperature", "value": 22, "type": "number"},
      {"name": "mode", "value": "cooling", "type": "string"}
    ],
    "status": {"online": true, "power": "on"}
  }'
```

### 2. 人员管理
```bash
# 创建人员数字孪生
curl -X POST http://localhost:8080/api/v1/things \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "type": "person",
    "description": "系统管理员",
    "properties": [
      {"name": "department", "value": "IT", "type": "string"},
      {"name": "role", "value": "admin", "type": "string"}
    ],
    "status": {"active": true, "lastLogin": "2024-01-01T09:00:00Z"}
  }'
```

### 3. 物品追踪
```bash
# 创建物品数字孪生
curl -X POST http://localhost:8080/api/v1/things \
  -H "Content-Type: application/json" \
  -d '{
    "name": "笔记本电脑001",
    "type": "object",
    "description": "公司配发笔记本",
    "properties": [
      {"name": "serialNumber", "value": "ABC123456", "type": "string"},
      {"name": "location", "value": "办公室A", "type": "string"}
    ],
    "status": {"inUse": true, "assignedTo": "张三"}
  }'
```

## 开发计划

- [ ] 添加认证和授权
- [ ] 支持更多数据库 (PostgreSQL, MySQL)
- [ ] 添加数据验证
- [ ] 实现批量操作
- [ ] 添加监控和日志
- [ ] 支持数据导入导出
- [ ] 添加 GraphQL 支持

## 许可证

MIT License

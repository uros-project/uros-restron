# ROSIX Java Implementation

ROSIX (Resource Operating System Interface eXtension) 的Java实现版本。

## 概述

这是ROSIX编程层的Java语言实现，提供与Go版本相同的API接口和功能，但使用Java生态系统的工具和框架。

## 技术栈

- **Java**: 17+
- **Spring Boot**: 3.2.0
- **Maven**: 项目构建
- **Lombok**: 减少样板代码
- **Jackson**: JSON处理
- **SLF4J**: 日志框架

## 项目结构

```
src/main/java/com/uros/rosix/
├── core/              # 核心接口和类型
│   ├── Resource.java
│   ├── ResourceDescriptor.java
│   ├── ResourcePath.java
│   ├── Context.java
│   └── ROSIX.java
│
├── resource/          # 资源层
│   ├── ResourceAdapter.java
│   ├── ThingAdapter.java
│   ├── ActorAdapter.java
│   └── ResourceRegistry.java
│
├── syscall/           # 系统调用实现
│   └── ROSIXSystem.java
│
├── ai/                # AI协同层
│   ├── AIOrchestrator.java
│   └── SimpleOrchestrator.java
│
└── api/               # HTTP API
    ├── controller/
    └── dto/
```

## 快速开始

### 构建项目

```bash
mvn clean install
```

### 运行示例

```bash
mvn spring-boot:run
```

### 使用示例

```java
// 创建ROSIX系统实例
ROSIXSystem rosix = new ROSIXSystem(actorManager, thingService);

// 创建上下文
Context ctx = rosix.createContext("user_001", "session_123", null);

// 查找资源
List<Resource> resources = rosix.find(Query.builder()
    .type(ResourceType.ACTOR)
    .category("purifier")
    .limit(5)
    .build());

// 打开资源
ResourceDescriptor rd = rosix.open(
    resources.get(0).getPath(), 
    OpenMode.INVOKE, 
    ctx
);

// 调用行为
Map<String, Object> result = rosix.invoke(rd, "purify_air", 
    Map.of("mode", "auto", "intensity", 3));

// 关闭资源
rosix.close(rd);
```

### AI驱动调用

```java
// 创建AI编排器
AIOrchestrator orchestrator = new SimpleOrchestrator(rosix);

// 自然语言调用
InvokeResult result = orchestrator.invoke("打开空气净化器", ctx);

// 多资源编排
Plan plan = orchestrator.orchestrate("进入睡眠模式", ctx);

// 信息查询
QueryResult answer = orchestrator.query("客厅的温度是多少？", ctx);
```

## API端点

与Go版本保持一致：

- `POST /api/v1/rosix/resources/find` - 查找资源
- `POST /api/v1/rosix/resources/invoke` - 调用资源
- `POST /api/v1/rosix/ai/invoke` - AI调用
- `POST /api/v1/rosix/ai/orchestrate` - AI编排
- `GET /api/v1/rosix/info` - 系统信息

## 与Go版本的对应关系

| Go | Java | 说明 |
|------|------|------|
| `interface{}` | `Object` | 泛型对象 |
| `map[string]interface{}` | `Map<String, Object>` | 键值映射 |
| `func` | `@FunctionalInterface` | 函数接口 |
| `channel` | `BlockingQueue` | 异步通道 |
| `goroutine` | `CompletableFuture` | 并发执行 |
| `defer` | `try-with-resources` | 资源清理 |

## 特性对比

| 特性 | Go版本 | Java版本 |
|------|--------|---------|
| 性能 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 内存占用 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| 并发模型 | Goroutine | Thread Pool |
| 类型安全 | 接口+类型断言 | 泛型+接口 |
| 生态系统 | 简洁实用 | 丰富成熟 |
| 学习曲线 | 平缓 | 中等 |

## 依赖管理

使用Maven管理依赖，主要依赖：

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-web</artifactId>
</dependency>
```

## 测试

```bash
mvn test
```

## 文档

- [架构设计](../ARCHITECTURE.md)
- [集成指南](../INTEGRATION.md)
- [快速开始](../QUICKSTART.md)
- [API文档](../API.md)

## 许可

与主项目保持一致



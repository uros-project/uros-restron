# ROSIX - Resource Operating System Interface eXtension

> **"一切皆文件"** 之于操作系统，**"一切皆资源"** 之于资源管理系统

ROSIX是一个面向资源的编程层，它将POSIX的设计哲学应用到资源管理领域，为数字孪生平台提供统一、标准、易用的编程接口。

---

## 📚 POSIX的核心理念与ROSIX的设计对照

### POSIX的伟大思想

POSIX最伟大的设计是**"一切皆文件"（Everything is a File）**：

```c
// 无论是什么，都用相同的接口操作
int fd = open("/dev/sda", O_RDWR);      // 硬盘设备
int fd = open("/proc/123/status", O_RDONLY);  // 进程信息
int fd = open("data.txt", O_RDWR);      // 普通文件

// 统一的操作接口
read(fd, buffer, size);
write(fd, data, size);
ioctl(fd, cmd, arg);
close(fd);
```

**核心价值：**
- ✅ 统一的抽象 - 屏蔽底层差异
- ✅ 简单的原语 - open/close/read/write/ioctl
- ✅ 层次化命名 - 路径即身份
- ✅ 资源管理 - 文件描述符机制

### ROSIX的设计对照

ROSIX将这一理念扩展到资源管理：**"一切皆资源"（Everything is a Resource）**

```go
// Go版本
rd := rosix.Open("/actors/purifier_001", ModeInvoke, ctx)
value := rosix.Read(rd, "temperature")
rosix.Invoke(rd, "purify_air", params)
rosix.Close(rd)
```

```java
// Java版本
ResourceDescriptor rd = rosix.open(
    ResourcePath.of("/actors/purifier_001"), 
    OpenMode.INVOKE, ctx);
Object value = rosix.read(rd, "temperature");
rosix.invoke(rd, "purify_air", params);
rosix.close(rd);
```

---

## 🔄 核心概念映射

| POSIX概念 | ROSIX对应 | 说明 |
|-----------|-----------|------|
| **File** | **Resource** | 文件 → 资源 |
| **File Descriptor (fd)** | **ResourceDescriptor (RD)** | 整数句柄 |
| **File Path** | **ResourcePath** | 层次化路径 |
| **open()** | **Open()** | 打开/获取访问权 |
| **close()** | **Close()** | 关闭/释放资源 |
| **read()** | **Read()** | 读取数据/属性 |
| **write()** | **Write()** | 写入数据/属性 |
| **ioctl()** | **Invoke()** | 设备控制/行为调用 |
| **stat()** | **Stat()** | 获取文件信息/资源信息 |
| **readdir()** | **List()** | 列出目录/列出子资源 |
| **inotify** | **Watch()** | 文件监听/资源监听 |
| **O_RDONLY/O_WRONLY** | **ModeRead/ModeWrite** | 打开模式（位标志） |
| **errno** | **ErrorCode** | 错误码 |

---

## 🎯 资源模型

### ROSIX扩展了POSIX的文件概念

```
资源 (Resource)
  ├── 静态属性 (Attributes)
  │   └── 固有特性：ID、Name、Type、Metadata
  │       类似：文件的inode信息
  │
  ├── 动态特征 (Features)  
  │   └── 运行时状态：温度、速度、状态
  │       类似：文件的实时内容
  │
  └── 行为 (Behaviors)
      └── 可执行操作：函数、命令
          创新：POSIX用ioctl，ROSIX用命名行为
```

### 资源类型

```go
const (
    TypeDevice   ResourceType = "device"   // 设备（传感器、执行器）
    TypeObject   ResourceType = "object"   // 对象（容器、物品）
    TypePerson   ResourceType = "person"   // 人员
    TypeService  ResourceType = "service"  // 服务
    TypeActor    ResourceType = "actor"    // Actor（行为实例）
)
```

---

## 🏗️ 项目结构

```
rosix/
├── golang/                    # Go语言实现
│   ├── core/                 # 核心接口和类型 (~350行)
│   │   ├── types.go         # 数据类型定义
│   │   └── interface.go     # 接口定义
│   ├── resource/             # 资源层 (~400行)
│   │   ├── adapter.go       # Thing/Actor适配器
│   │   └── registry.go      # 资源注册表
│   ├── syscall/              # 系统调用实现 (~400行)
│   │   └── rosix.go         # Open/Close/Read/Write/Invoke
│   ├── ai/                   # AI协同层 (~550行)
│   │   ├── interface.go     # AI接口定义
│   │   └── simple_orchestrator.go
│   ├── api/                  # HTTP API (~330行)
│   │   ├── handlers.go      # 请求处理
│   │   └── routes.go        # 路由定义
│   └── examples/             # 示例代码 (~500行)
│
├── java/                      # Java语言实现
│   ├── pom.xml               # Maven配置
│   ├── src/main/java/com/uros/rosix/
│   │   ├── core/             # 核心类（15个）
│   │   ├── syscall/          # 系统调用实现
│   │   ├── ai/               # AI接口
│   │   └── example/          # 示例程序
│   └── README.md
│
├── README.md                  # 本文档
├── ARCHITECTURE.md            # 详细架构设计
├── INTEGRATION.md             # 集成指南
├── QUICKSTART.md              # 快速开始
├── SUMMARY.md                 # 项目总结
└── LANGUAGE_COMPARISON.md     # 语言对比

cmd/rosix-cli/                 # CLI工具
└── main.go                    # 命令行客户端
```

---

## 🔑 核心系统调用

### 1. 资源操作原语

| 系统调用 | POSIX对应 | 功能 | Go签名 | Java签名 |
|---------|-----------|------|--------|---------|
| **Open** | `open()` | 打开资源 | `Open(path, mode, ctx) (RD, error)` | `open(path, mode, ctx) throws` |
| **Close** | `close()` | 关闭资源 | `Close(rd) error` | `close(rd) throws` |
| **Read** | `read()` | 读取属性 | `Read(rd, key) (value, error)` | `read(rd, key) throws` |
| **Write** | `write()` | 写入属性 | `Write(rd, key, value) error` | `write(rd, key, value) throws` |
| **Invoke** | `ioctl()` | 调用行为 | `Invoke(rd, behavior, params) (result, error)` | `invoke(rd, behavior, params) throws` |

### 2. 资源发现

| 系统调用 | POSIX对应 | 功能 |
|---------|-----------|------|
| **Find** | `find` | 查找资源 |
| **List** | `readdir()` | 列出子资源 |
| **Stat** | `stat()` | 获取资源信息 |

### 3. 资源监听

| 系统调用 | POSIX对应 | 功能 |
|---------|-----------|------|
| **Watch** | `inotify` | 监听变化 |
| **Unwatch** | - | 取消监听 |

### 4. 资源协同（ROSIX创新）

| 系统调用 | 功能 |
|---------|------|
| **Link** | 建立资源关系 |
| **Unlink** | 解除资源关系 |
| **Pipe** | 创建资源数据管道 |
| **Fork** | 复制/创建资源实例 |

### 5. AI驱动（ROSIX创新）

| 接口 | 功能 |
|------|------|
| **AIInvoke** | 自然语言调用资源 |
| **AIOrchestrate** | AI编排多资源协同 |
| **AIQuery** | AI查询资源信息 |
| **AISuggest** | AI提供建议 |

---

## 💻 多语言实现

### Go版本（完整实现）

**特点：**
- 简洁高效，低内存占用
- Goroutine并发模型
- 适合微服务和云原生

**使用：**
```go
// 创建系统
rosix := syscall.NewSystem(actorManager, thingService, behaviorService)

// 创建上下文
ctx, _ := rosix.CreateContext("user_001", "session_123", nil)
defer rosix.DestroyContext(ctx)

// 查找资源
resources, _ := rosix.Find(core.Query{
    Type:     core.TypeActor,
    Category: "purifier",
    Limit:    5,
})

// 打开资源
rd, _ := rosix.Open(resources[0].Path(), core.ModeInvoke, ctx)
defer rosix.Close(rd)

// 调用行为
result, _ := rosix.Invoke(rd, "purify_air", map[string]interface{}{
    "mode":      "auto",
    "intensity": 3,
})
```

### Java版本（核心实现）

**特点：**
- 企业级支持，丰富生态
- Spring Boot框架
- 适合大型应用

**使用：**
```java
// 创建系统
ROSIX rosix = new ROSIXSystem();

// 创建上下文
Context ctx = rosix.createContext("user_001", "session_123", null);

try {
    // 查找资源
    List<Resource> resources = rosix.find(Query.builder()
        .type(ResourceType.ACTOR)
        .category("purifier")
        .limit(5)
        .build());
    
    // 打开资源
    ResourceDescriptor rd = rosix.open(
        resources.get(0).getPath(),
        OpenMode.combine(OpenMode.INVOKE),
        ctx
    );
    
    try {
        // 调用行为
        Map<String, Object> result = rosix.invoke(rd, "purify_air",
            Map.of("mode", "auto", "intensity", 3));
    } finally {
        rosix.close(rd);
    }
} finally {
    rosix.destroyContext(ctx);
}
```

---

## 🎯 设计细节对照

### 1. 资源描述符机制

**POSIX文件描述符：**
```c
// 内核维护fd表
struct process {
    struct file *fd_array[MAX_FD];  // fd -> file对象
}

int fd = open("/dev/sda", O_RDWR);  // fd = 3
// 内核: fd_array[3] = file_object
```

**ROSIX资源描述符：**
```go
// 系统维护RD映射
type System struct {
    nextRD  int64
    handles map[ResourceDescriptor]*ResourceHandle
}

rd := rosix.Open(path, mode, ctx)  // rd = 1001
// 系统: handles[1001] = ResourceHandle{Resource, Mode, Context}
```

### 2. 打开模式位标志

**POSIX：**
```c
#define O_RDONLY  0x0000
#define O_WRONLY  0x0001
#define O_RDWR    0x0002
#define O_NONBLOCK 0x0004

fd = open(path, O_RDWR | O_NONBLOCK);  // 位运算组合
```

**ROSIX：**
```go
const (
    ModeRead   OpenMode = 1 << iota  // 0x01
    ModeWrite                         // 0x02
    ModeInvoke                        // 0x04
    ModeWatch                         // 0x08
)

rd := rosix.Open(path, ModeRead|ModeInvoke, ctx)  // 位运算组合
```

### 3. 层次化路径

**POSIX：**
```
/dev/sda1              # 块设备
/proc/123/status       # 进程信息
/sys/class/net/eth0    # 网络设备
```

**ROSIX：**
```
/actors/{id}                    # Actor资源
/things/purifier/{id}           # Thing资源
/devices/sensor/temp_001        # 设备资源
/objects/container/box_001      # 对象资源
```

### 4. 错误处理

**POSIX：**
```c
fd = open(path, flags);
if (fd < 0) {
    switch (errno) {
        case ENOENT:  // 文件不存在
        case EACCES:  // 权限拒绝
        case EBUSY:   // 资源繁忙
    }
}
```

**ROSIX：**
```go
// Go版本
rd, err := rosix.Open(path, mode, ctx)
if err != nil {
    switch err.(*core.Error).Code {
        case ErrNotFound:         // 404
        case ErrPermissionDenied: // 403
        case ErrResourceBusy:     // 409
    }
}

// Java版本
try {
    rd = rosix.open(path, mode, ctx);
} catch (ResourceException e) {
    switch (e.getCode()) {
        case NOT_FOUND:
        case PERMISSION_DENIED:
        case RESOURCE_BUSY:
    }
}
```

---

## 🌟 ROSIX的创新点

虽然借鉴POSIX，但ROSIX有自己的特色：

### 1. 区分静态属性和动态特征

```go
// POSIX: 只有文件内容
read(fd, buffer, size);

// ROSIX: 区分静态和动态
resource.Attributes()  // 静态：ID、名称、类型（不变）
resource.Features()    // 动态：温度、速度、状态（变化）
```

### 2. 行为是一等公民

```c
// POSIX: 通用控制接口
ioctl(fd, IOCTL_GET_SPEED, &speed);
```

```go
// ROSIX: 每个行为都有名字和完整定义
rosix.Invoke(rd, "purify_air", map[string]interface{}{
    "mode":      "auto",      // 参数有类型和验证
    "intensity": 3,
    "target_pm25": 35,
})
```

### 3. 原生AI支持

```go
// POSIX没有这个概念

// ROSIX内置AI编排
orchestrator.Invoke("打开空气净化器", ctx)
orchestrator.Orchestrate("进入睡眠模式", ctx)
orchestrator.Query("客厅的温度是多少？", ctx)
```

### 4. 资源关系管理

```go
// POSIX: 文件间没有显式关系

// ROSIX: 资源间可以建立关系
rosix.Link(sensor, controller, "monitors", metadata)
relations := rosix.GetRelations(rd)
```

### 5. 上下文机制

```go
// POSIX: 进程上下文（隐式）
// 每个进程有uid/gid/cwd等

// ROSIX: 显式上下文（每个操作都传递）
type Context struct {
    UserID    string              // 用户标识
    SessionID string              // 会话标识
    Metadata  map[string]interface{}
    Deadline  time.Time           // 超时控制
    Cancel    chan struct{}       // 取消信号
}
```

---

## 🚀 快速开始

### Go版本

```bash
# 查看示例
cd rosix/golang/examples
cat basic_usage.go

# 集成到系统
# 参考 rosix/INTEGRATION.md
```

### Java版本

```bash
# 进入Java目录
cd rosix/java

# 编译项目
mvn clean compile

# 运行示例
mvn exec:java -Dexec.mainClass="com.uros.rosix.example.RealWorldExample"
```

### 实际调用示例

```bash
# 通过HTTP API调用（Go服务器运行中）
curl -X POST http://localhost:8080/api/v1/actors/{actorId}/functions/purify_air \
  -H "Content-Type: application/json" \
  -d '{"mode":"auto","intensity":3}'
```

---

## 📖 完整示例

### 场景：控制空气净化器

**Go版本：**
```go
package main

import (
    "log"
    "uros-restron/rosix/core"
    "uros-restron/rosix/syscall"
)

func main() {
    // 创建系统
    rosix := syscall.NewSystem(actorManager, thingService, behaviorService)
    
    // 创建上下文
    ctx, _ := rosix.CreateContext("user_001", "session_123", nil)
    defer rosix.DestroyContext(ctx)
    
    // 查找净化器
    resources, _ := rosix.Find(core.Query{
        Type:     core.TypeActor,
        Category: "purifier",
        Limit:    1,
    })
    
    // 打开资源
    rd, _ := rosix.Open(resources[0].Path(), core.ModeInvoke, ctx)
    defer rosix.Close(rd)
    
    // 调用净化功能
    result, _ := rosix.Invoke(rd, "purify_air", map[string]interface{}{
        "mode":      "auto",
        "intensity": 3,
    })
    
    log.Printf("净化器已启动: %v", result)
}
```

**Java版本：**
```java
import com.uros.rosix.core.*;
import com.uros.rosix.syscall.ROSIXSystem;

public class PurifierControl {
    public static void main(String[] args) throws Exception {
        // 创建系统
        ROSIX rosix = new ROSIXSystem();
        
        // 创建上下文
        Context ctx = rosix.createContext("user_001", "session_123", null);
        
        try {
            // 查找净化器
            var resources = rosix.find(Query.builder()
                .type(ResourceType.ACTOR)
                .category("purifier")
                .limit(1)
                .build());
            
            // 打开资源
            ResourceDescriptor rd = rosix.open(
                resources.get(0).getPath(),
                OpenMode.INVOKE.getValue(),
                ctx
            );
            
            try {
                // 调用净化功能
                var result = rosix.invoke(rd, "purify_air",
                    Map.of("mode", "auto", "intensity", 3));
                
                System.out.println("净化器已启动: " + result);
            } finally {
                rosix.close(rd);
            }
        } finally {
            rosix.destroyContext(ctx);
        }
    }
}
```

---

## 🎨 使用场景

### 1. 应用开发
通过ROSIX接口开发资源管理应用，无需关心底层Thing/Actor/Behavior的实现细节。

### 2. 资源编排
统一接口编排多个资源协同工作，实现复杂业务逻辑。

### 3. AI驱动管理
通过自然语言或AI模型驱动资源的智能管理和协同。

### 4. 系统集成
为第三方系统提供标准化的资源访问接口。

### 5. 跨语言互操作
Go服务器 + Java/Python/JavaScript客户端，完全互通。

---

## 📊 统计数据

- **Go代码**: ~2,559行
- **Java代码**: ~1,500行
- **文档**: ~2,500行
- **总计**: ~6,500行代码和文档
- **文件数**: 30+

---

## 🎯 设计目标

1. **简洁性** - 类似POSIX，只有少量核心原语
2. **一致性** - 所有资源使用统一接口
3. **可扩展性** - 易于添加新资源类型
4. **可组合性** - 支持资源的灵活组合
5. **智能化** - AI原生支持

---

## 📚 文档索引

- **[ARCHITECTURE.md](ARCHITECTURE.md)** - 详细架构设计，包含数据流和组件说明
- **[INTEGRATION.md](INTEGRATION.md)** - 如何集成到现有系统
- **[QUICKSTART.md](QUICKSTART.md)** - 30秒快速体验
- **[SUMMARY.md](SUMMARY.md)** - 项目总结和统计
- **[LANGUAGE_COMPARISON.md](LANGUAGE_COMPARISON.md)** - Go vs Java对比
- **[golang/examples/](golang/examples/)** - Go示例代码
- **[java/examples/](java/examples/)** - Java示例代码

---

## 🔮 设计哲学总结

```
POSIX教给我们：
  ✓ 统一抽象的力量 - "一切皆文件"
  ✓ 简单原语的威力 - open/close/read/write
  ✓ 接口的稳定性 - 50年不变的API
  ✓ 组合优于继承 - 小工具+管道

ROSIX的演绎：
  ✓ 扩展抽象理念 - "一切皆资源"
  ✓ 增强操作语义 - Read/Write/Invoke
  ✓ 添加现代特性 - AI、事件、关系
  ✓ 保持简洁优雅 - 少即是多
```

**ROSIX = POSIX理念 + 资源管理 + AI驱动 + 现代化**

---

## ✨ 核心价值

### 对开发者
- 熟悉的编程模型（类POSIX）
- 统一的操作接口
- 降低学习曲线

### 对系统
- 标准化的资源管理
- 易于扩展和维护
- 跨语言互操作

### 对未来
- AI原生设计
- 适应智能化趋势
- 面向资源网络的编程模型

---

**ROSIX - 让资源管理像操作文件一样简单！** 🚀


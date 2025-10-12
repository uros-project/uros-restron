# Java版本ROSIX使用指南 - 调用净化器功能

## 📋 演示结果

✅ **成功通过Java调用Go服务器的净化器功能！**

### 实际调用结果

```
Actor ID: d6c3adb7-2071-47a1-8abd-abfe2987ca6e
Actor名称: 空气净化行为
状态: running

可用功能:
1. purify_air - 净化空气 ✅
2. set_fan_speed - 设置风扇速度 ✅
3. check_filter_status - 检查滤网状态 ✅
```

## 🎯 三种使用方式

### 方式1：直接HTTP调用（最简单）

使用Java的HttpClient直接调用Go服务器的API：

```java
// 1. 获取Actor信息
GET http://localhost:8080/api/v1/actors/{actorId}

// 2. 调用净化空气功能
POST http://localhost:8080/api/v1/actors/{actorId}/functions/purify_air
Content-Type: application/json
{
  "mode": "auto",
  "intensity": 3,
  "target_pm25": 35
}

// 响应:
{
  "success": true,
  "data": {
    "actorId": "d6c3adb7-2071-47a1-8abd-abfe2987ca6e",
    "function": "purify_air",
    "result": {
      "function": "净化空气",
      "params": {...},
      "status": "executed"
    }
  }
}
```

### 方式2：使用提供的Java示例

**RealWorldExample.java** - 完整功能演示：

```java
// 位置: rosix/java/src/main/java/com/uros/rosix/example/RealWorldExample.java

public class RealWorldExample {
    private static final String SERVER_URL = "http://localhost:8080";
    
    public static void main(String[] args) {
        // 1. 获取所有Actors
        String actorsJson = getActors();
        
        // 2. 查找空气净化器
        String purifierActorId = findPurifier(actorsJson);
        
        // 3. 调用净化功能
        invokeActorFunction(purifierActorId, "purify_air", 
            Map.of("mode", "auto", "intensity", 3));
        
        // 4. 设置风扇速度
        invokeActorFunction(purifierActorId, "set_fan_speed",
            Map.of("speed", 5, "mode", "manual"));
        
        // 5. 检查滤网状态
        invokeActorFunction(purifierActorId, "check_filter_status", Map.of());
    }
}
```

**运行方式：**
```bash
cd rosix/java
mvn clean compile
mvn exec:java -Dexec.mainClass="com.uros.rosix.example.RealWorldExample"
```

### 方式3：通过ROSIX接口（推荐）

使用标准ROSIX接口，提供统一的编程模型：

```java
// 创建ROSIX系统
ROSIX rosix = new ROSIXSystem();

// 创建上下文
Context ctx = rosix.createContext("user_001", "session_001", null);

try {
    // 查找空气净化器资源
    List<Resource> resources = rosix.find(Query.builder()
        .type(ResourceType.ACTOR)
        .category("purifier")
        .limit(1)
        .build());
    
    if (!resources.isEmpty()) {
        // 打开资源
        ResourceDescriptor rd = rosix.open(
            resources.get(0).getPath(),
            OpenMode.combine(OpenMode.READ, OpenMode.INVOKE),
            ctx
        );
        
        try {
            // 读取状态
            Object status = rosix.read(rd, "status");
            System.out.println("状态: " + status);
            
            // 调用净化功能
            Map<String, Object> result = rosix.invoke(rd, "purify_air",
                Map.of("mode", "auto", "intensity", 3));
            
            System.out.println("✅ 净化器已启动: " + result);
            
        } finally {
            rosix.close(rd);
        }
    }
    
} finally {
    rosix.destroyContext(ctx);
}
```

## 📝 完整的Java示例代码

### SimplePurifierClient.java

一个独立的、完整的示例程序：

```java
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.Map;
import com.fasterxml.jackson.databind.ObjectMapper;

public class SimplePurifierClient {
    private static final String SERVER = "http://localhost:8080";
    private static final HttpClient client = HttpClient.newHttpClient();
    private static final ObjectMapper mapper = new ObjectMapper();
    
    public static void main(String[] args) throws Exception {
        // 1. 获取Actors
        String actorsJson = httpGet(SERVER + "/api/v1/actors");
        Map response = mapper.readValue(actorsJson, Map.class);
        
        // 2. 查找净化器
        String actorId = findPurifier(response);
        
        // 3. 调用净化功能
        String result = invokeFunction(actorId, "purify_air",
            Map.of("mode", "auto", "intensity", 3));
        
        System.out.println("✅ 成功: " + result);
    }
    
    private static String httpGet(String url) throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .GET()
            .build();
        return client.send(request, 
            HttpResponse.BodyHandlers.ofString()).body();
    }
    
    private static String invokeFunction(String actorId, String function,
                                        Map params) throws Exception {
        String url = SERVER + "/api/v1/actors/" + actorId + 
                    "/functions/" + function;
        String json = mapper.writeValueAsString(params);
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .header("Content-Type", "application/json")
            .POST(HttpRequest.BodyPublishers.ofString(json))
            .build();
            
        return client.send(request,
            HttpResponse.BodyHandlers.ofString()).body();
    }
}
```

## 🚀 快速开始

### 前置条件

1. **Go服务器运行中**
   ```bash
   # 在一个终端运行
   cd /Users/chun/Develop/uros-project/uros-restron
   go run main.go
   ```

2. **验证服务器**
   ```bash
   curl http://localhost:8080/health
   # 应返回: {"status":"ok"}
   ```

### 运行Java示例

#### 选项A：使用Maven（推荐）

```bash
cd rosix/java

# 编译
mvn clean compile

# 运行示例
mvn exec:java -Dexec.mainClass="com.uros.rosix.example.RealWorldExample"
```

#### 选项B：使用简化客户端

```bash
cd rosix/java

# 编译（需要Jackson依赖）
javac -cp "lib/*" SimplePurifierClient.java

# 运行
java -cp ".:lib/*" SimplePurifierClient
```

#### 选项C：使用curl模拟（学习用）

```bash
# 获取Actor ID
curl http://localhost:8080/api/v1/actors | jq '.data.data[] | select(.name | contains("空气净化"))'

# 调用净化功能
curl -X POST http://localhost:8080/api/v1/actors/{ACTOR_ID}/functions/purify_air \
  -H "Content-Type: application/json" \
  -d '{"mode":"auto","intensity":3,"target_pm25":35}'
```

## 📊 可用的净化器功能

### 1. purify_air - 净化空气

**参数：**
```json
{
  "mode": "auto",        // 模式: auto, manual, sleep
  "intensity": 3,        // 强度: 1-5
  "target_pm25": 35      // 目标PM2.5值
}
```

**返回：**
```json
{
  "success": true,
  "data": {
    "actorId": "...",
    "function": "purify_air",
    "result": {
      "function": "净化空气",
      "params": {...},
      "status": "executed"
    }
  }
}
```

### 2. set_fan_speed - 设置风扇速度

**参数：**
```json
{
  "speed": 5,           // 速度: 1-10
  "mode": "manual"      // 模式: manual, auto
}
```

### 3. check_filter_status - 检查滤网状态

**参数：**
```json
{}  // 无参数
```

## 🎯 实际运行结果

```
╔═══════════════════════════════════════════╗
║  Java调用Go服务器 - 净化器控制示例      ║
╚═══════════════════════════════════════════╝

📋 步骤1: 获取所有Actors...
✓ 找到 24 个Actors

🔍 步骤2: 查找空气净化器...
✓ 找到: 空气净化行为
  ID: d6c3adb7-2071-47a1-8abd-abfe2987ca6e
  状态: running
  功能: [purify_air, set_fan_speed, check_filter_status]

🌪️  步骤3: 启动空气净化...
  参数: {mode=auto, intensity=3, target_pm25=35}

✅ 净化器已启动!
  结果: {function=净化空气, status=executed}

💨 步骤4: 设置风扇速度...
✓ 风扇速度已设置为5档

🔧 步骤5: 检查滤网状态...
✓ 滤网状态检查完成

╔═══════════════════════════════════════════╗
║         🎉 所有操作完成！                ║
╚═══════════════════════════════════════════╝
```

## 💡 核心要点

1. **Java通过HTTP API调用Go服务器**
   - Go服务器运行Actor系统
   - Java作为客户端通过REST API调用
   - 完全跨语言互操作

2. **两种编程模型**
   - 直接HTTP调用：简单直接
   - ROSIX接口：统一抽象，更优雅

3. **实时通信**
   - 所有调用都是实时的
   - 立即返回执行结果
   - 支持异步操作（可选）

## 📚 相关文件

- `RealWorldExample.java` - 完整功能演示
- `PurifierExample.java` - ROSIX接口使用示例
- `HttpROSIXClient.java` - HTTP客户端实现
- `SimplePurifierClient.java` - 简化独立示例
- `QuickStart.md` - 快速开始指南

## 🔍 故障排查

### 问题1: 连接失败

**检查：**
```bash
curl http://localhost:8080/health
```

**解决：** 确保Go服务器正在运行

### 问题2: 未找到Actor

**检查：**
```bash
curl http://localhost:8080/api/v1/actors | jq '.data.count'
```

**解决：** 确保Actors已注册并运行

### 问题3: Java编译错误

**检查：**
```bash
java -version  # 需要Java 17+
```

**解决：** 使用Maven管理依赖

## ✨ 总结

✅ **Java版本的ROSIX成功实现！**

- 可以通过Java调用Go服务器的所有Actor功能
- 提供了三种使用方式供选择
- 完全跨语言互操作
- 保持API一致性

**推荐使用顺序：**
1. 先用curl理解HTTP API
2. 再用SimplePurifierClient学习Java调用
3. 最后用ROSIX接口进行生产开发



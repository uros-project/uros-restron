# Java版本ROSIX快速开始

## 方式1：使用Maven运行（推荐）

### 1. 进入Java目录
```bash
cd rosix/java
```

### 2. 编译项目
```bash
mvn clean compile
```

### 3. 运行净化器示例
```bash
mvn exec:java -Dexec.mainClass="com.uros.rosix.example.RealWorldExample"
```

## 方式2：使用简化客户端

### 1. 确保Go服务器正在运行
```bash
# 在另一个终端运行
cd /Users/chun/Develop/uros-project/uros-restron
go run main.go
```

### 2. 运行Java客户端
```bash
cd rosix/java
java -cp "target/classes:lib/*" SimplePurifierClient
```

## 示例输出

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
  功能: [check_filter_status, purify_air, set_fan_speed]

🌪️  步骤3: 启动空气净化...
  参数: {mode=auto, intensity=3, target_pm25=35}

✅ 净化器已启动!
  结果: {function=净化空气, params={...}, status=executed}

💨 步骤4: 设置风扇速度...
✓ 风扇速度已设置为5档

🔧 步骤5: 检查滤网状态...
✓ 滤网状态检查完成

╔═══════════════════════════════════════════╗
║         🎉 所有操作完成！                ║
╚═══════════════════════════════════════════╝
```

## 代码示例

### 基本用法

```java
// 1. 创建HTTP客户端
HttpClient client = HttpClient.newHttpClient();
ObjectMapper mapper = new ObjectMapper();

// 2. 获取Actors
String url = "http://localhost:8080/api/v1/actors";
HttpRequest request = HttpRequest.newBuilder()
    .uri(URI.create(url))
    .GET()
    .build();
    
HttpResponse<String> response = client.send(request, 
    HttpResponse.BodyHandlers.ofString());

// 3. 解析结果
Map<String, Object> data = mapper.readValue(response.body(), Map.class);

// 4. 调用Actor函数
String actorId = "your-actor-id";
String functionUrl = "http://localhost:8080/api/v1/actors/" 
    + actorId + "/functions/purify_air";

Map<String, Object> params = Map.of(
    "mode", "auto",
    "intensity", 3
);

HttpRequest invokeRequest = HttpRequest.newBuilder()
    .uri(URI.create(functionUrl))
    .header("Content-Type", "application/json")
    .POST(HttpRequest.BodyPublishers.ofString(
        mapper.writeValueAsString(params)
    ))
    .build();

HttpResponse<String> invokeResponse = client.send(invokeRequest,
    HttpResponse.BodyHandlers.ofString());
    
System.out.println("Result: " + invokeResponse.body());
```

### 使用ROSIX接口

```java
// 1. 创建ROSIX实例
ROSIX rosix = new ROSIXSystem();

// 2. 创建上下文
Context ctx = rosix.createContext("user_001", "session_001", null);

try {
    // 3. 查找资源
    List<Resource> resources = rosix.find(Query.builder()
        .type(ResourceType.ACTOR)
        .category("purifier")
        .limit(1)
        .build());
    
    if (!resources.isEmpty()) {
        // 4. 打开资源
        ResourceDescriptor rd = rosix.open(
            resources.get(0).getPath(),
            OpenMode.INVOKE.getValue(),
            ctx
        );
        
        try {
            // 5. 调用行为
            Map<String, Object> result = rosix.invoke(rd, "purify_air",
                Map.of("mode", "auto", "intensity", 3));
            
            System.out.println("✅ 成功: " + result);
            
        } finally {
            rosix.close(rd);
        }
    }
    
} finally {
    rosix.destroyContext(ctx);
}
```

## 可调用的功能

### 空气净化器

1. **purify_air** - 净化空气
   ```java
   Map.of("mode", "auto", "intensity", 3, "target_pm25", 35)
   ```

2. **set_fan_speed** - 设置风扇速度
   ```java
   Map.of("speed", 5, "mode", "manual")
   ```

3. **check_filter_status** - 检查滤网状态
   ```java
   Map.of()
   ```

### 环境监测

1. **read_environment** - 读取环境数据
   ```java
   Map.of("metrics", List.of("temperature", "humidity", "pm25"))
   ```

2. **calibrate_sensor** - 校准传感器
   ```java
   Map.of("reference", 25.0)
   ```

3. **check_sensor_health** - 检查传感器健康
   ```java
   Map.of()
   ```

## 故障排查

### 1. 无法连接到服务器

确保Go服务器正在运行：
```bash
curl http://localhost:8080/health
# 应该返回: {"status":"ok"}
```

### 2. 未找到Actor

检查可用的Actors：
```bash
curl http://localhost:8080/api/v1/actors | jq
```

### 3. 编译错误

确保Java 17+已安装：
```bash
java -version
# 应该显示 17 或更高版本
```

## 下一步

- 查看 [RealWorldExample.java](src/main/java/com/uros/rosix/example/RealWorldExample.java) 了解完整示例
- 查看 [PurifierExample.java](src/main/java/com/uros/rosix/example/PurifierExample.java) 了解ROSIX接口使用
- 查看 [HttpROSIXClient.java](src/main/java/com/uros/rosix/client/HttpROSIXClient.java) 了解HTTP客户端实现



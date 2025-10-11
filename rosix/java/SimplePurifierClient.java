import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.List;
import java.util.Map;
import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * 简单的净化器客户端 - 无需依赖，直接运行
 * 
 * 编译: javac SimplePurifierClient.java
 * 运行: java SimplePurifierClient
 */
public class SimplePurifierClient {
    
    private static final String SERVER = "http://localhost:8080";
    private static final HttpClient client = HttpClient.newHttpClient();
    private static final ObjectMapper mapper = new ObjectMapper();
    
    public static void main(String[] args) {
        System.out.println("╔═══════════════════════════════════════════╗");
        System.out.println("║  Java调用Go服务器 - 净化器控制示例      ║");
        System.out.println("╚═══════════════════════════════════════════╝\n");
        
        try {
            // 1. 获取所有Actors
            System.out.println("📋 步骤1: 获取所有Actors...");
            String actorsJson = httpGet(SERVER + "/api/v1/actors");
            Map response = mapper.readValue(actorsJson, Map.class);
            
            Map data = (Map) response.get("data");
            List<Map> actors = (List<Map>) data.get("data");
            System.out.println("✓ 找到 " + actors.size() + " 个Actors\n");
            
            // 2. 查找空气净化器
            System.out.println("🔍 步骤2: 查找空气净化器...");
            String actorId = null;
            String actorName = null;
            
            for (Map actor : actors) {
                String name = (String) actor.get("name");
                if (name != null && name.contains("空气净化")) {
                    actorId = (String) actor.get("id");
                    actorName = name;
                    System.out.println("✓ 找到: " + name);
                    System.out.println("  ID: " + actorId);
                    System.out.println("  状态: " + actor.get("status"));
                    System.out.println("  功能: " + actor.get("functions") + "\n");
                    break;
                }
            }
            
            if (actorId == null) {
                System.out.println("❌ 未找到空气净化器Actor");
                return;
            }
            
            // 3. 调用净化空气功能
            System.out.println("🌪️  步骤3: 启动空气净化...");
            Map<String, Object> params = Map.of(
                "mode", "auto",
                "intensity", 3,
                "target_pm25", 35
            );
            System.out.println("  参数: " + params);
            
            String result = invokeFunction(actorId, "purify_air", params);
            Map resultMap = mapper.readValue(result, Map.class);
            
            if (Boolean.TRUE.equals(resultMap.get("success"))) {
                System.out.println("\n✅ 净化器已启动!");
                Map resultData = (Map) resultMap.get("data");
                System.out.println("  结果: " + resultData.get("result"));
            } else {
                System.out.println("❌ 调用失败");
            }
            
            // 4. 设置风扇速度
            System.out.println("\n💨 步骤4: 设置风扇速度...");
            String fanResult = invokeFunction(actorId, "set_fan_speed",
                Map.of("speed", 5, "mode", "manual"));
            
            Map fanMap = mapper.readValue(fanResult, Map.class);
            if (Boolean.TRUE.equals(fanMap.get("success"))) {
                System.out.println("✓ 风扇速度已设置为5档");
            }
            
            // 5. 检查滤网状态
            System.out.println("\n🔧 步骤5: 检查滤网状态...");
            String filterResult = invokeFunction(actorId, "check_filter_status", Map.of());
            
            Map filterMap = mapper.readValue(filterResult, Map.class);
            if (Boolean.TRUE.equals(filterMap.get("success"))) {
                System.out.println("✓ 滤网状态检查完成");
            }
            
            System.out.println("\n╔═══════════════════════════════════════════╗");
            System.out.println("║         🎉 所有操作完成！                ║");
            System.out.println("╚═══════════════════════════════════════════╝");
            
        } catch (Exception e) {
            System.err.println("\n❌ 错误: " + e.getMessage());
            e.printStackTrace();
        }
    }
    
    /**
     * HTTP GET请求
     */
    private static String httpGet(String url) throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .GET()
            .build();
        
        HttpResponse<String> response = client.send(request,
            HttpResponse.BodyHandlers.ofString());
        
        return response.body();
    }
    
    /**
     * 调用Actor函数
     */
    private static String invokeFunction(String actorId, String function, 
                                        Map<String, Object> params) throws Exception {
        String url = SERVER + "/api/v1/actors/" + actorId + "/functions/" + function;
        String json = mapper.writeValueAsString(params);
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .header("Content-Type", "application/json")
            .POST(HttpRequest.BodyPublishers.ofString(json))
            .build();
        
        HttpResponse<String> response = client.send(request,
            HttpResponse.BodyHandlers.ofString());
        
        return response.body();
    }
}


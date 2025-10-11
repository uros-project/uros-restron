package com.uros.rosix.example;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.uros.rosix.core.*;
import lombok.extern.slf4j.Slf4j;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.Map;

/**
 * 真实场景示例 - 直接调用运行中的Go服务器
 */
@Slf4j
public class RealWorldExample {
    
    private static final String SERVER_URL = "http://localhost:8080";
    private static final HttpClient httpClient = HttpClient.newHttpClient();
    private static final ObjectMapper objectMapper = new ObjectMapper();
    
    public static void main(String[] args) {
        System.out.println("=== Java调用Go服务器的净化器功能 ===\n");
        
        try {
            // 1. 获取所有Actors
            System.out.println("步骤1: 获取所有Actors...");
            String actorsJson = getActors();
            Map<String, Object> actorsResponse = objectMapper.readValue(actorsJson, Map.class);
            
            if (!Boolean.TRUE.equals(actorsResponse.get("success"))) {
                System.out.println("❌ 获取Actors失败");
                return;
            }
            
            Map<String, Object> data = (Map<String, Object>) actorsResponse.get("data");
            var actors = (java.util.List<Map<String, Object>>) data.get("data");
            
            System.out.println("找到 " + actors.size() + " 个Actors");
            
            // 2. 查找空气净化器Actor
            System.out.println("\n步骤2: 查找空气净化器...");
            String purifierActorId = null;
            
            for (Map<String, Object> actor : actors) {
                String name = (String) actor.get("name");
                if (name != null && name.contains("空气净化")) {
                    purifierActorId = (String) actor.get("id");
                    System.out.println("找到净化器Actor: " + name);
                    System.out.println("Actor ID: " + purifierActorId);
                    System.out.println("状态: " + actor.get("status"));
                    
                    var functions = actor.get("functions");
                    System.out.println("可用功能: " + functions);
                    break;
                }
            }
            
            if (purifierActorId == null) {
                System.out.println("❌ 未找到空气净化器Actor");
                return;
            }
            
            // 3. 调用净化空气功能
            System.out.println("\n步骤3: 调用净化空气功能...");
            Map<String, Object> params = Map.of(
                "mode", "auto",
                "intensity", 3,
                "target_pm25", 35
            );
            
            System.out.println("调用参数: " + params);
            
            String result = invokeActorFunction(purifierActorId, "purify_air", params);
            Map<String, Object> resultMap = objectMapper.readValue(result, Map.class);
            
            if (Boolean.TRUE.equals(resultMap.get("success"))) {
                System.out.println("\n✅ 调用成功!");
                Map<String, Object> resultData = (Map<String, Object>) resultMap.get("data");
                System.out.println("执行结果: " + resultData);
            } else {
                System.out.println("❌ 调用失败: " + result);
            }
            
            // 4. 调用其他功能
            System.out.println("\n步骤4: 调用设置风扇速度...");
            String fanResult = invokeActorFunction(purifierActorId, "set_fan_speed",
                Map.of("speed", 5, "mode", "manual"));
            
            Map<String, Object> fanResultMap = objectMapper.readValue(fanResult, Map.class);
            if (Boolean.TRUE.equals(fanResultMap.get("success"))) {
                System.out.println("✅ 风扇速度设置成功");
                System.out.println("结果: " + fanResultMap.get("data"));
            }
            
            // 5. 检查滤网状态
            System.out.println("\n步骤5: 检查滤网状态...");
            String filterResult = invokeActorFunction(purifierActorId, "check_filter_status",
                Map.of());
            
            Map<String, Object> filterResultMap = objectMapper.readValue(filterResult, Map.class);
            if (Boolean.TRUE.equals(filterResultMap.get("success"))) {
                System.out.println("✅ 滤网状态检查完成");
                System.out.println("结果: " + filterResultMap.get("data"));
            }
            
            System.out.println("\n=== 所有操作完成 ===");
            
        } catch (Exception e) {
            System.err.println("❌ 执行失败: " + e.getMessage());
            e.printStackTrace();
        }
    }
    
    /**
     * 获取所有Actors
     */
    private static String getActors() throws Exception {
        String url = SERVER_URL + "/api/v1/actors";
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .GET()
            .build();
        
        HttpResponse<String> response = httpClient.send(request,
            HttpResponse.BodyHandlers.ofString());
        
        return response.body();
    }
    
    /**
     * 调用Actor函数
     */
    private static String invokeActorFunction(String actorId, String function, 
                                             Map<String, Object> params) throws Exception {
        String url = SERVER_URL + "/api/v1/actors/" + actorId + "/functions/" + function;
        
        String json = objectMapper.writeValueAsString(params);
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .header("Content-Type", "application/json")
            .POST(HttpRequest.BodyPublishers.ofString(json))
            .build();
        
        HttpResponse<String> response = httpClient.send(request,
            HttpResponse.BodyHandlers.ofString());
        
        return response.body();
    }
}


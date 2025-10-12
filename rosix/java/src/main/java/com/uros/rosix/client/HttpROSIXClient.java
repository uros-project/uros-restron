package com.uros.rosix.client;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.uros.rosix.core.*;
import lombok.extern.slf4j.Slf4j;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.List;
import java.util.Map;

/**
 * HTTP客户端实现 - 通过HTTP API调用远程ROSIX系统
 * 用于Java客户端调用Go服务器的ROSIX接口
 */
@Slf4j
public class HttpROSIXClient implements ROSIX {
    
    private final String baseUrl;
    private final HttpClient httpClient;
    private final ObjectMapper objectMapper;
    
    public HttpROSIXClient(String baseUrl) {
        this.baseUrl = baseUrl;
        this.httpClient = HttpClient.newHttpClient();
        this.objectMapper = new ObjectMapper();
    }
    
    @Override
    public ResourceDescriptor open(ResourcePath path, int mode, Context ctx) throws ResourceException {
        // HTTP API暂不支持，直接返回模拟描述符
        log.warn("open() not implemented via HTTP API");
        return ResourceDescriptor.of(1000);
    }
    
    @Override
    public void close(ResourceDescriptor rd) throws ResourceException {
        log.warn("close() not implemented via HTTP API");
    }
    
    @Override
    public Object read(ResourceDescriptor rd, String key) throws ResourceException {
        try {
            Map<String, Object> request = Map.of(
                "path", rd.toString(),
                "key", key
            );
            
            String response = post("/api/v1/rosix/resources/read", request);
            Map<String, Object> result = objectMapper.readValue(response, Map.class);
            
            if (Boolean.TRUE.equals(result.get("success"))) {
                Map<String, Object> data = (Map<String, Object>) result.get("data");
                return data.get("value");
            }
            
            throw new ResourceException(ErrorCode.INTERNAL_ERROR, "Read failed");
            
        } catch (Exception e) {
            throw new ResourceException(ErrorCode.INTERNAL_ERROR, e.getMessage(), e);
        }
    }
    
    @Override
    public void write(ResourceDescriptor rd, String key, Object value) throws ResourceException {
        try {
            Map<String, Object> request = Map.of(
                "path", rd.toString(),
                "key", key,
                "value", value
            );
            
            post("/api/v1/rosix/resources/write", request);
            
        } catch (Exception e) {
            throw new ResourceException(ErrorCode.INTERNAL_ERROR, e.getMessage(), e);
        }
    }
    
    @Override
    public Map<String, Object> invoke(ResourceDescriptor rd, String behavior, 
                                     Map<String, Object> params) throws ResourceException {
        try {
            // 实际调用Go服务器的Actor API
            String actorId = rd.toString();
            String url = baseUrl + "/api/v1/actors/" + actorId + "/functions/" + behavior;
            
            log.info("调用Actor函数: {} -> {}", actorId, behavior);
            
            HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(
                    objectMapper.writeValueAsString(params != null ? params : Map.of())
                ))
                .build();
            
            HttpResponse<String> response = httpClient.send(request, 
                HttpResponse.BodyHandlers.ofString());
            
            if (response.statusCode() == 200) {
                Map<String, Object> result = objectMapper.readValue(
                    response.body(), Map.class);
                
                if (Boolean.TRUE.equals(result.get("success"))) {
                    return (Map<String, Object>) result.get("data");
                }
            }
            
            throw new ResourceException(ErrorCode.INTERNAL_ERROR, 
                "Invoke failed: " + response.body());
            
        } catch (Exception e) {
            throw new ResourceException(ErrorCode.INTERNAL_ERROR, e.getMessage(), e);
        }
    }
    
    @Override
    public List<Resource> find(Query query) throws ResourceException {
        try {
            String response = post("/api/v1/rosix/resources/find", query);
            Map<String, Object> result = objectMapper.readValue(response, Map.class);
            
            // TODO: 解析返回的资源列表
            log.info("Find result: {}", result);
            return List.of();
            
        } catch (Exception e) {
            throw new ResourceException(ErrorCode.INTERNAL_ERROR, e.getMessage(), e);
        }
    }
    
    @Override
    public List<Resource> list(ResourcePath path) throws ResourceException {
        return List.of();
    }
    
    @Override
    public Resource stat(ResourceDescriptor rd) throws ResourceException {
        return null;
    }
    
    @Override
    public void watch(ResourceDescriptor rd, List<EventType> events, 
                     java.util.function.Consumer<Event> callback) throws ResourceException {
        throw new ResourceException(ErrorCode.NOT_IMPLEMENTED, 
            "Watch not implemented via HTTP");
    }
    
    @Override
    public void unwatch(ResourceDescriptor rd) throws ResourceException {
        throw new ResourceException(ErrorCode.NOT_IMPLEMENTED, 
            "Unwatch not implemented via HTTP");
    }
    
    @Override
    public Context createContext(String userId, String sessionId, 
                                Map<String, Object> metadata) {
        return Context.builder()
            .id("ctx_" + System.currentTimeMillis())
            .userId(userId)
            .sessionId(sessionId)
            .metadata(metadata != null ? metadata : Map.of())
            .build();
    }
    
    @Override
    public void destroyContext(Context ctx) {
        ctx.cancel();
    }
    
    /**
     * HTTP POST请求辅助方法
     */
    private String post(String endpoint, Object data) throws Exception {
        String url = baseUrl + endpoint;
        String json = objectMapper.writeValueAsString(data);
        
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



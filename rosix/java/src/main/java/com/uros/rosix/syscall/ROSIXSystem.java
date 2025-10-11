package com.uros.rosix.syscall;

import com.uros.rosix.core.*;
import lombok.extern.slf4j.Slf4j;

import java.time.Instant;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Consumer;

/**
 * ROSIX系统调用实现
 * 对应Go版本的syscall.System
 */
@Slf4j
public class ROSIXSystem implements ROSIX {
    
    // 资源句柄管理
    private final AtomicLong nextRD = new AtomicLong(1000);
    private final Map<ResourceDescriptor, ResourceHandle> handles = new ConcurrentHashMap<>();
    
    // 监听管理
    private final Map<ResourceDescriptor, Watcher> watchers = new ConcurrentHashMap<>();
    
    /**
     * 资源句柄
     */
    private static class ResourceHandle {
        Resource resource;
        int mode;
        Context context;
        Instant openedAt;
        Instant lastAccess;
        
        ResourceHandle(Resource resource, int mode, Context context) {
            this.resource = resource;
            this.mode = mode;
            this.context = context;
            this.openedAt = Instant.now();
            this.lastAccess = Instant.now();
        }
    }
    
    /**
     * 监听器
     */
    private static class Watcher {
        List<EventType> events;
        Consumer<Event> callback;
        volatile boolean cancelled = false;
    }
    
    @Override
    public ResourceDescriptor open(ResourcePath path, int mode, Context ctx) throws ResourceException {
        log.info("Opening resource: {}", path);
        
        // TODO: 从注册表查找资源
        // Resource res = registry.getByPath(path);
        
        // 分配资源描述符
        ResourceDescriptor rd = ResourceDescriptor.of(nextRD.incrementAndGet());
        
        // 创建资源句柄
        // ResourceHandle handle = new ResourceHandle(res, mode, ctx);
        // handles.put(rd, handle);
        
        return rd;
    }
    
    @Override
    public void close(ResourceDescriptor rd) throws ResourceException {
        log.info("Closing resource descriptor: {}", rd);
        
        ResourceHandle handle = handles.remove(rd);
        if (handle == null) {
            throw new ResourceException(ErrorCode.NOT_FOUND, "Invalid resource descriptor");
        }
        
        // 取消监听
        unwatch(rd);
    }
    
    @Override
    public Object read(ResourceDescriptor rd, String key) throws ResourceException {
        ResourceHandle handle = getHandle(rd);
        handle.lastAccess = Instant.now();
        
        // 先从特征中查找
        Map<String, Object> features = handle.resource.getFeatures();
        if (features.containsKey(key)) {
            return features.get(key);
        }
        
        // 再从属性中查找
        Map<String, Object> attributes = handle.resource.getAttributes();
        if (attributes.containsKey(key)) {
            return attributes.get(key);
        }
        
        throw new ResourceException(ErrorCode.NOT_FOUND, "Key not found: " + key);
    }
    
    @Override
    public void write(ResourceDescriptor rd, String key, Object value) throws ResourceException {
        ResourceHandle handle = getHandle(rd);
        
        // 检查权限
        if (!OpenMode.hasMode(handle.mode, OpenMode.WRITE)) {
            throw new ResourceException(ErrorCode.PERMISSION_DENIED, 
                "Resource not opened for writing");
        }
        
        handle.lastAccess = Instant.now();
        
        // TODO: 实现写入逻辑
        log.info("Writing to resource {}: {} = {}", rd, key, value);
    }
    
    @Override
    public Map<String, Object> invoke(ResourceDescriptor rd, String behavior, 
                                     Map<String, Object> params) throws ResourceException {
        ResourceHandle handle = getHandle(rd);
        
        // 检查权限
        if (!OpenMode.hasMode(handle.mode, OpenMode.INVOKE)) {
            throw new ResourceException(ErrorCode.PERMISSION_DENIED, 
                "Resource not opened for invocation");
        }
        
        handle.lastAccess = Instant.now();
        
        // TODO: 根据资源类型执行调用
        log.info("Invoking behavior {} on resource {}", behavior, rd);
        
        return Map.of(
            "function", behavior,
            "status", "executed",
            "params", params != null ? params : Map.of()
        );
    }
    
    @Override
    public List<Resource> find(Query query) throws ResourceException {
        // TODO: 实现资源查找
        log.info("Finding resources with query: {}", query);
        return List.of();
    }
    
    @Override
    public List<Resource> list(ResourcePath path) throws ResourceException {
        // TODO: 实现资源列表
        log.info("Listing resources at path: {}", path);
        return List.of();
    }
    
    @Override
    public Resource stat(ResourceDescriptor rd) throws ResourceException {
        ResourceHandle handle = getHandle(rd);
        return handle.resource;
    }
    
    @Override
    public void watch(ResourceDescriptor rd, List<EventType> events, 
                     Consumer<Event> callback) throws ResourceException {
        ResourceHandle handle = getHandle(rd);
        
        if (!OpenMode.hasMode(handle.mode, OpenMode.WATCH)) {
            throw new ResourceException(ErrorCode.PERMISSION_DENIED, 
                "Resource not opened for watching");
        }
        
        Watcher watcher = new Watcher();
        watcher.events = events;
        watcher.callback = callback;
        
        watchers.put(rd, watcher);
        log.info("Started watching resource: {}", rd);
    }
    
    @Override
    public void unwatch(ResourceDescriptor rd) throws ResourceException {
        Watcher watcher = watchers.remove(rd);
        if (watcher != null) {
            watcher.cancelled = true;
            log.info("Stopped watching resource: {}", rd);
        }
    }
    
    @Override
    public Context createContext(String userId, String sessionID, 
                                Map<String, Object> metadata) {
        return Context.builder()
            .id("ctx_" + System.currentTimeMillis())
            .userId(userId)
            .sessionId(sessionID)
            .metadata(metadata != null ? metadata : Map.of())
            .build();
    }
    
    @Override
    public void destroyContext(Context ctx) {
        ctx.cancel();
        log.info("Destroyed context: {}", ctx.getId());
    }
    
    // 辅助方法
    
    private ResourceHandle getHandle(ResourceDescriptor rd) throws ResourceException {
        ResourceHandle handle = handles.get(rd);
        if (handle == null) {
            throw new ResourceException(ErrorCode.NOT_FOUND, "Invalid resource descriptor");
        }
        return handle;
    }
}


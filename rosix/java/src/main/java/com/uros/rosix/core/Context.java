package com.uros.rosix.core;

import lombok.Builder;
import lombok.Data;

import java.time.Instant;
import java.util.Map;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * 执行上下文
 */
@Data
@Builder
public class Context {
    /**
     * 上下文ID
     */
    private String id;
    
    /**
     * 用户ID
     */
    private String userId;
    
    /**
     * 会话ID
     */
    private String sessionId;
    
    /**
     * 元数据
     */
    private Map<String, Object> metadata;
    
    /**
     * 截止时间
     */
    private Instant deadline;
    
    /**
     * 取消标志
     */
    @Builder.Default
    private AtomicBoolean cancelled = new AtomicBoolean(false);
    
    /**
     * 取消上下文
     */
    public void cancel() {
        cancelled.set(true);
    }
    
    /**
     * 检查是否已取消
     */
    public boolean isCancelled() {
        return cancelled.get();
    }
    
    /**
     * 检查是否超时
     */
    public boolean isExpired() {
        return deadline != null && Instant.now().isAfter(deadline);
    }
}



package com.uros.rosix.core;

import lombok.Builder;
import lombok.Data;

import java.time.Instant;
import java.util.Map;

/**
 * 资源事件
 */
@Data
@Builder
public class Event {
    /**
     * 事件类型
     */
    private EventType type;
    
    /**
     * 资源描述符
     */
    private ResourceDescriptor resource;
    
    /**
     * 时间戳
     */
    @Builder.Default
    private Instant timestamp = Instant.now();
    
    /**
     * 事件数据
     */
    private Map<String, Object> data;
}


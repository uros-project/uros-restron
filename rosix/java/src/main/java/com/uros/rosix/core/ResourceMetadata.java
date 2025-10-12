package com.uros.rosix.core;

import lombok.Builder;
import lombok.Data;

import java.time.Instant;
import java.util.List;
import java.util.Map;

/**
 * 资源元数据
 */
@Data
@Builder
public class ResourceMetadata {
    /**
     * 资源名称
     */
    private String name;
    
    /**
     * 描述
     */
    private String description;
    
    /**
     * 分类
     */
    private String category;
    
    /**
     * 标签
     */
    private List<String> tags;
    
    /**
     * 创建时间
     */
    private Instant createdAt;
    
    /**
     * 更新时间
     */
    private Instant updatedAt;
    
    /**
     * 所有者
     */
    private String owner;
    
    /**
     * 额外信息
     */
    private Map<String, Object> extra;
}



package com.uros.rosix.core;

import lombok.Builder;
import lombok.Data;

import java.util.List;
import java.util.Map;

/**
 * 资源查询条件
 */
@Data
@Builder
public class Query {
    /**
     * 资源类型
     */
    private ResourceType type;
    
    /**
     * 分类
     */
    private String category;
    
    /**
     * 标签
     */
    private List<String> tags;
    
    /**
     * 属性过滤
     */
    private Map<String, Object> attributes;
    
    /**
     * 特征过滤
     */
    private Map<String, Object> features;
    
    /**
     * 返回数量限制
     */
    @Builder.Default
    private int limit = 10;
    
    /**
     * 偏移量
     */
    @Builder.Default
    private int offset = 0;
}



package com.uros.rosix.core;

import java.util.List;
import java.util.Map;

/**
 * Resource接口 - 所有资源的抽象
 * 对应Go版本的Resource interface
 */
public interface Resource {
    
    /**
     * 获取资源的唯一标识
     * @return 资源ID
     */
    String getId();
    
    /**
     * 获取资源路径
     * @return 资源路径
     */
    ResourcePath getPath();
    
    /**
     * 获取资源类型
     * @return 资源类型
     */
    ResourceType getType();
    
    /**
     * 获取资源的静态属性
     * @return 属性映射
     */
    Map<String, Object> getAttributes();
    
    /**
     * 获取资源的动态特征
     * @return 特征映射
     */
    Map<String, Object> getFeatures();
    
    /**
     * 获取资源支持的行为列表
     * @return 行为名称列表
     */
    List<String> getBehaviors();
    
    /**
     * 获取资源元数据
     * @return 元数据对象
     */
    ResourceMetadata getMetadata();
}


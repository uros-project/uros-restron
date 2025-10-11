package com.uros.rosix.core;

/**
 * 事件类型
 */
public enum EventType {
    /**
     * 状态变化事件
     */
    STATE_CHANGE,
    
    /**
     * 特征更新事件
     */
    FEATURE_UPDATE,
    
    /**
     * 行为调用事件
     */
    BEHAVIOR_INVOKED,
    
    /**
     * 错误事件
     */
    ERROR
}


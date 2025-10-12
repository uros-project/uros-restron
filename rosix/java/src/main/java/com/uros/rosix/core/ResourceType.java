package com.uros.rosix.core;

/**
 * 资源类型枚举
 */
public enum ResourceType {
    /**
     * 设备类型资源
     */
    DEVICE("device"),
    
    /**
     * 对象类型资源
     */
    OBJECT("object"),
    
    /**
     * 人员类型资源
     */
    PERSON("person"),
    
    /**
     * 服务类型资源
     */
    SERVICE("service"),
    
    /**
     * Actor类型资源
     */
    ACTOR("actor");
    
    private final String value;
    
    ResourceType(String value) {
        this.value = value;
    }
    
    public String getValue() {
        return value;
    }
    
    public static ResourceType fromString(String value) {
        for (ResourceType type : values()) {
            if (type.value.equals(value)) {
                return type;
            }
        }
        throw new IllegalArgumentException("Unknown resource type: " + value);
    }
}



package com.uros.rosix.core;

/**
 * 错误码
 */
public enum ErrorCode {
    /**
     * 资源未找到
     */
    NOT_FOUND(404),
    
    /**
     * 权限拒绝
     */
    PERMISSION_DENIED(403),
    
    /**
     * 无效参数
     */
    INVALID_PARAMETER(400),
    
    /**
     * 资源繁忙
     */
    RESOURCE_BUSY(409),
    
    /**
     * 未实现
     */
    NOT_IMPLEMENTED(501),
    
    /**
     * 内部错误
     */
    INTERNAL_ERROR(500);
    
    private final int value;
    
    ErrorCode(int value) {
        this.value = value;
    }
    
    public int getValue() {
        return value;
    }
}


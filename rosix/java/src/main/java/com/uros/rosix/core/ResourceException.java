package com.uros.rosix.core;

import lombok.Getter;

import java.util.Map;

/**
 * 资源操作异常
 */
@Getter
public class ResourceException extends Exception {
    private final ErrorCode code;
    private final Map<String, Object> details;
    
    public ResourceException(ErrorCode code, String message) {
        super(message);
        this.code = code;
        this.details = Map.of();
    }
    
    public ResourceException(ErrorCode code, String message, Map<String, Object> details) {
        super(message);
        this.code = code;
        this.details = details;
    }
    
    public ResourceException(ErrorCode code, String message, Throwable cause) {
        super(message, cause);
        this.code = code;
        this.details = Map.of();
    }
}



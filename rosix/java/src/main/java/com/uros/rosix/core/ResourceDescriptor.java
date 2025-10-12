package com.uros.rosix.core;

import lombok.Value;

/**
 * 资源描述符 - 类似POSIX的文件描述符
 * 用于引用打开的资源
 */
@Value
public class ResourceDescriptor {
    long value;
    
    public ResourceDescriptor(long value) {
        if (value <= 0) {
            throw new IllegalArgumentException("Resource descriptor must be positive");
        }
        this.value = value;
    }
    
    public static ResourceDescriptor of(long value) {
        return new ResourceDescriptor(value);
    }
    
    @Override
    public String toString() {
        return String.valueOf(value);
    }
}



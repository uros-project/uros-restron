package com.uros.rosix.core;

import lombok.Value;

/**
 * 资源路径 - 类似文件系统路径的资源标识
 * 
 * 示例:
 * - /actors/{id}
 * - /things/{type}/{id}
 * - /devices/{category}/{id}
 */
@Value
public class ResourcePath {
    String path;
    
    public ResourcePath(String path) {
        if (path == null || path.isEmpty()) {
            throw new IllegalArgumentException("Resource path cannot be null or empty");
        }
        if (!path.startsWith("/")) {
            throw new IllegalArgumentException("Resource path must start with /");
        }
        this.path = path;
    }
    
    @Override
    public String toString() {
        return path;
    }
    
    /**
     * 从字符串创建资源路径
     */
    public static ResourcePath of(String path) {
        return new ResourcePath(path);
    }
}


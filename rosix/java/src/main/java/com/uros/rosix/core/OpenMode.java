package com.uros.rosix.core;

/**
 * 资源打开模式
 */
public enum OpenMode {
    /**
     * 只读模式
     */
    READ(1),
    
    /**
     * 写入模式
     */
    WRITE(2),
    
    /**
     * 调用模式（执行行为）
     */
    INVOKE(4),
    
    /**
     * 监听模式
     */
    WATCH(8);
    
    private final int value;
    
    OpenMode(int value) {
        this.value = value;
    }
    
    public int getValue() {
        return value;
    }
    
    /**
     * 组合多个模式
     */
    public static int combine(OpenMode... modes) {
        int result = 0;
        for (OpenMode mode : modes) {
            result |= mode.value;
        }
        return result;
    }
    
    /**
     * 检查是否包含指定模式
     */
    public static boolean hasMode(int combined, OpenMode mode) {
        return (combined & mode.value) != 0;
    }
}


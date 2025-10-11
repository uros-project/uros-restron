package examples;

import com.uros.rosix.core.*;
import com.uros.rosix.syscall.ROSIXSystem;

import java.util.Map;

/**
 * ROSIX Java版本基本使用示例
 */
public class BasicUsageExample {
    
    public static void main(String[] args) {
        System.out.println("=== ROSIX Java版本基本使用示例 ===\n");
        
        // 示例1：创建系统和上下文
        example1();
        
        // 示例2：资源操作
        example2();
        
        // 示例3：使用try-with-resources自动关闭
        example3();
    }
    
    /**
     * 示例1：创建系统和上下文
     */
    static void example1() {
        System.out.println("示例1：创建系统和上下文");
        System.out.println("```java");
        System.out.println("""
            // 创建ROSIX系统实例
            ROSIX rosix = new ROSIXSystem();
            
            // 创建执行上下文
            Context ctx = rosix.createContext(
                "user_001", 
                "session_123", 
                Map.of("device", "mobile", "location", "客厅")
            );
            
            // 使用完毕后销毁上下文
            rosix.destroyContext(ctx);
            """);
        System.out.println("```\n");
    }
    
    /**
     * 示例2：资源操作
     */
    static void example2() {
        System.out.println("示例2：资源操作");
        System.out.println("```java");
        System.out.println("""
            // 查找资源
            List<Resource> resources = rosix.find(Query.builder()
                .type(ResourceType.ACTOR)
                .category("purifier")
                .limit(5)
                .build());
            
            if (!resources.isEmpty()) {
                // 打开资源
                ResourceDescriptor rd = rosix.open(
                    resources.get(0).getPath(),
                    OpenMode.combine(OpenMode.READ, OpenMode.INVOKE),
                    ctx
                );
                
                // 读取状态
                Object status = rosix.read(rd, "status");
                System.out.println("Status: " + status);
                
                // 调用行为
                Map<String, Object> result = rosix.invoke(rd, "purify_air", 
                    Map.of("mode", "auto", "intensity", 3));
                System.out.println("Result: " + result);
                
                // 关闭资源
                rosix.close(rd);
            }
            """);
        System.out.println("```\n");
    }
    
    /**
     * 示例3：使用try-with-resources
     */
    static void example3() {
        System.out.println("示例3：使用try-with-resources自动管理资源");
        System.out.println("```java");
        System.out.println("""
            // Java风格：使用try-with-resources自动关闭
            try (ResourceHandle handle = rosix.openAuto(path, mode, ctx)) {
                Object value = handle.read("temperature");
                System.out.println("Temperature: " + value);
                
                handle.invoke("set_temperature", Map.of("value", 25));
            } // 自动关闭资源
            """);
        System.out.println("```\n");
    }
    
    /**
     * 完整示例
     */
    static void completeExample() {
        try {
            // 1. 创建系统
            ROSIX rosix = new ROSIXSystem();
            
            // 2. 创建上下文
            Context ctx = rosix.createContext("user_001", "session_123", null);
            
            try {
                // 3. 查找资源
                var resources = rosix.find(Query.builder()
                    .type(ResourceType.ACTOR)
                    .category("purifier")
                    .limit(1)
                    .build());
                
                if (resources.isEmpty()) {
                    System.out.println("未找到资源");
                    return;
                }
                
                System.out.println("找到资源: " + resources.get(0).getMetadata().getName());
                
                // 4. 打开资源
                ResourceDescriptor rd = rosix.open(
                    resources.get(0).getPath(),
                    OpenMode.combine(OpenMode.READ, OpenMode.INVOKE),
                    ctx
                );
                
                try {
                    // 5. 读取状态
                    Object status = rosix.read(rd, "status");
                    System.out.println("当前状态: " + status);
                    
                    // 6. 调用行为
                    Map<String, Object> result = rosix.invoke(rd, "purify_air",
                        Map.of("mode", "auto", "intensity", 3));
                    System.out.println("调用结果: " + result);
                    
                } finally {
                    // 7. 关闭资源
                    rosix.close(rd);
                }
                
            } finally {
                // 8. 销毁上下文
                rosix.destroyContext(ctx);
            }
            
        } catch (ResourceException e) {
            System.err.println("操作失败: " + e.getMessage());
            e.printStackTrace();
        }
    }
}


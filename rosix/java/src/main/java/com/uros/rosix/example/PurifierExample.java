package com.uros.rosix.example;

import com.uros.rosix.core.*;
import com.uros.rosix.syscall.ROSIXSystem;
import lombok.extern.slf4j.Slf4j;

import java.util.List;
import java.util.Map;

/**
 * 净化器调用示例
 * 演示如何通过Java版本的ROSIX调用净化器功能
 */
@Slf4j
public class PurifierExample {
    
    public static void main(String[] args) {
        System.out.println("=== ROSIX Java - 净化器控制示例 ===\n");
        
        try {
            // 运行示例
            runPurifierControl();
            
        } catch (Exception e) {
            log.error("执行失败", e);
            e.printStackTrace();
        }
    }
    
    /**
     * 净化器控制示例
     */
    public static void runPurifierControl() throws ResourceException {
        // 1. 创建ROSIX系统实例
        ROSIX rosix = new ROSIXSystem();
        log.info("ROSIX系统已创建");
        
        // 2. 创建执行上下文
        Context ctx = rosix.createContext(
            "user_java_001",           // 用户ID
            "session_java_" + System.currentTimeMillis(), // 会话ID
            Map.of(
                "device", "java_client",
                "location", "客厅",
                "timestamp", System.currentTimeMillis()
            )
        );
        log.info("创建上下文: {}", ctx.getId());
        
        try {
            // 3. 查找空气净化器资源
            System.out.println("\n步骤1: 查找空气净化器资源...");
            List<Resource> resources = rosix.find(Query.builder()
                .type(ResourceType.ACTOR)
                .category("purifier")
                .limit(5)
                .build());
            
            if (resources.isEmpty()) {
                System.out.println("未找到空气净化器资源");
                System.out.println("\n提示：请确保Go版本的UROS系统正在运行");
                return;
            }
            
            System.out.println("找到 " + resources.size() + " 个空气净化器资源");
            Resource purifier = resources.get(0);
            System.out.println("选择资源: " + purifier.getMetadata().getName());
            System.out.println("资源路径: " + purifier.getPath());
            
            // 4. 打开资源
            System.out.println("\n步骤2: 打开资源...");
            int mode = OpenMode.combine(
                OpenMode.READ, 
                OpenMode.INVOKE
            );
            
            ResourceDescriptor rd = rosix.open(
                purifier.getPath(),
                mode,
                ctx
            );
            System.out.println("资源已打开，描述符: " + rd);
            
            try {
                // 5. 读取当前状态
                System.out.println("\n步骤3: 读取当前状态...");
                try {
                    Object status = rosix.read(rd, "status");
                    System.out.println("当前状态: " + status);
                } catch (ResourceException e) {
                    System.out.println("读取状态: " + e.getMessage());
                }
                
                // 6. 调用净化空气功能
                System.out.println("\n步骤4: 调用净化空气功能...");
                Map<String, Object> params = Map.of(
                    "mode", "auto",
                    "intensity", 3,
                    "target_pm25", 35
                );
                
                System.out.println("调用参数: " + params);
                
                Map<String, Object> result = rosix.invoke(
                    rd, 
                    "purify_air", 
                    params
                );
                
                System.out.println("\n✅ 调用成功!");
                System.out.println("执行结果: " + result);
                
                // 7. 调用其他功能（可选）
                demonstrateOtherFunctions(rosix, rd);
                
            } finally {
                // 8. 关闭资源
                System.out.println("\n步骤5: 关闭资源...");
                rosix.close(rd);
                System.out.println("资源已关闭");
            }
            
        } finally {
            // 9. 销毁上下文
            rosix.destroyContext(ctx);
            log.info("上下文已销毁");
        }
        
        System.out.println("\n=== 示例执行完成 ===");
    }
    
    /**
     * 演示其他功能
     */
    private static void demonstrateOtherFunctions(ROSIX rosix, ResourceDescriptor rd) {
        System.out.println("\n--- 演示其他功能 ---");
        
        try {
            // 设置风扇速度
            System.out.println("\n调用: set_fan_speed");
            Map<String, Object> fanResult = rosix.invoke(rd, "set_fan_speed",
                Map.of("speed", 5, "mode", "manual"));
            System.out.println("结果: " + fanResult);
            
            // 检查滤网状态
            System.out.println("\n调用: check_filter_status");
            Map<String, Object> filterResult = rosix.invoke(rd, "check_filter_status",
                Map.of());
            System.out.println("结果: " + filterResult);
            
        } catch (ResourceException e) {
            log.error("调用失败", e);
        }
    }
    
    /**
     * 简化版本 - 使用try-with-resources风格
     */
    public static void simplifiedExample() {
        System.out.println("=== 简化版本示例 ===\n");
        
        ROSIX rosix = new ROSIXSystem();
        Context ctx = rosix.createContext("user_001", "session_001", null);
        
        try {
            // 查找资源
            var resources = rosix.find(Query.builder()
                .type(ResourceType.ACTOR)
                .category("purifier")
                .limit(1)
                .build());
            
            if (!resources.isEmpty()) {
                // 打开、使用、关闭资源
                ResourceDescriptor rd = rosix.open(
                    resources.get(0).getPath(),
                    OpenMode.INVOKE.getValue(),
                    ctx
                );
                
                try {
                    // 调用净化功能
                    var result = rosix.invoke(rd, "purify_air",
                        Map.of("mode", "auto", "intensity", 3));
                    
                    System.out.println("✅ 净化器已启动: " + result);
                    
                } finally {
                    rosix.close(rd);
                }
            }
            
        } catch (ResourceException e) {
            System.err.println("❌ 操作失败: " + e.getMessage());
        } finally {
            rosix.destroyContext(ctx);
        }
    }
}



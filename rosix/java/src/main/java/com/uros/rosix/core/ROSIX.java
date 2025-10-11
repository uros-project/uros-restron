package com.uros.rosix.core;

import java.util.List;
import java.util.Map;
import java.util.function.Consumer;

/**
 * ROSIX主接口 - 类似POSIX的系统调用接口
 * 对应Go版本的ROSIX interface
 */
public interface ROSIX {
    
    // ==================== 资源操作原语 ====================
    
    /**
     * 打开资源，返回资源描述符
     * 类似 POSIX: open()
     * 
     * @param path 资源路径
     * @param mode 打开模式
     * @param ctx 执行上下文
     * @return 资源描述符
     * @throws ResourceException 资源操作异常
     */
    ResourceDescriptor open(ResourcePath path, int mode, Context ctx) throws ResourceException;
    
    /**
     * 关闭资源
     * 类似 POSIX: close()
     */
    void close(ResourceDescriptor rd) throws ResourceException;
    
    /**
     * 读取资源属性或特征
     * 类似 POSIX: read()
     */
    Object read(ResourceDescriptor rd, String key) throws ResourceException;
    
    /**
     * 写入资源属性或特征
     * 类似 POSIX: write()
     */
    void write(ResourceDescriptor rd, String key, Object value) throws ResourceException;
    
    /**
     * 调用资源行为
     * 类似 POSIX: ioctl()
     */
    Map<String, Object> invoke(ResourceDescriptor rd, String behavior, 
                               Map<String, Object> params) throws ResourceException;
    
    // ==================== 资源发现和查询 ====================
    
    /**
     * 查找资源
     */
    List<Resource> find(Query query) throws ResourceException;
    
    /**
     * 列出指定路径下的资源
     * 类似 POSIX: readdir()
     */
    List<Resource> list(ResourcePath path) throws ResourceException;
    
    /**
     * 获取资源信息
     * 类似 POSIX: stat()
     */
    Resource stat(ResourceDescriptor rd) throws ResourceException;
    
    // ==================== 资源监听 ====================
    
    /**
     * 监听资源变化
     * 类似 Linux: inotify
     */
    void watch(ResourceDescriptor rd, List<EventType> events, 
              Consumer<Event> callback) throws ResourceException;
    
    /**
     * 取消监听
     */
    void unwatch(ResourceDescriptor rd) throws ResourceException;
    
    // ==================== 上下文管理 ====================
    
    /**
     * 创建执行上下文
     */
    Context createContext(String userId, String sessionId, 
                         Map<String, Object> metadata);
    
    /**
     * 销毁执行上下文
     */
    void destroyContext(Context ctx);
}


package com.uros.rosix.ai;

import com.uros.rosix.core.Context;
import com.uros.rosix.core.ResourceException;
import lombok.Builder;
import lombok.Data;

import java.util.List;
import java.util.Map;

/**
 * AI编排器接口
 * 对应Go版本的ai.AIOrchestrator
 */
public interface AIOrchestrator {
    
    /**
     * 通过自然语言调用资源
     * 示例: "打开客厅的空气净化器"
     */
    InvokeResult invoke(String prompt, Context context) throws ResourceException;
    
    /**
     * 编排多资源协同
     * 示例: "晚上8点，关闭所有灯光，调低空调温度"
     */
    Plan orchestrate(String goal, Context context) throws ResourceException;
    
    /**
     * 查询资源信息
     * 示例: "客厅的温度是多少？"
     */
    QueryResult query(String question, Context context) throws ResourceException;
    
    /**
     * AI调用结果
     */
    @Data
    @Builder
    class InvokeResult {
        boolean success;
        String intent;
        List<String> resources;
        List<Action> actions;
        Map<String, Object> result;
        String message;
    }
    
    /**
     * 执行计划
     */
    @Data
    @Builder
    class Plan {
        String id;
        String goal;
        List<PlanStep> steps;
        List<String> resources;
        int estimatedTime;
        Map<String, Object> metadata;
    }
    
    /**
     * 计划步骤
     */
    @Data
    @Builder
    class PlanStep {
        int order;
        String description;
        String resource;
        Action action;
        List<Integer> dependsOn;
        String condition;
    }
    
    /**
     * 动作定义
     */
    @Data
    @Builder
    class Action {
        ActionType type;
        String behavior;
        Map<String, Object> parameters;
    }
    
    /**
     * 动作类型
     */
    enum ActionType {
        READ, WRITE, INVOKE, WAIT
    }
    
    /**
     * 查询结果
     */
    @Data
    @Builder
    class QueryResult {
        String answer;
        Map<String, Object> data;
        List<String> resources;
        double confidence;
    }
}


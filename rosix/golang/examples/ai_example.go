package main

import (
	"fmt"
	"log"
	"uros-restron/rosix/ai"
	"uros-restron/rosix/core"
)

// AI驱动的资源管理示例
func main() {
	fmt.Println("=== AI驱动的资源管理示例 ===\n")

	// 这些示例展示了如何使用AI接口进行资源管理
	aiExamples()
}

func aiExamples() {
	// 示例场景
	scenarios := []struct {
		name        string
		description string
		prompt      string
		goal        string
	}{
		{
			name:        "场景1：环境控制",
			description: "通过自然语言控制室内环境",
			prompt:      "客厅太热了，帮我降低温度",
			goal:        "",
		},
		{
			name:        "场景2：空气质量管理",
			description: "智能管理室内空气质量",
			prompt:      "启动空气净化，PM2.5有点高",
			goal:        "",
		},
		{
			name:        "场景3：多资源协同",
			description: "编排多个资源协同工作",
			prompt:      "",
			goal:        "晚上8点，进入睡眠模式：关闭灯光，调低空调，启动净化器",
		},
		{
			name:        "场景4：信息查询",
			description: "查询资源状态和信息",
			prompt:      "客厅的温度和湿度是多少？",
			goal:        "",
		},
		{
			name:        "场景5：智能建议",
			description: "获取资源使用建议",
			prompt:      "如何降低室内PM2.5？",
			goal:        "",
		},
	}

	for _, scenario := range scenarios {
		fmt.Printf("### %s\n", scenario.name)
		fmt.Printf("描述: %s\n", scenario.description)

		if scenario.prompt != "" {
			fmt.Printf("\n用户: \"%s\"\n", scenario.prompt)
			fmt.Println("\n系统执行流程:")
			fmt.Println("1. AI识别用户意图")
			fmt.Println("2. 查找相关资源")
			fmt.Println("3. 生成执行计划")
			fmt.Println("4. 调用资源行为")
			fmt.Println("5. 返回执行结果")
		}

		if scenario.goal != "" {
			fmt.Printf("\n目标: \"%s\"\n", scenario.goal)
			fmt.Println("\n系统生成计划:")
			fmt.Println("步骤1: 查找所有灯光资源 -> 调用turn_off行为")
			fmt.Println("步骤2: 查找空调资源 -> 调用set_temperature(24)行为")
			fmt.Println("步骤3: 查找净化器资源 -> 调用purify_air(silent)行为")
		}

		fmt.Println("\n" + "---" + "\n")
	}
}

// 实际的AI使用示例（需要完整的系统初始化）
func aiUsageExample(rosix core.ROSIX) {
	// 创建AI编排器
	orchestrator := ai.NewSimpleOrchestrator(rosix)

	// 创建执行上下文
	ctx, err := rosix.CreateContext("user_001", "session_123", map[string]interface{}{
		"location": "客厅",
		"time":     "20:00",
	})
	if err != nil {
		log.Fatalf("创建上下文失败: %v", err)
	}
	defer rosix.DestroyContext(ctx)

	// 示例1：通过自然语言调用
	fmt.Println("=== 示例1：自然语言调用 ===")
	result, err := orchestrator.Invoke("打开空气净化器", ctx)
	if err != nil {
		log.Printf("调用失败: %v", err)
	} else {
		fmt.Printf("识别意图: %s\n", result.Intent)
		fmt.Printf("涉及资源: %v\n", result.Resources)
		fmt.Printf("执行动作: %v\n", result.Actions)
		fmt.Printf("结果: %s\n", result.Message)
	}

	// 示例2：多资源编排
	fmt.Println("\n=== 示例2：多资源编排 ===")
	plan, err := orchestrator.Orchestrate("进入睡眠模式", ctx)
	if err != nil {
		log.Printf("编排失败: %v", err)
	} else {
		fmt.Printf("计划ID: %s\n", plan.ID)
		fmt.Printf("目标: %s\n", plan.Goal)
		fmt.Printf("预计时间: %d秒\n", plan.EstimatedTime)
		fmt.Println("执行步骤:")
		for _, step := range plan.Steps {
			fmt.Printf("  %d. %s\n", step.Order, step.Description)
			fmt.Printf("     资源: %s\n", step.Resource)
			fmt.Printf("     行为: %s\n", step.Action.Behavior)
			if len(step.DependsOn) > 0 {
				fmt.Printf("     依赖: 步骤 %v\n", step.DependsOn)
			}
		}
	}

	// 示例3：信息查询
	fmt.Println("\n=== 示例3：信息查询 ===")
	queryResult, err := orchestrator.Query("客厅的温度是多少？", ctx)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("回答: %s\n", queryResult.Answer)
		fmt.Printf("置信度: %.2f\n", queryResult.Confidence)
		fmt.Printf("相关资源: %v\n", queryResult.Resources)
		fmt.Printf("数据: %v\n", queryResult.Data)
	}

	// 示例4：获取建议
	fmt.Println("\n=== 示例4：获取建议 ===")
	suggestion, err := orchestrator.Suggest("如何改善室内空气质量？", ctx)
	if err != nil {
		log.Printf("获取建议失败: %v", err)
	} else {
		fmt.Printf("问题: %s\n", suggestion.Question)
		fmt.Println("建议:")
		for i, item := range suggestion.Suggestions {
			fmt.Printf("  %d. %s\n", i+1, item.Title)
			fmt.Printf("     描述: %s\n", item.Description)
			fmt.Printf("     优先级: %d\n", item.Priority)
			fmt.Printf("     置信度: %.2f\n", item.Confidence)
		}
	}

	// 示例5：学习用户行为
	fmt.Println("\n=== 示例5：学习用户行为 ===")
	behavior := &ai.UserBehavior{
		UserID:    "user_001",
		Action:    "turn_on_purifier",
		Resources: []string{"purifier_001"},
		Timestamp: 1704096000,
		Context: map[string]interface{}{
			"time":     "morning",
			"location": "客厅",
			"pm25":     65,
		},
	}
	err = orchestrator.Learn(behavior)
	if err != nil {
		log.Printf("学习失败: %v", err)
	} else {
		fmt.Println("用户行为已记录，系统将学习用户偏好")
	}

	// 示例6：预测资源状态
	fmt.Println("\n=== 示例6：预测资源状态 ===")
	prediction, err := orchestrator.Predict("purifier_001", 3600)
	if err != nil {
		log.Printf("预测失败: %v", err)
	} else {
		fmt.Printf("资源: %s\n", prediction.ResourceID)
		fmt.Printf("预测时长: %d秒\n", prediction.Duration)
		fmt.Printf("预测结果: %v\n", prediction.Predictions)
		fmt.Printf("置信度: %.2f\n", prediction.Confidence)
	}
}

// AI协同场景示例
func collaborationScenarios() {
	scenarios := []string{
		"场景：回家模式",
		"  用户: '我回家了'",
		"  AI识别:",
		"    - 打开玄关灯光",
		"    - 启动空调到舒适温度",
		"    - 打开空气净化器",
		"    - 播放欢迎音乐",
		"",
		"场景：会议模式",
		"  用户: '准备开会'",
		"  AI识别:",
		"    - 关闭音乐设备",
		"    - 调整灯光亮度",
		"    - 启动投影仪",
		"    - 设置勿扰模式",
		"",
		"场景：节能模式",
		"  用户: '启动节能模式'",
		"  AI识别:",
		"    - 关闭无人区域灯光",
		"    - 调整空调到节能温度",
		"    - 降低净化器功率",
		"    - 关闭非必要设备",
		"",
		"场景：健康提醒",
		"  系统主动:",
		"    检测到PM2.5超标 -> 建议开启净化器",
		"    检测到CO2浓度高 -> 建议通风",
		"    检测到温度过高 -> 建议调整空调",
	}

	fmt.Println("=== AI协同场景示例 ===\n")
	for _, line := range scenarios {
		fmt.Println(line)
	}
}

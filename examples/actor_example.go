package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"uros-restron/internal/actor"
	"uros-restron/internal/models"
)

func RunActorExample() {
	fmt.Println("=== UROS Actor系统示例 ===")

	// 创建模拟的Behavior数据
	behavior := &models.Behavior{
		ID:          "purifier-001",
		Name:        "空气净化行为",
		Type:        "purifier",
		Description: "执行空气净化功能",
		Category:    "device",
		Functions: map[string]models.Function{
			"purify_air": models.Function{
				Name:        "净化空气",
				Description: "执行空气净化操作",
				InputParams: map[string]models.Parameter{
					"air_quality": {
						Type:        "number",
						Description: "当前空气质量指数",
						Required:    true,
						Min:         &[]float64{0}[0],
						Max:         &[]float64{500}[0],
					},
					"target_quality": {
						Type:        "number",
						Description: "目标空气质量指数",
						Required:    false,
						Default:     50,
					},
				},
				OutputParams: map[string]models.Parameter{
					"success": {
						Type:        "boolean",
						Description: "净化操作是否成功",
					},
					"final_quality": {
						Type:        "number",
						Description: "净化后的空气质量指数",
					},
				},
				Implementation: models.FunctionImplementation{
					Steps: []models.ImplementationStep{
						{
							Step:        1,
							Action:      "check_air_quality",
							Description: "检查当前空气质量",
						},
						{
							Step:        2,
							Action:      "start_fan",
							Description: "启动风扇",
							Condition:   "air_quality > target_quality",
						},
						{
							Step:        3,
							Action:      "activate_filter",
							Description: "激活过滤器",
						},
						{
							Step:        4,
							Action:      "monitor_progress",
							Description: "监控净化进度",
						},
					},
				},
			},
		},
		Parameters: map[string]interface{}{
			"fan_speed":   "auto",
			"filter_type": "hepa",
			"auto_mode":   true,
		},
	}

	// 创建BehaviorActor
	behaviorActor := actor.NewBehaviorActor(behavior)

	// 启动Actor
	ctx := context.Background()
	if err := behaviorActor.Start(ctx); err != nil {
		log.Fatalf("启动Actor失败: %v", err)
	}

	fmt.Printf("✅ Actor已启动: %s (%s)\n", behaviorActor.GetBehavior().Name, behaviorActor.ID())

	// 等待一下让Actor完全启动
	time.Sleep(100 * time.Millisecond)

	// 获取可用函数列表
	functions := behaviorActor.GetAvailableFunctions()
	fmt.Printf("📋 可用函数: %v\n", functions)

	// 调用净化空气函数
	fmt.Println("\n=== 调用净化空气函数 ===")
	params := map[string]interface{}{
		"air_quality":    150.0,
		"target_quality": 50.0,
	}

	result, err := behaviorActor.CallFunction("purify_air", params)
	if err != nil {
		log.Printf("❌ 函数调用失败: %v", err)
	} else {
		fmt.Printf("✅ 函数调用成功，结果: %+v\n", result)
	}

	// 获取函数信息
	fmt.Println("\n=== 获取函数信息 ===")
	funcInfo, err := behaviorActor.GetFunctionInfo("purify_air")
	if err != nil {
		log.Printf("❌ 获取函数信息失败: %v", err)
	} else {
		funcInfoJSON, _ := json.MarshalIndent(funcInfo, "", "  ")
		fmt.Printf("📄 函数信息:\n%s\n", funcInfoJSON)
	}

	// 发送消息到Actor
	fmt.Println("\n=== 发送消息到Actor ===")
	msg := actor.NewFunctionCallMessage("example", behaviorActor.ID(), "purify_air", params)
	msg.SetCorrelationID("test-001")

	if err := behaviorActor.Send(msg); err != nil {
		log.Printf("❌ 发送消息失败: %v", err)
	} else {
		fmt.Printf("✅ 消息发送成功: %s\n", msg.Type)
	}

	// 等待一下处理消息
	time.Sleep(200 * time.Millisecond)

	// 停止Actor
	fmt.Println("\n=== 停止Actor ===")
	if err := behaviorActor.Stop(); err != nil {
		log.Printf("❌ 停止Actor失败: %v", err)
	} else {
		fmt.Printf("✅ Actor已停止: %s\n", behaviorActor.ID())
	}

	fmt.Println("\n=== 示例完成 ===")
}

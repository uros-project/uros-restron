package main

import (
	"fmt"
	"log"
	"uros-restron/rosix/ai"
	"uros-restron/rosix/core"
)

// 演示ROSIX的基本使用
func main() {
	// 注意：这是一个示例，实际使用时需要完整初始化系统
	fmt.Println("=== ROSIX 基本使用示例 ===\n")

	// 示例1：资源发现和打开
	example1()

	// 示例2：读取资源属性
	example2()

	// 示例3：调用资源行为
	example3()

	// 示例4：AI驱动的资源操作
	example4()
}

// 示例1：资源发现和打开
func example1() {
	fmt.Println("示例1：资源发现和打开")
	fmt.Println("```go")
	fmt.Println(`// 创建执行上下文
ctx, _ := rosix.CreateContext("user_001", "session_123", nil)

// 查找环境监测类型的资源
resources, err := rosix.Find(core.Query{
    Type:     core.TypeDevice,
    Category: "environment",
    Limit:    5,
})

fmt.Printf("找到 %d 个环境监测资源\n", len(resources))

// 打开第一个资源
if len(resources) > 0 {
    rd, err := rosix.Open(resources[0].Path(), 
        core.ModeRead | core.ModeInvoke, ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer rosix.Close(rd)
    
    fmt.Printf("成功打开资源: %s\n", resources[0].ID())
}`)
	fmt.Println("```\n")
}

// 示例2：读取资源属性
func example2() {
	fmt.Println("示例2：读取资源属性")
	fmt.Println("```go")
	fmt.Println(`// 假设已经打开了资源，获得了资源描述符 rd

// 读取资源的温度特征
temperature, err := rosix.Read(rd, "temperature")
if err != nil {
    log.Printf("读取温度失败: %v", err)
} else {
    fmt.Printf("当前温度: %.1f°C\n", temperature)
}

// 读取资源的湿度特征
humidity, err := rosix.Read(rd, "humidity")
if err != nil {
    log.Printf("读取湿度失败: %v", err)
} else {
    fmt.Printf("当前湿度: %.1f%%\n", humidity)
}

// 读取资源的静态属性
name, _ := rosix.Read(rd, "name")
fmt.Printf("资源名称: %s\n", name)`)
	fmt.Println("```\n")
}

// 示例3：调用资源行为
func example3() {
	fmt.Println("示例3：调用资源行为")
	fmt.Println("```go")
	fmt.Println(`// 假设已经打开了空气净化器资源

// 调用净化空气行为
result, err := rosix.Invoke(rd, "purify_air", map[string]interface{}{
    "mode":        "auto",
    "intensity":   3,
    "target_pm25": 35,
})

if err != nil {
    log.Printf("调用失败: %v", err)
} else {
    fmt.Printf("净化器已启动: %v\n", result)
}

// 调用设置风扇速度行为
result, err = rosix.Invoke(rd, "set_fan_speed", map[string]interface{}{
    "speed": 5,
    "mode":  "manual",
})

if err != nil {
    log.Printf("调用失败: %v", err)
} else {
    fmt.Printf("风扇速度已设置: %v\n", result)
}`)
	fmt.Println("```\n")
}

// 示例4：AI驱动的资源操作
func example4() {
	fmt.Println("示例4：AI驱动的资源操作")
	fmt.Println("```go")
	fmt.Println(`// 创建AI编排器
orchestrator := ai.NewSimpleOrchestrator(rosix)

// 通过自然语言调用资源
result, err := orchestrator.Invoke("打开客厅的空气净化器", ctx)
if err != nil {
    log.Printf("AI调用失败: %v", err)
} else {
    fmt.Printf("AI识别意图: %s\n", result.Intent)
    fmt.Printf("执行结果: %s\n", result.Message)
}

// 编排多资源协同
plan, err := orchestrator.Orchestrate("进入睡眠模式", ctx)
if err != nil {
    log.Printf("编排失败: %v", err)
} else {
    fmt.Printf("生成计划: %s\n", plan.ID)
    fmt.Printf("计划步骤:\n")
    for _, step := range plan.Steps {
        fmt.Printf("  %d. %s\n", step.Order, step.Description)
    }
}

// 查询资源信息
queryResult, err := orchestrator.Query("客厅的温度是多少？", ctx)
if err != nil {
    log.Printf("查询失败: %v", err)
} else {
    fmt.Printf("回答: %s\n", queryResult.Answer)
}`)
	fmt.Println("```\n")
}

// 完整的使用示例
func completeExample(rosix core.ROSIX) {
	fmt.Println("=== 完整使用示例 ===\n")

	// 1. 创建上下文
	ctx, err := rosix.CreateContext("user_001", "session_123", map[string]interface{}{
		"device": "mobile",
		"ip":     "192.168.1.100",
	})
	if err != nil {
		log.Fatalf("创建上下文失败: %v", err)
	}
	defer rosix.DestroyContext(ctx)

	// 2. 查找资源
	resources, err := rosix.Find(core.Query{
		Category: "purifier",
		Limit:    1,
	})
	if err != nil || len(resources) == 0 {
		log.Printf("未找到空气净化器")
		return
	}

	fmt.Printf("找到资源: %s\n", resources[0].Metadata().Name)

	// 3. 打开资源
	rd, err := rosix.Open(
		resources[0].Path(),
		core.ModeRead|core.ModeWrite|core.ModeInvoke,
		ctx,
	)
	if err != nil {
		log.Fatalf("打开资源失败: %v", err)
	}
	defer rosix.Close(rd)

	// 4. 读取资源状态
	status, err := rosix.Read(rd, "status")
	if err != nil {
		log.Printf("读取状态失败: %v", err)
	} else {
		fmt.Printf("当前状态: %v\n", status)
	}

	// 5. 调用资源行为
	result, err := rosix.Invoke(rd, "purify_air", map[string]interface{}{
		"mode":      "auto",
		"intensity": 3,
	})
	if err != nil {
		log.Printf("调用失败: %v", err)
	} else {
		fmt.Printf("调用结果: %v\n", result)
	}

	// 6. 使用AI编排器
	orchestrator := ai.NewSimpleOrchestrator(rosix)
	aiResult, err := orchestrator.Invoke("启动空气净化", ctx)
	if err != nil {
		log.Printf("AI调用失败: %v", err)
	} else {
		fmt.Printf("AI执行结果: %s\n", aiResult.Message)
	}
}

// 资源监听示例
func watchExample(rosix core.ROSIX) {
	fmt.Println("=== 资源监听示例 ===\n")

	ctx, _ := rosix.CreateContext("user_001", "session_123", nil)
	defer rosix.DestroyContext(ctx)

	// 打开资源（监听模式）
	resources, _ := rosix.Find(core.Query{Limit: 1})
	if len(resources) == 0 {
		return
	}

	rd, err := rosix.Open(
		resources[0].Path(),
		core.ModeRead|core.ModeWatch,
		ctx,
	)
	if err != nil {
		log.Fatalf("打开资源失败: %v", err)
	}
	defer rosix.Close(rd)

	// 设置监听回调
	callback := func(event core.Event) error {
		fmt.Printf("收到事件: %s, 资源: %v, 数据: %v\n",
			event.Type, event.Resource, event.Data)
		return nil
	}

	// 开始监听
	err = rosix.Watch(rd, []core.EventType{
		core.EventStateChange,
		core.EventFeatureUpdate,
	}, callback)

	if err != nil {
		log.Printf("监听失败: %v", err)
	}

	fmt.Println("正在监听资源变化...")
	// 这里应该保持程序运行，实际使用中可能需要阻塞或使用其他机制
}

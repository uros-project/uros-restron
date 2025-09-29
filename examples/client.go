package examples

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Thing 数字孪生结构
type Thing struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Properties  []Property             `json:"properties"`
	Status      map[string]interface{} `json:"status"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

type Property struct {
	ID      string      `json:"id"`
	ThingID string      `json:"thingId"`
	Name    string      `json:"name"`
	Value   interface{} `json:"value"`
	Type    string      `json:"type"`
}

// DittoClient 客户端
type DittoClient struct {
	baseURL string
	client  *http.Client
}

// NewDittoClient 创建新的客户端
func NewDittoClient(baseURL string) *DittoClient {
	return &DittoClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// CreateThing 创建数字孪生
func (c *DittoClient) CreateThing(thing *Thing) (*Thing, error) {
	jsonData, err := json.Marshal(thing)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(c.baseURL+"/api/v1/things", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("创建失败: %s", string(body))
	}

	var result Thing
	err = json.Unmarshal(body, &result)
	return &result, err
}

// GetThing 获取数字孪生
func (c *DittoClient) GetThing(id string) (*Thing, error) {
	resp, err := c.client.Get(c.baseURL + "/api/v1/things/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取失败: %s", string(body))
	}

	var result Thing
	err = json.Unmarshal(body, &result)
	return &result, err
}

// UpdateProperty 更新属性
func (c *DittoClient) UpdateProperty(thingID, propertyName string, value interface{}) error {
	requestData := map[string]interface{}{
		"value": value,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", c.baseURL+"/api/v1/things/"+thingID+"/properties/"+propertyName, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("更新属性失败: %s", string(body))
	}

	return nil
}

// UpdateStatus 更新状态
func (c *DittoClient) UpdateStatus(thingID string, status map[string]interface{}) error {
	jsonData, err := json.Marshal(status)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", c.baseURL+"/api/v1/things/"+thingID+"/status", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("更新状态失败: %s", string(body))
	}

	return nil
}

func RunClientExample() {
	// 创建客户端
	client := NewDittoClient("http://localhost:8080")

	// 示例1: 创建智能设备
	fmt.Println("=== 创建智能设备 ===")
	device := &Thing{
		Name:        "智能空调001",
		Type:        "machine",
		Description: "客厅空调设备",
		Properties: []Property{
			{Name: "temperature", Value: 22.5, Type: "number"},
			{Name: "mode", Value: "cooling", Type: "string"},
		},
		Status: map[string]interface{}{
			"online": true,
			"power":  "on",
		},
	}

	createdDevice, err := client.CreateThing(device)
	if err != nil {
		fmt.Printf("创建设备失败: %v\n", err)
		return
	}
	fmt.Printf("设备创建成功: %s (ID: %s)\n", createdDevice.Name, createdDevice.ID)

	// 示例2: 创建人员
	fmt.Println("\n=== 创建人员 ===")
	person := &Thing{
		Name:        "张三",
		Type:        "person",
		Description: "系统管理员",
		Properties: []Property{
			{Name: "department", Value: "IT", Type: "string"},
			{Name: "role", Value: "admin", Type: "string"},
		},
		Status: map[string]interface{}{
			"active":    true,
			"lastLogin": time.Now().Format(time.RFC3339),
		},
	}

	createdPerson, err := client.CreateThing(person)
	if err != nil {
		fmt.Printf("创建人员失败: %v\n", err)
		return
	}
	fmt.Printf("人员创建成功: %s (ID: %s)\n", createdPerson.Name, createdPerson.ID)

	// 示例3: 更新设备属性
	fmt.Println("\n=== 更新设备属性 ===")
	err = client.UpdateProperty(createdDevice.ID, "temperature", 24.0)
	if err != nil {
		fmt.Printf("更新属性失败: %v\n", err)
	} else {
		fmt.Println("设备温度已更新为 24.0°C")
	}

	// 示例4: 更新设备状态
	fmt.Println("\n=== 更新设备状态 ===")
	newStatus := map[string]interface{}{
		"online":   true,
		"power":    "on",
		"mode":     "auto",
		"fanSpeed": "medium",
	}
	err = client.UpdateStatus(createdDevice.ID, newStatus)
	if err != nil {
		fmt.Printf("更新状态失败: %v\n", err)
	} else {
		fmt.Println("设备状态已更新")
	}

	// 示例5: 获取更新后的设备信息
	fmt.Println("\n=== 获取设备信息 ===")
	updatedDevice, err := client.GetThing(createdDevice.ID)
	if err != nil {
		fmt.Printf("获取设备信息失败: %v\n", err)
	} else {
		fmt.Printf("设备信息: %+v\n", updatedDevice)
	}

	fmt.Println("\n=== 示例完成 ===")
}

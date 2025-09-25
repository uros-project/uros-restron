package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Hub 维护活跃的客户端连接
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMessage
	mutex      sync.RWMutex
}

// Client 表示一个 WebSocket 客户端
type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan []byte
	thingID string // 可选的：只订阅特定事物的更新
}

// BroadcastMessage 广播消息结构
type BroadcastMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Message 客户端消息结构
type Message struct {
	Type    string      `json:"type"`
	ThingID string      `json:"thingId,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该更严格
	},
}

// NewHub 创建新的 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMessage),
	}
}

// Run 启动 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			logrus.Info("Client connected")

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()
			logrus.Info("Client disconnected")

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				// 如果客户端订阅了特定事物，只发送相关消息
				if client.thingID != "" && message.Type != "thing_created" {
					// 检查消息是否与客户端订阅的事物相关
					if !h.isMessageRelevant(client.thingID, message) {
						continue
					}
				}

				select {
				case client.send <- h.encodeMessage(message):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// isMessageRelevant 检查消息是否与特定事物相关
func (h *Hub) isMessageRelevant(thingID string, message BroadcastMessage) bool {
	switch message.Type {
	case "thing_updated", "property_updated", "status_updated":
		if data, ok := message.Data.(map[string]interface{}); ok {
			if id, exists := data["thingId"]; exists {
				return id == thingID
			}
		}
		return false
	case "thing_deleted":
		if data, ok := message.Data.(map[string]interface{}); ok {
			if id, exists := data["id"]; exists {
				return id == thingID
			}
		}
		return false
	default:
		return true
	}
}

// encodeMessage 编码消息为 JSON
func (h *Hub) encodeMessage(message BroadcastMessage) []byte {
	data, err := json.Marshal(message)
	if err != nil {
		logrus.Error("Failed to encode message:", err)
		return []byte{}
	}
	return data
}

// Broadcast 广播消息给所有客户端
func (h *Hub) Broadcast(messageType string, data interface{}) {
	message := BroadcastMessage{
		Type: messageType,
		Data: data,
	}
	select {
	case h.broadcast <- message:
	default:
		logrus.Warn("Broadcast channel is full, dropping message")
	}
}

// handleWebSocket 处理 WebSocket 连接
func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error("Failed to upgrade connection:", err)
		return
	}

	thingID := c.Query("thingId") // 可选的：只订阅特定事物

	client := &Client{
		hub:     s.hub,
		conn:    conn,
		send:    make(chan []byte, 256),
		thingID: thingID,
	}

	client.hub.register <- client

	// 启动 goroutine 处理客户端
	go client.writePump()
	go client.readPump()
}

// readPump 读取客户端消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var message Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Error("WebSocket error:", err)
			}
			break
		}

		// 处理客户端消息
		c.handleMessage(message)
	}
}

// writePump 向客户端发送消息
func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			logrus.Error("Failed to write message:", err)
			return
		}
	}
}

// handleMessage 处理客户端消息
func (c *Client) handleMessage(message Message) {
	switch message.Type {
	case "subscribe":
		// 客户端请求订阅特定事物
		if thingID, ok := message.Data.(string); ok {
			c.thingID = thingID
			logrus.Info("Client subscribed to thing:", thingID)
		}
	case "ping":
		// 心跳检测
		response := Message{Type: "pong"}
		c.conn.WriteJSON(response)
	default:
		logrus.Warn("Unknown message type:", message.Type)
	}
}

package services

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Client WebSocket 클라이언트
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   uint
	RoomID   uint
}

// Hub WebSocket 연결 관리
type Hub struct {
	// 채팅방별 클라이언트 관리
	Rooms map[uint]map[*Client]bool

	// 브로드캐스트 채널
	Broadcast chan *Message

	// 클라이언트 등록/해제
	Register   chan *Client
	Unregister chan *Client

	mu sync.RWMutex
}

// Message WebSocket 메시지 구조
type Message struct {
	Type       string      `json:"type"` // message, read, member_join, member_leave
	RoomID     uint        `json:"room_id"`
	UserID     uint        `json:"user_id"`
	Data       interface{} `json:"data"`
}

// NewHub Hub 생성
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[uint]map[*Client]bool),
		Broadcast:  make(chan *Message, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run Hub 실행
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, ok := h.Rooms[client.RoomID]; !ok {
				h.Rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.Rooms[client.RoomID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered: UserID=%d, RoomID=%d", client.UserID, client.RoomID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Rooms[client.RoomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Rooms, client.RoomID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: UserID=%d, RoomID=%d", client.UserID, client.RoomID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			if clients, ok := h.Rooms[message.RoomID]; ok {
				messageBytes, err := json.Marshal(message)
				if err != nil {
					log.Printf("Error marshaling message: %v", err)
					h.mu.RUnlock()
					continue
				}

				for client := range clients {
					select {
					case client.Send <- messageBytes:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// ReadPump 클라이언트로부터 메시지 읽기
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// 클라이언트로부터 받은 메시지 처리
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// 메시지 유효성 검증
		if msg.RoomID != c.RoomID {
			log.Printf("Invalid room ID: expected=%d, got=%d", c.RoomID, msg.RoomID)
			continue
		}

		// 브로드캐스트 (서버에서 처리 후 다시 보내는 방식이므로 여기서는 무시)
		// 실제 메시지는 HTTP API를 통해 저장되고, 저장 후 Hub를 통해 브로드캐스트됨
	}
}

// WritePump 클라이언트로 메시지 쓰기
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		// 대기 중인 메시지들도 한번에 전송
		n := len(c.Send)
		for i := 0; i < n; i++ {
			w.Write([]byte{'\n'})
			w.Write(<-c.Send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}

// BroadcastMessage 메시지 브로드캐스트
func (h *Hub) BroadcastMessage(roomID uint, msgType string, userID uint, data interface{}) {
	message := &Message{
		Type:   msgType,
		RoomID: roomID,
		UserID: userID,
		Data:   data,
	}
	h.Broadcast <- message
}

// 전역 Hub 인스턴스
var GlobalHub *Hub

// InitHub Hub 초기화
func InitHub() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
	log.Println("WebSocket Hub initialized")
}

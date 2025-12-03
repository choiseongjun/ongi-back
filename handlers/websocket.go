package handlers

import (
	"log"
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// WebSocketHandler WebSocket 연결 핸들러
func WebSocketHandler(c *fiber.Ctx) error {
	// WebSocket 업그레이드
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// HandleWebSocket WebSocket 연결 처리
func HandleWebSocket(c *websocket.Conn) {
	// URL 파라미터에서 roomID와 userID 가져오기
	roomIDStr := c.Params("roomId")
	userIDStr := c.Query("user_id")

	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid room ID: %v", err)
		c.Close()
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		c.Close()
		return
	}

	// 채팅방 멤버 확인
	var membership models.ChatRoomMember
	if err := database.DB.Where("chat_room_id = ? AND user_id = ?", roomID, userID).First(&membership).Error; err != nil {
		log.Printf("User is not a member of this room: roomID=%d, userID=%d", roomID, userID)
		c.Close()
		return
	}

	// 클라이언트 생성
	client := &services.Client{
		Hub:    services.GlobalHub,
		Conn:   c,
		Send:   make(chan []byte, 256),
		UserID: uint(userID),
		RoomID: uint(roomID),
	}

	// Hub에 등록
	client.Hub.Register <- client

	// 입장 메시지 브로드캐스트
	client.Hub.BroadcastMessage(
		uint(roomID),
		"member_online",
		uint(userID),
		fiber.Map{
			"user_id": userID,
			"status":  "online",
		},
	)

	// 고루틴 시작
	go client.WritePump()
	client.ReadPump()

	// 퇴장 메시지 브로드캐스트
	client.Hub.BroadcastMessage(
		uint(roomID),
		"member_offline",
		uint(userID),
		fiber.Map{
			"user_id": userID,
			"status":  "offline",
		},
	)
}

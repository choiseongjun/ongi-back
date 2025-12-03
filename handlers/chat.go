package handlers

import (
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateChatRoomRequest 채팅방 생성 요청
type CreateChatRoomRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ClubID      *uint  `json:"club_id"`      // 클럽 채팅방인 경우
	RoomType    string `json:"room_type"`    // group, club, direct
	MemberIDs   []uint `json:"member_ids"`   // 초대할 멤버 ID 목록
}

// CreateChatRoom 채팅방 생성
// POST /chat/rooms
func CreateChatRoom(c *fiber.Ctx) error {
	var req CreateChatRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// 생성자 ID (실제로는 JWT에서 가져와야 함)
	// 여기서는 임시로 request body에서 받거나 기본값 사용
	createdBy := uint(1) // TODO: JWT에서 user_id 추출

	// 기본값 설정
	if req.RoomType == "" {
		req.RoomType = "group"
	}

	// 채팅방 생성
	chatRoom := models.ChatRoom{
		Name:        req.Name,
		Description: req.Description,
		ClubID:      req.ClubID,
		RoomType:    req.RoomType,
		CreatedBy:   createdBy,
		MemberCount: len(req.MemberIDs) + 1, // 생성자 포함
	}

	if err := database.DB.Create(&chatRoom).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create chat room",
			"details": err.Error(),
		})
	}

	// 생성자를 admin으로 추가
	creatorMember := models.ChatRoomMember{
		ChatRoomID: chatRoom.ID,
		UserID:     createdBy,
		Role:       "admin",
		JoinedAt:   time.Now(),
	}
	database.DB.Create(&creatorMember)

	// 멤버 추가
	for _, memberID := range req.MemberIDs {
		if memberID == createdBy {
			continue // 생성자는 이미 추가됨
		}
		member := models.ChatRoomMember{
			ChatRoomID: chatRoom.ID,
			UserID:     memberID,
			Role:       "member",
			JoinedAt:   time.Now(),
		}
		database.DB.Create(&member)
	}

	// 생성된 채팅방 정보 조회 (멤버 정보 포함)
	database.DB.Preload("Members.User").Preload("Creator").Preload("Club").First(&chatRoom, chatRoom.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Chat room created successfully",
		"data":    chatRoom,
	})
}

// GetChatRooms 사용자의 채팅방 목록 조회
// GET /chat/rooms
func GetChatRooms(c *fiber.Ctx) error {
	// 사용자 ID (실제로는 JWT에서 가져와야 함)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "user_id is required",
		})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user_id",
		})
	}

	// 사용자가 속한 채팅방 멤버십 조회
	var memberships []models.ChatRoomMember
	if err := database.DB.Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch chat rooms",
		})
	}

	// 채팅방 ID 목록 추출
	var roomIDs []uint
	for _, membership := range memberships {
		roomIDs = append(roomIDs, membership.ChatRoomID)
	}

	if len(roomIDs) == 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    []models.ChatRoom{},
		})
	}

	// 채팅방 정보 조회
	var chatRooms []models.ChatRoom
	if err := database.DB.
		Preload("Creator").
		Preload("Club").
		Where("id IN ?", roomIDs).
		Order("last_message_at DESC NULLS LAST, created_at DESC").
		Find(&chatRooms).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch chat rooms",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    chatRooms,
	})
}

// GetChatRoom 특정 채팅방 조회
// GET /chat/rooms/:id
func GetChatRoom(c *fiber.Ctx) error {
	roomID := c.Params("id")

	var chatRoom models.ChatRoom
	if err := database.DB.
		Preload("Members.User").
		Preload("Creator").
		Preload("Club").
		First(&chatRoom, roomID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Chat room not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    chatRoom,
	})
}

// SendMessageRequest 메시지 전송 요청
type SendMessageRequest struct {
	UserID      uint   `json:"user_id" validate:"required"`
	Message     string `json:"message" validate:"required"`
	MessageType string `json:"message_type"` // text, image, file, system
	FileURL     string `json:"file_url"`
}

// SendMessage 메시지 전송
// POST /chat/rooms/:id/messages
func SendMessage(c *fiber.Ctx) error {
	roomID := c.Params("id")
	var req SendMessageRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// 채팅방 존재 확인
	var chatRoom models.ChatRoom
	if err := database.DB.First(&chatRoom, roomID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Chat room not found",
		})
	}

	// 사용자가 채팅방 멤버인지 확인
	var membership models.ChatRoomMember
	if err := database.DB.Where("chat_room_id = ? AND user_id = ?", roomID, req.UserID).First(&membership).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "User is not a member of this chat room",
		})
	}

	// 기본값 설정
	if req.MessageType == "" {
		req.MessageType = "text"
	}

	var fileURL *string
	if req.FileURL != "" {
		fileURL = &req.FileURL
	}

	// 메시지 생성
	message := models.ChatMessage{
		ChatRoomID:  chatRoom.ID,
		UserID:      req.UserID,
		Message:     req.Message,
		MessageType: req.MessageType,
		FileURL:     fileURL,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to send message",
			"details": err.Error(),
		})
	}

	// 채팅방의 last_message 및 last_message_at 업데이트
	now := time.Now()
	database.DB.Model(&chatRoom).Updates(map[string]interface{}{
		"last_message":    req.Message,
		"last_message_at": now,
	})

	// 다른 멤버들의 unread_count 증가
	database.DB.Model(&models.ChatRoomMember{}).
		Where("chat_room_id = ? AND user_id != ?", roomID, req.UserID).
		UpdateColumn("unread_count", database.DB.Raw("unread_count + 1"))

	// 메시지 정보 조회 (사용자 정보 포함)
	database.DB.Preload("User").First(&message, message.ID)

	// WebSocket으로 실시간 브로드캐스트
	if services.GlobalHub != nil {
		services.GlobalHub.BroadcastMessage(
			chatRoom.ID,
			"message",
			req.UserID,
			message,
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Message sent successfully",
		"data":    message,
	})
}

// GetMessages 채팅방의 메시지 목록 조회
// GET /chat/rooms/:id/messages
func GetMessages(c *fiber.Ctx) error {
	roomID := c.Params("id")

	// 페이지네이션
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	// 채팅방 존재 확인
	var chatRoom models.ChatRoom
	if err := database.DB.First(&chatRoom, roomID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Chat room not found",
		})
	}

	// 메시지 조회
	var messages []models.ChatMessage
	if err := database.DB.
		Preload("User").
		Where("chat_room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch messages",
		})
	}

	// 전체 메시지 수
	var total int64
	database.DB.Model(&models.ChatMessage{}).Where("chat_room_id = ?", roomID).Count(&total)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"messages": messages,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
		},
	})
}

// MarkAsRead 메시지 읽음 처리
// POST /chat/rooms/:id/read
func MarkAsRead(c *fiber.Ctx) error {
	roomID := c.Params("id")

	type ReadRequest struct {
		UserID uint `json:"user_id" validate:"required"`
	}

	var req ReadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// 멤버십 확인 및 업데이트
	var membership models.ChatRoomMember
	if err := database.DB.Where("chat_room_id = ? AND user_id = ?", roomID, req.UserID).First(&membership).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Membership not found",
		})
	}

	// last_read_at 업데이트 및 unread_count 초기화
	now := time.Now()
	database.DB.Model(&membership).Updates(map[string]interface{}{
		"last_read_at":  now,
		"unread_count": 0,
	})

	// WebSocket으로 읽음 처리 브로드캐스트
	roomIDUint, _ := strconv.ParseUint(roomID, 10, 32)
	if services.GlobalHub != nil {
		services.GlobalHub.BroadcastMessage(
			uint(roomIDUint),
			"read",
			req.UserID,
			fiber.Map{
				"user_id":      req.UserID,
				"last_read_at": now,
			},
		)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Messages marked as read",
	})
}

// AddMember 채팅방에 멤버 추가
// POST /chat/rooms/:id/members
func AddChatRoomMember(c *fiber.Ctx) error {
	roomID := c.Params("id")

	type AddMemberRequest struct {
		UserID uint `json:"user_id" validate:"required"`
	}

	var req AddMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// 채팅방 존재 확인
	var chatRoom models.ChatRoom
	if err := database.DB.First(&chatRoom, roomID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Chat room not found",
		})
	}

	// 이미 멤버인지 확인
	var existingMember models.ChatRoomMember
	if err := database.DB.Where("chat_room_id = ? AND user_id = ?", roomID, req.UserID).First(&existingMember).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error":   "User is already a member",
		})
	}

	// 멤버 추가
	member := models.ChatRoomMember{
		ChatRoomID: chatRoom.ID,
		UserID:     req.UserID,
		Role:       "member",
		JoinedAt:   time.Now(),
	}

	if err := database.DB.Create(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to add member",
		})
	}

	// 멤버 수 증가
	database.DB.Model(&chatRoom).UpdateColumn("member_count", database.DB.Raw("member_count + 1"))

	// 멤버 정보 조회
	database.DB.Preload("User").First(&member, member.ID)

	// WebSocket으로 멤버 추가 브로드캐스트
	if services.GlobalHub != nil {
		services.GlobalHub.BroadcastMessage(
			chatRoom.ID,
			"member_join",
			req.UserID,
			member,
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Member added successfully",
		"data":    member,
	})
}

// RemoveMember 채팅방에서 멤버 제거
// DELETE /chat/rooms/:id/members/:userId
func RemoveChatRoomMember(c *fiber.Ctx) error {
	roomID := c.Params("id")
	userID := c.Params("userId")

	// 멤버십 확인
	var member models.ChatRoomMember
	if err := database.DB.Where("chat_room_id = ? AND user_id = ?", roomID, userID).First(&member).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Member not found",
		})
	}

	// 멤버 삭제
	if err := database.DB.Delete(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to remove member",
		})
	}

	// 멤버 수 감소
	var chatRoom models.ChatRoom
	database.DB.First(&chatRoom, roomID)
	database.DB.Model(&chatRoom).UpdateColumn("member_count", database.DB.Raw("member_count - 1"))

	// WebSocket으로 멤버 제거 브로드캐스트
	roomIDUint, _ := strconv.ParseUint(roomID, 10, 32)
	userIDUint, _ := strconv.ParseUint(userID, 10, 32)
	if services.GlobalHub != nil {
		services.GlobalHub.BroadcastMessage(
			uint(roomIDUint),
			"member_leave",
			uint(userIDUint),
			fiber.Map{
				"user_id": userIDUint,
			},
		)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Member removed successfully",
	})
}

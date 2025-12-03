package models

import "time"

// ChatRoom 채팅방 (그룹 채팅)
type ChatRoom struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Name        string           `json:"name" gorm:"not null"`               // 채팅방 이름
	Description string           `json:"description" gorm:"type:text"`       // 채팅방 설명
	ClubID      *uint            `json:"club_id" gorm:"index"`               // 클럽 채팅방인 경우 (nullable)
	Club        *Club            `json:"club,omitempty" gorm:"foreignKey:ClubID"`
	RoomType    string           `json:"room_type" gorm:"default:'group'"`   // group, club, direct
	CreatedBy   uint             `json:"created_by" gorm:"not null"`         // 생성자 ID
	Creator     User             `json:"creator" gorm:"foreignKey:CreatedBy"`
	MemberCount int              `json:"member_count" gorm:"default:0"`      // 멤버 수
	LastMessage *string          `json:"last_message"`                       // 마지막 메시지
	LastMessageAt *time.Time     `json:"last_message_at"`                    // 마지막 메시지 시간
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Members     []ChatRoomMember `json:"members,omitempty" gorm:"foreignKey:ChatRoomID"`
	Messages    []ChatMessage    `json:"messages,omitempty" gorm:"foreignKey:ChatRoomID"`
}

// ChatRoomMember 채팅방 멤버
type ChatRoomMember struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	ChatRoomID    uint       `json:"chat_room_id" gorm:"not null;index"`
	ChatRoom      ChatRoom   `json:"-" gorm:"foreignKey:ChatRoomID"`
	UserID        uint       `json:"user_id" gorm:"not null;index"`
	User          User       `json:"user" gorm:"foreignKey:UserID"`
	Role          string     `json:"role" gorm:"default:'member'"`          // admin, member
	JoinedAt      time.Time  `json:"joined_at"`
	LastReadAt    *time.Time `json:"last_read_at"`                          // 마지막으로 읽은 시간
	UnreadCount   int        `json:"unread_count" gorm:"default:0"`         // 읽지 않은 메시지 수
	CreatedAt     time.Time  `json:"created_at"`
}

// ChatMessage 채팅 메시지
type ChatMessage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ChatRoomID uint      `json:"chat_room_id" gorm:"not null;index"`
	ChatRoom   ChatRoom  `json:"-" gorm:"foreignKey:ChatRoomID"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Message    string    `json:"message" gorm:"type:text;not null"`         // 메시지 내용
	MessageType string   `json:"message_type" gorm:"default:'text'"`        // text, image, file, system
	FileURL    *string   `json:"file_url"`                                  // 파일/이미지 URL (nullable)
	IsRead     bool      `json:"is_read" gorm:"default:false"`              // 읽음 여부
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

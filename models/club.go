package models

import "time"

type Club struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Category    string    `json:"category"`     // 클럽 카테고리 (운동, 문화, 학습 등)
	ImageURL    string    `json:"image_url"`
	MemberCount int       `json:"member_count"` // 멤버 수
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Members     []ClubMember `json:"members" gorm:"foreignKey:ClubID"`
}

type ClubMember struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ClubID    uint      `json:"club_id" gorm:"not null"`
	Club      Club      `json:"-" gorm:"foreignKey:ClubID"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	JoinedAt  time.Time `json:"joined_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Meeting struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	ClubID      uint      `json:"club_id"`
	Club        Club      `json:"club" gorm:"foreignKey:ClubID"`
	Location    string    `json:"location"`
	ScheduledAt time.Time `json:"scheduled_at"`
	MaxMembers  int       `json:"max_members"`
	Category    string    `json:"category"` // 모임 카테고리
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

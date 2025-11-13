package models

import (
	"time"
)

// GuestSession - 비회원 설문 세션
type GuestSession struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	SocialityScore  float64   `json:"sociality_score"`
	ActivityScore   float64   `json:"activity_score"`
	IntimacyScore   float64   `json:"intimacy_score"`
	ImmersionScore  float64   `json:"immersion_score"`
	FlexibilityScore float64  `json:"flexibility_score"`
	ProfileType     string    `json:"profile_type"`
	ResultSummary   string    `json:"result_summary" gorm:"type:text"`
	IsLinked        bool      `json:"is_linked" gorm:"default:false"` // 계정 연동 여부
	LinkedUserID    *uint     `json:"linked_user_id"`                 // 연동된 사용자 ID (nullable)
	ExpiresAt       time.Time `json:"expires_at"`                     // 세션 만료 시간
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// GuestAnswer - 비회원 답변
type GuestAnswer struct {
	ID         uint         `json:"id" gorm:"primaryKey"`
	SessionID  string       `json:"session_id" gorm:"not null"`
	Session    GuestSession `json:"-" gorm:"foreignKey:SessionID"`
	QuestionID uint         `json:"question_id" gorm:"not null"`
	Question   Question     `json:"question" gorm:"foreignKey:QuestionID"`
	OptionID   uint         `json:"option_id" gorm:"not null"`
	Option     Option       `json:"option" gorm:"foreignKey:OptionID"`
	CreatedAt  time.Time    `json:"created_at"`
}

// SessionVector - 세션의 벡터 표현 (빠른 유사도 계산용)
type SessionVector struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SessionID string    `json:"session_id" gorm:"uniqueIndex;not null"`
	UserID    *uint     `json:"user_id" gorm:"index"` // nullable, 회원인 경우
	Vector    []float64 `json:"vector" gorm:"type:jsonb;serializer:json"` // [sociality, activity, intimacy, immersion, flexibility]
	Magnitude float64   `json:"magnitude"` // 벡터 크기 (미리 계산)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

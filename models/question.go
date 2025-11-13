package models

import "time"

type Question struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	QuestionText string   `json:"question_text" gorm:"not null;type:text"`
	Order       int       `json:"order" gorm:"not null"` // 질문 순서 (1-10)
	Category    string    `json:"category"`              // 측정 카테고리 (sociality, activity, intimacy, immersion, flexibility)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Options     []Option  `json:"options" gorm:"foreignKey:QuestionID"`
}

type Option struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	QuestionID uint      `json:"question_id" gorm:"not null"`
	OptionText string    `json:"option_text" gorm:"not null;type:text"`
	Score      int       `json:"score"`       // 점수 (1-5)
	Weight     string    `json:"weight"`      // 가중치 카테고리
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AnswerPayload - 답변 제출용 DTO
type AnswerPayload struct {
	QuestionID uint `json:"question_id"`
	OptionID   uint `json:"option_id"`
}

type UserAnswer struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	User       User      `json:"-" gorm:"foreignKey:UserID"`
	QuestionID uint      `json:"question_id" gorm:"not null"`
	Question   Question  `json:"question" gorm:"foreignKey:QuestionID"`
	OptionID   uint      `json:"option_id" gorm:"not null"`
	Option     Option    `json:"option" gorm:"foreignKey:OptionID"`
	CreatedAt  time.Time `json:"created_at"`
}

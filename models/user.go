package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Password  string    `json:"-" gorm:"default:null"` // 비밀번호 (카카오 로그인 시 null)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfile struct {
	ID              uint    `json:"id" gorm:"primaryKey"`
	UserID          uint    `json:"user_id" gorm:"uniqueIndex;not null"`
	User            User    `json:"-" gorm:"foreignKey:UserID"`
	SocialityScore  float64 `json:"sociality_score"`   // 사교성
	ActivityScore   float64 `json:"activity_score"`    // 활동성
	IntimacyScore   float64 `json:"intimacy_score"`    // 친밀도
	ImmersionScore  float64 `json:"immersion_score"`   // 몰입도
	FlexibilityScore float64 `json:"flexibility_score"` // 유연성
	ResultSummary   string  `json:"result_summary" gorm:"type:text"`
	ProfileType     string  `json:"profile_type"` // 성향 유형
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

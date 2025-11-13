package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/utils"
	"time"

	"gorm.io/gorm"
)

// GenerateSessionID - 고유한 세션 ID 생성
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateGuestSession - 비회원 세션 생성
func CreateGuestSession() (*models.GuestSession, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, err
	}

	session := &models.GuestSession{
		ID:        sessionID,
		IsLinked:  false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7일 후 만료
	}

	err = database.DB.Create(session).Error
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetGuestSession - 세션 조회
func GetGuestSession(sessionID string) (*models.GuestSession, error) {
	var session models.GuestSession
	err := database.DB.Where("id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}

	// 만료된 세션 체크
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	return &session, nil
}

// SubmitGuestAnswers - 비회원 답변 제출
func SubmitGuestAnswers(sessionID string, answers []models.AnswerPayload) error {
	// 기존 답변 삭제
	database.DB.Where("session_id = ?", sessionID).Delete(&models.GuestAnswer{})

	// 새 답변 저장
	for _, ans := range answers {
		answer := models.GuestAnswer{
			SessionID:  sessionID,
			QuestionID: ans.QuestionID,
			OptionID:   ans.OptionID,
		}

		if err := database.DB.Create(&answer).Error; err != nil {
			return err
		}
	}

	return nil
}

// CalculateGuestScores - 비회원 세션 점수 계산
func CalculateGuestScores(sessionID string) (*ScoreResult, error) {
	var answers []models.GuestAnswer

	err := database.DB.Preload("Option").
		Where("session_id = ?", sessionID).
		Find(&answers).Error

	if err != nil {
		return nil, err
	}

	if len(answers) == 0 {
		return nil, fmt.Errorf("no answers found for session")
	}

	scores := &ScoreResult{}
	categoryScores := make(map[string][]int)

	// 각 답변의 점수를 카테고리별로 분류
	for _, answer := range answers {
		weight := answer.Option.Weight
		score := answer.Option.Score

		if weight != "" {
			categoryScores[weight] = append(categoryScores[weight], score)
		}
	}

	// 카테고리별 평균 계산
	scores.SocialityScore = calculateAverage(categoryScores["sociality"])
	scores.ActivityScore = calculateAverage(categoryScores["activity"])
	scores.IntimacyScore = calculateAverage(categoryScores["intimacy"])
	scores.ImmersionScore = calculateAverage(categoryScores["immersion"])
	scores.FlexibilityScore = calculateAverage(categoryScores["flexibility"])

	return scores, nil
}

// SaveGuestResult - 비회원 세션 결과 저장
func SaveGuestResult(sessionID string, scores *ScoreResult, profileType string, summary string) error {
	return database.DB.Model(&models.GuestSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"sociality_score":   scores.SocialityScore,
			"activity_score":    scores.ActivityScore,
			"intimacy_score":    scores.IntimacyScore,
			"immersion_score":   scores.ImmersionScore,
			"flexibility_score": scores.FlexibilityScore,
			"profile_type":      profileType,
			"result_summary":    summary,
		}).Error
}

// CreateSessionVector - 세션 벡터 생성/업데이트
func CreateSessionVector(sessionID string, userID *uint, scores *ScoreResult) error {
	vector := []float64{
		scores.SocialityScore,
		scores.ActivityScore,
		scores.IntimacyScore,
		scores.ImmersionScore,
		scores.FlexibilityScore,
	}

	v := utils.FromSlice(vector)
	magnitude := v.Magnitude()

	sessionVector := models.SessionVector{
		SessionID: sessionID,
		UserID:    userID,
		Vector:    vector,
		Magnitude: magnitude,
	}

	// Upsert
	var existing models.SessionVector
	result := database.DB.Where("session_id = ?", sessionID).First(&existing)

	if result.Error == nil {
		// 업데이트
		return database.DB.Model(&existing).Updates(sessionVector).Error
	} else {
		// 생성
		return database.DB.Create(&sessionVector).Error
	}
}

// LinkSessionToUser - 세션을 사용자 계정과 연동
func LinkSessionToUser(sessionID string, userID uint) error {
	// 트랜잭션으로 처리
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 세션을 사용자와 연결
		err := tx.Model(&models.GuestSession{}).
			Where("id = ?", sessionID).
			Updates(map[string]interface{}{
				"is_linked":      true,
				"linked_user_id": userID,
			}).Error
		if err != nil {
			return err
		}

		// 2. 세션의 답변을 사용자 답변으로 복사
		var guestAnswers []models.GuestAnswer
		err = tx.Preload("Option").Where("session_id = ?", sessionID).Find(&guestAnswers).Error
		if err != nil {
			return err
		}

		// 기존 사용자 답변 삭제
		tx.Where("user_id = ?", userID).Delete(&models.UserAnswer{})

		// 새 답변 생성
		for _, ga := range guestAnswers {
			userAnswer := models.UserAnswer{
				UserID:     userID,
				QuestionID: ga.QuestionID,
				OptionID:   ga.OptionID,
			}
			if err := tx.Create(&userAnswer).Error; err != nil {
				return err
			}
		}

		// 3. 세션 점수를 UserProfile로 복사
		var session models.GuestSession
		err = tx.Where("id = ?", sessionID).First(&session).Error
		if err != nil {
			return err
		}

		profile := models.UserProfile{
			UserID:           userID,
			SocialityScore:   session.SocialityScore,
			ActivityScore:    session.ActivityScore,
			IntimacyScore:    session.IntimacyScore,
			ImmersionScore:   session.ImmersionScore,
			FlexibilityScore: session.FlexibilityScore,
			ProfileType:      session.ProfileType,
			ResultSummary:    session.ResultSummary,
		}

		// UserProfile upsert
		var existingProfile models.UserProfile
		result := tx.Where("user_id = ?", userID).First(&existingProfile)

		if result.Error == nil {
			err = tx.Model(&existingProfile).Updates(profile).Error
		} else {
			err = tx.Create(&profile).Error
		}

		if err != nil {
			return err
		}

		// 4. SessionVector 업데이트
		return tx.Model(&models.SessionVector{}).
			Where("session_id = ?", sessionID).
			Update("user_id", userID).Error
	})
}

// CleanExpiredSessions - 만료된 세션 정리
func CleanExpiredSessions() error {
	now := time.Now()

	// 만료된 세션의 SessionVector 삭제
	err := database.DB.
		Where("session_id IN (?)",
			database.DB.Model(&models.GuestSession{}).
				Select("id").
				Where("expires_at < ? AND is_linked = false", now),
		).
		Delete(&models.SessionVector{}).Error

	if err != nil {
		return err
	}

	// 만료된 세션의 GuestAnswer 삭제
	err = database.DB.
		Where("session_id IN (?)",
			database.DB.Model(&models.GuestSession{}).
				Select("id").
				Where("expires_at < ? AND is_linked = false", now),
		).
		Delete(&models.GuestAnswer{}).Error

	if err != nil {
		return err
	}

	// 만료된 세션 삭제 (연동되지 않은 것만)
	return database.DB.
		Where("expires_at < ? AND is_linked = false", now).
		Delete(&models.GuestSession{}).Error
}

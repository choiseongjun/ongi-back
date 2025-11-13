package handlers

import (
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/services"
	"ongi-back/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CreateGuestSession - 비회원 세션 생성
func CreateGuestSession(c *fiber.Ctx) error {
	session, err := services.CreateGuestSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create guest session",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"session_id": session.ID,
			"expires_at": session.ExpiresAt,
		},
		"message": "Guest session created. Save this session_id to retrieve your results later.",
	})
}

// SubmitGuestAnswers - 비회원 답변 제출
func SubmitGuestAnswers(c *fiber.Ctx) error {
	var req struct {
		SessionID string                `json:"session_id"`
		Answers   []models.AnswerPayload `json:"answers"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 세션 확인
	_, err := services.GetGuestSession(req.SessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Session not found or expired",
		})
	}

	// 답변 저장
	err = services.SubmitGuestAnswers(req.SessionID, req.Answers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save answers",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Answers submitted successfully",
	})
}

// GetGuestResult - 비회원 결과 조회
func GetGuestResult(c *fiber.Ctx) error {
	sessionID := c.Params("sessionId")

	// 세션 확인
	session, err := services.GetGuestSession(sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Session not found or expired",
		})
	}

	// 이미 계산된 결과가 있는지 확인
	if session.ProfileType != "" {
		// 캐시된 결과 반환
		return returnGuestResult(c, sessionID, session)
	}

	// 점수 계산
	scores, err := services.CalculateGuestScores(sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate scores: " + err.Error(),
		})
	}

	// 프로필 타입 및 설명 생성
	profileType := services.DetermineProfileType(scores)
	descriptions := services.GenerateDescriptions(scores)
	summary := strings.Join(descriptions, " ")

	// 결과 저장
	err = services.SaveGuestResult(sessionID, scores, profileType, summary)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save result",
		})
	}

	// 벡터 생성
	err = services.CreateSessionVector(sessionID, nil, scores)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session vector",
		})
	}

	// 세션 다시 조회
	session, _ = services.GetGuestSession(sessionID)
	return returnGuestResult(c, sessionID, session)
}

func returnGuestResult(c *fiber.Ctx, sessionID string, session *models.GuestSession) error {
	scores := &services.ScoreResult{
		SocialityScore:   session.SocialityScore,
		ActivityScore:    session.ActivityScore,
		IntimacyScore:    session.IntimacyScore,
		ImmersionScore:   session.ImmersionScore,
		FlexibilityScore: session.FlexibilityScore,
	}

	descriptions := strings.Split(session.ResultSummary, " ")
	if len(descriptions) == 1 && descriptions[0] == "" {
		descriptions = services.GenerateDescriptions(scores)
	}

	// 추천 데이터 가져오기
	recommendedClubs, _ := services.GetRecommendedClubsForSession(sessionID, 5)
	similarClubs, _ := services.GetClubsWithSimilarMembersForSession(sessionID, 5)
	recommendedMeetings, _ := services.GetRecommendedMeetingsForSession(sessionID, 5)
	similarProfiles, _ := services.GetSimilarProfilesFast(sessionID, 5)

	result := fiber.Map{
		"session_id":  sessionID,
		"is_linked":   session.IsLinked,
		"scores":      scores,
		"profile_type": session.ProfileType,
		"descriptions": descriptions,
		"recommendations": fiber.Map{
			"clubs":           recommendedClubs,
			"similar_clubs":   similarClubs,
			"meetings":        recommendedMeetings,
			"similar_profiles": similarProfiles,
		},
		"expires_at": session.ExpiresAt,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// LinkSessionToAccount - 세션을 계정과 연동
func LinkSessionToAccount(c *fiber.Ctx) error {
	var req struct {
		SessionID string `json:"session_id"`
		UserID    uint   `json:"user_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 세션 확인
	session, err := services.GetGuestSession(req.SessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Session not found or expired",
		})
	}

	// 이미 연동된 세션인지 확인
	if session.IsLinked {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session is already linked to an account",
		})
	}

	// 사용자 확인
	var user models.User
	err = database.DB.First(&user, req.UserID).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// 연동 처리
	err = services.LinkSessionToUser(req.SessionID, req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to link session: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Session successfully linked to user account",
		"data": fiber.Map{
			"user_id":    req.UserID,
			"session_id": req.SessionID,
		},
	})
}

// GetSessionInfo - 세션 정보 조회
func GetSessionInfo(c *fiber.Ctx) error {
	sessionID := c.Params("sessionId")

	session, err := services.GetGuestSession(sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Session not found or expired",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"session_id":      session.ID,
			"is_linked":       session.IsLinked,
			"linked_user_id":  session.LinkedUserID,
			"has_results":     session.ProfileType != "",
			"profile_type":    session.ProfileType,
			"expires_at":      session.ExpiresAt,
			"created_at":      session.CreatedAt,
		},
	})
}

// GetCompatibility - 두 프로필 간 궁합 계산
func GetCompatibility(c *fiber.Ctx) error {
	var req struct {
		SessionID1 string `json:"session_id_1"`
		SessionID2 string `json:"session_id_2"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 두 세션의 벡터 가져오기
	var vector1, vector2 models.SessionVector
	err1 := database.DB.Where("session_id = ?", req.SessionID1).First(&vector1).Error
	err2 := database.DB.Where("session_id = ?", req.SessionID2).First(&vector2).Error

	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "One or both sessions not found",
		})
	}

	// 벡터로 변환
	v1 := utils.FromSlice(vector1.Vector)
	v2 := utils.FromSlice(vector2.Vector)

	if v1 == nil || v2 == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid vector data",
		})
	}

	// 궁합 계산
	compatibility := services.CalculateProfileCompatibility(v1, v2)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    compatibility,
	})
}

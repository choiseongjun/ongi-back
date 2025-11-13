package handlers

import (
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/services"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// 분석 결과 생성 및 조회
func GetAnalysisResult(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// 점수 계산
	scores, err := services.CalculateScores(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate scores: " + err.Error(),
		})
	}

	// 프로필 타입 및 설명 생성
	profileType := services.DetermineProfileType(scores)
	descriptions := services.GenerateDescriptions(scores)

	// UserProfile 저장 또는 업데이트
	profile := models.UserProfile{
		UserID:           uint(userID),
		SocialityScore:   scores.SocialityScore,
		ActivityScore:    scores.ActivityScore,
		IntimacyScore:    scores.IntimacyScore,
		ImmersionScore:   scores.ImmersionScore,
		FlexibilityScore: scores.FlexibilityScore,
		ProfileType:      profileType,
		ResultSummary:    strings.Join(descriptions, " "),
	}

	// upsert (존재하면 업데이트, 없으면 생성)
	var existingProfile models.UserProfile
	result := database.DB.Where("user_id = ?", userID).First(&existingProfile)

	if result.Error == nil {
		// 업데이트
		database.DB.Model(&existingProfile).Updates(profile)
	} else {
		// 생성
		database.DB.Create(&profile)
	}

	// 추천 데이터 가져오기
	recommendedClubs, _ := services.GetRecommendedClubs(uint(userID), 5)
	similarClubs, _ := services.GetClubsWithSimilarMembers(uint(userID), 5)
	recommendedMeetings, _ := services.GetRecommendedMeetings(uint(userID), 5)
	similarUsers, _ := services.GetSimilarUsers(uint(userID), 5)

	analysisResult := fiber.Map{
		"scores":       scores,
		"profile_type": profileType,
		"descriptions": descriptions,
		"recommendations": fiber.Map{
			"clubs":          recommendedClubs,
			"similar_clubs":  similarClubs,
			"meetings":       recommendedMeetings,
			"similar_users":  similarUsers,
		},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    analysisResult,
	})
}

// 사용자의 모든 답변 조회
func GetUserAnswers(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var answers []models.UserAnswer
	err := database.DB.Preload("Question").Preload("Option").
		Where("user_id = ?", userID).
		Find(&answers).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch answers",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    answers,
	})
}

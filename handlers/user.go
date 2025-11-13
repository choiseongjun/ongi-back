package handlers

import (
	"fmt"
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/services"

	"github.com/gofiber/fiber/v2"
)

// 사용자 생성
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user := models.User{
		Email: req.Email,
		Name:  req.Name,
	}

	err := database.DB.Create(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// 사용자 조회
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// 모든 사용자 조회
func GetUsers(c *fiber.Ctx) error {
	var users []models.User

	err := database.DB.Find(&users).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

// 성향 분석 함수
func analyzeTendency(score float64) string {
	if score >= 80 {
		return "매우 높음"
	} else if score >= 60 {
		return "높음"
	} else if score >= 40 {
		return "보통"
	} else if score >= 20 {
		return "낮음"
	}
	return "매우 낮음"
}

// 사용자 프로필 조회
func GetUserProfile(c *fiber.Ctx) error {
	userID := c.Params("id")

	var profile models.UserProfile
	err := database.DB.Preload("User").Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Profile not found",
		})
	}

	// 사용자 성향 분석
	tendencies := fiber.Map{
		"sociality":   fiber.Map{"score": profile.SocialityScore, "level": analyzeTendency(profile.SocialityScore)},
		"activity":    fiber.Map{"score": profile.ActivityScore, "level": analyzeTendency(profile.ActivityScore)},
		"intimacy":    fiber.Map{"score": profile.IntimacyScore, "level": analyzeTendency(profile.IntimacyScore)},
		"immersion":   fiber.Map{"score": profile.ImmersionScore, "level": analyzeTendency(profile.ImmersionScore)},
		"flexibility": fiber.Map{"score": profile.FlexibilityScore, "level": analyzeTendency(profile.FlexibilityScore)},
	}

	// 유사 사용자 추천 (70% 이상 유사도)
	var uid uint
	if _, err := fmt.Sscanf(userID, "%d", &uid); err == nil {
		similarUsers, _ := services.GetSimilarUsers(uid, 20) // 상위 20명

		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"profile":      profile,
				"tendencies":   tendencies,
				"similar_users": similarUsers,
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"profile":    profile,
			"tendencies": tendencies,
		},
	})
}

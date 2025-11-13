package handlers

import (
	"ongi-back/database"
	"ongi-back/models"

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

	return c.JSON(fiber.Map{
		"success": true,
		"data":    profile,
	})
}

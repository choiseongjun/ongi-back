package handlers

import (
	"ongi-back/database"
	"ongi-back/models"

	"github.com/gofiber/fiber/v2"
)

// 모든 클럽 조회
func GetClubs(c *fiber.Ctx) error {
	var clubs []models.Club

	err := database.DB.Preload("Members").Find(&clubs).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch clubs",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    clubs,
	})
}

// 특정 클럽 조회
func GetClub(c *fiber.Ctx) error {
	id := c.Params("id")

	var club models.Club
	err := database.DB.Preload("Members.User").First(&club, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Club not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    club,
	})
}

// 클럽 생성
type CreateClubRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
}

func CreateClub(c *fiber.Ctx) error {
	var req CreateClubRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	club := models.Club{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		MemberCount: 0,
	}

	err := database.DB.Create(&club).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create club",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    club,
	})
}

// 클럽 가입
type JoinClubRequest struct {
	UserID uint `json:"user_id"`
	ClubID uint `json:"club_id"`
}

func JoinClub(c *fiber.Ctx) error {
	var req JoinClubRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 이미 가입했는지 확인
	var existing models.ClubMember
	result := database.DB.Where("user_id = ? AND club_id = ?", req.UserID, req.ClubID).First(&existing)
	if result.Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Already a member of this club",
		})
	}

	member := models.ClubMember{
		UserID: req.UserID,
		ClubID: req.ClubID,
	}

	err := database.DB.Create(&member).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join club",
		})
	}

	// 클럽의 멤버 수 증가
	database.DB.Model(&models.Club{}).Where("id = ?", req.ClubID).
		Update("member_count", database.DB.Raw("member_count + 1"))

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully joined club",
		"data":    member,
	})
}

// 모든 모임 조회
func GetMeetings(c *fiber.Ctx) error {
	var meetings []models.Meeting

	err := database.DB.Preload("Club").Find(&meetings).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch meetings",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    meetings,
	})
}

// 특정 모임 조회
func GetMeeting(c *fiber.Ctx) error {
	id := c.Params("id")

	var meeting models.Meeting
	err := database.DB.Preload("Club").First(&meeting, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Meeting not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    meeting,
	})
}

// 모임 생성
type CreateMeetingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ClubID      uint   `json:"club_id"`
	Location    string `json:"location"`
	ScheduledAt string `json:"scheduled_at"`
	MaxMembers  int    `json:"max_members"`
	Category    string `json:"category"`
}

func CreateMeeting(c *fiber.Ctx) error {
	var req CreateMeetingRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	meeting := models.Meeting{
		Title:       req.Title,
		Description: req.Description,
		ClubID:      req.ClubID,
		Location:    req.Location,
		MaxMembers:  req.MaxMembers,
		Category:    req.Category,
	}

	err := database.DB.Create(&meeting).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create meeting",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    meeting,
	})
}

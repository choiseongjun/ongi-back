package handlers

import (
	"ongi-back/database"
	"ongi-back/models"

	"github.com/gofiber/fiber/v2"
)

// 모든 질문 가져오기
func GetQuestions(c *fiber.Ctx) error {
	var questions []models.Question

	err := database.DB.Preload("Options").Order("\"order\" ASC").Find(&questions).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch questions",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    questions,
	})
}

// 특정 질문 가져오기
func GetQuestion(c *fiber.Ctx) error {
	id := c.Params("id")

	var question models.Question
	err := database.DB.Preload("Options").First(&question, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Question not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    question,
	})
}

// 사용자 답변 제출
type SubmitAnswerRequest struct {
	UserID     uint `json:"user_id"`
	QuestionID uint `json:"question_id"`
	OptionID   uint `json:"option_id"`
}

func SubmitAnswer(c *fiber.Ctx) error {
	var req SubmitAnswerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 답변 저장
	answer := models.UserAnswer{
		UserID:     req.UserID,
		QuestionID: req.QuestionID,
		OptionID:   req.OptionID,
	}

	err := database.DB.Create(&answer).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save answer",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Answer submitted successfully",
		"data":    answer,
	})
}

// 여러 답변 한번에 제출
type SubmitAnswersRequest struct {
	UserID  uint              `json:"user_id"`
	Answers []AnswerSubmission `json:"answers"`
}

type AnswerSubmission struct {
	QuestionID uint `json:"question_id"`
	OptionID   uint `json:"option_id"`
}

func SubmitAnswers(c *fiber.Ctx) error {
	var req SubmitAnswersRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 기존 답변 삭제 (재시험 가능하도록)
	database.DB.Where("user_id = ?", req.UserID).Delete(&models.UserAnswer{})

	// 새 답변들 저장
	for _, ans := range req.Answers {
		answer := models.UserAnswer{
			UserID:     req.UserID,
			QuestionID: ans.QuestionID,
			OptionID:   ans.OptionID,
		}

		if err := database.DB.Create(&answer).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save answers",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "All answers submitted successfully",
	})
}

package routes

import (
	"ongi-back/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// Guest/Session routes (비회원 설문)
	guest := api.Group("/guest")
	guest.Post("/session", handlers.CreateGuestSession)           // 세션 생성
	guest.Post("/answers", handlers.SubmitGuestAnswers)            // 답변 제출
	guest.Get("/result/:sessionId", handlers.GetGuestResult)       // 결과 조회
	guest.Get("/session/:sessionId", handlers.GetSessionInfo)      // 세션 정보
	guest.Post("/link", handlers.LinkSessionToAccount)             // 계정 연동
	guest.Post("/compatibility", handlers.GetCompatibility)        // 궁합 계산

	// User routes
	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Post("/", handlers.CreateUser)
	users.Get("/:id", handlers.GetUser)
	users.Get("/:id/profile", handlers.GetUserProfile)

	// Question routes
	questions := api.Group("/questions")
	questions.Get("/", handlers.GetQuestions)
	questions.Get("/:id", handlers.GetQuestion)

	// Answer routes
	answers := api.Group("/answers")
	answers.Post("/", handlers.SubmitAnswer)
	answers.Post("/batch", handlers.SubmitAnswers)
	answers.Get("/user/:userId", handlers.GetUserAnswers)

	// Result routes
	results := api.Group("/results")
	results.Get("/:userId", handlers.GetAnalysisResult)

	// Club routes
	clubs := api.Group("/clubs")
	clubs.Get("/", handlers.GetClubs)
	clubs.Post("/", handlers.CreateClub)
	clubs.Get("/:id", handlers.GetClub)
	clubs.Post("/join", handlers.JoinClub)

	// Meeting routes
	meetings := api.Group("/meetings")
	meetings.Get("/", handlers.GetMeetings)
	meetings.Post("/", handlers.CreateMeeting)
	meetings.Get("/:id", handlers.GetMeeting)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Server is running",
		})
	})
}

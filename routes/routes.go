package routes

import (
	"ongi-back/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// Auth routes (인증)
	auth := api.Group("/auth")
	auth.Post("/kakao/login", handlers.KakaoLogin)

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
	users.Post("/profile", handlers.CreateOrUpdateUserProfile)
	users.Get("/:id/profile", handlers.GetUserProfile)
	users.Post("/:id/auto-match", handlers.AutoMatchClubs)
	users.Post("/:id/auto-match-group", handlers.AutoMatchWithSimilarUsers)

	// Matching routes - 전체 사용자 그룹 매칭
	api.Post("/match-all", handlers.MatchAllUsersToClubs)

	// Chat routes (그룹 채팅)
	chat := api.Group("/chat")
	chat.Post("/rooms", handlers.CreateChatRoom)                         // 채팅방 생성
	chat.Get("/rooms", handlers.GetChatRooms)                            // 채팅방 목록 조회
	chat.Get("/rooms/:id", handlers.GetChatRoom)                         // 채팅방 상세 조회
	chat.Post("/rooms/:id/messages", handlers.SendMessage)               // 메시지 전송
	chat.Get("/rooms/:id/messages", handlers.GetMessages)                // 메시지 목록 조회
	chat.Post("/rooms/:id/read", handlers.MarkAsRead)                    // 메시지 읽음 처리
	chat.Post("/rooms/:id/members", handlers.AddChatRoomMember)          // 멤버 추가
	chat.Delete("/rooms/:id/members/:userId", handlers.RemoveChatRoomMember) // 멤버 제거

	// WebSocket route (실시간 채팅)
	app.Use("/ws", handlers.WebSocketHandler)
	app.Get("/ws/chat/:roomId", websocket.New(handlers.HandleWebSocket))

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

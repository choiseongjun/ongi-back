package handlers

import (
	"fmt"
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/services"
	"ongi-back/utils"

	"github.com/gofiber/fiber/v2"
)

// KakaoLoginRequest 카카오 로그인 요청
type KakaoLoginRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

// KakaoLoginResponse 카카오 로그인 응답
type KakaoLoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	User    models.User `json:"user"`
	IsNewUser bool `json:"is_new_user"`
}

// KakaoLogin 카카오 로그인 처리
// POST /auth/kakao/login
func KakaoLogin(c *fiber.Ctx) error {
	var req KakaoLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.AccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Access token is required",
		})
	}

	// 1. 카카오 API 호출 (Access Token 검증)
	kakaoUserInfo, err := services.ValidateKakaoToken(req.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid kakao access token",
			"details": err.Error(),
		})
	}

	// 2. 회원가입/로그인 처리
	email := kakaoUserInfo.KakaoAccount.Email
	if email == "" {
		// 이메일이 없는 경우 카카오 ID 사용
		email = fmt.Sprintf("kakao_%d@kakao.com", kakaoUserInfo.ID)
	}

	name := kakaoUserInfo.KakaoAccount.Profile.Nickname
	if name == "" {
		name = fmt.Sprintf("User_%d", kakaoUserInfo.ID)
	}

	// 기존 사용자 확인
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)

	isNewUser := false
	if result.Error != nil {
		// 신규 사용자 - 회원가입
		user = models.User{
			Email: email,
			Name:  name,
		}

		if err := database.DB.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to create user",
				"details": err.Error(),
			})
		}
		isNewUser = true
	}

	// 3. JWT 토큰 발급
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate token",
			"details": err.Error(),
		})
	}

	// 4. 응답 반환
	return c.Status(fiber.StatusOK).JSON(KakaoLoginResponse{
		Success:   true,
		Token:     token,
		User:      user,
		IsNewUser: isNewUser,
	})
}

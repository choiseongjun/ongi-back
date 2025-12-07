package handlers

import (
	"fmt"
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/services"
	"ongi-back/utils"

	"github.com/gofiber/fiber/v2"
)

// RegisterRequest 회원가입 요청
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

// LoginRequest 로그인 요청
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse 인증 응답
type AuthResponse struct {
	Success   bool        `json:"success"`
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	IsNewUser bool        `json:"is_new_user,omitempty"`
}

// KakaoLoginRequest 카카오 로그인 요청
type KakaoLoginRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

// KakaoLoginResponse 카카오 로그인 응답
type KakaoLoginResponse struct {
	Success   bool        `json:"success"`
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	IsNewUser bool        `json:"is_new_user"`
}

// KakaoLogin 카카오 로그인 처리 (클라이언트사이드 OAuth)
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
	user, isNewUser, err := processKakaoUser(kakaoUserInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to process user",
			"details": err.Error(),
		})
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

// KakaoCallback 카카오 OAuth 콜백 처리 (서버사이드 OAuth)
// GET /auth/kakao/callback
func KakaoCallback(c *fiber.Ctx) error {
	// 1. Authorization Code 가져오기
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Authorization code is required",
		})
	}

	// 에러 체크
	errorParam := c.Query("error")
	if errorParam != "" {
		errorDescription := c.Query("error_description")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   errorParam,
			"details": errorDescription,
		})
	}

	// 2. Authorization Code를 Access Token으로 교환
	tokenResp, err := services.ExchangeCodeForToken(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to exchange code for token",
			"details": err.Error(),
		})
	}

	// 3. Access Token으로 사용자 정보 가져오기
	kakaoUserInfo, err := services.ValidateKakaoToken(tokenResp.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to get user info",
			"details": err.Error(),
		})
	}

	// 4. 회원가입/로그인 처리
	user, isNewUser, err := processKakaoUser(kakaoUserInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to process user",
			"details": err.Error(),
		})
	}

	// 5. JWT 토큰 발급
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate token",
			"details": err.Error(),
		})
	}

	// 6. 응답 반환
	return c.Status(fiber.StatusOK).JSON(KakaoLoginResponse{
		Success:   true,
		Token:     token,
		User:      user,
		IsNewUser: isNewUser,
	})
}

// processKakaoUser 카카오 사용자 정보로 회원가입/로그인 처리
func processKakaoUser(kakaoUserInfo *services.KakaoUserInfo) (models.User, bool, error) {
	// 이메일 결정
	email := kakaoUserInfo.KakaoAccount.Email
	if email == "" {
		// 이메일이 없는 경우 카카오 ID 사용
		email = fmt.Sprintf("kakao_%d@kakao.com", kakaoUserInfo.ID)
	}

	// 이름 결정
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
			return user, false, fmt.Errorf("failed to create user: %w", err)
		}
		isNewUser = true
	}

	return user, isNewUser, nil
}

// Register 일반 회원가입
// POST /auth/register
func Register(c *fiber.Ctx) error {
	var req RegisterRequest

	// 요청 파싱
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// 필수 필드 검증
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Email, password, and name are required",
		})
	}

	// 비밀번호 길이 검증
	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Password must be at least 6 characters",
		})
	}

	// 이메일 중복 확인
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error":   "Email already exists",
		})
	}

	// 비밀번호 해싱
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to hash password",
		})
	}

	// 사용자 생성
	user := models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create user",
			"details": err.Error(),
		})
	}

	// JWT 토큰 발급
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate token",
		})
	}

	// 응답 반환
	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Success:   true,
		Token:     token,
		User:      user,
		IsNewUser: true,
	})
}

// Login 일반 로그인
// POST /auth/login
func Login(c *fiber.Ctx) error {
	var req LoginRequest

	// 요청 파싱
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// 필수 필드 검증
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Email and password are required",
		})
	}

	// 사용자 조회
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid email or password",
		})
	}

	// 비밀번호가 설정되지 않은 경우 (카카오 로그인 사용자)
	if user.Password == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "This account uses social login. Please use Kakao login.",
		})
	}

	// 비밀번호 검증
	if !utils.CheckPassword(user.Password, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid email or password",
		})
	}

	// JWT 토큰 발급
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate token",
		})
	}

	// 응답 반환
	return c.Status(fiber.StatusOK).JSON(AuthResponse{
		Success: true,
		Token:   token,
		User:    user,
	})
}

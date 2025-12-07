package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// KakaoUserInfo 카카오 사용자 정보
type KakaoUserInfo struct {
	ID           int64              `json:"id"`
	ConnectedAt  string             `json:"connected_at"`
	KakaoAccount KakaoAccount       `json:"kakao_account"`
	Properties   map[string]string  `json:"properties"`
}

// KakaoAccount 카카오 계정 정보
type KakaoAccount struct {
	Profile           KakaoProfile `json:"profile"`
	Email             string       `json:"email"`
	EmailNeedsAgreement bool       `json:"email_needs_agreement"`
}

// KakaoProfile 카카오 프로필 정보
type KakaoProfile struct {
	Nickname              string `json:"nickname"`
	ProfileImageURL       string `json:"profile_image_url"`
	ThumbnailImageURL     string `json:"thumbnail_image_url"`
}

// KakaoTokenResponse 카카오 토큰 응답
type KakaoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// ExchangeCodeForToken 카카오 Authorization Code를 Access Token으로 교환
func ExchangeCodeForToken(code string) (*KakaoTokenResponse, error) {
	tokenURL := "https://kauth.kakao.com/oauth/token"

	// 환경 변수에서 카카오 설정 가져오기
	clientID := os.Getenv("KAKAO_CLIENT_ID")
	clientSecret := os.Getenv("KAKAO_CLIENT_SECRET")
	redirectURI := os.Getenv("KAKAO_REDIRECT_URI")

	if clientID == "" {
		return nil, fmt.Errorf("KAKAO_CLIENT_ID is not set")
	}
	if redirectURI == "" {
		return nil, fmt.Errorf("KAKAO_REDIRECT_URI is not set")
	}

	// POST 요청 파라미터 구성
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)

	// client_secret은 선택사항 (보안 강화를 위해 사용)
	if clientSecret != "" {
		data.Set("client_secret", clientSecret)
	}

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kakao token API returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp KakaoTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// ValidateKakaoToken 카카오 Access Token 검증 및 사용자 정보 가져오기
func ValidateKakaoToken(accessToken string) (*KakaoUserInfo, error) {
	// 카카오 사용자 정보 가져오기 API
	apiURL := "https://kapi.kakao.com/v2/user/me"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call kakao API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kakao API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo KakaoUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse kakao response: %w", err)
	}

	return &userInfo, nil
}

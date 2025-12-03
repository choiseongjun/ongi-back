package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// ValidateKakaoToken 카카오 Access Token 검증 및 사용자 정보 가져오기
func ValidateKakaoToken(accessToken string) (*KakaoUserInfo, error) {
	// 카카오 사용자 정보 가져오기 API
	url := "https://kapi.kakao.com/v2/user/me"

	req, err := http.NewRequest("GET", url, nil)
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

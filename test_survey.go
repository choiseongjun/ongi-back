package main

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"log"
//	"math/rand"
//	"net/http"
//	"time"
//)
//
//const baseURL = "http://localhost:5000/api/v1"
//
//// 응답 구조체
//type Response struct {
//	Success bool            `json:"success"`
//	Data    json.RawMessage `json:"data"`
//	Message string          `json:"message,omitempty"`
//}
//
//type User struct {
//	ID    uint   `json:"id"`
//	Email string `json:"email"`
//	Name  string `json:"name"`
//}
//
//type Question struct {
//	ID           uint     `json:"id"`
//	QuestionText string   `json:"question_text"`
//	Order        int      `json:"order"`
//	Category     string   `json:"category"`
//	Options      []Option `json:"options"`
//}
//
//type Option struct {
//	ID         uint   `json:"id"`
//	QuestionID uint   `json:"question_id"`
//	OptionText string `json:"option_text"`
//	Score      int    `json:"score"`
//	Weight     string `json:"weight"`
//}
//
//type AnswerRequest struct {
//	UserID  uint            `json:"user_id"`
//	Answers []AnswerPayload `json:"answers"`
//}
//
//type AnswerPayload struct {
//	QuestionID uint `json:"question_id"`
//	OptionID   uint `json:"option_id"`
//}
//
//type Result struct {
//	Scores struct {
//		SocialityScore   float64 `json:"sociality_score"`
//		ActivityScore    float64 `json:"activity_score"`
//		IntimacyScore    float64 `json:"intimacy_score"`
//		ImmersionScore   float64 `json:"immersion_score"`
//		FlexibilityScore float64 `json:"flexibility_score"`
//	} `json:"scores"`
//	ProfileType     string   `json:"profile_type"`
//	Descriptions    []string `json:"descriptions"`
//	Recommendations struct {
//		Clubs        []Club `json:"clubs"`
//		SimilarClubs []Club `json:"similar_clubs"`
//		SimilarUsers []User `json:"similar_users"`
//	} `json:"recommendations"`
//}
//
//type Club struct {
//	ID          uint   `json:"id"`
//	Name        string `json:"name"`
//	Description string `json:"description"`
//	Category    string `json:"category"`
//	MemberCount int    `json:"member_count"`
//}
//
//// 테스트 유저 데이터
//var testUsers = []struct {
//	Email string
//	Name  string
//}{
//	{"user1@test.com", "김민수"},
//	{"user2@test.com", "이영희"},
//	{"user3@test.com", "박철수"},
//	{"user4@test.com", "정수진"},
//	{"user5@test.com", "최지훈"},
//	{"user6@test.com", "강서연"},
//	{"user7@test.com", "윤태영"},
//	{"user8@test.com", "임나영"},
//	{"user9@test.com", "한동욱"},
//	{"user10@test.com", "송미래"},
//}
//
//// 성향별 답변 패턴 (다양성을 위해)
//var answerPatterns = [][]int{
//	// 패턴 1: 사교적이고 활동적인 성향
//	{5, 5, 5, 3, 4, 5, 5, 5, 4, 5},
//	// 패턴 2: 내향적이고 몰입형
//	{1, 2, 1, 5, 2, 1, 2, 1, 1, 1},
//	// 패턴 3: 균형잡힌 성향
//	{3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
//	// 패턴 4: 유연하고 적응적인 성향
//	{3, 4, 3, 3, 5, 3, 3, 3, 3, 3},
//	// 패턴 5: 활동적이지만 소규모 선호
//	{3, 2, 2, 4, 4, 4, 3, 3, 4, 2},
//	// 패턴 6: 사교적이지만 깊이있는 관계 선호
//	{5, 3, 2, 4, 3, 4, 5, 5, 3, 3},
//	// 패턴 7: 몰입형이면서 유연한 성향
//	{2, 2, 3, 5, 5, 2, 2, 2, 2, 2},
//	// 패턴 8: 매우 외향적
//	{5, 5, 5, 2, 3, 5, 5, 5, 5, 5},
//	// 패턴 9: 내향적이지만 유연한 성향
//	{2, 2, 2, 3, 5, 2, 2, 2, 3, 2},
//	// 패턴 10: 도전적인 탐험가
//	{4, 3, 4, 4, 4, 5, 4, 4, 5, 4},
//}
//
//func main() {
//	rand.Seed(time.Now().UnixNano())
//
//	fmt.Println("=== 설문조사 테스트 시작 ===")
//	fmt.Println()
//
//	// 1. 설문 질문 조회
//	fmt.Println("1. 설문 질문 조회 중...")
//	questions, err := getQuestions()
//	if err != nil {
//		log.Fatalf("질문 조회 실패: %v", err)
//	}
//	fmt.Printf("✓ 총 %d개의 질문을 불러왔습니다.\n\n", len(questions))
//
//	// 2. 10명의 유저 생성 및 설문 진행
//	var createdUsers []User
//	for i, userData := range testUsers {
//		fmt.Printf("=== 유저 %d: %s ===\n", i+1, userData.Name)
//
//		// 유저 생성
//		user, err := createUser(userData.Email, userData.Name)
//		if err != nil {
//			log.Printf("유저 생성 실패: %v", err)
//			continue
//		}
//		createdUsers = append(createdUsers, user)
//		fmt.Printf("✓ 유저 생성 완료 (ID: %d)\n", user.ID)
//
//		// 답변 제출 (패턴 사용)
//		pattern := answerPatterns[i]
//		answers := generateAnswers(questions, pattern)
//		if err := submitAnswers(user.ID, answers); err != nil {
//			log.Printf("답변 제출 실패: %v", err)
//			continue
//		}
//		fmt.Printf("✓ 설문 답변 제출 완료\n")
//
//		// 잠시 대기
//		time.Sleep(500 * time.Millisecond)
//	}
//
//	fmt.Println()
//	fmt.Println("=== 모든 유저 설문 완료 ===")
//	fmt.Println()
//	time.Sleep(1 * time.Second)
//
//	// 3. 각 유저의 결과 조회 및 출력
//	fmt.Println("=== 설문 결과 분석 ===")
//	fmt.Println()
//
//	for i, user := range createdUsers {
//		fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
//		fmt.Printf("유저 %d: %s (%s)\n", i+1, user.Name, user.Email)
//		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
//
//		result, err := getResult(user.ID)
//		if err != nil {
//			log.Printf("결과 조회 실패: %v", err)
//			continue
//		}
//
//		printResult(result)
//		time.Sleep(300 * time.Millisecond)
//	}
//
//	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
//	fmt.Println("=== 테스트 완료 ===")
//	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
//}
//
//// 질문 조회
//func getQuestions() ([]Question, error) {
//	resp, err := http.Get(baseURL + "/questions")
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	var response Response
//	if err := json.Unmarshal(body, &response); err != nil {
//		return nil, err
//	}
//
//	var questions []Question
//	if err := json.Unmarshal(response.Data, &questions); err != nil {
//		return nil, err
//	}
//
//	return questions, nil
//}
//
//// 유저 생성
//func createUser(email, name string) (User, error) {
//	payload := map[string]string{
//		"email": email,
//		"name":  name,
//	}
//
//	jsonData, err := json.Marshal(payload)
//	if err != nil {
//		return User{}, err
//	}
//
//	resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(jsonData))
//	if err != nil {
//		return User{}, err
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return User{}, err
//	}
//
//	var response Response
//	if err := json.Unmarshal(body, &response); err != nil {
//		return User{}, err
//	}
//
//	var user User
//	if err := json.Unmarshal(response.Data, &user); err != nil {
//		return User{}, err
//	}
//
//	return user, nil
//}
//
//// 답변 생성 (패턴 기반)
//func generateAnswers(questions []Question, pattern []int) []AnswerPayload {
//	var answers []AnswerPayload
//
//	for i, question := range questions {
//		if i >= len(pattern) {
//			break
//		}
//
//		// 패턴에 따라 옵션 선택
//		optionIndex := pattern[i] - 1 // 1-based to 0-based
//		if optionIndex < 0 {
//			optionIndex = 0
//		}
//		if optionIndex >= len(question.Options) {
//			optionIndex = len(question.Options) - 1
//		}
//
//		answers = append(answers, AnswerPayload{
//			QuestionID: question.ID,
//			OptionID:   question.Options[optionIndex].ID,
//		})
//	}
//
//	return answers
//}
//
//// 답변 제출
//func submitAnswers(userID uint, answers []AnswerPayload) error {
//	payload := AnswerRequest{
//		UserID:  userID,
//		Answers: answers,
//	}
//
//	jsonData, err := json.Marshal(payload)
//	if err != nil {
//		return err
//	}
//
//	resp, err := http.Post(baseURL+"/answers/batch", "application/json", bytes.NewBuffer(jsonData))
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
//		body, _ := io.ReadAll(resp.Body)
//		return fmt.Errorf("답변 제출 실패: %s", string(body))
//	}
//
//	return nil
//}
//
//// 결과 조회
//func getResult(userID uint) (Result, error) {
//	url := fmt.Sprintf("%s/results/%d", baseURL, userID)
//	resp, err := http.Get(url)
//	if err != nil {
//		return Result{}, err
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return Result{}, err
//	}
//
//	var response Response
//	if err := json.Unmarshal(body, &response); err != nil {
//		return Result{}, err
//	}
//
//	var result Result
//	if err := json.Unmarshal(response.Data, &result); err != nil {
//		return Result{}, err
//	}
//
//	return result, nil
//}
//
//// 결과 출력
//func printResult(result Result) {
//	fmt.Printf("\n📊 성향 점수:\n")
//	fmt.Printf("  • 사교성: %.1f점\n", result.Scores.SocialityScore)
//	fmt.Printf("  • 활동성: %.1f점\n", result.Scores.ActivityScore)
//	fmt.Printf("  • 친밀도: %.1f점\n", result.Scores.IntimacyScore)
//	fmt.Printf("  • 몰입도: %.1f점\n", result.Scores.ImmersionScore)
//	fmt.Printf("  • 유연성: %.1f점\n", result.Scores.FlexibilityScore)
//
//	fmt.Printf("\n🎭 프로필 타입: %s\n", result.ProfileType)
//
//	fmt.Printf("\n📝 성향 설명:\n")
//	for i, desc := range result.Descriptions {
//		fmt.Printf("  %d. %s\n", i+1, desc)
//	}
//
//	fmt.Printf("\n🎯 추천 클럽 (%d개):\n", len(result.Recommendations.Clubs))
//	for i, club := range result.Recommendations.Clubs {
//		if i >= 3 {
//			break
//		}
//		fmt.Printf("  • %s - %s\n", club.Name, club.Description)
//	}
//
//	fmt.Printf("\n👥 유사한 성향의 사용자 (%d명):\n", len(result.Recommendations.SimilarUsers))
//	for i, user := range result.Recommendations.SimilarUsers {
//		if i >= 3 {
//			break
//		}
//		fmt.Printf("  • %s\n", user.Name)
//	}
//
//	fmt.Printf("\n🌟 유사 사용자가 많은 클럽 (%d개):\n", len(result.Recommendations.SimilarClubs))
//	for i, club := range result.Recommendations.SimilarClubs {
//		if i >= 3 {
//			break
//		}
//		fmt.Printf("  • %s - %s\n", club.Name, club.Description)
//	}
//}

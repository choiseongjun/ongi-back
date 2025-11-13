package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const baseURL = "http://localhost:5000/api/v1"

// 응답 구조체
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message,omitempty"`
}

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Question struct {
	ID           uint     `json:"id"`
	QuestionText string   `json:"question_text"`
	Order        int      `json:"order"`
	Category     string   `json:"category"`
	Options      []Option `json:"options"`
}

type Option struct {
	ID         uint   `json:"id"`
	QuestionID uint   `json:"question_id"`
	OptionText string `json:"option_text"`
	Score      int    `json:"score"`
	Weight     string `json:"weight"`
}

type AnswerRequest struct {
	UserID  uint            `json:"user_id"`
	Answers []AnswerPayload `json:"answers"`
}

type AnswerPayload struct {
	QuestionID uint `json:"question_id"`
	OptionID   uint `json:"option_id"`
}

type Result struct {
	Scores struct {
		SocialityScore   float64 `json:"sociality_score"`
		ActivityScore    float64 `json:"activity_score"`
		IntimacyScore    float64 `json:"intimacy_score"`
		ImmersionScore   float64 `json:"immersion_score"`
		FlexibilityScore float64 `json:"flexibility_score"`
	} `json:"scores"`
	ProfileType     string   `json:"profile_type"`
	Descriptions    []string `json:"descriptions"`
	Recommendations struct {
		Clubs        []Club `json:"clubs"`
		SimilarClubs []Club `json:"similar_clubs"`
		SimilarUsers []User `json:"similar_users"`
	} `json:"recommendations"`
}

type Club struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	MemberCount int    `json:"member_count"`
}

// 한국 성씨와 이름
var surnames = []string{"김", "이", "박", "최", "정", "강", "조", "윤", "장", "임", "한", "오", "서", "신", "권", "황", "안", "송", "류", "홍"}
var firstNames = []string{"민준", "서연", "예준", "하은", "도윤", "서윤", "시우", "지우", "수호", "지민", "현우", "수빈", "준서", "은서", "하준", "윤서", "건우", "채원", "우진", "다인"}

// 성향별 답변 패턴
var answerPatterns = [][]int{
	{5, 5, 5, 3, 4, 5, 5, 5, 4, 5}, // 사교적이고 활동적
	{1, 2, 1, 5, 2, 1, 2, 1, 1, 1}, // 내향적이고 몰입형
	{3, 3, 3, 3, 3, 3, 3, 3, 3, 3}, // 균형잡힌 성향
	{3, 4, 3, 3, 5, 3, 3, 3, 3, 3}, // 유연하고 적응적
	{3, 2, 2, 4, 4, 4, 3, 3, 4, 2}, // 활동적이지만 소규모 선호
	{5, 3, 2, 4, 3, 4, 5, 5, 3, 3}, // 사교적이지만 깊이있는 관계
	{2, 2, 3, 5, 5, 2, 2, 2, 2, 2}, // 몰입형이면서 유연
	{5, 5, 5, 2, 3, 5, 5, 5, 5, 5}, // 매우 외향적
	{2, 2, 2, 3, 5, 2, 2, 2, 3, 2}, // 내향적이지만 유연
	{4, 3, 4, 4, 4, 5, 4, 4, 5, 4}, // 도전적인 탐험가
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== 설문조사 테스트 시작 (500명) ===")
	fmt.Println()

	// 1. 설문 질문 조회
	fmt.Println("1. 설문 질문 조회 중...")
	questions, err := getQuestions()
	if err != nil {
		log.Fatalf("질문 조회 실패: %v", err)
	}
	fmt.Printf("✓ 총 %d개의 질문을 불러왔습니다.\n\n", len(questions))

	// 2. 500명의 유저 생성 및 설문 진행
	fmt.Println("2. 500명의 유저 생성 및 설문 진행 중...")
	var createdUsers []User
	startTime := time.Now()

	for i := 0; i < 500; i++ {
		// 랜덤 이름 생성
		surname := surnames[rand.Intn(len(surnames))]
		firstName := firstNames[rand.Intn(len(firstNames))]
		name := surname + firstName
		email := fmt.Sprintf("user%d@test.com", i+1)

		// 유저 생성
		user, err := createUser(email, name)
		if err != nil {
			log.Printf("유저 %d 생성 실패: %v", i+1, err)
			continue
		}
		createdUsers = append(createdUsers, user)

		// 랜덤 패턴 선택
		pattern := answerPatterns[rand.Intn(len(answerPatterns))]
		answers := generateAnswers(questions, pattern)

		// 답변 제출
		if err := submitAnswers(user.ID, answers); err != nil {
			log.Printf("유저 %d 답변 제출 실패: %v", i+1, err)
			continue
		}

		// 진행 상황 표시
		if (i+1)%50 == 0 {
			fmt.Printf("  ✓ %d명 완료... (경과 시간: %.1f초)\n", i+1, time.Since(startTime).Seconds())
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n✓ 총 %d명의 유저 생성 및 설문 완료 (총 소요 시간: %.1f초)\n\n", len(createdUsers), elapsed.Seconds())

	// 3. 샘플로 처음 10명의 결과 출력
	fmt.Println("=== 샘플 결과 출력 (처음 10명) ===")
	fmt.Println()

	sampleCount := 10
	if len(createdUsers) < sampleCount {
		sampleCount = len(createdUsers)
	}

	for i := 0; i < sampleCount; i++ {
		user := createdUsers[i]

		fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("유저 %d: %s (ID: %d)\n", i+1, user.Name, user.ID)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		result, err := getResult(user.ID)
		if err != nil {
			log.Printf("결과 조회 실패: %v", err)
			continue
		}

		printResult(result)
	}

	// 4. 매칭 통계 출력
	fmt.Println("\n\n=== 매칭 통계 분석 중... ===")
	fmt.Println()

	// 프로필 타입별 분포
	profileCounts := make(map[string]int)
	totalMatches := 0

	for _, user := range createdUsers {
		result, err := getResult(user.ID)
		if err != nil {
			continue
		}
		profileCounts[result.ProfileType]++
		totalMatches += len(result.Recommendations.SimilarUsers)
	}

	avgMatches := float64(totalMatches) / float64(len(createdUsers))

	fmt.Println("📊 프로필 타입별 분포:")
	for profileType, count := range profileCounts {
		percentage := float64(count) / float64(len(createdUsers)) * 100
		fmt.Printf("  • %s: %d명 (%.1f%%)\n", profileType, count, percentage)
	}

	fmt.Printf("\n📈 평균 유사 사용자 매칭 수: %.1f명\n", avgMatches)

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("=== 테스트 완료 (총 %d명 생성) ===\n", len(createdUsers))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func getQuestions() ([]Question, error) {
	resp, err := http.Get(baseURL + "/questions")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var questions []Question
	if err := json.Unmarshal(response.Data, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

func createUser(email, name string) (User, error) {
	payload := map[string]string{
		"email": email,
		"name":  name,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return User{}, err
	}

	resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return User{}, err
	}

	var user User
	if err := json.Unmarshal(response.Data, &user); err != nil {
		return User{}, err
	}

	return user, nil
}

func generateAnswers(questions []Question, pattern []int) []AnswerPayload {
	var answers []AnswerPayload

	for i, question := range questions {
		if i >= len(pattern) {
			break
		}

		optionIndex := pattern[i] - 1
		if optionIndex < 0 {
			optionIndex = 0
		}
		if optionIndex >= len(question.Options) {
			optionIndex = len(question.Options) - 1
		}

		answers = append(answers, AnswerPayload{
			QuestionID: question.ID,
			OptionID:   question.Options[optionIndex].ID,
		})
	}

	return answers
}

func submitAnswers(userID uint, answers []AnswerPayload) error {
	payload := AnswerRequest{
		UserID:  userID,
		Answers: answers,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(baseURL+"/answers/batch", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("답변 제출 실패: %s", string(body))
	}

	return nil
}

func getResult(userID uint) (Result, error) {
	url := fmt.Sprintf("%s/results/%d", baseURL, userID)
	resp, err := http.Get(url)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Result{}, err
	}

	var result Result
	if err := json.Unmarshal(response.Data, &result); err != nil {
		return Result{}, err
	}

	return result, nil
}

func printResult(result Result) {
	fmt.Printf("\n📊 성향 점수:\n")
	fmt.Printf("  • 사교성: %.1f점\n", result.Scores.SocialityScore)
	fmt.Printf("  • 활동성: %.1f점\n", result.Scores.ActivityScore)
	fmt.Printf("  • 친밀도: %.1f점\n", result.Scores.IntimacyScore)
	fmt.Printf("  • 몰입도: %.1f점\n", result.Scores.ImmersionScore)
	fmt.Printf("  • 유연성: %.1f점\n", result.Scores.FlexibilityScore)

	fmt.Printf("\n🎭 프로필 타입: %s\n", result.ProfileType)

	fmt.Printf("\n📝 성향 설명:\n")
	for i, desc := range result.Descriptions {
		fmt.Printf("  %d. %s\n", i+1, desc)
	}

	fmt.Printf("\n🎯 추천 클럽 (상위 3개):\n")
	for i, club := range result.Recommendations.Clubs {
		if i >= 3 {
			break
		}
		fmt.Printf("  • %s - %s\n", club.Name, club.Description)
	}

	fmt.Printf("\n👥 유사한 성향의 사용자 (상위 5명):\n")
	for i, user := range result.Recommendations.SimilarUsers {
		if i >= 5 {
			break
		}
		fmt.Printf("  • %s\n", user.Name)
	}
}

package services

import (
	"fmt"
	"math"
	"ongi-back/database"
	"ongi-back/models"
)

type ScoreResult struct {
	SocialityScore   float64 `json:"sociality_score"`
	ActivityScore    float64 `json:"activity_score"`
	IntimacyScore    float64 `json:"intimacy_score"`
	ImmersionScore   float64 `json:"immersion_score"`
	FlexibilityScore float64 `json:"flexibility_score"`
}

type AnalysisResult struct {
	Scores        ScoreResult `json:"scores"`
	ProfileType   string      `json:"profile_type"`
	Descriptions  []string    `json:"descriptions"`
	SuggestedClubs    []models.Club   `json:"suggested_clubs"`
	SuggestedMeetings []models.Meeting `json:"suggested_meetings"`
	SimilarUsers      []models.User    `json:"similar_users"`
}

func CalculateScores(userID uint) (*ScoreResult, error) {
	var answers []models.UserAnswer

	err := database.DB.Preload("Option").
		Where("user_id = ?", userID).
		Find(&answers).Error

	if err != nil {
		return nil, err
	}

	if len(answers) == 0 {
		return nil, fmt.Errorf("no answers found for user")
	}

	scores := &ScoreResult{}
	categoryScores := make(map[string][]int)

	// 각 답변의 점수를 카테고리별로 분류
	for _, answer := range answers {
		weight := answer.Option.Weight
		score := answer.Option.Score

		if weight != "" {
			categoryScores[weight] = append(categoryScores[weight], score)
		}
	}

	// 카테고리별 평균 계산 (0-100 스케일로 변환)
	scores.SocialityScore = calculateAverage(categoryScores["sociality"])
	scores.ActivityScore = calculateAverage(categoryScores["activity"])
	scores.IntimacyScore = calculateAverage(categoryScores["intimacy"])
	scores.ImmersionScore = calculateAverage(categoryScores["immersion"])
	scores.FlexibilityScore = calculateAverage(categoryScores["flexibility"])

	return scores, nil
}

func calculateAverage(scores []int) float64 {
	if len(scores) == 0 {
		return 0
	}

	sum := 0
	for _, score := range scores {
		sum += score
	}

	avg := float64(sum) / float64(len(scores))
	// 1-5 점수를 0-100으로 변환
	return math.Round((avg / 5.0 * 100) * 10) / 10
}

func GenerateDescriptions(scores *ScoreResult) []string {
	descriptions := []string{}

	// 유연성 기반 설명
	if scores.FlexibilityScore >= 60 {
		descriptions = append(descriptions, "당신은 상황에 따라 유연하게 대처하며, 내향과 외향의 균형을 잘 맞춥니다.")
	} else if scores.FlexibilityScore >= 40 {
		descriptions = append(descriptions, "당신은 때로는 계획적이고 때로는 즉흥적인 성향을 보입니다.")
	} else {
		descriptions = append(descriptions, "당신은 계획적이고 체계적인 접근을 선호합니다.")
	}

	// 사교성과 활동성 기반 설명
	if scores.SocialityScore >= 60 && scores.ActivityScore >= 60 {
		descriptions = append(descriptions, "다양한 활동을 즐기고, 사람들과의 조화를 중요하게 생각합니다.")
	} else if scores.SocialityScore >= 60 {
		descriptions = append(descriptions, "사람들과 함께하는 시간을 즐기며, 깊은 대화를 선호합니다.")
	} else if scores.ActivityScore >= 60 {
		descriptions = append(descriptions, "적극적으로 새로운 활동에 참여하며, 도전을 즐깁니다.")
	} else {
		descriptions = append(descriptions, "조용하고 안정적인 환경에서 집중하는 것을 선호합니다.")
	}

	// 친밀도와 몰입도 기반 설명
	if scores.IntimacyScore >= 60 && scores.ImmersionScore >= 60 {
		descriptions = append(descriptions, "깊이 있는 관계를 형성하고, 한 가지 일에 오랫동안 몰두하는 성향이 있습니다.")
	} else if scores.IntimacyScore >= 60 {
		descriptions = append(descriptions, "소수의 사람들과 깊은 유대감을 형성하는 것을 중요하게 여깁니다.")
	} else if scores.ImmersionScore >= 60 {
		descriptions = append(descriptions, "관심사에 깊이 몰입하며, 전문성을 추구합니다.")
	}

	return descriptions
}

func DetermineProfileType(scores *ScoreResult) string {
	// 가장 높은 점수를 가진 카테고리 조합으로 프로필 타입 결정
	if scores.SocialityScore >= 70 && scores.ActivityScore >= 70 {
		return "열정적인 사교가"
	} else if scores.SocialityScore >= 70 && scores.IntimacyScore >= 70 {
		return "따뜻한 조력자"
	} else if scores.ActivityScore >= 70 && scores.ImmersionScore >= 70 {
		return "도전적인 탐험가"
	} else if scores.ImmersionScore >= 70 && scores.IntimacyScore >= 70 {
		return "깊이있는 전문가"
	} else if scores.FlexibilityScore >= 70 {
		return "유연한 적응형"
	} else if scores.SocialityScore >= 60 {
		return "친근한 외향형"
	} else if scores.ImmersionScore >= 60 {
		return "집중하는 몰입형"
	} else {
		return "균형잡힌 조화형"
	}
}

func GetCompleteAnalysis(userID uint) (*AnalysisResult, error) {
	scores, err := CalculateScores(userID)
	if err != nil {
		return nil, err
	}

	profileType := DetermineProfileType(scores)
	descriptions := GenerateDescriptions(scores)

	result := &AnalysisResult{
		Scores:       *scores,
		ProfileType:  profileType,
		Descriptions: descriptions,
	}

	return result, nil
}

package services

import (
	"math"
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/utils"
	"runtime"
	"sort"
)

// SimilarProfile - 유사한 프로필 정보
type SimilarProfile struct {
	SessionID  string        `json:"session_id,omitempty"`
	UserID     *uint         `json:"user_id,omitempty"`
	User       *models.User  `json:"user,omitempty"`
	Similarity float64       `json:"similarity"`
	Vector     *utils.Vector5D `json:"vector,omitempty"`
}

// GetSimilarProfilesFast - 고속 유사 프로필 검색 (벡터 연산 최적화)
func GetSimilarProfilesFast(sessionID string, limit int) ([]SimilarProfile, error) {
	// 1. 현재 세션의 벡터 가져오기
	var currentVector models.SessionVector
	err := database.DB.Where("session_id = ?", sessionID).First(&currentVector).Error
	if err != nil {
		return nil, err
	}

	currentV := utils.FromSlice(currentVector.Vector)
	if currentV == nil {
		return nil, err
	}

	// 2. 모든 벡터 가져오기 (자기 자신 제외)
	var allVectors []models.SessionVector
	err = database.DB.
		Preload("User").
		Where("session_id != ?", sessionID).
		Find(&allVectors).Error

	if err != nil {
		return nil, err
	}

	if len(allVectors) == 0 {
		return []SimilarProfile{}, nil
	}

	// 3. 벡터 변환
	vectors := make([]*utils.Vector5D, len(allVectors))
	for i, v := range allVectors {
		vectors[i] = utils.FromSlice(v.Vector)
	}

	// 4. 병렬 유사도 계산 (CPU 코어 수만큼 워커 사용)
	workers := runtime.NumCPU()
	results := utils.BatchSimilarity(currentV, vectors, workers)

	// 5. 결과를 SimilarProfile로 변환
	profiles := make([]SimilarProfile, len(results))
	for i, result := range results {
		profile := SimilarProfile{
			SessionID:  allVectors[result.Index].SessionID,
			UserID:     allVectors[result.Index].UserID,
			Similarity: result.Similarity,
			Vector:     vectors[result.Index],
		}

		// 회원인 경우 User 정보 포함
		if profile.UserID != nil {
			var user models.User
			database.DB.First(&user, *profile.UserID)
			profile.User = &user
		}

		profiles[i] = profile
	}

	// 6. 유사도 높은 순으로 정렬
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Similarity > profiles[j].Similarity
	})

	// 7. 상위 N개 반환
	if len(profiles) > limit {
		profiles = profiles[:limit]
	}

	return profiles, nil
}

// GetRecommendedClubsForSession - 세션 기반 클럽 추천
func GetRecommendedClubsForSession(sessionID string, limit int) ([]models.Club, error) {
	var sessionVector models.SessionVector
	err := database.DB.Where("session_id = ?", sessionID).First(&sessionVector).Error
	if err != nil {
		return nil, err
	}

	v := utils.FromSlice(sessionVector.Vector)
	if v == nil {
		return []models.Club{}, nil
	}

	var clubs []models.Club
	query := database.DB.Preload("Members")

	// 사교성이 높은 사람에게는 멤버가 많은 클럽 추천
	if v.Sociality >= 70 {
		query = query.Order("member_count DESC")
	} else if v.Intimacy >= 60 {
		// 친밀도가 높은 사람에게는 적당한 규모의 클럽 추천
		query = query.Where("member_count <= ?", 50).Order("member_count ASC")
	} else {
		// 균형잡힌 사람에게는 중간 규모
		query = query.Where("member_count BETWEEN ? AND ?", 10, 100)
	}

	err = query.Limit(limit).Find(&clubs).Error
	if err != nil {
		return nil, err
	}

	return clubs, nil
}

// GetClubsWithSimilarMembersForSession - 유사한 사람들이 많은 클럽 추천
func GetClubsWithSimilarMembersForSession(sessionID string, limit int) ([]models.Club, error) {
	// 1. 유사한 프로필 찾기
	similarProfiles, err := GetSimilarProfilesFast(sessionID, 20)
	if err != nil {
		return nil, err
	}

	if len(similarProfiles) == 0 {
		return GetRecommendedClubsForSession(sessionID, limit)
	}

	// 2. 유사한 사용자들의 ID 수집 (회원만)
	var userIDs []uint
	for _, profile := range similarProfiles {
		if profile.UserID != nil {
			userIDs = append(userIDs, *profile.UserID)
		}
	}

	if len(userIDs) == 0 {
		return GetRecommendedClubsForSession(sessionID, limit)
	}

	// 3. 클럽별 유사 사용자 수 계산
	type ClubCount struct {
		ClubID uint
		Count  int64
	}

	var clubCounts []ClubCount
	err = database.DB.Model(&models.ClubMember{}).
		Select("club_id, COUNT(*) as count").
		Where("user_id IN ?", userIDs).
		Group("club_id").
		Order("count DESC").
		Limit(limit).
		Scan(&clubCounts).Error

	if err != nil {
		return nil, err
	}

	// 4. 클럽 정보 가져오기
	var clubIDs []uint
	for _, cc := range clubCounts {
		clubIDs = append(clubIDs, cc.ClubID)
	}

	var clubs []models.Club
	if len(clubIDs) > 0 {
		err = database.DB.Where("id IN ?", clubIDs).Find(&clubs).Error
		if err != nil {
			return nil, err
		}

		// 원래 순서대로 정렬 (count 높은 순)
		clubMap := make(map[uint]models.Club)
		for _, club := range clubs {
			clubMap[club.ID] = club
		}

		sortedClubs := make([]models.Club, 0, len(clubIDs))
		for _, id := range clubIDs {
			if club, ok := clubMap[id]; ok {
				sortedClubs = append(sortedClubs, club)
			}
		}
		clubs = sortedClubs
	}

	return clubs, nil
}

// GetRecommendedMeetingsForSession - 세션 기반 모임 추천
func GetRecommendedMeetingsForSession(sessionID string, limit int) ([]models.Meeting, error) {
	var sessionVector models.SessionVector
	err := database.DB.Where("session_id = ?", sessionID).First(&sessionVector).Error
	if err != nil {
		return nil, err
	}

	v := utils.FromSlice(sessionVector.Vector)
	if v == nil {
		return []models.Meeting{}, nil
	}

	var meetings []models.Meeting
	query := database.DB.Preload("Club")

	// 활동성이 높은 사람에게는 다양한 모임 추천
	if v.Activity >= 70 {
		query = query.Order("scheduled_at ASC")
	} else if v.Intimacy >= 60 {
		// 친밀도가 높은 사람에게는 소규모 모임
		query = query.Where("max_members <= ?", 20).Order("max_members ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	err = query.Limit(limit).Find(&meetings).Error
	if err != nil {
		return nil, err
	}

	return meetings, nil
}

// CalculateProfileCompatibility - 두 프로필 간 궁합 점수 계산
func CalculateProfileCompatibility(v1, v2 *utils.Vector5D) map[string]interface{} {
	similarity := utils.SimilarityScore(v1, v2)

	// 차원별 궁합 분석
	compatibility := map[string]interface{}{
		"overall_score": similarity,
		"details": map[string]interface{}{
			"sociality_match":   100 - math.Abs(v1.Sociality-v2.Sociality),
			"activity_match":    100 - math.Abs(v1.Activity-v2.Activity),
			"intimacy_match":    100 - math.Abs(v1.Intimacy-v2.Intimacy),
			"immersion_match":   100 - math.Abs(v1.Immersion-v2.Immersion),
			"flexibility_match": 100 - math.Abs(v1.Flexibility-v2.Flexibility),
		},
	}

	// 궁합 평가
	if similarity >= 80 {
		compatibility["rating"] = "최고의 궁합"
		compatibility["description"] = "매우 비슷한 성향으로 서로 잘 맞을 것입니다"
	} else if similarity >= 70 {
		compatibility["rating"] = "좋은 궁합"
		compatibility["description"] = "비슷한 성향으로 편안한 관계를 형성할 수 있습니다"
	} else if similarity >= 60 {
		compatibility["rating"] = "보통 궁합"
		compatibility["description"] = "서로 다른 점이 있지만 조화롭게 지낼 수 있습니다"
	} else if similarity >= 50 {
		compatibility["rating"] = "상호보완적"
		compatibility["description"] = "다른 성향으로 서로에게 새로운 자극이 될 수 있습니다"
	} else {
		compatibility["rating"] = "흥미로운 조합"
		compatibility["description"] = "매우 다른 성향이지만 배울 점이 많을 것입니다"
	}

	return compatibility
}

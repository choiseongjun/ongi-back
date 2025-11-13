package services

import (
	"math"
	"ongi-back/database"
	"ongi-back/models"
	"sort"
)

type UserSimilarity struct {
	User       models.User
	Similarity float64
}

// 유클리드 거리 기반 유사도 계산
func calculateSimilarity(profile1, profile2 *models.UserProfile) float64 {
	diff1 := profile1.SocialityScore - profile2.SocialityScore
	diff2 := profile1.ActivityScore - profile2.ActivityScore
	diff3 := profile1.IntimacyScore - profile2.IntimacyScore
	diff4 := profile1.ImmersionScore - profile2.ImmersionScore
	diff5 := profile1.FlexibilityScore - profile2.FlexibilityScore

	distance := math.Sqrt(
		diff1*diff1 + diff2*diff2 + diff3*diff3 + diff4*diff4 + diff5*diff5,
	)

	// 거리를 유사도로 변환 (0-100)
	maxDistance := math.Sqrt(5 * 100 * 100) // 최대 거리
	similarity := (1 - (distance / maxDistance)) * 100

	return similarity
}

func GetSimilarUsers(userID uint, limit int) ([]models.User, error) {
	var userProfile models.UserProfile
	err := database.DB.Where("user_id = ?", userID).First(&userProfile).Error
	if err != nil {
		return nil, err
	}

	var allProfiles []models.UserProfile
	err = database.DB.Preload("User").
		Where("user_id != ?", userID).
		Find(&allProfiles).Error
	if err != nil {
		return nil, err
	}

	// 유사도 계산
	similarities := []UserSimilarity{}
	for _, profile := range allProfiles {
		similarity := calculateSimilarity(&userProfile, &profile)
		similarities = append(similarities, UserSimilarity{
			User:       profile.User,
			Similarity: similarity,
		})
	}

	// 유사도 높은 순으로 정렬
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	// 상위 N명 반환
	var users []models.User
	maxResults := limit
	if len(similarities) < maxResults {
		maxResults = len(similarities)
	}

	for i := 0; i < maxResults; i++ {
		users = append(users, similarities[i].User)
	}

	return users, nil
}

func GetRecommendedClubs(userID uint, limit int) ([]models.Club, error) {
	var userProfile models.UserProfile
	err := database.DB.Where("user_id = ?", userID).First(&userProfile).Error
	if err != nil {
		return nil, err
	}

	var clubs []models.Club
	query := database.DB.Preload("Members")

	// 사교성이 높은 사람에게는 멤버가 많은 클럽 추천
	if userProfile.SocialityScore >= 70 {
		query = query.Order("member_count DESC")
	} else {
		// 친밀도가 높은 사람에게는 적당한 규모의 클럽 추천
		query = query.Order("member_count ASC")
	}

	err = query.Limit(limit).Find(&clubs).Error
	if err != nil {
		return nil, err
	}

	return clubs, nil
}

func GetClubsWithSimilarMembers(userID uint, limit int) ([]models.Club, error) {
	// 유사한 사용자들이 많이 가입한 클럽 찾기
	similarUsers, err := GetSimilarUsers(userID, 20)
	if err != nil {
		return nil, err
	}

	if len(similarUsers) == 0 {
		return GetRecommendedClubs(userID, limit)
	}

	var userIDs []uint
	for _, user := range similarUsers {
		userIDs = append(userIDs, user.ID)
	}

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
	}

	return clubs, nil
}

func GetRecommendedMeetings(userID uint, limit int) ([]models.Meeting, error) {
	var userProfile models.UserProfile
	err := database.DB.Where("user_id = ?", userID).First(&userProfile).Error
	if err != nil {
		return nil, err
	}

	var meetings []models.Meeting
	query := database.DB.Preload("Club")

	// 활동성이 높은 사람에게는 다양한 모임 추천
	if userProfile.ActivityScore >= 70 {
		query = query.Order("scheduled_at ASC")
	} else {
		query = query.Order("max_members ASC")
	}

	err = query.Limit(limit).Find(&meetings).Error
	if err != nil {
		return nil, err
	}

	return meetings, nil
}

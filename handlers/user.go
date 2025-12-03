package handlers

import (
	"fmt"
	"math/rand"
	"ongi-back/database"
	"ongi-back/models"
	"ongi-back/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// 사용자 생성
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user := models.User{
		Email: req.Email,
		Name:  req.Name,
	}

	err := database.DB.Create(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// 사용자 조회
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// 모든 사용자 조회
func GetUsers(c *fiber.Ctx) error {
	var users []models.User

	err := database.DB.Find(&users).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

// 성향 분석 함수
func analyzeTendency(score float64) string {
	if score >= 80 {
		return "매우 높음"
	} else if score >= 60 {
		return "높음"
	} else if score >= 40 {
		return "보통"
	} else if score >= 20 {
		return "낮음"
	}
	return "매우 낮음"
}

// CreateOrUpdateUserProfile 사용자 프로필 생성/수정
// POST /users/profile
type CreateUserProfileRequest struct {
	UserID           uint    `json:"user_id" validate:"required"`
	SocialityScore   float64 `json:"sociality_score"`
	ActivityScore    float64 `json:"activity_score"`
	IntimacyScore    float64 `json:"intimacy_score"`
	ImmersionScore   float64 `json:"immersion_score"`
	FlexibilityScore float64 `json:"flexibility_score"`
	ResultSummary    string  `json:"result_summary"`
	ProfileType      string  `json:"profile_type"`
}

func CreateOrUpdateUserProfile(c *fiber.Ctx) error {
	var req CreateUserProfileRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// 사용자 존재 확인
	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	// 기존 프로필 확인
	var profile models.UserProfile
	result := database.DB.Where("user_id = ?", req.UserID).First(&profile)

	if result.Error != nil {
		// 신규 프로필 생성
		profile = models.UserProfile{
			UserID:           req.UserID,
			SocialityScore:   req.SocialityScore,
			ActivityScore:    req.ActivityScore,
			IntimacyScore:    req.IntimacyScore,
			ImmersionScore:   req.ImmersionScore,
			FlexibilityScore: req.FlexibilityScore,
			ResultSummary:    req.ResultSummary,
			ProfileType:      req.ProfileType,
		}

		if err := database.DB.Create(&profile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to create profile",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"message": "Profile created successfully",
			"data":    profile,
		})
	}

	// 기존 프로필 업데이트
	profile.SocialityScore = req.SocialityScore
	profile.ActivityScore = req.ActivityScore
	profile.IntimacyScore = req.IntimacyScore
	profile.ImmersionScore = req.ImmersionScore
	profile.FlexibilityScore = req.FlexibilityScore
	profile.ResultSummary = req.ResultSummary
	profile.ProfileType = req.ProfileType

	if err := database.DB.Save(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
		"data":    profile,
	})
}

// 사용자 프로필 조회
func GetUserProfile(c *fiber.Ctx) error {
	userID := c.Params("id")

	var profile models.UserProfile
	err := database.DB.Preload("User").Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Profile not found",
		})
	}

	// 사용자 성향 분석
	tendencies := fiber.Map{
		"sociality":   fiber.Map{"score": profile.SocialityScore, "level": analyzeTendency(profile.SocialityScore)},
		"activity":    fiber.Map{"score": profile.ActivityScore, "level": analyzeTendency(profile.ActivityScore)},
		"intimacy":    fiber.Map{"score": profile.IntimacyScore, "level": analyzeTendency(profile.IntimacyScore)},
		"immersion":   fiber.Map{"score": profile.ImmersionScore, "level": analyzeTendency(profile.ImmersionScore)},
		"flexibility": fiber.Map{"score": profile.FlexibilityScore, "level": analyzeTendency(profile.FlexibilityScore)},
	}

	// 유사 사용자 추천 (70% 이상 유사도)
	var uid uint
	if _, err := fmt.Sscanf(userID, "%d", &uid); err == nil {
		similarUsers, _ := services.GetSimilarUsers(uid, 20) // 상위 20명

		// 클럽 추천 (유사한 멤버들이 있는 클럽 우선)
		recommendedClubs, _ := services.GetClubsWithSimilarMembers(uid, 10)

		// 만약 유사 멤버 기반 클럽이 부족하면 성향 기반 클럽 추가
		if len(recommendedClubs) < 5 {
			additionalClubs, _ := services.GetRecommendedClubs(uid, 10)
			recommendedClubs = append(recommendedClubs, additionalClubs...)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"profile":           profile,
				"tendencies":        tendencies,
				"similar_users":     similarUsers,
				"recommended_clubs": recommendedClubs,
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"profile":    profile,
			"tendencies": tendencies,
		},
	})
}

// 자동 매칭 - 추천 모임 중 랜덤으로 자동 가입
func AutoMatchClubs(c *fiber.Ctx) error {
	userID := c.Params("id")

	var uid uint
	if _, err := fmt.Sscanf(userID, "%d", &uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// 사용자가 설문을 완료했는지 확인
	var profile models.UserProfile
	err := database.DB.Where("user_id = ?", uid).First(&profile).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User profile not found. Please complete the survey first.",
		})
	}

	// 추천 클럽 가져오기 (유사 멤버 기반 + 성향 기반)
	recommendedClubs, err := services.GetClubsWithSimilarMembers(uid, 20)
	if err != nil || len(recommendedClubs) < 5 {
		// 추가 클럽 가져오기
		additionalClubs, _ := services.GetRecommendedClubs(uid, 20)
		recommendedClubs = append(recommendedClubs, additionalClubs...)
	}

	if len(recommendedClubs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No recommended clubs available",
		})
	}

	// 랜덤 인원수 결정 (1~5명, 추천 클럽 수를 초과하지 않도록)
	rand.Seed(time.Now().UnixNano())
	maxClubs := len(recommendedClubs)
	if maxClubs > 5 {
		maxClubs = 5
	}
	numClubsToJoin := rand.Intn(maxClubs) + 1 // 1부터 maxClubs까지

	// 추천 클럽 셔플
	rand.Shuffle(len(recommendedClubs), func(i, j int) {
		recommendedClubs[i], recommendedClubs[j] = recommendedClubs[j], recommendedClubs[i]
	})

	// 랜덤으로 선택된 클럽들에 자동 가입
	var joinedClubs []models.Club
	var alreadyMember []models.Club

	for i := 0; i < numClubsToJoin && i < len(recommendedClubs); i++ {
		club := recommendedClubs[i]

		// 이미 가입된 클럽인지 확인
		var existingMember models.ClubMember
		checkErr := database.DB.Where("club_id = ? AND user_id = ?", club.ID, uid).First(&existingMember).Error

		if checkErr == nil {
			// 이미 가입된 경우
			alreadyMember = append(alreadyMember, club)
			continue
		}

		// 클럽 가입
		newMember := models.ClubMember{
			ClubID:   club.ID,
			UserID:   uid,
			JoinedAt: time.Now(),
		}

		if err := database.DB.Create(&newMember).Error; err != nil {
			continue // 에러 발생시 다음 클럽으로
		}

		// 멤버 수 증가
		database.DB.Model(&models.Club{}).Where("id = ?", club.ID).Update("member_count", club.MemberCount+1)

		joinedClubs = append(joinedClubs, club)
	}

	// 가입된 클럽이 없는 경우 (모두 이미 가입되어 있었던 경우)
	if len(joinedClubs) == 0 && len(alreadyMember) > 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Already member of selected clubs",
			"data": fiber.Map{
				"joined_clubs":         []models.Club{},
				"already_member_clubs": alreadyMember,
				"attempted_count":      numClubsToJoin,
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Successfully joined %d club(s)", len(joinedClubs)),
		"data": fiber.Map{
			"joined_clubs":         joinedClubs,
			"already_member_clubs": alreadyMember,
			"attempted_count":      numClubsToJoin,
			"total_recommended":    len(recommendedClubs),
		},
	})
}

// 그룹 자동 매칭 - 유사한 성향의 사용자들과 함께 클럽에 가입
func AutoMatchWithSimilarUsers(c *fiber.Ctx) error {
	userID := c.Params("id")

	var uid uint
	if _, err := fmt.Sscanf(userID, "%d", &uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// 사용자가 설문을 완료했는지 확인
	var profile models.UserProfile
	err := database.DB.Where("user_id = ?", uid).First(&profile).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User profile not found. Please complete the survey first.",
		})
	}

	// 유사한 사용자 찾기 (유사도 70% 이상, 최대 20명)
	similarUsers, err := services.GetSimilarUsers(uid, 20)
	if err != nil || len(similarUsers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No similar users found for group matching",
		})
	}

	// 랜덤으로 2-4명의 유사 사용자 선택
	rand.Seed(time.Now().UnixNano())
	groupSize := rand.Intn(3) + 2 // 2~4명
	if groupSize > len(similarUsers) {
		groupSize = len(similarUsers)
	}

	// 유사 사용자 셔플 후 선택
	rand.Shuffle(len(similarUsers), func(i, j int) {
		similarUsers[i], similarUsers[j] = similarUsers[j], similarUsers[i]
	})

	selectedUsers := make([]uint, groupSize+1)
	selectedUsers[0] = uid // 본인 포함
	for i := 0; i < groupSize; i++ {
		selectedUsers[i+1] = similarUsers[i].User.ID
	}

	// 그룹이 함께 들어갈 추천 클럽 찾기
	recommendedClubs, err := services.GetClubsWithSimilarMembers(uid, 20)
	if err != nil || len(recommendedClubs) < 5 {
		additionalClubs, _ := services.GetRecommendedClubs(uid, 20)
		recommendedClubs = append(recommendedClubs, additionalClubs...)
	}

	if len(recommendedClubs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No recommended clubs available",
		})
	}

	// 추천 클럽 셔플
	rand.Shuffle(len(recommendedClubs), func(i, j int) {
		recommendedClubs[i], recommendedClubs[j] = recommendedClubs[j], recommendedClubs[i]
	})

	// 1~3개의 클럽에 그룹 전체를 가입시키기
	numClubsToJoin := rand.Intn(3) + 1 // 1~3개
	if numClubsToJoin > len(recommendedClubs) {
		numClubsToJoin = len(recommendedClubs)
	}

	type JoinResult struct {
		Club         models.Club `json:"club"`
		JoinedUsers  []uint      `json:"joined_users"`
		SkippedUsers []uint      `json:"skipped_users"` // 이미 가입된 사용자
	}

	var results []JoinResult

	for clubIdx := 0; clubIdx < numClubsToJoin; clubIdx++ {
		club := recommendedClubs[clubIdx]
		var joinedUsers []uint
		var skippedUsers []uint

		// 그룹의 각 사용자를 클럽에 가입시키기
		for _, userID := range selectedUsers {
			// 이미 가입된 클럽인지 확인
			var existingMember models.ClubMember
			checkErr := database.DB.Where("club_id = ? AND user_id = ?", club.ID, userID).First(&existingMember).Error

			if checkErr == nil {
				skippedUsers = append(skippedUsers, userID)
				continue
			}

			// 클럽 가입
			newMember := models.ClubMember{
				ClubID:   club.ID,
				UserID:   userID,
				JoinedAt: time.Now(),
			}

			if err := database.DB.Create(&newMember).Error; err != nil {
				skippedUsers = append(skippedUsers, userID)
				continue
			}

			joinedUsers = append(joinedUsers, userID)
		}

		// 멤버 수 업데이트
		if len(joinedUsers) > 0 {
			database.DB.Model(&models.Club{}).Where("id = ?", club.ID).
				UpdateColumn("member_count", club.MemberCount+len(joinedUsers))

			results = append(results, JoinResult{
				Club:         club,
				JoinedUsers:  joinedUsers,
				SkippedUsers: skippedUsers,
			})
		}
	}

	if len(results) == 0 {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "All users were already members of selected clubs",
			"data": fiber.Map{
				"group_users":     selectedUsers,
				"attempted_clubs": numClubsToJoin,
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Successfully matched group to %d club(s)", len(results)),
		"data": fiber.Map{
			"group_size":      len(selectedUsers),
			"group_users":     selectedUsers,
			"matched_clubs":   results,
			"total_attempted": numClubsToJoin,
		},
	})
}

// 전체 사용자 그룹 매칭 - 비슷한 성향의 사용자들을 그룹화하여 클럽에 매칭
func MatchAllUsersToClubs(c *fiber.Ctx) error {
	err := services.MatchUsersToClubs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to match users to clubs",
			"details": err.Error(),
		})
	}

	// 매칭 결과 통계
	var totalMembers int64
	database.DB.Model(&models.ClubMember{}).Count(&totalMembers)

	var totalClubs int64
	database.DB.Model(&models.Club{}).Where("member_count > 0").Count(&totalClubs)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully matched users to clubs based on similar profiles",
		"data": fiber.Map{
			"total_memberships": totalMembers,
			"active_clubs":      totalClubs,
		},
	})
}

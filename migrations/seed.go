package migrations

import (
	"log"
	"ongi-back/database"
	"ongi-back/models"
)

// SeedQuestions - 초기 설문 데이터 생성 (예/아니오 형식)
func SeedQuestions() error {
	questions := []models.Question{
		// 사교성 관련 질문 (2개)
		{
			QuestionText: "처음 만난 사람에게 먼저 말을 거는 편인가요?",
			Order:        1,
			Category:     "sociality",
			Options: []models.Option{
				{OptionText: "예", Score: 5, Weight: "sociality"},
				{OptionText: "아니오", Score: 1, Weight: "intimacy"},
			},
		},
		{
			QuestionText: "큰 모임(10명 이상)보다 소규모 모임을 선호하나요?",
			Order:        2,
			Category:     "sociality",
			Options: []models.Option{
				{OptionText: "예", Score: 1, Weight: "intimacy"},
				{OptionText: "아니오", Score: 5, Weight: "sociality"},
			},
		},
		{
			QuestionText: "친구들과의 관계에서 중요하게 생각하는 것은?",
			Order:        3,
			Category:     "intimacy",
			Options: []models.Option{
				{OptionText: "많은 사람들과 넓은 인맥을 유지하는 것", Score: 5, Weight: "sociality"},
				{OptionText: "다양한 친구들과 즐거운 시간을 보내는 것", Score: 4, Weight: "activity"},
				{OptionText: "상황에 맞게 다양한 관계를 유지하는 것", Score: 3, Weight: "flexibility"},
				{OptionText: "소수의 친구들과 깊은 대화를 나누는 것", Score: 2, Weight: "intimacy"},
				{OptionText: "오랜 시간 함께한 친구와의 신뢰", Score: 1, Weight: "intimacy"},
			},
		},
		{
			QuestionText: "한 가지 일에 집중할 때 당신은?",
			Order:        4,
			Category:     "immersion",
			Options: []models.Option{
				{OptionText: "시간 가는 줄 모르고 완전히 몰입한다", Score: 5, Weight: "immersion"},
				{OptionText: "집중해서 완성도 높게 끝낸다", Score: 4, Weight: "immersion"},
				{OptionText: "집중하다가 다른 일도 병행한다", Score: 3, Weight: "flexibility"},
				{OptionText: "여러 가지를 동시에 처리하는 편이다", Score: 2, Weight: "activity"},
				{OptionText: "짧게 집중하고 다른 일로 전환한다", Score: 1, Weight: "activity"},
			},
		},
		{
			QuestionText: "갑자기 계획이 변경되면 당신은?",
			Order:        5,
			Category:     "flexibility",
			Options: []models.Option{
				{OptionText: "즐겁게 받아들이고 새로운 계획을 세운다", Score: 5, Weight: "flexibility"},
				{OptionText: "빠르게 적응하고 대처한다", Score: 4, Weight: "flexibility"},
				{OptionText: "약간 당황하지만 적응한다", Score: 3, Weight: "flexibility"},
				{OptionText: "원래 계획을 유지하려고 노력한다", Score: 2, Weight: "immersion"},
				{OptionText: "스트레스를 받고 불편해한다", Score: 1, Weight: "immersion"},
			},
		},
		{
			QuestionText: "주말에 선호하는 활동은?",
			Order:        6,
			Category:     "activity",
			Options: []models.Option{
				{OptionText: "여러 명이 모여 액티비티를 즐긴다", Score: 5, Weight: "activity"},
				{OptionText: "친구들과 외출하거나 새로운 곳을 탐험한다", Score: 4, Weight: "activity"},
				{OptionText: "기분에 따라 외출하거나 집에서 쉰다", Score: 3, Weight: "flexibility"},
				{OptionText: "소수의 친구들과 조용히 만난다", Score: 2, Weight: "intimacy"},
				{OptionText: "집에서 혼자만의 시간을 보낸다", Score: 1, Weight: "intimacy"},
			},
		},
		{
			QuestionText: "팀 프로젝트에서 당신의 스타일은?",
			Order:        7,
			Category:     "sociality",
			Options: []models.Option{
				{OptionText: "리더 역할을 맡아 팀을 이끈다", Score: 5, Weight: "sociality"},
				{OptionText: "적극적으로 의견을 제시하고 참여한다", Score: 4, Weight: "sociality"},
				{OptionText: "상황에 따라 역할을 조절한다", Score: 3, Weight: "flexibility"},
				{OptionText: "맡은 부분에 집중하며 완성도를 높인다", Score: 2, Weight: "immersion"},
				{OptionText: "조용히 자신의 일을 처리한다", Score: 1, Weight: "immersion"},
			},
		},
		{
			QuestionText: "새로운 사람을 만날 때 당신은?",
			Order:        8,
			Category:     "sociality",
			Options: []models.Option{
				{OptionText: "먼저 다가가 대화를 시작한다", Score: 5, Weight: "sociality"},
				{OptionText: "자연스럽게 소통한다", Score: 4, Weight: "sociality"},
				{OptionText: "상황을 보고 행동한다", Score: 3, Weight: "flexibility"},
				{OptionText: "상대방이 먼저 말을 걸기를 기다린다", Score: 2, Weight: "intimacy"},
				{OptionText: "낯을 많이 가리는 편이다", Score: 1, Weight: "intimacy"},
			},
		},
		{
			QuestionText: "취미 활동을 할 때 중요하게 생각하는 것은?",
			Order:        9,
			Category:     "immersion",
			Options: []models.Option{
				{OptionText: "다양한 경험과 새로운 시도", Score: 5, Weight: "activity"},
				{OptionText: "여러 가지를 조금씩 경험하는 것", Score: 4, Weight: "activity"},
				{OptionText: "상황에 따라 깊이와 폭을 조절", Score: 3, Weight: "flexibility"},
				{OptionText: "한 가지를 깊이 있게 파고드는 것", Score: 2, Weight: "immersion"},
				{OptionText: "전문가 수준까지 도달하는 것", Score: 1, Weight: "immersion"},
			},
		},
		{
			QuestionText: "모임의 규모는 어느 정도가 가장 편한가요?",
			Order:        10,
			Category:     "intimacy",
			Options: []models.Option{
				{OptionText: "10명 이상의 큰 모임", Score: 5, Weight: "sociality"},
				{OptionText: "5-10명 정도의 모임", Score: 4, Weight: "sociality"},
				{OptionText: "상황에 따라 다르다", Score: 3, Weight: "flexibility"},
				{OptionText: "2-4명의 소규모 모임", Score: 2, Weight: "intimacy"},
				{OptionText: "1:1 만남이 가장 편하다", Score: 1, Weight: "intimacy"},
			},
		},
	}

	for _, question := range questions {
		var existingQuestion models.Question
		result := database.DB.Where("\"order\" = ?", question.Order).First(&existingQuestion)

		if result.Error != nil {
			// 질문이 없으면 생성
			if err := database.DB.Create(&question).Error; err != nil {
				return err
			}
			log.Printf("Created question %d: %s", question.Order, question.QuestionText)
		} else {
			log.Printf("Question %d already exists, skipping", question.Order)
		}
	}

	log.Println("Question seeding completed")
	return nil
}

// SeedSampleClubs - 샘플 클럽 데이터 생성
func SeedSampleClubs() error {
	clubs := []models.Club{
		{
			Name:        "러닝 크루",
			Description: "함께 달리며 건강을 챙기는 모임",
			Category:    "운동",
			MemberCount: 0,
		},
		{
			Name:        "독서 모임",
			Description: "매달 한 권의 책을 읽고 토론하는 모임",
			Category:    "문화",
			MemberCount: 0,
		},
		{
			Name:        "사진 동호회",
			Description: "함께 출사를 다니며 사진을 공유하는 모임",
			Category:    "취미",
			MemberCount: 0,
		},
		{
			Name:        "요리 클래스",
			Description: "다양한 요리를 배우고 함께 만들어보는 모임",
			Category:    "취미",
			MemberCount: 0,
		},
		{
			Name:        "프로그래밍 스터디",
			Description: "함께 코딩을 공부하고 프로젝트를 진행하는 모임",
			Category:    "학습",
			MemberCount: 0,
		},
	}

	for _, club := range clubs {
		var existing models.Club
		result := database.DB.Where("name = ?", club.Name).First(&existing)

		if result.Error != nil {
			if err := database.DB.Create(&club).Error; err != nil {
				return err
			}
			log.Printf("Created club: %s", club.Name)
		} else {
			log.Printf("Club '%s' already exists, skipping", club.Name)
		}
	}

	log.Println("Club seeding completed")
	return nil
}

// SeedAll - 모든 초기 데이터 생성
func SeedAll() error {
	log.Println("Starting database seeding...")

	if err := SeedQuestions(); err != nil {
		return err
	}

	if err := SeedSampleClubs(); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

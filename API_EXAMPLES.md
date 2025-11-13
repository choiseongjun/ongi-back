# API 사용 예제

## 전체 플로우

### 1단계: 사용자 생성

```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "테스트유저"
  }'
```

**응답:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "test@example.com",
    "name": "테스트유저",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2단계: 설문 질문 조회

```bash
curl http://localhost:3000/api/v1/questions
```

**응답:** 10개의 질문과 각 질문의 5개 옵션

### 3단계: 설문 답변 제출

```bash
curl -X POST http://localhost:3000/api/v1/answers/batch \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "answers": [
      {"question_id": 1, "option_id": 3},
      {"question_id": 2, "option_id": 4},
      {"question_id": 3, "option_id": 2},
      {"question_id": 4, "option_id": 5},
      {"question_id": 5, "option_id": 3},
      {"question_id": 6, "option_id": 2},
      {"question_id": 7, "option_id": 4},
      {"question_id": 8, "option_id": 3},
      {"question_id": 9, "option_id": 2},
      {"question_id": 10, "option_id": 3}
    ]
  }'
```

### 4단계: 분석 결과 및 추천 조회

```bash
curl http://localhost:3000/api/v1/results/1
```

**응답:**
```json
{
  "success": true,
  "data": {
    "scores": {
      "sociality_score": 65.0,
      "activity_score": 75.0,
      "intimacy_score": 55.0,
      "immersion_score": 80.0,
      "flexibility_score": 60.0
    },
    "profile_type": "도전적인 탐험가",
    "descriptions": [
      "당신은 상황에 따라 유연하게 대처하며, 내향과 외향의 균형을 잘 맞춥니다.",
      "다양한 활동을 즐기고, 사람들과의 조화를 중요하게 생각합니다.",
      "관심사에 깊이 몰입하며, 전문성을 추구합니다."
    ],
    "recommendations": {
      "clubs": [
        {
          "id": 1,
          "name": "러닝 크루",
          "description": "함께 달리며 건강을 챙기는 모임",
          "category": "운동",
          "member_count": 0
        }
      ],
      "similar_clubs": [...],
      "meetings": [...],
      "similar_users": [...]
    }
  }
}
```

## 클럽 관련 API

### 클럽 생성

```bash
curl -X POST http://localhost:3000/api/v1/clubs \
  -H "Content-Type: application/json" \
  -d '{
    "name": "배드민턴 동호회",
    "description": "주말마다 함께 배드민턴을 치는 모임",
    "category": "운동",
    "image_url": "https://example.com/image.jpg"
  }'
```

### 클럽 가입

```bash
curl -X POST http://localhost:3000/api/v1/clubs/join \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "club_id": 1
  }'
```

### 모든 클럽 조회

```bash
curl http://localhost:3000/api/v1/clubs
```

### 특정 클럽 상세 조회

```bash
curl http://localhost:3000/api/v1/clubs/1
```

## 모임 관련 API

### 모임 생성

```bash
curl -X POST http://localhost:3000/api/v1/meetings \
  -H "Content-Type: application/json" \
  -d '{
    "title": "주말 러닝",
    "description": "한강에서 10km 달리기",
    "club_id": 1,
    "location": "한강공원 반포지구",
    "scheduled_at": "2024-12-21T09:00:00Z",
    "max_members": 20,
    "category": "운동"
  }'
```

### 모든 모임 조회

```bash
curl http://localhost:3000/api/v1/meetings
```

### 특정 모임 상세 조회

```bash
curl http://localhost:3000/api/v1/meetings/1
```

## 사용자 관련 API

### 모든 사용자 조회

```bash
curl http://localhost:3000/api/v1/users
```

### 특정 사용자 조회

```bash
curl http://localhost:3000/api/v1/users/1
```

### 사용자 프로필 조회

```bash
curl http://localhost:3000/api/v1/users/1/profile
```

**응답:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "sociality_score": 65.0,
    "activity_score": 75.0,
    "intimacy_score": 55.0,
    "immersion_score": 80.0,
    "flexibility_score": 60.0,
    "result_summary": "당신은 상황에 따라 유연하게...",
    "profile_type": "도전적인 탐험가",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 사용자의 답변 조회

```bash
curl http://localhost:3000/api/v1/answers/user/1
```

## 건강 체크

```bash
curl http://localhost:3000/api/v1/health
```

**응답:**
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

## Postman Collection

Postman에서 사용하려면 다음과 같이 환경 변수를 설정하세요:

- `base_url`: `http://localhost:3000`
- `user_id`: `1` (생성된 사용자 ID)

## 주의사항

1. 데이터베이스가 실행 중이어야 합니다
2. 먼저 `go run cmd/seed/main.go`를 실행하여 초기 데이터를 생성하세요
3. 설문 답변은 10개 질문 모두에 대해 제출해야 정확한 결과를 얻을 수 있습니다
4. 사용자 추천은 최소 2명 이상의 사용자가 설문을 완료해야 동작합니다

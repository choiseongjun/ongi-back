# 비회원 설문 API 가이드

## 개요

비회원도 설문을 진행하고 결과를 확인할 수 있습니다. 나중에 계정을 생성하면 설문 결과를 계정과 연동할 수 있습니다.

## 비회원 설문 플로우

```
1. 세션 생성 → 2. 설문 조회 → 3. 답변 제출 → 4. 결과 확인 → 5. (선택) 계정 연동
```

## API 엔드포인트

### 1. 세션 생성

비회원 설문을 시작하기 위해 임시 세션을 생성합니다.

**요청:**
```bash
curl -X POST http://localhost:3000/api/v1/guest/session
```

**응답:**
```json
{
  "success": true,
  "data": {
    "session_id": "a1b2c3d4e5f6...",
    "expires_at": "2024-12-28T10:00:00Z"
  },
  "message": "Guest session created. Save this session_id to retrieve your results later."
}
```

⚠️ **중요**: `session_id`를 반드시 저장하세요! (로컬 스토리지, 쿠키 등)
- 세션은 7일 후 만료됩니다
- 세션 ID가 없으면 결과를 다시 볼 수 없습니다

### 2. 설문 질문 조회

```bash
curl http://localhost:3000/api/v1/questions
```

응답: 10개의 질문과 각 질문의 5개 옵션

### 3. 답변 제출

```bash
curl -X POST http://localhost:3000/api/v1/guest/answers \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "a1b2c3d4e5f6...",
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

**응답:**
```json
{
  "success": true,
  "message": "Answers submitted successfully"
}
```

### 4. 결과 조회

```bash
curl http://localhost:3000/api/v1/guest/result/a1b2c3d4e5f6...
```

**응답:**
```json
{
  "success": true,
  "data": {
    "session_id": "a1b2c3d4e5f6...",
    "is_linked": false,
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
      "clubs": [...],
      "similar_clubs": [...],
      "meetings": [...],
      "similar_profiles": [
        {
          "session_id": "xyz123...",
          "user_id": 5,
          "user": {
            "id": 5,
            "name": "홍길동",
            "email": "hong@example.com"
          },
          "similarity": 87.5
        }
      ]
    },
    "expires_at": "2024-12-28T10:00:00Z"
  }
}
```

### 5. 세션 정보 조회

세션의 상태를 확인합니다.

```bash
curl http://localhost:3000/api/v1/guest/session/a1b2c3d4e5f6...
```

**응답:**
```json
{
  "success": true,
  "data": {
    "session_id": "a1b2c3d4e5f6...",
    "is_linked": false,
    "linked_user_id": null,
    "has_results": true,
    "profile_type": "도전적인 탐험가",
    "expires_at": "2024-12-28T10:00:00Z",
    "created_at": "2024-12-21T10:00:00Z"
  }
}
```

## 계정 연동

### 6. 계정과 연동

비회원으로 진행한 설문 결과를 나중에 생성한 계정과 연동합니다.

**1단계: 계정 생성**
```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "홍길동"
  }'
```

**응답:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "홍길동"
  }
}
```

**2단계: 세션 연동**
```bash
curl -X POST http://localhost:3000/api/v1/guest/link \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "a1b2c3d4e5f6...",
    "user_id": 1
  }'
```

**응답:**
```json
{
  "success": true,
  "message": "Session successfully linked to user account",
  "data": {
    "user_id": 1,
    "session_id": "a1b2c3d4e5f6..."
  }
}
```

연동 후 효과:
- 설문 답변이 사용자 계정으로 복사됨
- UserProfile이 생성/업데이트됨
- 세션이 만료되지 않고 영구 보존됨
- 회원 전용 기능 이용 가능

**3단계: 회원 결과 조회**
```bash
curl http://localhost:3000/api/v1/results/1
```

## 궁합 계산

두 프로필 간의 궁합을 계산합니다. (비회원끼리도 가능)

```bash
curl -X POST http://localhost:3000/api/v1/guest/compatibility \
  -H "Content-Type: application/json" \
  -d '{
    "session_id_1": "a1b2c3d4e5f6...",
    "session_id_2": "x7y8z9w0v1u2..."
  }'
```

**응답:**
```json
{
  "success": true,
  "data": {
    "overall_score": 87.5,
    "rating": "최고의 궁합",
    "description": "매우 비슷한 성향으로 서로 잘 맞을 것입니다",
    "details": {
      "sociality_match": 92.0,
      "activity_match": 85.0,
      "intimacy_match": 88.0,
      "immersion_match": 90.0,
      "flexibility_match": 83.0
    }
  }
}
```

**궁합 등급:**
- 80점 이상: "최고의 궁합" - 매우 비슷한 성향
- 70-80점: "좋은 궁합" - 비슷한 성향
- 60-70점: "보통 궁합" - 조화로운 관계
- 50-60점: "상호보완적" - 다른 성향, 자극적
- 50점 미만: "흥미로운 조합" - 배울 점이 많음

## 성능 최적화

### 벡터 연산
- **병렬 처리**: CPU 코어 수만큼 워커를 사용하여 유사도 계산
- **사전 계산**: 벡터 크기(Magnitude)를 미리 계산하여 DB에 저장
- **고속 알고리즘**:
  - 유클리드 거리 (Euclidean Distance)
  - 코사인 유사도 (Cosine Similarity)
  - 맨해튼 거리 (Manhattan Distance)

### Go의 장점
```go
// 병렬 유사도 계산 예제
workers := runtime.NumCPU()  // CPU 코어 수
results := utils.BatchSimilarity(target, vectors, workers)
```

- **동시성**: 고루틴을 활용한 병렬 처리
- **성능**: C/C++ 수준의 빠른 연산
- **메모리 효율**: 효율적인 메모리 관리

## 프론트엔드 통합 예제

### React/JavaScript
```javascript
// 1. 세션 생성
const createSession = async () => {
  const res = await fetch('http://localhost:3000/api/v1/guest/session', {
    method: 'POST'
  });
  const data = await res.json();

  // 로컬 스토리지에 저장
  localStorage.setItem('session_id', data.data.session_id);
  return data.data.session_id;
};

// 2. 답변 제출
const submitAnswers = async (answers) => {
  const sessionId = localStorage.getItem('session_id');

  await fetch('http://localhost:3000/api/v1/guest/answers', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      session_id: sessionId,
      answers: answers
    })
  });
};

// 3. 결과 조회
const getResult = async () => {
  const sessionId = localStorage.getItem('session_id');

  const res = await fetch(`http://localhost:3000/api/v1/guest/result/${sessionId}`);
  return await res.json();
};

// 4. 계정 연동
const linkAccount = async (userId) => {
  const sessionId = localStorage.getItem('session_id');

  await fetch('http://localhost:3000/api/v1/guest/link', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      session_id: sessionId,
      user_id: userId
    })
  });

  // 연동 후 세션 ID 삭제
  localStorage.removeItem('session_id');
};
```

## 주의사항

1. **세션 ID 관리**
   - 클라이언트에서 안전하게 보관 (로컬 스토리지/쿠키)
   - 세션 ID를 잃어버리면 결과 복구 불가능
   - 7일 후 자동 만료 (연동 안 된 경우)

2. **보안**
   - 세션 ID는 추측 불가능한 랜덤 문자열 (16바이트 hex)
   - HTTPS 사용 권장
   - 민감한 정보는 세션에 저장하지 않음

3. **성능**
   - 유사도 계산은 CPU 집약적
   - 많은 사용자가 있을 때는 캐싱 권장
   - 정기적으로 만료된 세션 정리 필요

## 만료된 세션 정리

서버에서 주기적으로 실행:
```go
// 크론 작업 등으로 주기적 실행
services.CleanExpiredSessions()
```

또는 수동 실행:
```bash
# 별도 스크립트 작성 필요
go run cmd/cleanup/main.go
```

# 인증 및 프로필 API 문서

## 목차
1. [카카오 로그인 API](#카카오-로그인-api)
2. [프로필 생성/수정 API](#프로필-생성수정-api)
3. [프로필 조회 API](#프로필-조회-api)

---

## 카카오 로그인 API

### POST /api/v1/auth/kakao/login

카카오 Access Token을 검증하고 회원가입/로그인 처리 후 JWT 토큰을 발급합니다.

#### Request

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "access_token": "카카오_액세스_토큰"
}
```

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| access_token | string | O | 카카오 OAuth로 받은 Access Token |

#### Response

**성공 (200 OK):**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@kakao.com",
    "name": "홍길동",
    "created_at": "2024-11-13T10:00:00Z",
    "updated_at": "2024-11-13T10:00:00Z"
  },
  "is_new_user": true
}
```

| 필드 | 타입 | 설명 |
|------|------|------|
| success | boolean | 요청 성공 여부 |
| token | string | JWT 토큰 (7일 유효) |
| user | object | 사용자 정보 |
| is_new_user | boolean | 신규 가입 여부 (true: 신규, false: 기존 회원) |

**실패 (400 Bad Request):**
```json
{
  "success": false,
  "error": "Access token is required"
}
```

**실패 (401 Unauthorized):**
```json
{
  "success": false,
  "error": "Invalid kakao access token",
  "details": "kakao API returned status 401: ..."
}
```

#### 처리 흐름

1. **카카오 API 호출 (Access Token 검증)**
   - 카카오 사용자 정보 API (`https://kapi.kakao.com/v2/user/me`) 호출
   - Access Token 유효성 검증
   - 사용자 정보 (ID, 이메일, 닉네임) 가져오기

2. **회원가입/로그인 처리**
   - 이메일로 기존 사용자 조회
   - 신규 사용자인 경우: DB에 사용자 정보 저장
   - 기존 사용자인 경우: 기존 정보 사용

3. **JWT 토큰 발급**
   - 사용자 ID와 이메일을 포함한 JWT 생성
   - 토큰 만료 시간: 7일
   - 서명 알고리즘: HS256

#### cURL 예제

```bash
curl -X POST http://localhost:3000/api/v1/auth/kakao/login \
  -H "Content-Type: application/json" \
  -d '{
    "access_token": "YOUR_KAKAO_ACCESS_TOKEN"
  }'
```

---

## 프로필 생성/수정 API

### POST /api/v1/users/profile

사용자의 성향 프로필을 생성하거나 수정합니다.

#### Request

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}  (선택사항)
```

**Body:**
```json
{
  "user_id": 1,
  "sociality_score": 75.5,
  "activity_score": 80.0,
  "intimacy_score": 65.5,
  "immersion_score": 70.0,
  "flexibility_score": 85.5,
  "result_summary": "당신은 활발하고 적극적인 성향을 가지고 있습니다.",
  "profile_type": "열정적인 사교가"
}
```

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| user_id | uint | O | 사용자 ID |
| sociality_score | float64 | X | 사교성 점수 (0-100) |
| activity_score | float64 | X | 활동성 점수 (0-100) |
| intimacy_score | float64 | X | 친밀도 점수 (0-100) |
| immersion_score | float64 | X | 몰입도 점수 (0-100) |
| flexibility_score | float64 | X | 유연성 점수 (0-100) |
| result_summary | string | X | 결과 요약 |
| profile_type | string | X | 프로필 타입 |

#### Response

**성공 - 신규 생성 (201 Created):**
```json
{
  "success": true,
  "message": "Profile created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "sociality_score": 75.5,
    "activity_score": 80.0,
    "intimacy_score": 65.5,
    "immersion_score": 70.0,
    "flexibility_score": 85.5,
    "result_summary": "당신은 활발하고 적극적인 성향을 가지고 있습니다.",
    "profile_type": "열정적인 사교가",
    "created_at": "2024-11-13T10:00:00Z",
    "updated_at": "2024-11-13T10:00:00Z"
  }
}
```

**성공 - 업데이트 (200 OK):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "sociality_score": 75.5,
    "activity_score": 80.0,
    "intimacy_score": 65.5,
    "immersion_score": 70.0,
    "flexibility_score": 85.5,
    "result_summary": "당신은 활발하고 적극적인 성향을 가지고 있습니다.",
    "profile_type": "열정적인 사교가",
    "created_at": "2024-11-13T10:00:00Z",
    "updated_at": "2024-11-13T11:30:00Z"
  }
}
```

**실패 (400 Bad Request):**
```json
{
  "success": false,
  "error": "Invalid request body"
}
```

**실패 (404 Not Found):**
```json
{
  "success": false,
  "error": "User not found"
}
```

#### cURL 예제

```bash
# 프로필 생성/수정
curl -X POST http://localhost:3000/api/v1/users/profile \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "sociality_score": 75.5,
    "activity_score": 80.0,
    "intimacy_score": 65.5,
    "immersion_score": 70.0,
    "flexibility_score": 85.5,
    "result_summary": "당신은 활발하고 적극적인 성향을 가지고 있습니다.",
    "profile_type": "열정적인 사교가"
  }'
```

---

## 프로필 조회 API

### GET /api/v1/users/:userId/profile

사용자의 프로필, 성향 분석, 유사 사용자, 추천 클럽 정보를 조회합니다.

#### Request

**Headers:**
```
Authorization: Bearer {JWT_TOKEN}  (선택사항)
```

**Path Parameters:**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| userId | uint | O | 사용자 ID |

#### Response

**성공 (200 OK):**
```json
{
  "success": true,
  "data": {
    "profile": {
      "id": 1,
      "user_id": 1,
      "sociality_score": 75.5,
      "activity_score": 80.0,
      "intimacy_score": 65.5,
      "immersion_score": 70.0,
      "flexibility_score": 85.5,
      "result_summary": "당신은 활발하고 적극적인 성향을 가지고 있습니다.",
      "profile_type": "열정적인 사교가",
      "created_at": "2024-11-13T10:00:00Z",
      "updated_at": "2024-11-13T10:00:00Z"
    },
    "tendencies": {
      "sociality": {
        "score": 75.5,
        "level": "높음"
      },
      "activity": {
        "score": 80.0,
        "level": "매우 높음"
      },
      "intimacy": {
        "score": 65.5,
        "level": "높음"
      },
      "immersion": {
        "score": 70.0,
        "level": "높음"
      },
      "flexibility": {
        "score": 85.5,
        "level": "매우 높음"
      }
    },
    "similar_users": [
      {
        "user": {
          "id": 2,
          "email": "user2@example.com",
          "name": "김철수"
        },
        "similarity": 0.92
      }
    ],
    "recommended_clubs": [
      {
        "id": 1,
        "name": "등산 동호회",
        "description": "주말마다 산을 오르는 모임",
        "category": "운동",
        "member_count": 15
      }
    ]
  }
}
```

**실패 (404 Not Found):**
```json
{
  "error": "Profile not found"
}
```

#### 성향 레벨 기준

| 점수 | 레벨 |
|------|------|
| 80 이상 | 매우 높음 |
| 60~79 | 높음 |
| 40~59 | 보통 |
| 20~39 | 낮음 |
| 20 미만 | 매우 낮음 |

#### cURL 예제

```bash
# 사용자 프로필 조회
curl -X GET http://localhost:3000/api/v1/users/1/profile
```

---

## JWT 토큰 사용법

발급받은 JWT 토큰은 이후 API 요청 시 Authorization 헤더에 포함하여 사용합니다.

```bash
curl -X GET http://localhost:3000/api/v1/users/1/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### JWT 토큰 정보

- **유효 기간**: 7일
- **알고리즘**: HS256
- **포함 정보**:
  - `user_id`: 사용자 ID
  - `email`: 사용자 이메일
  - `exp`: 만료 시간
  - `iat`: 발급 시간
  - `iss`: 발급자 (ongi-back)

---

## 전체 플로우 예시

### 1. 카카오 로그인
```bash
# 1-1. 프론트엔드에서 카카오 OAuth로 Access Token 받기
# 1-2. Access Token으로 백엔드 로그인
curl -X POST http://localhost:3000/api/v1/auth/kakao/login \
  -H "Content-Type: application/json" \
  -d '{"access_token": "KAKAO_ACCESS_TOKEN"}'

# Response: JWT 토큰 받음
# {"success": true, "token": "JWT_TOKEN", "user": {...}, "is_new_user": true}
```

### 2. 프로필 생성
```bash
# 설문 결과를 바탕으로 프로필 생성
curl -X POST http://localhost:3000/api/v1/users/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer JWT_TOKEN" \
  -d '{
    "user_id": 1,
    "sociality_score": 75.5,
    "activity_score": 80.0,
    "intimacy_score": 65.5,
    "immersion_score": 70.0,
    "flexibility_score": 85.5,
    "result_summary": "활발한 성향",
    "profile_type": "열정적인 사교가"
  }'
```

### 3. 프로필 조회
```bash
# 프로필 및 추천 정보 조회
curl -X GET http://localhost:3000/api/v1/users/1/profile \
  -H "Authorization: Bearer JWT_TOKEN"
```

---

## 에러 코드

| HTTP 상태 코드 | 설명 |
|----------------|------|
| 200 OK | 요청 성공 |
| 201 Created | 리소스 생성 성공 |
| 400 Bad Request | 잘못된 요청 (필수 파라미터 누락, 형식 오류) |
| 401 Unauthorized | 인증 실패 (유효하지 않은 토큰) |
| 404 Not Found | 리소스를 찾을 수 없음 |
| 500 Internal Server Error | 서버 내부 오류 |

---

## 환경 변수 설정

`.env` 파일에 다음 설정이 필요합니다:

```env
# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production

# Note: 카카오 Access Token은 클라이언트(프론트엔드)에서 제공
# 백엔드에서는 토큰 검증만 수행하므로 별도의 Kakao API Key 불필요
```

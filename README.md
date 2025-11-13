# Ongi Backend - 성향 기반 취미 매칭 플랫폼

Fiber와 PostgreSQL을 사용한 성향 기반 취미 매칭 플랫폼 백엔드 API

## 주요 기능

- **성향 설문**: 10개의 질문을 통한 사용자 성향 분석
- **결과 분석**: 사교성, 활동성, 친밀도, 몰입도, 유연성 점수 계산
- **추천 시스템**:
  - 비슷한 성향의 사용자 추천
  - 맞춤형 클럽/모임 추천
  - 유사 사용자가 많은 클럽 추천
- **클럽 관리**: 클럽 생성, 가입, 멤버 관리
- **모임 관리**: 모임 생성 및 일정 관리

## 기술 스택

- **언어**: Go 1.21+
- **웹 프레임워크**: Fiber v2
- **ORM**: GORM
- **데이터베이스**: PostgreSQL
- **환경 변수 관리**: godotenv

## 프로젝트 구조

```
ongi-back/
├── cmd/
│   ├── api/          # 메인 API 서버
│   └── seed/         # 데이터베이스 시드
├── config/           # 설정 관리
├── database/         # 데이터베이스 연결
├── handlers/         # HTTP 핸들러
├── models/           # 데이터 모델
├── routes/           # 라우트 정의
├── services/         # 비즈니스 로직
├── migrations/       # 마이그레이션 및 시드
└── utils/            # 유틸리티 함수
```

## 설치 및 실행

### 1. 의존성 설치

```bash
go mod tidy
```

### 2. 환경 변수 설정

`.env.example` 파일을 `.env`로 복사하고 값을 수정합니다:

```bash
cp .env.example .env
```

```env
PORT=3000
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ongi_db
DB_SSLMODE=disable
```

### 3. PostgreSQL 설정

PostgreSQL이 설치되어 있어야 합니다. 데이터베이스를 생성합니다:

```sql
CREATE DATABASE ongi_db;
```

### 4. 데이터베이스 시드 (초기 데이터 생성)

```bash
go run cmd/seed/main.go
```

이 명령어는 다음을 수행합니다:
- 테이블 마이그레이션
- 10개의 설문 질문 및 옵션 생성
- 샘플 클럽 데이터 생성

### 5. 서버 실행

```bash
go run cmd/api/main.go
```

서버가 `http://localhost:3000`에서 실행됩니다.

## API 엔드포인트

### Health Check
- `GET /api/v1/health` - 서버 상태 확인

### Users
- `GET /api/v1/users` - 모든 사용자 조회
- `POST /api/v1/users` - 사용자 생성
- `GET /api/v1/users/:id` - 특정 사용자 조회
- `GET /api/v1/users/:id/profile` - 사용자 프로필 조회

### Questions (설문)
- `GET /api/v1/questions` - 모든 질문 조회
- `GET /api/v1/questions/:id` - 특정 질문 조회

### Answers (답변)
- `POST /api/v1/answers` - 단일 답변 제출
- `POST /api/v1/answers/batch` - 여러 답변 한번에 제출
- `GET /api/v1/answers/user/:userId` - 사용자의 모든 답변 조회

### Results (결과)
- `GET /api/v1/results/:userId` - 사용자 분석 결과 및 추천 조회

### Clubs (클럽)
- `GET /api/v1/clubs` - 모든 클럽 조회
- `POST /api/v1/clubs` - 클럽 생성
- `GET /api/v1/clubs/:id` - 특정 클럽 조회
- `POST /api/v1/clubs/join` - 클럽 가입

### Meetings (모임)
- `GET /api/v1/meetings` - 모든 모임 조회
- `POST /api/v1/meetings` - 모임 생성
- `GET /api/v1/meetings/:id` - 특정 모임 조회

## 사용 예제

### 1. 사용자 생성

```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "홍길동"
  }'
```

### 2. 설문 질문 조회

```bash
curl http://localhost:3000/api/v1/questions
```

### 3. 답변 제출 (일괄)

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

### 4. 분석 결과 조회

```bash
curl http://localhost:3000/api/v1/results/1
```

응답 예시:
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
      "clubs": [...],
      "similar_clubs": [...],
      "meetings": [...],
      "similar_users": [...]
    }
  }
}
```

## 성향 분석 기준

### 점수 카테고리
- **사교성 (Sociality)**: 사람들과의 상호작용 선호도
- **활동성 (Activity)**: 새로운 활동과 도전에 대한 적극성
- **친밀도 (Intimacy)**: 깊은 관계 형성에 대한 선호도
- **몰입도 (Immersion)**: 한 가지에 집중하는 경향
- **유연성 (Flexibility)**: 상황 변화에 대한 적응력

### 프로필 타입
- 열정적인 사교가
- 따뜻한 조력자
- 도전적인 탐험가
- 깊이있는 전문가
- 유연한 적응형
- 친근한 외향형
- 집중하는 몰입형
- 균형잡힌 조화형

## 추천 알고리즘

### 유사 사용자 찾기
- 유클리드 거리 기반 5차원 벡터 유사도 계산
- 5가지 성향 점수를 기반으로 가장 유사한 사용자 추천

### 클럽/모임 추천
- 사용자 성향에 따른 맞춤형 추천
- 유사 사용자가 많이 가입한 클럽 우선 추천
- 사교성 높음 → 멤버가 많은 클럽
- 친밀도 높음 → 소규모 클럽

## 개발

### 테스트
```bash
go test ./...
```

### 빌드
```bash
# API 서버
go build -o bin/server cmd/api/main.go

# 시드 프로그램
go build -o bin/seed cmd/seed/main.go
```

## 라이센스

MIT License

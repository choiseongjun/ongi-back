# Ongi Backend - 프로젝트 요약

## 프로젝트 개요

성향 기반 취미 매칭 플랫폼의 백엔드 API 서버입니다.
사용자의 성향을 10개 질문으로 분석하여 맞춤형 클럽/모임/사용자를 추천합니다.

## 기술 스택

- **언어**: Go 1.21+
- **프레임워크**: Fiber v2 (고성능 웹 프레임워크)
- **ORM**: GORM (PostgreSQL 지원)
- **데이터베이스**: PostgreSQL 16
- **환경 설정**: godotenv

## 핵심 기능

### 1. 성향 분석 시스템
- 10개 질문, 각 5개 선택지
- 5가지 성향 점수 계산:
  - 사교성 (Sociality)
  - 활동성 (Activity)
  - 친밀도 (Intimacy)
  - 몰입도 (Immersion)
  - 유연성 (Flexibility)

### 2. 프로필 타입 분류
8가지 프로필 타입:
- 열정적인 사교가
- 따뜻한 조력자
- 도전적인 탐험가
- 깊이있는 전문가
- 유연한 적응형
- 친근한 외향형
- 집중하는 몰입형
- 균형잡힌 조화형

### 3. 추천 시스템
- **유사 사용자 추천**: 유클리드 거리 기반 5차원 벡터 유사도
- **클럽 추천**: 성향 맞춤형 클럽 추천
- **모임 추천**: 활동성/친밀도 기반 모임 추천
- **인기 클럽 추천**: 유사 사용자가 많이 가입한 클럽

## 프로젝트 구조

```
ongi-back/
├── cmd/
│   ├── api/main.go          # 메인 서버
│   └── seed/main.go         # 데이터베이스 시드
├── config/
│   └── config.go            # 설정 관리
├── database/
│   └── database.go          # DB 연결 및 마이그레이션
├── handlers/
│   ├── user.go              # 사용자 핸들러
│   ├── question.go          # 설문 핸들러
│   ├── result.go            # 결과 핸들러
│   └── club.go              # 클럽/모임 핸들러
├── models/
│   ├── user.go              # User, UserProfile
│   ├── question.go          # Question, Option, UserAnswer
│   └── club.go              # Club, ClubMember, Meeting
├── services/
│   ├── analysis.go          # 성향 분석 로직
│   └── recommendation.go    # 추천 알고리즘
├── routes/
│   └── routes.go            # API 라우트 설정
├── migrations/
│   └── seed.go              # 초기 데이터 생성
└── utils/                   # 유틸리티 함수
```

## 데이터베이스 스키마

### 주요 테이블
1. **users**: 사용자 기본 정보
2. **user_profiles**: 사용자 성향 분석 결과
3. **questions**: 설문 질문
4. **options**: 질문 선택지
5. **user_answers**: 사용자 답변
6. **clubs**: 클럽 정보
7. **club_members**: 클럽 멤버십
8. **meetings**: 모임 정보

## API 엔드포인트

### Users
- `GET /api/v1/users` - 모든 사용자
- `POST /api/v1/users` - 사용자 생성
- `GET /api/v1/users/:id` - 사용자 조회
- `GET /api/v1/users/:id/profile` - 프로필 조회

### Questions
- `GET /api/v1/questions` - 모든 질문
- `GET /api/v1/questions/:id` - 특정 질문

### Answers
- `POST /api/v1/answers` - 답변 제출
- `POST /api/v1/answers/batch` - 일괄 답변 제출
- `GET /api/v1/answers/user/:userId` - 사용자 답변 조회

### Results
- `GET /api/v1/results/:userId` - 분석 결과 및 추천

### Clubs
- `GET /api/v1/clubs` - 모든 클럽
- `POST /api/v1/clubs` - 클럽 생성
- `GET /api/v1/clubs/:id` - 클럽 조회
- `POST /api/v1/clubs/join` - 클럽 가입

### Meetings
- `GET /api/v1/meetings` - 모든 모임
- `POST /api/v1/meetings` - 모임 생성
- `GET /api/v1/meetings/:id` - 모임 조회

## 시작하기

### 1. Docker로 PostgreSQL 시작

```bash
docker-compose up -d
```

또는 로컬 PostgreSQL 사용:
```sql
CREATE DATABASE ongi_db;
```

### 2. 환경 설정

`.env` 파일 확인 및 수정

### 3. 의존성 설치

```bash
make install
# 또는
go mod tidy
```

### 4. 데이터베이스 시드

```bash
make seed
# 또는
go run cmd/seed/main.go
```

### 5. 서버 실행

```bash
make run
# 또는
go run cmd/api/main.go
```

서버: http://localhost:3000

## 주요 알고리즘

### 점수 계산
- 각 질문의 선택지는 1-5점
- 카테고리별 평균 계산
- 0-100 스케일로 변환

### 유사도 계산
```go
distance = sqrt(
  (s1-s2)² + (a1-a2)² + (i1-i2)² + (im1-im2)² + (f1-f2)²
)
similarity = (1 - distance/maxDistance) * 100
```

### 추천 로직
- 사교성 높음 → 대규모 클럽 추천
- 친밀도 높음 → 소규모 클럽 추천
- 활동성 높음 → 다양한 모임 추천
- 유사 사용자 기반 클럽 필터링

## 빌드 및 배포

### 로컬 빌드
```bash
make build
```

### 실행 파일
```bash
./bin/server  # 서버 실행
./bin/seed    # 시드 실행
```

## 개발 도구

- **Makefile**: 자주 사용하는 명령어 단축
- **docker-compose.yml**: PostgreSQL 컨테이너
- **.env**: 환경 변수 관리
- **API_EXAMPLES.md**: API 사용 예제

## 다음 단계

### 추가 구현 고려사항
1. JWT 인증 시스템
2. 실시간 채팅 (WebSocket)
3. 이미지 업로드 (S3)
4. 이메일 알림
5. 모임 참가 기능
6. 리뷰/평점 시스템
7. 검색 기능 (Elasticsearch)
8. 캐싱 (Redis)
9. API Rate Limiting
10. 로깅 및 모니터링

### 성능 최적화
- 데이터베이스 인덱싱
- 쿼리 최적화
- 캐싱 레이어
- 페이지네이션

### 테스트
- 유닛 테스트
- 통합 테스트
- E2E 테스트
- 부하 테스트

## 라이센스

MIT License

## 기여

이슈와 PR은 언제나 환영합니다!

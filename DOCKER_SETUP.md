# Docker Setup Guide

이 가이드는 Ongi 백엔드를 Docker Compose로 실행하는 방법을 설명합니다.

## 사전 요구사항

- Docker Desktop 설치 필수
- Docker Compose 포함 (Docker Desktop에 기본 포함)

## 빠른 시작

### 1. 프로젝트 클론 (이미 완료된 경우 건너뛰기)

```bash
git clone <repository-url>
cd ongi-back
```

### 2. Docker Compose로 실행

```bash
# 백그라운드에서 실행
docker-compose up -d

# 또는 로그를 보면서 실행
docker-compose up
```

이 명령으로 다음이 자동으로 실행됩니다:
- PostgreSQL 데이터베이스 (포트 5432)
- Ongi 백엔드 서버 (포트 5000)
- 데이터베이스 테이블 자동 생성 (migration)

### 3. 서버 확인

서버가 정상적으로 실행되었는지 확인:

```bash
# 로그 확인
docker-compose logs -f backend

# 헬스체크 (브라우저 또는 curl)
curl http://localhost:5000/health
```

### 4. 중지 및 재시작

```bash
# 중지
docker-compose down

# 중지 + 데이터 삭제
docker-compose down -v

# 재시작
docker-compose restart

# 특정 서비스만 재시작
docker-compose restart backend
```

## 환경 변수 설정 (선택사항)

기본 설정으로 실행되지만, 커스터마이징이 필요한 경우:

1. `.env` 파일 생성:
```bash
cp .env.example .env
```

2. `.env` 파일 수정 (필요시)

**참고**: Docker Compose 사용 시 환경 변수는 `docker-compose.yml`에 이미 설정되어 있습니다.

## API 엔드포인트

서버 실행 후 다음 주소로 접근:

- 백엔드 API: `http://localhost:5000`
- PostgreSQL: `localhost:5432`

## 문제 해결

### 포트가 이미 사용 중인 경우

`docker-compose.yml` 파일에서 포트 번호 변경:

```yaml
services:
  backend:
    ports:
      - "5001:5000"  # 5000 대신 5001 사용

  postgres:
    ports:
      - "5433:5432"  # 5432 대신 5433 사용
```

### 컨테이너 재빌드

코드 변경 후 이미지를 다시 빌드해야 하는 경우:

```bash
docker-compose up -d --build
```

### 로그 확인

```bash
# 모든 서비스 로그
docker-compose logs -f

# 백엔드만
docker-compose logs -f backend

# PostgreSQL만
docker-compose logs -f postgres
```

### 데이터베이스 초기화

데이터베이스를 완전히 초기화하려면:

```bash
docker-compose down -v
docker-compose up -d
```

`-v` 옵션은 볼륨(데이터)도 함께 삭제합니다.

## 로컬 개발 (Docker 없이)

Docker 없이 로컬에서 개발하려면:

1. PostgreSQL을 별도로 실행
2. `.env` 파일 생성 및 설정
3. Go 서버 실행:

```bash
go run cmd/api/main.go
```

## 프론트엔드 개발자를 위한 안내

프론트엔드 개발 시:

1. 이 저장소를 클론
2. `docker-compose up -d` 명령 실행
3. `http://localhost:5000`으로 API 요청

그게 전부입니다! 별도의 설정이나 Go 설치가 필요 없습니다.

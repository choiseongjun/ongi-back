.PHONY: help build run seed test clean dev install

help: ## 도움말 표시
	@echo "사용 가능한 명령어:"
	@echo "  make install   - Go 모듈 의존성 설치"
	@echo "  make seed      - 데이터베이스 시드 (초기 데이터 생성)"
	@echo "  make run       - 서버 실행"
	@echo "  make dev       - 개발 모드로 서버 실행"
	@echo "  make build     - 프로덕션 빌드"
	@echo "  make test      - 테스트 실행"
	@echo "  make clean     - 빌드 파일 정리"

install: ## Go 모듈 의존성 설치
	go mod download
	go mod tidy

seed: ## 데이터베이스 시드
	go run cmd/seed/main.go

run: ## 서버 실행
	go run cmd/api/main.go

dev: ## 개발 모드 (hot reload with air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air가 설치되어 있지 않습니다. 'go install github.com/air-verse/air@latest' 로 설치하세요"; \
		go run cmd/api/main.go; \
	fi

build: ## 프로덕션 빌드
	@echo "Building server..."
	@mkdir -p bin
	go build -o bin/server cmd/api/main.go
	@echo "Building seed..."
	go build -o bin/seed cmd/seed/main.go
	@echo "Build complete!"

test: ## 테스트 실행
	go test -v ./...

test-coverage: ## 테스트 커버리지
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## 빌드 파일 정리
	rm -rf bin/
	rm -f coverage.out

migrate: ## 데이터베이스 마이그레이션만 실행
	@echo "Running migrations..."
	@go run -exec 'go run' cmd/seed/main.go

docker-up: ## Docker Compose로 PostgreSQL 시작
	docker-compose up -d

docker-down: ## Docker Compose 중지
	docker-compose down

docker-logs: ## Docker 로그 확인
	docker-compose logs -f

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	Environment   string
}

var AppConfig *Config

func Load() {
	// .env 파일 로드 (없어도 에러 발생 안함)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	log.Println("Configuration loaded")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

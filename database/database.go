package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ongi-back/models"
)

var DB *gorm.DB

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Connect() error {
	config := DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "ongi_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Database connection established")
	return nil
}

func Migrate() error {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Question{},
		&models.Option{},
		&models.UserAnswer{},
		&models.Club{},
		&models.ClubMember{},
		&models.Meeting{},
		&models.GuestSession{},
		&models.GuestAnswer{},
		&models.SessionVector{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

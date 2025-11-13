package main

import (
	"log"
	"ongi-back/config"
	"ongi-back/database"
	"ongi-back/migrations"
)

func main() {
	log.Println("Starting seed process...")

	// Load configuration
	config.Load()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations first
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed data
	if err := migrations.SeedAll(); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	log.Println("Seed process completed successfully!")
}

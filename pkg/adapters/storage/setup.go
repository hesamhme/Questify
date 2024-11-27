package storage

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"Questify/config" 
	"Questify/pkg/adapters/storage/entities"
)

var DB *gorm.DB

// SetupDatabase initializes the database connection and performs migrations
func SetupDatabase(cfg *config.DatabaseConfig) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Ensure UUID extension is enabled
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}

	// Run migrations
	if err := DB.AutoMigrate(&entities.User{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database setup completed successfully.")
}

package storage

import (
	"fmt"
	"log"

	"Questify/pkg/adapters/storage/entities"
	"Questify/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupDatabase initializes the database connection and performs migrations.
func SetupDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DBName,
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Println("Database connection established.")

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database setup and migrations completed successfully.")

	return db, nil
}

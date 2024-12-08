package storage

import (
	"Questify/config"
	"Questify/pkg/adapters/storage/entities"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresGormConnection(dbConfig config.DB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host, dbConfig.User, dbConfig.Pass, dbConfig.DBName, dbConfig.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func AddExtension(db *gorm.DB) error {
	return db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
}

func Migrate(db *gorm.DB) error {
	migrator := db.Migrator()

	// Ensure the User table exists
	if !migrator.HasTable(&entities.User{}) {
		if err := migrator.CreateTable(&entities.User{}); err != nil {
			return err
		}
	}

	// Ensure the column exists and set a default value
	if !migrator.HasColumn(&entities.User{}, "national_code") {
		err := db.Exec(`ALTER TABLE "users" ADD COLUMN "national_code" char(10) DEFAULT '0000000000'`).Error
		if err != nil {
			return err
		}
		// Apply NOT NULL constraint
		err = db.Exec(`ALTER TABLE "users" ALTER COLUMN "national_code" SET NOT NULL`).Error
		if err != nil {
			return err
		}
	}

	// Perform other migrations for roles, permissions, and user_roles
	err := migrator.AutoMigrate(
		&entities.User{},
		&entities.Role{},
		&entities.Permission{},
		&entities.UserRole{},
		&entities.Survey{},
		&entities.SurveyRequirements{},
		&entities.City{},
		&entities.Question{},
		&entities.Answer{},
		&entities.QuestionChoices{},
		&entities.SurveyUserRole{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Seed permissions after migrations
	err = SeedPermissions(db)
	if err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	return nil
}

// SeedPermissions seeds initial permissions into the database
func SeedPermissions(db *gorm.DB) error {
	permissions := []entities.Permission{
		{ID: 1, Description: "View Survey"},
		{ID: 2, Description: "see selected users vote"},
		{ID: 3, Description: "vote permission"},
		{ID: 4, Description: "edit survey"},
		{ID: 5, Description: "role permissions"},
		{ID: 6, Description: "see reports"},
	}

	for _, perm := range permissions {
		// Use FirstOrCreate to avoid duplicate entries
		err := db.FirstOrCreate(&perm, entities.Permission{ID: perm.ID}).Error
		if err != nil {
			return fmt.Errorf("failed to seed permission ID %d: %w", perm.ID, err)
		}
	}
	return nil

}

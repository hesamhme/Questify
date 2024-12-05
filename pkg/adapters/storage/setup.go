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

	// Ensure the table exists
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

	// Perform other migrations
	return migrator.AutoMigrate(
		&entities.User{},
		&entities.Survey{},
		&entities.SurveyRequirements{},
		&entities.City{},
		&entities.Question{},

		&entities.Answer{},
		&entities.QuestionChoices{},
	)

}

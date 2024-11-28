package service

import (
	"Questify/config"
	"Questify/pkg/adapters/storage"
	"gorm.io/gorm"
)

// AppContainer holds all initialized services and shared dependencies.
type AppContainer struct {
	DB          *gorm.DB
	AuthService *AuthService
}

// NewAppContainer initializes the application container with all services.
func NewAppContainer(cfg *config.Config) (*AppContainer, error) {
	// Step 1: Setup Database
	db, err := storage.SetupDatabase(cfg.DB)
	if err != nil {
		return nil, err
	}

	// Step 2: Initialize Services
	authService := NewAuthService(db, cfg.JWT.Secret)

	// Step 3: Return the AppContainer
	return &AppContainer{
		DB:          db,
		AuthService: authService,
	}, nil
}

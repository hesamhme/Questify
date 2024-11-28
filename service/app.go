package service

import (
	"log"

	"Questify/config"
	"Questify/pkg/adapters/storage"
	"gorm.io/gorm"
)

// AppContainer holds shared dependencies for the application.
type AppContainer struct {
	DB          *gorm.DB
	AuthService *AuthService
	UserService *UserService
}

// NewAppContainer initializes the container with all shared dependencies.
func NewAppContainer(cfg *config.Config) (*AppContainer, error) {
	// Setup database
	db, err := storage.SetupDatabase(&cfg.DB)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := &storage.UserRepository{DB: db}

	// Initialize services
	authService := NewAuthService(db, cfg.JWT.Secret)
	userService := NewUserService(userRepo)

	// Log success and return the container
	log.Println("App container initialized successfully.")
	return &AppContainer{
		DB:          db,
		AuthService: authService,
		UserService: userService,
	}, nil
}

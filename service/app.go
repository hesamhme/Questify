package service

import (
	"errors"
	"Questify/pkg/adapters/storage"
	"Questify/pkg/adapters/storage/entities"

)

// UserInput represents the user input for registration
type UserInput struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	NationalCode string `json:"national_code"`
}

// UserService handles user-related business logic
var UserService = &userService{}

type userService struct{}

// Signup validates the input and registers a new user
func (s *userService) Signup(input UserInput) error {
	// Validate the National ID
	if !IsValidIranianNationalCode(input.NationalCode) {
		return errors.New("invalid National ID")
	}

	// Check for duplicate National Code
	var existingUser entities.User
	if err := storage.DB.Where("national_code = ?", input.NationalCode).First(&existingUser).Error; err == nil {
		return errors.New("user with this National ID already exists")
	}

	// Hash the password (this step will be implemented in the next phase)
	hashedPassword := input.Password // Placeholder for hashing

	// Create and save the user
	user := entities.User{
		Email:        input.Email,
		Password:     hashedPassword,
		NationalCode: input.NationalCode,
	}
	if err := storage.DB.Create(&user).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

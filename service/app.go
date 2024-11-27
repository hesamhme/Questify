package service

import (
	"errors"
	"net/mail"
	"Questify/config"
	"Questify/pkg/adapters/storage"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var appConfig *config.Config

// InitializeServices sets up the service with the necessary configuration
func InitializeServices(cfg *config.Config) {
	appConfig = cfg
}

// UserInput represents the user input for registration
type UserInput struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	NationalCode string `json:"national_code"`
}

// SignupResponse contains the JWT token and user details
type SignupResponse struct {
	Token        string `json:"token"`
	NationalCode string `json:"national_code"`
	Email        string `json:"email"`
}

// UserService handles user-related business logic
var UserService = &userService{}

type userService struct{}

// Signup validates input, registers a new user, and returns a JWT
func (s *userService) Signup(input UserInput) (*SignupResponse, error) {
	// Validate the National ID
	if !IsValidIranianNationalCode(input.NationalCode) {
		return nil, errors.New("invalid National ID")
	}

	// Validate email format
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return nil, errors.New("invalid email format")
	}

	// Check for duplicate National Code
	var existingUser entities.User
	if err := storage.DB.Where("national_code = ?", input.NationalCode).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this National ID already exists")
	}

	// Hash the password
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create and save the user
	user := entities.User{
		Email:        input.Email,
		Password:     hashedPassword,
		NationalCode: input.NationalCode,
	}
	if err := storage.DB.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := jwt.GenerateJWT(user.ID.String(), appConfig.JWT.Secret)
	if err != nil {
		return nil, errors.New("failed to generate JWT")
	}

	return &SignupResponse{
		Token:        token,
		NationalCode: user.NationalCode,
		Email:        user.Email,
	}, nil
}

func (s *userService) Login(email, password string) (string, error) {
	// Find the user by email
	var user entities.User
	if err := storage.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid email or password")
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate a JWT token
	token, err := jwt.GenerateJWT(user.ID.String(), appConfig.JWT.Secret)
	if err != nil {
		return "", errors.New("failed to generate JWT")
	}

	return token, nil
}

// hashPassword hashes the user's password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

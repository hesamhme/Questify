package service

import (
	"Questify/internal/user"
	"Questify/pkg/adapters/storage"
	"Questify/pkg/jwt"
	"errors"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)



// AuthService handles authentication-related logic.
type AuthService struct {
	UserRepo *storage.UserRepository
	JWTSecret string
}

// NewAuthService initializes the AuthService with dependencies.
func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		UserRepo:  &storage.UserRepository{DB: db},
		JWTSecret: jwtSecret,
	}
}

// SignUp creates a new user and returns a JWT token.
func (s *AuthService) SignUp(email, password, nid string) (string, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// Create the user
	newUser := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		NID:      nid,
	}

	createdUser, err := s.UserRepo.CreateUser(newUser)
	if err != nil {
		return "", err
	}

	// Generate a JWT
	token, err := jwt.GenerateToken(createdUser.ID, s.JWTSecret, jwt.DefaultExpiry)
	if err != nil {
		return "", err
	}

	return token, nil
}

// SignIn authenticates a user and returns a JWT token.
func (s *AuthService) SignIn(email, password string) (string, error) {
	// Fetch the user by email
	existingUser, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate a JWT
	token, err := jwt.GenerateToken(existingUser.ID, s.JWTSecret, jwt.DefaultExpiry)
	if err != nil {
		return "", err
	}

	return token, nil
}

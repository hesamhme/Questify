package storage

import (
	userpkg "Questify/internal/user" 
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles database operations for the User entity.
type UserRepository struct {
	DB *gorm.DB
}

// GetUserByID fetches a user by their ID.
func (repo *UserRepository) GetUserByID(userID string) (*userpkg.User, error) {
	var user userpkg.User
	if err := repo.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database.
func (repo *UserRepository) CreateUser(user *userpkg.User) (*userpkg.User, error) {
	// Validate email format
	if !userpkg.IsValidEmail(user.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate Iranian National ID
	if !userpkg.IsValidIranianNationalCode(user.NID) {
		return nil, errors.New("invalid Iranian National ID")
	}

	// Generate a new UUID if not already set
	if user.ID == (uuid.UUID{}) {
		user.ID = uuid.New()
	}

	// Check if email already exists
	existingUser, err := repo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	} else if err != nil && err.Error() != "user not found" {
		// Only propagate errors that aren't "user not found"
		return nil, err
	}

	// Create the user in the database
	if err := repo.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail fetches a user by their email.
func (repo *UserRepository) GetUserByEmail(email string) (*userpkg.User, error) {
	if !userpkg.IsValidEmail(email) {
		return nil, errors.New("invalid email format")
	}
	var user userpkg.User
	if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

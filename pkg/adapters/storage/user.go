package storage

import (
	"Questify/internal/user"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles database operations for the User entity.
type UserRepository struct {
	DB *gorm.DB
}

// GetUserByID fetches a user by their ID.
func (repo *UserRepository) GetUserByID(userID string) (*user.User, error) {
	var user user.User
	if err := repo.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
func (repo *UserRepository) CreateUser(user *user.User) (*user.User, error) {
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
func (repo *UserRepository) GetUserByEmail(email string) (*user.User, error) {
	var user user.User
	if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

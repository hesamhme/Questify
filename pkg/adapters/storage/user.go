package storage

import (
	"Questify/internal/user"
	"errors"

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

// CreateUser creates a new user in the database.
func (repo *UserRepository) CreateUser(user *user.User) (*user.User, error) {
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
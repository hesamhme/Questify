package user

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq" // For PostgreSQL error handling
	"gorm.io/gorm"
)

// Repository handles user-related database operations.
type Repository struct {
	DB *gorm.DB
}

// CreateUser creates a new user in the database.
func (repo *Repository) CreateUser(user *User) (*User, error) {
	user.ID = uuid.New() // Generate a new UUID for the user
	if err := repo.DB.Create(user).Error; err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Handle unique constraint violation
			return nil, errors.New("user already exists")
		}
		return nil, err
	}
	return user, nil
}

// GetUserByEmail fetches a user by their email.
func (repo *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID fetches a user by their ID.
func (repo *Repository) GetUserByID(userID uuid.UUID) (*User, error) {
	var user User
	if err := repo.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user's details in the database.
func (repo *Repository) UpdateUser(user *User) (*User, error) {
	if err := repo.DB.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser deletes a user from the database by ID.
func (repo *Repository) DeleteUser(userID uuid.UUID) error {
	if err := repo.DB.Delete(&User{}, "id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}

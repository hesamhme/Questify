package storage

import (
	"Questify/internal/user"
	"Questify/pkg/adapters/storage/mappers"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (repo *UserRepository) CreateUser(domainUser *user.User) (*user.User, error) {
	entityUser := mappers.ToEntityUser(domainUser)
	if err := repo.DB.Create(entityUser).Error; err != nil {
		return nil, err
	}
	return mappers.ToDomainUser(entityUser), nil
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


// GetUserByID fetches a user by their UUID.
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


// UpdateUser updates an existing user in the database.
func (repo *UserRepository) UpdateUser(user *user.User) (*user.User, error) {
	if err := repo.DB.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

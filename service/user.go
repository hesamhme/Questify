package service

import (
	"Questify/internal/user"
	"Questify/pkg/adapters/storage"
	"errors"
)

// UserService handles user-related business logic.
type UserService struct {
	UserRepo *storage.UserRepository
}

// NewUserService initializes the UserService with a UserRepository.
func NewUserService(repo *storage.UserRepository) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

// GetUserByID fetches a user by their ID.
func (s *UserService) GetUserByID(userID string) (*user.User, error) {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

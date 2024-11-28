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

// NewUserService initializes the UserService with dependencies.
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

// UpdateUser updates the user's profile fields.
func (s *UserService) UpdateUser(userID string, updates map[string]interface{}) (*user.User, error) {
	// Fetch the existing user
	existingUser, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if email, ok := updates["email"].(string); ok {
		existingUser.Email = email
	}
	if nid, ok := updates["nid"].(string); ok {
		existingUser.NID = nid
	}

	// Save the updated user
	updatedUser, err := s.UserRepo.UpdateUser(existingUser)
	if err != nil {
		return nil, errors.New("failed to update user")
	}
	return updatedUser, nil
}

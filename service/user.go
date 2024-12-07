package service

import (
	"Questify/internal/user"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUnAuthorizedUser = errors.New("in order to get list of users you need to login")
)

type UserService struct {
	userOps *user.Ops
}

func NewUserService(userOps *user.Ops) *UserService {
	return &UserService{
		userOps: userOps,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]user.User, int64, error) {
	_, err := s.userOps.GetUserByID(ctx, userID)
	if err != nil {
		return nil, 0, ErrUnAuthorizedUser
	}

	return s.userOps.GetUsers(ctx, page, pageSize)
}

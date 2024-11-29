package service

import (

	"github.com/hesamhme/Qustify/internal/user"


)

type UserService struct {
	userOps                *user.Ops

}

func NewUserService(userOps *user.Ops) *UserService {
	return &UserService{
		userOps:                userOps,
		
	}
}
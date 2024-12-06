package presenter

import (
	"Questify/internal/user"
	"Questify/pkg/fp"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	NationalCode string    `json:"nationalcode"`
}

func RegisterRequestToUser(registerrequest *RegisterRequest) *user.User {

	return &user.User{
		Email:        registerrequest.Email,
		Password:     registerrequest.Password,
		NationalCode: registerrequest.NationalCode,
	}

}

type UserLoginReq struct{
	Email string  `json:"email"`
	Password string `json:"password"`
}

type UserGet struct {
	ID       uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
}

func UserToUserGet(user user.User) UserGet {
	return UserGet{
		ID:       user.ID,
		Email:    user.Email,
	}
}

func BatchUsersToUserGet(users []user.User) []UserGet {
	return fp.Map(users, UserToUserGet)
}
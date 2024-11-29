package presenter

import (
	"github.com/hesamhme/Questify/internal/user"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	ID           uuid.UUID  `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	NationalCode string `json:"nationcode"`
}


func RegisterRequestToUser( registerrequest *RegisterRequest ) *user.User{

	return &user.User{
		Email:        registerrequest.Email,
		Password:      registerrequest.Password,
		NationalCode:  registerrequest.NationalCode,
	}
	
}
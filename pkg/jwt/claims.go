package jwt

import (
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type UserClaims struct {
	jwt2.RegisteredClaims
	UserID   uuid.UUID
	Role     string
	Sections []string
}

type TFAClaims struct {
	jwt2.RegisteredClaims
	Email   string    `json:"email"`
	TFACode string    `json:"tfa_code"`
	Expires time.Time `json:"expires"`
}

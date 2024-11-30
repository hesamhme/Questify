package jwt

import (
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	jwt2.RegisteredClaims
	UserID   uuid.UUID
	Role     string
	Sections []string
}
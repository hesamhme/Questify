package presenter

import "github.com/google/uuid"

// UserClaims represents the claims for an authenticated user.
type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

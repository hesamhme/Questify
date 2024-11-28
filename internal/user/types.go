package user

import "github.com/google/uuid"

// User represents the domain model for a user.
type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	NID      string    `json:"nid"`
}

// SignupInput represents the input for a user sign-up request.
type SignupInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	NID      string `json:"nid" validate:"required"`
}

// LoginInput represents the input for a user login request.
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the output for user-related responses.
type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	NID   string `json:"nid"`
}

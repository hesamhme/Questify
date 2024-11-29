package user

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

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
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	NID   string    `json:"nid"`
}

// IsValidEmail checks if the given email is in a valid format.
func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// IsValidPassword checks password complexity (e.g., min length, special chars).
func IsValidPassword(password string) bool {
	return len(password) >= 8 // Add more checks as needed.
}

// IsValidIranianNationalCode validates an Iranian National ID.
func IsValidIranianNationalCode(input string) bool {
	if len(input) != 10 {
		return false
	}
	for i := 0; i < 10; i++ {
		if input[i] < '0' || input[i] > '9' {
			return false
		}
	}
	check := int(input[9] - '0')
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(input[i]-'0') * (10 - i)
	}
	sum %= 11
	return (sum < 2 && check == sum) || (sum >= 2 && check+sum == 11)
}

// ExtractToken extracts the token from the Authorization header.
func ExtractToken(authorizationHeader string) (string, error) {
	if authorizationHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}

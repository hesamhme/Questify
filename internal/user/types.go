package user

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrDuplicatedUser        = errors.New("email/national ID already exists")
	ErrInvalidNationalCode   = errors.New("national code is invalid")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidPassword       = errors.New("invalid password format")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrInvalidAuthentication = errors.New("email and password don't match")
	ErrInvalidTFA            = errors.New("invalid or expired TFA code")
	ErrHashPasswordFailed    = errors.New("failed to hash password")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters long")
	ErrPasswordMissingNumber = errors.New("password must include at least one number")
	ErrPasswordMissingUpper  = errors.New("password must include at least one uppercase letter")
	ErrPasswordMissingLower  = errors.New("password must include at least one lowercase letter")
	ErrPasswordMissingSymbol = errors.New("password must include at least one special character (!@#$%^&*)")
)

type Repo interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error // New method for updating the user
}

type User struct {
	ID           uuid.UUID
	Email        string
	Password     string
	NationalCode string
	Role         string
	TfaCode      string    // Temporary TFA code
	TfaExpiresAt time.Time // Expiration time for TFA code
	WalletID     uuid.UUID
}

// IsValidIranianNationalCode validates an Iranian National Code.
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

// IsValidEmail validates the format of an email using a regex
func IsValidEmail(email string) bool {
	// Define a regex for validating email
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// LowerCaseEmail converts an email to lowercase
func LowerCaseEmail(email string) string {
	return strings.ToLower(email)
}

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashPasswordFailed
	}
	return string(hashed), nil
}

// SetPassword hashes the password and sets it to the User struct
func (u *User) SetPassword(password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// CheckPasswordHash compares a plaintext password with a hashed password
func CheckPasswordHash(password, hashedPassword string) error {
	// Compare the password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrInvalidAuthentication
	}
	return nil
}

// ValidatePasswordWithFeedback validates a password for strength and provides feedback
func ValidatePasswordWithFeedback(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasNumber, hasUpper, hasLower, hasSymbol bool
	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsSymbol(char) || unicode.IsPunct(char):
			hasSymbol = true
		}
	}

	if !hasNumber {
		return ErrPasswordMissingNumber
	}
	if !hasUpper {
		return ErrPasswordMissingUpper
	}
	if !hasLower {
		return ErrPasswordMissingLower
	}
	if !hasSymbol {
		return ErrPasswordMissingSymbol
	}

	return nil
}

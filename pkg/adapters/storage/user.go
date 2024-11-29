package storage

import (
	"Questify/internal/user"
	"errors"
	"regexp"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles database operations for the User entity.
type UserRepository struct {
	DB *gorm.DB
}

// GetUserByID fetches a user by their ID.
func (repo *UserRepository) GetUserByID(userID string) (*user.User, error) {
	var user user.User
	if err := repo.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
func (repo *UserRepository) CreateUser(user *user.User) (*user.User, error) {
	// Validate email format
	if !isValidEmail(user.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate Iranian National ID
	if !isValidIranianNID(user.NID) {
		return nil, errors.New("invalid Iranian National ID")
	}

	// Generate a new UUID if not already set
	if user.ID == (uuid.UUID{}) {
		user.ID = uuid.New()
	}

	// Check if email already exists
	existingUser, err := repo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	} else if err != nil && err.Error() != "user not found" {
		// Only propagate errors that aren't "user not found"
		return nil, err
	}

	// Create the user in the database
	if err := repo.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// isValidEmail checks if the provided email is valid.
func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// isValidIranianNID validates an Iranian National ID.
func isValidIranianNID(nid string) bool {
	if len(nid) != 10 {
		return false
	}
	for i := 0; i < 10; i++ {
		if nid[i] < '0' || nid[i] > '9' {
			return false
		}
	}
	check := int(nid[9] - '0')
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(nid[i]-'0') * (10 - i)
	}
	sum %= 11
	return (sum < 2 && check == sum) || (sum >= 2 && check+sum == 11)
}

// GetUserByEmail fetches a user by their email.
func (repo *UserRepository) GetUserByEmail(email string) (*user.User, error) {
	var user user.User
	if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

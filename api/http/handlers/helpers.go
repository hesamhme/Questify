package handlers

import (
	"errors"
	"regexp"
	"strings"
)

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

// IsValidIranianNationalCode checks if the NID follows the Iranian validation rules.
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

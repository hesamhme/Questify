package service

import (
	"Questify/internal/user"
	"Questify/pkg/jwt"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userOps                *user.Ops
	secret                 []byte
	tokenExpiration        uint
	refreshTokenExpiration uint
}

func NewAuthService(userOps *user.Ops, secret []byte,
	tokenExpiration uint, refreshTokenExpiration uint) *AuthService {
	return &AuthService{
		userOps:                userOps,
		secret:                 secret,
		tokenExpiration:        tokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}
}

type UserToken struct {
	AuthorizationToken string
	RefreshToken       string
	ExpiresAt          int64
}

func (s *AuthService) CreateUser(ctx context.Context, user *user.User) error {
	// Generate TFA Code
	tfaCode := fmt.Sprintf("%04d", rand.Intn(10000)) // Generate once
	user.TfaCode = tfaCode
	user.TfaExpiresAt = time.Now().Add(2 * time.Minute)

	// Save user to the database
	err := s.userOps.Create(ctx, user)
	if err != nil {
		return err
	}

	// Return the generated TFA code in the response
	return nil
}


func (s *AuthService) Login(ctx context.Context, email, pass string) (*UserToken, error) {
	fetchedUser, err := s.userOps.GetUserByEmailAndPassword(ctx, email, pass)
	if err != nil {
		return nil, err
	}

	// Calculate expiration time values
	var (
		authExp    = time.Now().Add(time.Minute * time.Duration(s.tokenExpiration))
		refreshExp = time.Now().Add(time.Minute * time.Duration(s.refreshTokenExpiration))
	)

	authToken, err := jwt.CreateToken(s.secret, s.userClaims(fetchedUser, authExp))
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.CreateToken(s.secret, s.userClaims(fetchedUser, refreshExp))
	if err != nil {
		return nil, err
	}

	return &UserToken{
		AuthorizationToken: authToken,
		RefreshToken:       refreshToken,
		ExpiresAt:          authExp.Unix(),
	}, nil
}

func (s *AuthService) RefreshAuth(ctx context.Context, refreshToken string) (*UserToken, error) {
	claim, err := jwt.ParseToken(refreshToken, s.secret)
	if err != nil {
		return nil, err
	}

	u, err := s.userOps.GetUserByID(ctx, claim.UserID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, user.ErrUserNotFound
	}

	// Calculate expiration time values
	var (
		authExp = time.Now().Add(time.Minute * time.Duration(s.tokenExpiration))
	)

	authToken, err := jwt.CreateToken(s.secret, s.userClaims(u, authExp))
	if err != nil {
		return nil, err
	}

	return &UserToken{
		AuthorizationToken: authToken,
		RefreshToken:       refreshToken,
		ExpiresAt:          authExp.UnixMilli(),
	}, nil
}

func (s *AuthService) userClaims(user *user.User, exp time.Time) *jwt.UserClaims {
	return &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: &jwt2.NumericDate{
				Time: exp,
			},
		},
		UserID: user.ID,
		Role:   user.Role,
	}
}

func (a *AuthService) Generate2FA(ctx context.Context, email string) (string, error) {
	// Generate a random 4-digit code
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	// Create the 2FA JWT
	tfaToken, err := jwt.CreateTFAToken(a.secret, email, code, 2*time.Minute)
	if err != nil {
		return "", err
	}

	// Send the code via email
	err = a.SendTFAEmail(ctx, email, code)
	if err != nil {
		return "", err
	}

	return tfaToken, nil
}

func (a *AuthService) SendTFAEmail(ctx context.Context, email, code string) error {
	return a.userOps.Send2FACodeEmail(ctx, email, code)
}

func (a *AuthService) ConfirmTFA(ctx context.Context, email, code string) error {
	// Normalize the input code (e.g., trim spaces)
	code = strings.TrimSpace(code)

	// Retrieve the user from the database
	user, err := a.userOps.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// Validate if the user exists
	if user == nil {
		return fmt.Errorf("user not found with email: %s", email)
	}

	// Check if the TFA code matches
	if user.TfaCode != code {
		return fmt.Errorf("invalid TFA code provided for email: %s", email)
	}

	// Check if the TFA code is expired
	if time.Now().After(user.TfaExpiresAt) {
		return fmt.Errorf("TFA code expired for email: %s", email)
	}

	// Clear TFA code after successful confirmation
	user.TfaCode = ""
	user.TfaExpiresAt = time.Time{}
	if err := a.userOps.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user during TFA confirmation: %w", err)
	}

	return nil
}

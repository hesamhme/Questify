package jwt

import (
	"errors"
	"strings"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

const (
	UserClaimKey = "User-Claims"
)

func CreateToken(secret []byte, claims *UserClaims) (string, error) {
	claims.Role = "user"
	return jwt2.NewWithClaims(jwt2.SigningMethodHS512, claims).SignedString(secret)
}

func ParseToken(tokenString string, secret []byte) (*UserClaims, error) {
	if strings.Count(tokenString, ".") != 2 {
		return nil, errors.New("token contains an invalid number of segments")
	}
	token, err := jwt2.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt2.Token) (interface{}, error) {
		return secret, nil
	})

	var claim *UserClaims
	if token.Claims != nil {
		cc, ok := token.Claims.(*UserClaims)
		if ok {
			claim = cc
		}
	}

	if err != nil {
		return claim, err
	}

	if !token.Valid {
		return claim, errors.New("token is not valid")
	}

	return claim, nil
}

func CreateTFAToken(secret []byte, email, code string, duration time.Duration) (string, error) {
	claims := TFAClaims{
		Email:   email,
		TFACode: code,
		Expires: time.Now().Add(duration),
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(duration)),
		},
	}

	return jwt2.NewWithClaims(jwt2.SigningMethodHS512, claims).SignedString(secret)
}

func ParseTFAToken(tokenString string, secret []byte) (*TFAClaims, error) {
	token, err := jwt2.ParseWithClaims(tokenString, &TFAClaims{}, func(t *jwt2.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TFAClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired TFA token")
	}

	return claims, nil
}

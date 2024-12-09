package middlewares

import (
	"Questify/api/http/handlers"
	jw2 "Questify/pkg/jwt"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Auth(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" {
			fmt.Println("Authorization header missing")
			return handlers.SendError(c, errors.New("authorization header missing"), fiber.StatusUnauthorized)
		}

		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Invalid authorization token format")
			return handlers.SendError(c, errors.New("invalid authorization token format"), fiber.StatusUnauthorized)
		}

		pureToken := parts[1]
		claims, err := jw2.ParseToken(pureToken, secret)
		if err != nil {
			fmt.Printf("Token parsing failed: %s\n", err.Error())
			return handlers.SendError(c, err, fiber.StatusUnauthorized)
		}

		fmt.Printf("Token parsed successfully: %+v\n", claims)

		c.Locals(jw2.UserClaimKey, claims)
		fmt.Printf("Claims before calling next: %+v\n", c.Locals(jw2.UserClaimKey))
		return c.Next()
	}
}

func RoleChecker(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals(jw2.UserClaimKey).(*jw2.UserClaims)
		hasAccess := false
		for _, role := range roles {
			if claims.Role == role {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			return handlers.SendError(c, errors.New("you don't have access to this section"), fiber.StatusForbidden)
		}

		return c.Next()
	}
}

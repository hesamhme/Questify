package handlers

import (
	"errors"
	"fmt"
	presenter "github.com/hesamhme/Questify/api/http/handlers/presentor"
	"github.com/hesamhme/Questify/internal/user"
	"github.com/hesamhme/Questify/service"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req presenter.RegisterRequest

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		u := presenter.UserRegisterToUserDomain(&req)

		newUser, err := authService.CreateUser(c.Context(), u)
		if err != nil {
			if errors.Is(err, user.ErrInvalidEmail) || errors.Is(err, user.ErrInvalidPassword) {
				return presenter.BadRequest(c, err)
			}
			if errors.Is(err, user.ErrEmailAlreadyExists) {
				return presenter.Conflict(c, err)
			}

			return presenter.InternalServerError(c, err)
		}

		return presenter.Created(c, "user successfully registered", fiber.Map{
			"user_id": newUser.ID,
		})
	}
}


func LoginUser(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req presenter.UserLoginReq

		if err := c.BodyParser(&req); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err := BodyValidator(req)
		if err != nil {
			return presenter.BadRequest(c, err)
		}

		c.Cookie(&fiber.Cookie{
			Name:        "X-Session-ID",
			Value:       fmt.Sprint(time.Now().UnixNano()),
			HTTPOnly:    true,
			SessionOnly: true,
		})

		authToken, err := authService.Login(c.Context(), req.Email, req.Password)
		if err != nil {

			return presenter.BadRequest(c, err)
		}
		return SendUserToken(c, authToken)
	}
}


func RefreshToken(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refToken := c.GetReqHeaders()["Authorization"][0]
		if len(refToken) == 0 {
			return SendError(c, errors.New("token should be provided"), fiber.StatusBadRequest)
		}
		pureToken := strings.Split(refToken, " ")[1]
		authToken, err := authService.RefreshAuth(c.UserContext(), pureToken)
		if err != nil {

			return presenter.Unauthorized(c, err)
		}

		return SendUserToken(c, authToken)
	}
}
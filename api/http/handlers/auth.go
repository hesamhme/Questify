package handlers

import (
	"Questify/api/http/handlers/presenter"
	"Questify/internal/user"
	"Questify/service"
	"errors"
	"strings"

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

		u := presenter.RegisterRequestToUser(&req)

		err = authService.CreateUser(c.Context(), u)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		// Response with the same TFA code stored in the database
		return c.JSON(fiber.Map{
			"success": true,
			"message": "TFA code sent to email.",
			"data": fiber.Map{
				"message":  "User successfully registered. Please verify your TFA code.",
				"tfa_code": u.TfaCode, // Ensure this matches the database
			},
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

		// c.Cookie(&fiber.Cookie{
		// 	Name:        "X-Session-ID",
		// 	Value:       fmt.Sprint(time.Now().UnixNano()),
		// 	HTTPOnly:    true,
		// 	SessionOnly: true,
		// })

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

func ConfirmTFA(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		}

		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
		}

		err := authService.ConfirmTFA(c.Context(), req.Email, req.Code)
		if err != nil {
			if errors.Is(err, user.ErrInvalidTFA) {
				return presenter.BadRequest(c, err)
			}
			return presenter.InternalServerError(c, err)
		}

		return c.JSON(fiber.Map{"message": "TFA confirmed. Registration complete."})
	}
}

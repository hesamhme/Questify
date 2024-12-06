package handlers

import (
	"Questify/api/http/handlers/presenter"
	"Questify/pkg/jwt"
	"Questify/pkg/valuecontext"
	"Questify/service"
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const UserClaimKey = jwt.UserClaimKey

var (
	errWrongClaimType = errors.New("wrong claim type")
)

type ServiceFactory[T any] func(context.Context) T

func SendError(c *fiber.Ctx, err error, status int) error {
	if status == 0 {
		status = fiber.StatusInternalServerError
	}

	c.Locals(valuecontext.IsTxError, err)

	return c.Status(status).JSON(map[string]any{
		"error_msg": err.Error(),
	})
}

func SendUserToken(c *fiber.Ctx, authToken *service.UserToken) error {

	return presenter.OK(c, "User successfully logged in", fiber.Map{
		"auth_token":    authToken.AuthorizationToken,
		"refresh_token": authToken.RefreshToken,
	})
}

func PageAndPageSize(c *fiber.Ctx) (int, int) {
	page, pageSize := c.QueryInt("page"), c.QueryInt("page_size")
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 20
	}

	return page, pageSize
}

func BodyValidator[T any](req T) error {
	myValidator := presenter.GetValidator()
	if errs := myValidator.Validate(req); len(errs) > 0 {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, err.Error)
		}

		return errors.New(strings.Join(errMsgs, "and"))
	}
	return nil
}

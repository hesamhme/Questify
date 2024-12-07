package handlers

import (
	"Questify/api/http/handlers/presenter"
	"Questify/pkg/jwt"
	"Questify/service"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func GetAllVerifiedUsers(userService *service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals(UserClaimKey).(*jwt.UserClaims)
		if !ok {
			return SendError(c, errWrongClaimType, fiber.StatusBadRequest)
		}
		//query parameter
		page, pageSize := PageAndPageSize(c)

		users, total, err := userService.GetAllUsers(c.UserContext(), userClaims.UserID, page, pageSize)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, service.ErrUnAuthorizedUser) {
				status = fiber.StatusUnauthorized
			}
			return SendError(c, err, status)
		}
		data := presenter.NewPagination(
			presenter.BatchUsersToUserGet(users),
			uint(page),
			uint(pageSize),
			uint(total),
		)
		return presenter.OK(c, "users successfully fetched.", data)
	}
}

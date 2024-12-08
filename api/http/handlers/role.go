package handlers

import (
	"Questify/api/http/handlers/presenter"
	"Questify/internal/role"
	"Questify/service"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AssignRoleToSurveyUser assigns a role to a user for a specific survey.
func AssignRoleToSurveyUser(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			UserID  string `json:"user_id"`
			RoleID  string `json:"role_id"`
			Timeout int    `json:"timeout"` // in minutes
		}

		surveyID := c.Params("surveyId")
		if surveyID == "" {
			return presenter.BadRequest(c, errors.New("survey ID is required"))
		}
		surveyUUID, err := uuid.Parse(surveyID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("invalid survey ID format"))
		}

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		userUUID, err := uuid.Parse(req.UserID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("invalid user ID format"))
		}

		roleUUID, err := uuid.Parse(req.RoleID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("invalid role ID format"))
		}

		var timeout *time.Duration
		if req.Timeout > 0 {
			duration := time.Duration(req.Timeout) * time.Minute
			timeout = &duration
		}

		// Validate if the user has the right permissions or is the survey owner
		claims := c.Locals("user_claims").(*presenter.UserClaims)
		if claims.UserID == uuid.Nil {
			return presenter.Unauthorized(c, errors.New("user not authenticated"))
		}

		isOwner, err := roleService.CheckSurveyPermission(c.Context(), surveyUUID, claims.UserID, role.PermissionIDManageRoles)
		if err != nil || !isOwner {
			return presenter.Forbidden(c, errors.New("user does not have the necessary permissions"))
		}

		err = roleService.AssignRoleToSurveyUser(c.Context(), surveyUUID, userUUID, roleUUID, timeout)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.Created(c, "Role assigned to survey user successfully", nil)
	}
}

// CheckSurveyPermission checks if a user has a specific permission for a survey.
func CheckSurveyPermission(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			UserID       string `json:"user_id"`
			PermissionID int    `json:"permission_id"`
		}

		surveyID := c.Params("surveyId")
		if surveyID == "" {
			return presenter.BadRequest(c, errors.New("survey ID is required"))
		}
		surveyUUID, err := uuid.Parse(surveyID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("invalid survey ID format"))
		}

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		userUUID, err := uuid.Parse(req.UserID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("invalid user ID format"))
		}

		hasPermission, err := roleService.CheckSurveyPermission(c.Context(), surveyUUID, userUUID, req.PermissionID)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Permission check completed", fiber.Map{
			"has_permission": hasPermission,
		})
	}
}

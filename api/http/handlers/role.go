package handlers

import (
	"Questify/api/http/handlers/presenter"
	"Questify/internal/survey"
	"Questify/internal/user"
	jw2 "Questify/pkg/jwt"
	"Questify/service"
	"errors"
	"fmt"
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
			Timeout int    `json:"timeout"`
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

		claims, ok := c.Locals(UserClaimKey).(*jw2.UserClaims)
		if !ok || claims == nil {
			fmt.Println("User claims are not found in context or are invalid")
			return presenter.Unauthorized(c, errors.New("user not authenticated"))
		}
		fmt.Printf("User claims retrieved successfully: %+v\n", claims)

		isOwner, err := roleService.CheckSurveyPermission(c.Context(), surveyUUID, claims.UserID, 0)
		if errors.Is(err, user.ErrUserNotFound) {
			return presenter.BadRequest(c, err)
		}
		if errors.Is(err, survey.ErrSurveyNotFound) {
			return presenter.BadRequest(c, err)
		}
		if err != nil || !isOwner || errors.Is(err, service.ErrNotOwner) {
			return presenter.Forbidden(c, err)
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


// CreateRole creates a new role with specified permissions.
func CreateRole(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Name         string `json:"name"`
			PermissionIDs []int `json:"permissions"`
		}

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		if req.Name == "" || len(req.PermissionIDs) == 0 {
			return presenter.BadRequest(c, errors.New("invalid role data"))
		}

		role, err := roleService.CreateRole(c.Context(), req.Name, req.PermissionIDs)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.Created(c, "Role created successfully", role)
	}
}

// GetAllRoles retrieves all roles in the system.
func GetAllRoles(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, err := roleService.GetAllRoles(c.Context())
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Roles retrieved successfully", roles)
	}
}

func DeleteRoles(roleService *service.RoleService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var req struct {
            RoleIDs []string `json:"role_ids"`
        }

        // Parse the request body
        if err := c.BodyParser(&req); err != nil {
            return presenter.BadRequest(c, errors.New("invalid request body"))
        }

        // Validate request body
        if len(req.RoleIDs) == 0 {
            return presenter.BadRequest(c, errors.New("role IDs are required"))
        }

        // Get user claims
        claims, ok := c.Locals(UserClaimKey).(*jw2.UserClaims)
        if !ok || claims == nil {
            return presenter.Unauthorized(c, errors.New("user not authenticated"))
        }

        // Iterate over RoleIDs and delete them if the user is the creator
        for _, roleID := range req.RoleIDs {
            roleUUID, err := uuid.Parse(roleID)
            if err != nil {
                return presenter.BadRequest(c, fmt.Errorf("invalid role ID format for: %s", roleID))
            }

            // Fetch the role to verify the creator
            _, err = roleService.GetRoleByID(c.Context(), roleUUID)
            if err != nil {
                return presenter.BadRequest(c, fmt.Errorf("role not found for ID: %s", roleID))
            }


            // Delete the role
            err = roleService.DeleteRole(c.Context(), roleUUID)
            if err != nil {
                return presenter.InternalServerError(c, fmt.Errorf("failed to delete role: %s", roleID))
            }
        }

        return presenter.OK(c, "Roles deleted successfully", nil)
    }
}


package handlers

import (
	"Questify/api/http/handlers/presenter"
	"Questify/service"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

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

// AssignRoleToUser assigns a role to a user.
func AssignRoleToUser(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			UserID string `json:"user_id"`
			RoleID string `json:"role_id"`
			Timeout int   `json:"timeout"` // in minutes
		}

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("invalid user ID"))
		}

		roleID, err := uuid.Parse(req.RoleID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("invalid role ID"))
		}

		var timeout *time.Duration
		if req.Timeout > 0 {
			duration := time.Duration(req.Timeout) * time.Minute
			timeout = &duration
		}

		err = roleService.AssignRoleToUser(c.Context(), userID, roleID, timeout)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.Created(c, "Role assigned to user successfully", nil)
	}
}

// CheckUserPermission checks if a user has a specific permission.
func CheckUserPermission(roleService *service.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			UserID       string `json:"user_id"`
			PermissionID int    `json:"permission_id"`
		}

		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, err)
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("invalid user ID"))
		}

		hasPermission, err := roleService.CheckPermission(c.Context(), userID, req.PermissionID)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.OK(c, "Permission check completed", fiber.Map{
			"has_permission": hasPermission,
		})
	}
}

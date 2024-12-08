package service

import (
	"Questify/internal/role"
	"context"
	"time"

	"github.com/google/uuid"
)

// RoleService provides methods to manage roles and permissions.
type RoleService struct {
	roleOps *role.Ops
}

// NewRoleService creates a new RoleService.
func NewRoleService(roleOps *role.Ops) *RoleService {
	return &RoleService{
		roleOps: roleOps,
	}
}

// CreateRole creates a new role with specified permissions.
func (s *RoleService) CreateRole(ctx context.Context, name string, permissionIDs []int) (*role.Role, error) {
	return s.roleOps.CreateRole(ctx, name, permissionIDs)
}

// AssignRoleToUser assigns a role to a user with an optional timeout.
func (s *RoleService) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, timeout *time.Duration) error {
	return s.roleOps.AssignRoleToUser(ctx, userID, roleID, timeout)
}

// GetRolesByUser retrieves all roles assigned to a user.
func (s *RoleService) GetRolesByUser(ctx context.Context, userID uuid.UUID) ([]role.Role, error) {
	return s.roleOps.GetRolesByUser(ctx, userID)
}

// RemoveRoleFromUser removes a role assignment from a user.
func (s *RoleService) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	return s.roleOps.RemoveRoleFromUser(ctx, userID, roleID)
}

// CheckPermission checks if a user has a specific permission.
func (s *RoleService) CheckPermission(ctx context.Context, userID uuid.UUID, permissionID int) (bool, error) {
	return s.roleOps.CheckPermission(ctx, userID, permissionID)
}

// GetAllRoles fetches all roles in the system.
func (s *RoleService) GetAllRoles(ctx context.Context) ([]role.Role, error) {
	return s.roleOps.GetAllRoles(ctx)
}

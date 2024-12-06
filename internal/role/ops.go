package role

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repository
}

// NewOps creates a new role operation handler.
func NewOps(repo Repository) *Ops {
	return &Ops{repo: repo}
}

func (o *Ops) CreateRole(ctx context.Context, name string, permissionIDs []int) (*Role, error) {
	// Create a new role instance
	role := &Role{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	// Fetch all permissions for the given IDs
	permissions, err := o.repo.GetPermissionsByIDs(ctx, permissionIDs) // Updated to fetch permissions using int IDs
	if err != nil {
		return nil, ErrPermissionNotFound
	}

	// Attach the fetched permissions to the role
	role.Permissions = permissions

	// Save the role to the repository
	if err := o.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// AssignRoleToUser assigns a role to a user with an optional timeout.
func (o *Ops) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, timeout *time.Duration) error {
	userRole := &UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedAt: time.Now(),
	}

	if timeout != nil {
		expiration := time.Now().Add(*timeout)
		userRole.ExpiresAt = &expiration
	}

	return o.repo.AssignRoleToUser(ctx, userRole)
}

// GetRolesByUser retrieves all roles assigned to a user.
func (o *Ops) GetRolesByUser(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	return o.repo.GetRolesByUserID(ctx, userID)
}

// RemoveRoleFromUser removes a role assignment from a user.
func (o *Ops) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	return o.repo.RemoveUserRole(ctx, userID, roleID)
}

// CheckPermission checks if a user has a specific permission.
func (o *Ops) CheckPermission(ctx context.Context, userID uuid.UUID, permissionID int) (bool, error) {
	roles, err := o.repo.GetRolesByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		for _, perm := range role.Permissions {
			if perm.ID == permissionID {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetAllRoles fetches all roles from the repository.
func (o *Ops) GetAllRoles(ctx context.Context) ([]Role, error) {
	return o.repo.GetAllRoles(ctx)
}

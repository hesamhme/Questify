package role

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Errors
var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrUserRoleAssignment = errors.New("failed to assign role to user")
	ErrUserRoleNotFound   = errors.New("user role not found")
)

// Role represents a role in the system.
type Role struct {
	ID          uuid.UUID
	Name        string
	CreatedAt   time.Time
	Permissions []Permission
}

// Permission represents a specific permission in the system.
type Permission struct {
	ID          int
	Description string
}

// UserRole represents the association between a user and a role.
type UserRole struct {
	UserID     uuid.UUID
	RoleID     uuid.UUID
	AssignedAt time.Time
	ExpiresAt  *time.Time
}

// Repository defines methods for role and permission management.
type Repository interface {
    CreateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error
    GetRoleByID(ctx context.Context, roleID uuid.UUID) (*Role, error)
    GetRoleByName(ctx context.Context, name string) (*Role, error)
    GetAllRoles(ctx context.Context) ([]Role, error)
    GetPermissionsByIDs(ctx context.Context, ids []int) ([]Permission, error)
    AssignRoleToUser(ctx context.Context, userRole *UserRole) error
    GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]Role, error)
    RemoveUserRole(ctx context.Context, userID, roleID uuid.UUID) error
	AssignRoleToSurveyUser(ctx context.Context, surveyUserRole *SurveyUserRole) error
	GetRolesBySurveyAndUser(ctx context.Context, surveyID uuid.UUID, userID uuid.UUID) ([]Role, error)
}


// SurveyUserRole represents a user's role in a specific survey.
type SurveyUserRole struct {
	ID        uuid.UUID
	SurveyID  uuid.UUID
	UserID    uuid.UUID
	RoleID    uuid.UUID
	AssignedAt time.Time
	ExpiresAt *time.Time
}


// Permission IDs (constants for predefined permissions)
const (
	PermissionIDViewSurvey     = 1 // View Survey permission ID
	PermissionIDVote           = 2 // Vote permission ID
	PermissionIDEditSurvey     = 3 // Edit Survey permission ID
	PermissionIDManageRoles    = 4 // Manage Roles permission ID
	PermissionIDViewReports    = 5 // View Reports permission ID
)




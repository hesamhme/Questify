package service

import (
	"Questify/internal/role"
	"Questify/internal/survey"
	"Questify/internal/user"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotOwner = errors.New("you do not have permission to make roles, only owner of the survey can create roles")
)

// RoleService provides methods to manage roles and permissions.
type RoleService struct {
	roleOps   *role.Ops
	surveyOps *survey.Ops
	userOps   *user.Ops
}

// NewRoleService creates a new RoleService.
func NewRoleService(roleOps *role.Ops, surveyOps *survey.Ops, userOps *user.Ops) *RoleService {
	return &RoleService{
		roleOps:   roleOps,
		surveyOps: surveyOps,
		userOps:   userOps,
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

// AssignRoleToSurveyUser assigns a role to a user for a specific survey.
func (s *RoleService) AssignRoleToSurveyUser(ctx context.Context, surveyID uuid.UUID, userID uuid.UUID, roleID uuid.UUID, timeout *time.Duration) error {
	// TODO: Add logic to validate if the current user is the survey owner or has "role permissions"
	return s.roleOps.AssignRoleToSurveyUser(ctx, surveyID, userID, roleID, timeout)
}

// GetSurveyRolesByUser retrieves roles assigned to a user for a specific survey.
func (s *RoleService) GetSurveyRolesByUser(ctx context.Context, surveyID uuid.UUID, userID uuid.UUID) ([]role.Role, error) {
	return s.roleOps.GetRolesBySurveyAndUser(ctx, surveyID, userID)
}

// CheckSurveyPermission checks if a user has a specific permission for a survey.
func (s *RoleService) CheckSurveyPermission(ctx context.Context, surveyID uuid.UUID, userID uuid.UUID, permissionID int) (bool, error) {
	if permissionID == 0 {
		_, err := s.userOps.GetUserByID(ctx, userID)
		if err != nil {
			return false, err
		}
		su, err := s.surveyOps.GetByID(ctx, surveyID)

		if err != nil {
			return false, err
		}

		if su.OwnerID != userID {
			return false, ErrNotOwner
		}
		return true, nil
	} else {
		return s.roleOps.CheckSurveyPermission(ctx, surveyID, userID, permissionID)
	}
}

// CheckSurveyPermission checks if a user has a specific permission for a survey.
func (s *RoleService) GetRolesBySurveyAndUser(ctx context.Context, surveyID uuid.UUID, userID uuid.UUID) ([]role.Role, error) {
	return s.roleOps.GetRolesBySurveyAndUser(ctx, surveyID, userID)
}

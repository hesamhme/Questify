package storage

import (
	"Questify/internal/role"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/adapters/storage/mappers"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roleRepo struct {
	db *gorm.DB
}

// NewRoleRepo creates a new instance of roleRepo.
func NewRoleRepo(db *gorm.DB) role.Repository {
	return &roleRepo{db: db}
}

// CreateRole creates a new role in the database.
func (r *roleRepo) CreateRole(ctx context.Context, role *role.Role) error {
	roleEntity := mappers.RoleDomainToEntity(*role)
	if err := r.db.WithContext(ctx).Create(&roleEntity).Error; err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	role.ID = roleEntity.ID
	return nil
}

// GetPermissionsByIDs retrieves permissions for given permission IDs.
func (r *roleRepo) GetPermissionsByIDs(ctx context.Context, ids []int) ([]role.Permission, error) {
	var permissionEntities []entities.Permission
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&permissionEntities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve permissions: %w", err)
	}
	return mappers.BatchPermissionEntityToDomain(permissionEntities), nil
}

// AssignRoleToUser assigns a role to a user.
func (r *roleRepo) AssignRoleToUser(ctx context.Context, userRole *role.UserRole) error {
	userRoleEntity := mappers.UserRoleDomainToEntity(*userRole)
	if err := r.db.WithContext(ctx).Create(&userRoleEntity).Error; err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}
	return nil
}

// GetRolesByUserID retrieves all roles assigned to a user.
func (r *roleRepo) GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]role.Role, error) {
	var userRoleEntities []entities.UserRole
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Role.Permissions").Find(&userRoleEntities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user roles: %w", err)
	}

	roles := make([]role.Role, 0, len(userRoleEntities))
	for _, userRole := range userRoleEntities {
		roles = append(roles, mappers.RoleEntityToDomain(userRole.Role))
	}

	return roles, nil
}

// RemoveUserRole removes a role assignment from a user.
func (r *roleRepo) RemoveUserRole(ctx context.Context, userID, roleID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&entities.UserRole{}).Error
	if err != nil {
		return fmt.Errorf("failed to remove user role: %w", err)
	}
	return nil
}
func (r *roleRepo) GetRoleByName(ctx context.Context, name string) (*role.Role, error) {
	var roleEntity entities.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).Preload("Permissions").First(&roleEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, role.ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}
	domainRole := mappers.RoleEntityToDomain(roleEntity)
	return &domainRole, nil
}

func (r *roleRepo) GetAllRoles(ctx context.Context) ([]role.Role, error) {
	var roleEntities []entities.Role

	// Fetch all roles along with their permissions using GORM's Preload
	err := r.db.WithContext(ctx).Preload("Permissions").Find(&roleEntities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}

	// Convert entities to domain models
	roles := mappers.BatchRoleEntityToDomain(roleEntities)

	return roles, nil
}

// GetRoleByID retrieves a role by its ID, including its permissions.
func (r *roleRepo) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*role.Role, error) {
	var roleEntity entities.Role
	err := r.db.WithContext(ctx).Preload("Permissions").Where("id = ?", roleID).First(&roleEntity).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve role by ID: %w", err)
	}
	domainRole := mappers.RoleEntityToDomain(roleEntity)
	return &domainRole, nil
}

// AssignRoleToSurveyUser assigns a role to a user for a specific survey.
func (r *roleRepo) AssignRoleToSurveyUser(ctx context.Context, surveyUserRole *role.SurveyUserRole) error {
	surveyUserRoleEntity := mappers.SurveyUserRoleDomainToEntity(*surveyUserRole)
	if err := r.db.WithContext(ctx).Create(&surveyUserRoleEntity).Error; err != nil {
		return fmt.Errorf("failed to assign role to survey user: %w", err)
	}
	return nil
}

// GetSurveyUserRoles retrieves all roles for a user within a specific survey.
func (r *roleRepo) GetSurveyUserRoles(ctx context.Context, surveyID, userID uuid.UUID) ([]role.Role, error) {
	var surveyUserRoleEntities []entities.SurveyUserRole
	err := r.db.WithContext(ctx).Where("survey_id = ? AND user_id = ?", surveyID, userID).Preload("Role.Permissions").Find(&surveyUserRoleEntities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve survey user roles: %w", err)
	}

	roles := make([]role.Role, 0, len(surveyUserRoleEntities))
	for _, surveyUserRole := range surveyUserRoleEntities {
		roles = append(roles, mappers.RoleEntityToDomain(surveyUserRole.Role))
	}

	return roles, nil
}

// GetRoleBySurveyAndUser retrieves roles assigned to a specific user for a specific survey.
func (r *roleRepo) GetRolesBySurveyAndUser(ctx context.Context, surveyID, userID uuid.UUID) ([]role.Role, error) {
	var surveyUserRoleEntities []entities.SurveyUserRole

	err := r.db.WithContext(ctx).
		Where("survey_id = ? AND user_id = ?", surveyID, userID).
		Preload("Role.Permissions").
		Find(&surveyUserRoleEntities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve roles for user %s in survey %s: %w", userID, surveyID, err)
	}

	roles := make([]role.Role, len(surveyUserRoleEntities))
	for i, surveyUserRole := range surveyUserRoleEntities {
		roles[i] = mappers.RoleEntityToDomain(surveyUserRole.Role)
	}

	return roles, nil
}

// DeleteRole removes a role by its ID.
func (r *roleRepo) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
    // Attempt to delete the role by ID
    result := r.db.WithContext(ctx).Where("id = ?", roleID).Delete(&entities.Role{})
    if result.Error != nil {
        return fmt.Errorf("failed to delete role: %w", result.Error)
    }

    // Check if a record was deleted
    if result.RowsAffected == 0 {
        return role.ErrRoleNotFound
    }

    return nil
}

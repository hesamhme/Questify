package mappers

import (
	"Questify/internal/role"
	"Questify/pkg/adapters/storage/entities"
)

// Map Role Entity to Domain
func RoleEntityToDomain(entity entities.Role) role.Role {
	return role.Role{
		ID:          entity.ID,             // UUID
		Name:        entity.Name,           // String
		CreatedAt:   entity.CreatedAt,      // Time
		Permissions: BatchPermissionEntityToDomain(entity.Permissions),
	}
}

// Map Role Domain to Entity
func RoleDomainToEntity(domain role.Role) entities.Role {
	return entities.Role{
		ID:          domain.ID,             // UUID
		Name:        domain.Name,           // String
		CreatedAt:   domain.CreatedAt,      // Time
		Permissions: BatchPermissionDomainToEntity(domain.Permissions),
	}
}

// Batch Role Entity to Domain
func BatchRoleEntityToDomain(entities []entities.Role) []role.Role {
	domains := make([]role.Role, len(entities))
	for i, entity := range entities {
		domains[i] = RoleEntityToDomain(entity)
	}
	return domains
}

// Batch Role Domain to Entity
func BatchRoleDomainToEntity(domains []role.Role) []entities.Role {
	entities := make([]entities.Role, len(domains))
	for i, domain := range domains {
		entities[i] = RoleDomainToEntity(domain)
	}
	return entities
}

// Map Permission Entity to Domain
func PermissionEntityToDomain(entity entities.Permission) role.Permission {
	return role.Permission{
		ID:          entity.ID,             // Integer
		Description: entity.Description,    // String
	}
}

// Map Permission Domain to Entity
func PermissionDomainToEntity(domain role.Permission) entities.Permission {
	return entities.Permission{
		ID:          domain.ID,             // Integer
		Description: domain.Description,    // String
	}
}

// Batch Permission Entity to Domain
func BatchPermissionEntityToDomain(entities []entities.Permission) []role.Permission {
	domains := make([]role.Permission, len(entities))
	for i, entity := range entities {
		domains[i] = PermissionEntityToDomain(entity)
	}
	return domains
}

// Batch Permission Domain to Entity
func BatchPermissionDomainToEntity(domains []role.Permission) []entities.Permission {
	entities := make([]entities.Permission, len(domains))
	for i, domain := range domains {
		entities[i] = PermissionDomainToEntity(domain)
	}
	return entities
}

// Map UserRole Entity to Domain
func UserRoleEntityToDomain(entity entities.UserRole) role.UserRole {
	return role.UserRole{
		UserID:     entity.UserID,    // UUID
		RoleID:     entity.RoleID,    // UUID
		AssignedAt: entity.AssignedAt, // Time
		ExpiresAt:  entity.ExpiresAt, // Nullable Time
	}
}

// Map UserRole Domain to Entity
func UserRoleDomainToEntity(domain role.UserRole) entities.UserRole {
	return entities.UserRole{
		UserID:     domain.UserID,    // UUID
		RoleID:     domain.RoleID,    // UUID
		AssignedAt: domain.AssignedAt, // Time
		ExpiresAt:  domain.ExpiresAt, // Nullable Time
	}
}

// Map SurveyUserRole Entity to Domain
func SurveyUserRoleEntityToDomain(entity entities.SurveyUserRole) role.SurveyUserRole {
	return role.SurveyUserRole{
		UserID:     entity.UserID,
		SurveyID:   entity.SurveyID,
		RoleID:     entity.RoleID,
		AssignedAt: entity.AssignedAt,
		ExpiresAt:  entity.ExpiresAt,
	}
}

// Map SurveyUserRole Domain to Entity
func SurveyUserRoleDomainToEntity(domain role.SurveyUserRole) entities.SurveyUserRole {
	return entities.SurveyUserRole{
		UserID:     domain.UserID,
		SurveyID:   domain.SurveyID,
		RoleID:     domain.RoleID,
		AssignedAt: domain.AssignedAt,
		ExpiresAt:  domain.ExpiresAt,
	}
}

package mappers

import (
	"Questify/internal/user"
	"Questify/pkg/adapters/storage/entities"
)

// ToDomainUser converts an entities.User to internal/user.User.
func ToDomainUser(entity *entities.User) *user.User {
	if entity == nil {
		return nil
	}

	return &user.User{
		ID:       entity.ID,
		Email:    entity.Email,
		Password: entity.Password,
		NID:      entity.NID,
	}
}

// ToEntityUser converts an internal/user.User to entities.User.
func ToEntityUser(domain *user.User) *entities.User {
	if domain == nil {
		return nil
	}

	return &entities.User{
		ID:       domain.ID,
		Email:    domain.Email,
		Password: domain.Password,
		NID:      domain.NID,
	}
}

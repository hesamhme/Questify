package mappers

import (
	"Questify/internal/user"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/fp"
)

func UserEntityToDomain(entity entities.User) user.User {
	return user.User{
		ID:           entity.ID,
		Email:        entity.Email,
		Password:     entity.Password,
		NationalCode: entity.NationalCode,
		TfaCode:      entity.TfaCode,
		TfaExpiresAt: entity.TfaExpiresAt,
		WalletID:     entity.WalletID,
	}
}

func BatchUserEntityToDomain(entities []entities.User) []user.User {
	return fp.Map(entities, UserEntityToDomain)
}

func UserDomainToEntity(domainUser *user.User) *entities.User {
	return &entities.User{
		Email:        domainUser.Email,
		Password:     domainUser.Password,
		NationalCode: domainUser.NationalCode,
		TfaCode:      domainUser.TfaCode,
		TfaExpiresAt: domainUser.TfaExpiresAt,
		WalletID:     domainUser.WalletID,
	}
}

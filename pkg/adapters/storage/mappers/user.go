package mappers

import (
	"github.com/hesamhme/Questify/internal/user"
	"github.com/hesamhme/Questify/pkg/adapters/storage/entities"
	"github.com/hesamhme/Questify/pkg/fp"
)

func UserEntityToDomain(entity entities.User) user.User{
	return user.User{
		ID:           entity.ID,
		Email:        entity.Email,
		Password:     entity.Password,
		NationalCode: entity.NationalCode,
	}
}


func BatchUserEntityToDomain(entities []entities.User) []user.User {
	return fp.Map(entities, UserEntityToDomain)
}

func UserDomainToEntity(domainUser *user.User) *entities.User {
	return &entities.User{
		Email:        domainUser.Email,
		Password:      domainUser.Password,
		NationalCode:  domainUser.NationalCode,
	}
}
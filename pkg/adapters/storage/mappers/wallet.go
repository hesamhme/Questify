package mappers

import (
	"Questify/internal/wallet"
	"Questify/pkg/adapters/storage/entities"
)

func WalletEntityToDomain(entity entities.Wallet) wallet.Wallet {
	return wallet.Wallet{
		ID:     entity.ID,
		Credit: entity.Credit,
	}
}

func WalletDomainToEntity(domain *wallet.Wallet) *entities.Wallet {
	return &entities.Wallet{
		ID:     domain.ID,
		Credit: domain.Credit,
	}
}

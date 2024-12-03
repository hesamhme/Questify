package storage

import (
	"Questify/internal/wallet"
	"Questify/pkg/adapters/storage/mappers"
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) wallet.Repo {
	return &walletRepo{
		db: db,
	}
}

func (r *walletRepo) Create(ctx context.Context, wallet *wallet.Wallet) error {
	wl := mappers.WalletDomainToEntity(wallet)
	err := r.db.WithContext(ctx).Create(&wl).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return err
	}
	wallet.ID = wl.ID
	return nil

}

func (r *walletRepo) Update(ctx context.Context, wallet *wallet.Wallet) error {
	return nil
}

func (r *walletRepo) GetById(ctx context.Context, id uuid.UUID) (*wallet.Wallet, error) {
	return nil, nil
}

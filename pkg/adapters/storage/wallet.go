package storage

import (
	"Questify/internal/wallet"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/adapters/storage/mappers"
	"context"
	"errors"
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
	w := mappers.WalletDomainToEntity(wallet)
	err := r.db.WithContext(ctx).Model(&entities.Wallet{}).Where("id = ?", wallet.ID).Updates(w).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *walletRepo) GetById(ctx context.Context, id uuid.UUID) (*wallet.Wallet, error) {
	var w entities.Wallet
	err := r.db.WithContext(ctx).Model(&entities.Wallet{}).Where("id = ?", id).First(&w).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	res := mappers.WalletEntityToDomain(w)
	return &res, nil
}

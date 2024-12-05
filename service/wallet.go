package service

import (
	"Questify/internal/wallet"
	"context"

	"github.com/google/uuid"
)

type WalletService struct {
	walletOps *wallet.Ops
}

func NewWalletService(walletOps *wallet.Ops) *WalletService {
	return &WalletService{walletOps: walletOps}
}

func (s *WalletService) CreateWallet(ctx context.Context, wallet *wallet.Wallet) (error, uuid.UUID) {
	err, id := s.walletOps.Create(ctx, wallet)
	if err != nil {
		return err, id
	}
	return nil, id
}

func (s *WalletService) Delete(ctx context.Context, id uuid.UUID) {
	s.walletOps.Delete(ctx, id)
}

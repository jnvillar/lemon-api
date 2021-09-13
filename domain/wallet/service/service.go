package walletservice

import (
	"context"

	walletmodel "lemonapp/domain/wallet/model"
)

type Service interface {
	CreateWallets(ctx context.Context, userID string) ([]walletmodel.Wallet, error)
	GetByID(ctx context.Context, userID string) (walletmodel.Wallet, error)
	GetUserWallets(ctx context.Context, userID string) ([]walletmodel.Wallet, error)
	GetUserWallet(ctx context.Context, userID string, walletID string) (walletmodel.Wallet, error)
	UpdateBalance(ctx context.Context, walletID string, balance int64) error
}

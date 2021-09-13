package walletrepository

import (
	"context"

	walletmodel "lemonapp/domain/wallet/model"
)

type Repository interface {
	CreateWallets(ctx context.Context, wallets []walletmodel.Wallet) ([]walletmodel.Wallet, error)
	GetUserWallets(ctx context.Context, walletID string) ([]walletmodel.Wallet, error)
	GetByID(ctx context.Context, walletID string) (walletmodel.Wallet, error)
	UpdateBalance(ctx context.Context, walletID string, amount int64) error
}

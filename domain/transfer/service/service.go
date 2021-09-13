package transferservice

import (
	"context"

	walletmodel "lemonapp/domain/wallet/model"
	trasnfermodel "lemonapp/domain/transfer/model"
)

type Service interface {
	Deposit(ctx context.Context, transfer *trasnfermodel.Transfer, to walletmodel.Wallet) error
	Extraction(ctx context.Context, transfer *trasnfermodel.Transfer, from walletmodel.Wallet) error
	Transfer(ctx context.Context, transfer *trasnfermodel.Transfer, from, to walletmodel.Wallet) error
}

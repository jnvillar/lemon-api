package transferrepository

import (
	"context"

	walletmodel "lemonapp/domain/wallet/model"
	transactionmodel "lemonapp/domain/transaction/model"
)

type Repository interface {
	Transfer(ctx context.Context,
		outgoingTransaction transactionmodel.Transaction,
		incomingTransaction transactionmodel.Transaction,
		from walletmodel.Wallet,
		to walletmodel.Wallet)  error

	Deposit(ctx context.Context, depositTransaction transactionmodel.Transaction, to walletmodel.Wallet) error
	Extraction(ctx context.Context, extractionTransaction transactionmodel.Transaction, from walletmodel.Wallet) error
}

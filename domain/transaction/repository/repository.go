package transactionrepository

import (
	"context"

	transactionmodel "lemonapp/domain/transaction/model"
)

type Repository interface {
	CreateTransactions(ctx context.Context, transactions []transactionmodel.Transaction) ([]transactionmodel.Transaction, error)
	Search(ctx context.Context, searchParams *transactionmodel.SearchParams) ([]transactionmodel.Transaction, error)
}

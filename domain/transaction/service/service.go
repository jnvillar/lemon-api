package transactionservice

import (
	"context"

	transactionmodel "lemonapp/domain/transaction/model"
)

type Service interface {
	Search(ctx context.Context, searchParams *transactionmodel.SearchParams) ([]transactionmodel.Transaction, error)
}

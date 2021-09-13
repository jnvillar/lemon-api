package transactionservice

import (
	"context"

	transactionmodel "lemonapp/domain/transaction/model"
	transactionrepository "lemonapp/domain/transaction/repository"
)

type ServiceImpl struct {
	repository transactionrepository.Repository
}

func (s *ServiceImpl) Search(ctx context.Context, searchParams *transactionmodel.SearchParams) ([]transactionmodel.Transaction, error) {
	return s.repository.Search(ctx, searchParams)
}

func NewServiceImpl(transactionRepository transactionrepository.Repository) Service {
	return &ServiceImpl{
		repository: transactionRepository,
	}
}

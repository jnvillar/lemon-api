package transferservice

import (
	"context"
	"fmt"

	walletmodel "lemonapp/domain/wallet/model"
	transactionmodel "lemonapp/domain/transaction/model"
	trasnfermodel "lemonapp/domain/transfer/model"
	transferrepository "lemonapp/domain/transfer/repository"
)

type ServiceImpl struct {
	repository transferrepository.Repository
}

func (s *ServiceImpl) Deposit(ctx context.Context, transfer *trasnfermodel.Transfer, to walletmodel.Wallet) error {
	deposit := s.getDepositTransaction(transfer, to)
	return s.repository.Deposit(ctx, deposit, to)
}

func (s *ServiceImpl) Extraction(ctx context.Context, transfer *trasnfermodel.Transfer, from walletmodel.Wallet) error {
	if from.GetBaseWallet().Balance < transfer.Amount {
		return fmt.Errorf("not enough balance")
	}
	extraction := s.getExtractionTransaction(transfer, from)
	return s.repository.Extraction(ctx, extraction, from)
}

func (s *ServiceImpl) Transfer(ctx context.Context, transfer *trasnfermodel.Transfer, from, to walletmodel.Wallet) error {
	currencyCheck := from.GetBaseWallet().Currency == to.GetBaseWallet().Currency && transfer.Currency == from.GetBaseWallet().Currency
	if !currencyCheck {
		return fmt.Errorf("currencies doesn't match")
	}

	if from.GetBaseWallet().Balance < transfer.Amount {
		return fmt.Errorf("not enough balance")
	}

	outgoingTransaction := s.getOutgoingTransaction(transfer, from, to)
	incomingTransaction := s.getIncomingTransaction(transfer, from, to)

	return s.repository.Transfer(ctx, outgoingTransaction, incomingTransaction, from, to)
}

func (s *ServiceImpl) getIncomingTransaction(transfer *trasnfermodel.Transfer, from, to walletmodel.Wallet) *transactionmodel.IncomingTransaction {
	return transactionmodel.NewIncomingTransaction(
		from.GetBaseWallet().UserID,
		to.GetBaseWallet().UserID,
		from.GetBaseWallet().ID,
		to.GetBaseWallet().ID,
		transfer.Amount,
		transfer.Currency,
	)
}

func (s *ServiceImpl) getOutgoingTransaction(transfer *trasnfermodel.Transfer, from, to walletmodel.Wallet) *transactionmodel.OutgoingTransaction {
	return transactionmodel.NewOutgoingTransaction(
		from.GetBaseWallet().UserID,
		to.GetBaseWallet().UserID,
		from.GetBaseWallet().ID,
		to.GetBaseWallet().ID,
		transfer.Amount,
		transfer.Currency,
	)
}

func (s *ServiceImpl) getExtractionTransaction(transfer *trasnfermodel.Transfer, from walletmodel.Wallet) *transactionmodel.ExtractionTransaction {
	return transactionmodel.NewExtractionTransaction(
		from.GetBaseWallet().UserID,
		from.GetBaseWallet().ID,
		transfer.Amount,
		from.GetBaseWallet().Currency,
	)
}

func (s *ServiceImpl) getDepositTransaction(transfer *trasnfermodel.Transfer, to walletmodel.Wallet) *transactionmodel.DepositTransaction {
	return transactionmodel.NewDepositTransaction(
		to.GetBaseWallet().UserID,
		to.GetBaseWallet().ID,
		transfer.Amount,
		to.GetBaseWallet().Currency,
	)
}

func NewServiceImpl(transferRepository transferrepository.Repository) Service {
	return &ServiceImpl{
		repository: transferRepository,
	}
}

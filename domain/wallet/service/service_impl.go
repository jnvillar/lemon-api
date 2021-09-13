package walletservice

import (
	"context"
	"fmt"

	walletmodel "lemonapp/domain/wallet/model"
	walletRepository "lemonapp/domain/wallet/repository"
	"lemonapp/errors"
)

type ServiceImpl struct {
	repository walletRepository.Repository
}

func NewServiceImpl(walletRepository walletRepository.Repository) Service {
	return &ServiceImpl{
		repository: walletRepository,
	}
}

func (s *ServiceImpl) UpdateBalance(ctx context.Context, walletID string, balance int64) error {
	return s.repository.UpdateBalance(ctx, walletID, balance)
}

func (s *ServiceImpl) GetUserWallet(ctx context.Context, userID string, walletID string) (walletmodel.Wallet, error) {
	wallets, err := s.GetUserWallets(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, wallet := range wallets {
		if wallet.GetBaseWallet().ID == walletID {
			return wallet, nil
		}
	}
	return nil, errors.NewNotFoundError(fmt.Errorf("wallet %v not found", walletID))
}

func (s *ServiceImpl) GetByID(ctx context.Context, walletID string) (walletmodel.Wallet, error) {
	return s.repository.GetByID(ctx, walletID)
}

func (s *ServiceImpl) CreateWallets(ctx context.Context, userID string) ([]walletmodel.Wallet, error) {
	return s.repository.CreateWallets(ctx, createInitialWallets(userID))
}

func createInitialWallets(userID string) []walletmodel.Wallet {
	return []walletmodel.Wallet{
		walletmodel.NewArsWallet(userID),
		walletmodel.NewUSDTWallet(userID),
		walletmodel.NewBTCWallet(userID),
	}
}

func (s *ServiceImpl) GetUserWallets(ctx context.Context, ID string) ([]walletmodel.Wallet, error) {
	return s.repository.GetUserWallets(ctx, ID)
}

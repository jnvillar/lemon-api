package transferrepository

import (
	"context"
	"database/sql"
	"fmt"

	transactionmodel "lemonapp/domain/transaction/model"
	transactionrepository "lemonapp/domain/transaction/repository"
	walletmodel "lemonapp/domain/wallet/model"
	walletRepository "lemonapp/domain/wallet/repository"
)

type repositoryImpl struct {
	walletRepository      walletRepository.Repository
	transactionRepository transactionrepository.Repository
	database              *sql.DB
}

func (r *repositoryImpl) Deposit(ctx context.Context, depositTransaction transactionmodel.Transaction, to walletmodel.Wallet) error {
	tx, err := r.database.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}

	if err := r.walletRepository.UpdateBalance(ctx, to.GetBaseWallet().ID, depositTransaction.GetBaseTransaction().Amount); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("could not rollback error %v", err)
		}
		return err
	}

	if _, err := r.transactionRepository.CreateTransactions(ctx, []transactionmodel.Transaction{depositTransaction}); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("could not rollback error %v", err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit changes: %v", err)
	}
	return nil
}

func (r *repositoryImpl) Extraction(ctx context.Context, extractionTransaction transactionmodel.Transaction, from walletmodel.Wallet) error {
	tx, err := r.database.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}

	if err := r.walletRepository.UpdateBalance(ctx, from.GetBaseWallet().ID, -extractionTransaction.GetBaseTransaction().Amount); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("could not rollback error %v", err)
		}
		return err
	}

	if _, err := r.transactionRepository.CreateTransactions(ctx, []transactionmodel.Transaction{extractionTransaction}); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("could not rollback error %v", err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit changes: %v", err)
	}
	return nil
}

func (r *repositoryImpl) Transfer(
	ctx context.Context,
	outgoingTransaction transactionmodel.Transaction,
	incomingTransaction transactionmodel.Transaction,
	from walletmodel.Wallet,
	to walletmodel.Wallet) error {

	tx, err := r.database.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}

	if err := r.walletRepository.UpdateBalance(ctx, from.GetBaseWallet().ID, -outgoingTransaction.GetBaseTransaction().Amount); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("could not rollback error %v", err)
		}
		return err
	}

	if err := r.walletRepository.UpdateBalance(ctx, to.GetBaseWallet().ID, incomingTransaction.GetBaseTransaction().Amount); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("could not rollback error %v", err)
		}
		return err
	}

	if _, err := r.transactionRepository.CreateTransactions(ctx, []transactionmodel.Transaction{outgoingTransaction, incomingTransaction}); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("could not rollback error %v", err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit changes: %v", err)
	}

	return nil
}

func NewRepository(
	db *sql.DB,
	walletRepository walletRepository.Repository,
	transactionRepository transactionrepository.Repository,
) Repository {
	return &repositoryImpl{
		database:              db,
		walletRepository:      walletRepository,
		transactionRepository: transactionRepository,
	}
}

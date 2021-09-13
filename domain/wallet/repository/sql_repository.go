package walletrepository

import (
	"context"
	"database/sql"
	"fmt"

	"lemonapp/currency"
	walletmodel "lemonapp/domain/wallet/model"
	"lemonapp/errors"
	"lemonapp/storage/mysql"

	sq "github.com/Masterminds/squirrel"
)

type sqlRepository struct {
	sqlBuilder sq.StatementBuilderType
	db         *sql.DB
}

const (
	walletTable    = "wallet"
	idColumn       = "id"
	currencyColumn = "currency"
	balanceColumn  = "balance"
	userIDColumn   = "user_id"
)

var walletColumns = []string{
	idColumn,
	currencyColumn,
	balanceColumn,
	userIDColumn,
}

func (r *sqlRepository) UpdateBalance(ctx context.Context, walletID string, amount int64) error {
	wallet, err := r.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	query, args, err := r.sqlBuilder.
		Update(walletTable).
		Set(balanceColumn, sq.Expr("balance + ?", amount)).
		Where(sq.Eq{idColumn: wallet.GetBaseWallet().ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error creating wallet query: %v", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating wallet balance %v", err)
	}

	return nil
}

func (r *sqlRepository) GetByID(ctx context.Context, walletID string) (walletmodel.Wallet, error) {
	query, args, err := r.sqlBuilder.
		Select(walletColumns...).
		From(walletTable).
		Where(sq.Eq{idColumn: walletID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating wallet query: %v", err)
	}
	row := r.db.QueryRowContext(ctx, query, args...)
	return r.scanWallet(row)
}

func (r *sqlRepository) scanWallet(row mysql.Scannable) (walletmodel.Wallet, error) {
	baseWallet := &walletmodel.BaseWallet{}

	scanArgs := []interface{}{
		&baseWallet.ID,
		&baseWallet.Currency,
		&baseWallet.Balance,
		&baseWallet.UserID,
	}

	if err := row.Scan(scanArgs...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(fmt.Errorf("wallet not found: %v", err))
		}
		return nil, fmt.Errorf(fmt.Sprintf("error scanning wallet: %v", err))
	}

	switch baseWallet.Currency {
	case currency.ARSCurrency:
		return &walletmodel.ARSWallet{BaseWallet: baseWallet}, nil
	case currency.BTCCurrency:
		return &walletmodel.BTCWallets{BaseWallet: baseWallet}, nil
	case currency.USDTCurrency:
		return &walletmodel.USDTWallets{BaseWallet: baseWallet}, nil
	default:
		return nil, fmt.Errorf("invalid wallet")
	}
}

func (r *sqlRepository) CreateWallets(ctx context.Context, wallets []walletmodel.Wallet) ([]walletmodel.Wallet, error) {
	builder := r.sqlBuilder.Insert(walletTable).Columns(walletColumns...)
	for _, wallet := range wallets {
		builder = builder.Values(
			wallet.GetBaseWallet().ID,
			wallet.GetBaseWallet().Currency,
			wallet.GetBaseWallet().Balance,
			wallet.GetBaseWallet().UserID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building create wallets query")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error creating wallets: %v", err)
	}

	return wallets, nil
}

func (r *sqlRepository) GetUserWallets(ctx context.Context, ID string) ([]walletmodel.Wallet, error) {
	builder := r.sqlBuilder.
		Select(walletColumns...).
		From(walletTable).
		Where(sq.Eq{userIDColumn: ID}).
		OrderBy(fmt.Sprintf("%s DESC", currencyColumn))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list users query")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing user wallets %v", err)
	}

	wallets := make([]walletmodel.Wallet, 0)
	for rows.Next() {
		wallet, err := r.scanWallet(rows)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func NewRepository(db *sql.DB) Repository {
	return &sqlRepository{
		db:         db,
		sqlBuilder: sq.StatementBuilder,
	}
}

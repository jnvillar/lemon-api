package transactionrepository

import (
	"context"
	"database/sql"
	"fmt"

	transactionmodel "lemonapp/domain/transaction/model"
	"lemonapp/errors"
	"lemonapp/storage/mysql"

	sq "github.com/Masterminds/squirrel"
)

type sqlRepository struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
}

const (
	defaultPageSizeLimit = 10
	maxPageSizeLimit     = 100
)

const (
	transactionTable  = "transaction"
	idColumn          = "id"
	currencyColumn    = "currency"
	amountColumn      = "amount"
	typeColumn        = "type"
	userIDColumn      = "user_id"
	userFromColumn    = "user_from"
	userToColumn      = "user_to"
	walletIDColumn    = "wallet_id"
	walletFromColumn  = "wallet_from"
	walletToColumn    = "wallet_to"
	dateCreatedColumn = "date_created"
)

var transactionColumns = []string{
	idColumn,
	currencyColumn,
	amountColumn,
	typeColumn,
	userIDColumn,
	userFromColumn,
	userToColumn,
	walletIDColumn,
	walletFromColumn,
	walletToColumn,
	dateCreatedColumn,
}

func (r *sqlRepository) GetWalletTransactions(ctx context.Context, walletID string) ([]transactionmodel.Transaction, error) {
	builder := r.sqlBuilder.
		Select(transactionColumns...).
		From(transactionTable).
		Where(sq.Eq{walletIDColumn: walletID}).
		OrderBy(fmt.Sprintf("%s DESC", dateCreatedColumn))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list transactions query")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing transactions: %v", err)
	}

	transactions := make([]transactionmodel.Transaction, 0)
	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *sqlRepository) CreateTransactions(ctx context.Context, transactions []transactionmodel.Transaction) ([]transactionmodel.Transaction, error) {
	builder := r.sqlBuilder.Insert(transactionTable).Columns(transactionColumns...)
	for _, transaction := range transactions {
		builder = builder.Values(
			transaction.GetBaseTransaction().ID,
			transaction.GetBaseTransaction().Currency,
			transaction.GetBaseTransaction().Amount,
			transaction.GetType(),
			transaction.GetBaseTransaction().UserID,
			transaction.GetBaseTransaction().UserFrom,
			transaction.GetBaseTransaction().UserTo,
			transaction.GetBaseTransaction().WalletID,
			transaction.GetBaseTransaction().WalletFrom,
			transaction.GetBaseTransaction().WalletTo,
			transaction.GetBaseTransaction().DateCreated,
		)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building create transactions query")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error creating transactions: %v", err)
	}

	return transactions, nil
}

func (r *sqlRepository) getLimit(limit int) int {
	if limit <= 0 {
		return defaultPageSizeLimit
	}
	if limit > maxPageSizeLimit {
		return maxPageSizeLimit
	}
	return limit
}

func (r *sqlRepository) getOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}

func (r *sqlRepository) Search(ctx context.Context, searchParams *transactionmodel.SearchParams) ([]transactionmodel.Transaction, error) {

	builder := r.sqlBuilder.
		Select(transactionColumns...).
		From(transactionTable).
		Limit(uint64(r.getLimit(searchParams.Limit))).
		Offset(uint64(r.getOffset(searchParams.Offset))).
		OrderBy(fmt.Sprintf("%s DESC", dateCreatedColumn))

	if searchParams.UserID != "" {
		builder = builder.Where(sq.Eq{userIDColumn: searchParams.UserID})
	}

	if searchParams.WalletID != "" {
		builder = builder.Where(sq.Eq{walletIDColumn: searchParams.WalletID})
	}

	if searchParams.TransactionType != "" {
		builder = builder.Where(sq.Eq{typeColumn: searchParams.TransactionType})
	}

	if searchParams.Currency != "" {
		builder = builder.Where(sq.Eq{currencyColumn: searchParams.Currency})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building list transactions query")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing transactions: %v", err)
	}

	transactions := make([]transactionmodel.Transaction, 0)
	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *sqlRepository) scanTransaction(row mysql.Scannable) (transactionmodel.Transaction, error) {
	baseTransaction := &transactionmodel.BaseTransaction{}

	scanArgs := []interface{}{
		&baseTransaction.ID,
		&baseTransaction.Currency,
		&baseTransaction.Amount,
		&baseTransaction.Type,
		&baseTransaction.UserID,
		&baseTransaction.UserFrom,
		&baseTransaction.UserTo,
		&baseTransaction.WalletID,
		&baseTransaction.WalletFrom,
		&baseTransaction.WalletTo,
		&baseTransaction.DateCreated,
	}

	if err := row.Scan(scanArgs...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError(fmt.Errorf("wallet not found: %v", err))
		}
		return nil, fmt.Errorf(fmt.Sprintf("error scanning wallet: %v", err))
	}

	switch baseTransaction.Type {
	case transactionmodel.OutgoingTransactionType:
		return &transactionmodel.OutgoingTransaction{
			BaseTransaction: baseTransaction,
		}, nil
	case transactionmodel.IncomingTransactionType:
		return &transactionmodel.IncomingTransaction{
			BaseTransaction: baseTransaction,
		}, nil
	case transactionmodel.DepositTransactionType:
		return &transactionmodel.DepositTransaction{
			BaseTransaction: baseTransaction,
		}, nil
	case transactionmodel.ExtractionTransactionType:
		return &transactionmodel.ExtractionTransaction{
			BaseTransaction: baseTransaction,
		}, nil
	default:
		return nil, fmt.Errorf("invalid transaction")
	}
}

func NewRepository(db *sql.DB) Repository {
	return &sqlRepository{
		db:         db,
		sqlBuilder: sq.StatementBuilder,
	}
}

package transactionmodel

import (
	"encoding/json"
	"time"

	"lemonapp/currency"

	"github.com/google/uuid"
)

type TransactionType = string

const (
	OutgoingTransactionType   TransactionType = "outgoing"
	IncomingTransactionType   TransactionType = "incoming"
	DepositTransactionType    TransactionType = "deposit"
	ExtractionTransactionType TransactionType = "extraction"
)

var transactionsTypes = []TransactionType{
	OutgoingTransactionType,
	IncomingTransactionType,
	DepositTransactionType,
	ExtractionTransactionType,
}

func ValidTransaction(transactionType TransactionType) bool {
	for _, t := range transactionsTypes {
		if transactionType == t {
			return true
		}
	}
	return false
}

type Transaction interface {
	GetBaseTransaction() *BaseTransaction
	GetType() TransactionType
}

type BaseTransaction struct {
	ID          string            `json:"id"`
	Currency    currency.Currency `json:"currency"`
	Amount      int64             `json:"-"`
	UserID      string            `json:"user_id"`
	UserFrom    string            `json:"user_from"`
	UserTo      string            `json:"user_to"`
	WalletID    string            `json:"wallet_id"`
	WalletTo    string            `json:"wallet_to"`
	WalletFrom  string            `json:"wallet_from"`
	Type        string            `json:"type"`
	DateCreated time.Time         `json:"date_created"`
}

func NewBaseTransaction(userID, userFrom, userTo, walletID, walletFrom, walletTo string, amount int64, currency currency.Currency) *BaseTransaction {
	return &BaseTransaction{
		ID:          uuid.NewString(),
		Currency:    currency,
		Amount:      amount,
		UserID:      userID,
		UserFrom:    userFrom,
		UserTo:      userTo,
		WalletID:    walletID,
		WalletTo:    walletTo,
		WalletFrom:  walletFrom,
		Type:        IncomingTransactionType,
		DateCreated: time.Now().UTC(),
	}
}

func (b *BaseTransaction) GetBaseTransaction() *BaseTransaction {
	return b
}

func (b BaseTransaction) MarshalJSON() ([]byte, error) {
	type Alias BaseTransaction
	t := struct {
		Alias
		Amount string `json:"amount"`
	}{
		Alias:  Alias(b),
		Amount: currency.FormatCurrency(b.Amount, b.Currency),
	}
	return json.Marshal(t)
}

type IncomingTransaction struct {
	*BaseTransaction
}

func NewIncomingTransaction(userFrom, userTo, walletFrom, walletTo string, amount int64, currency currency.Currency) *IncomingTransaction {
	return &IncomingTransaction{NewBaseTransaction(userTo, userFrom, userTo, walletTo, walletFrom, walletTo, amount, currency)}
}

func (t *IncomingTransaction) GetType() TransactionType {
	return IncomingTransactionType
}

type OutgoingTransaction struct {
	*BaseTransaction
}

func NewOutgoingTransaction(userFrom, userTo, walletFrom, walletTo string, amount int64, currency currency.Currency) *OutgoingTransaction {
	return &OutgoingTransaction{
		NewBaseTransaction(userFrom, userFrom, userTo, walletFrom, walletFrom, walletTo, amount, currency)}
}

func (t *OutgoingTransaction) GetType() TransactionType {
	return OutgoingTransactionType
}

type DepositTransaction struct {
	*BaseTransaction
}

func (t *DepositTransaction) GetType() TransactionType {
	return DepositTransactionType
}

func NewDepositTransaction(userID, walletID string, amount int64, currency currency.Currency) *DepositTransaction {
	return &DepositTransaction{
		NewBaseTransaction(userID, "", userID, walletID, "", walletID, amount, currency)}
}

type ExtractionTransaction struct {
	*BaseTransaction
}

func (t *ExtractionTransaction) GetType() TransactionType {
	return ExtractionTransactionType
}

func NewExtractionTransaction(userID, walletID string, amount int64, currency currency.Currency) *ExtractionTransaction {
	return &ExtractionTransaction{
		NewBaseTransaction(userID, userID, "", walletID, walletID, "", amount, currency)}
}

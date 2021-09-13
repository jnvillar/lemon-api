package transactionmodel

import "lemonapp/currency"

type SearchParams struct {
	UserID          string
	WalletID        string
	Limit           int
	Offset          int
	TransactionType TransactionType
	Currency        currency.Currency
}



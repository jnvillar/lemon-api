package trasnfermodel

import "lemonapp/currency"

type Transfer struct {
	Amount     int64             `json:"amount"`
	WalletFrom string            `json:"wallet_from"`
	WalletTo   string            `json:"wallet_to"`
	Currency   currency.Currency `json:"currency"`
}

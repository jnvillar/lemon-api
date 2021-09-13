package walletmodel

import (
	"encoding/json"

	"lemonapp/currency"

	"github.com/google/uuid"
)

type Wallet interface {
	GetBaseWallet() *BaseWallet
}

type BaseWallet struct {
	ID               string            `json:"id"`
	Currency         currency.Currency `json:"currency"`
	Balance          int64             `json:"-"`
	FormattedBalance string            `json:"balance"`
	UserID           string            `json:"user_id"`
}

func (b BaseWallet) MarshalJSON() ([]byte, error) {
	type Alias BaseWallet
	t := struct {
		Alias
		Balance string `json:"balance"`
	}{
		Alias:   Alias(b),
		Balance: currency.FormatCurrency(b.Balance, b.Currency),
	}
	return json.Marshal(t)
}

func (b *BaseWallet) GetBaseWallet() *BaseWallet {
	return b
}

type ARSWallet struct {
	*BaseWallet
}

func NewArsWallet(userID string) *ARSWallet {
	return &ARSWallet{
		BaseWallet: &BaseWallet{ID: uuid.NewString(), Currency: currency.ARSCurrency, UserID: userID},
	}
}

type BTCWallets struct {
	*BaseWallet
}

func NewBTCWallet(userID string) *BTCWallets {
	return &BTCWallets{
		BaseWallet: &BaseWallet{ID: uuid.NewString(), Currency: currency.BTCCurrency, UserID: userID},
	}
}

type USDTWallets struct {
	*BaseWallet
}

func NewUSDTWallet(userID string) *USDTWallets {
	return &USDTWallets{
		BaseWallet: &BaseWallet{ID: uuid.NewString(), Currency: currency.USDTCurrency, UserID: userID},
	}
}

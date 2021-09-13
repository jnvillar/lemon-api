package currency

import "fmt"

type Currency = string

const (
	ARSCurrency  Currency = "ARS"
	USDTCurrency Currency = "USDT"
	BTCCurrency  Currency = "BTC"
)

var currencies = []Currency{ARSCurrency, USDTCurrency, BTCCurrency}

func ValidCurrency(currency Currency) bool {
	for _, c := range currencies {
		if currency == c {
			return true
		}
	}
	return false
}

func FormatCurrency(amount int64, currency Currency) string {
	switch currency {
	case BTCCurrency:
		return fmt.Sprintf("%.8f", float64(amount/100000000))
	default:
		return fmt.Sprintf("%.2f", float64(amount/100))
	}
}

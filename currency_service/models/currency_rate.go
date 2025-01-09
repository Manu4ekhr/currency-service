package models

type CurrencyRate struct {
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
}

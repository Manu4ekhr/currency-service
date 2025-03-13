package models

import "gorm.io/gorm"

type CurrencyRate struct {
	gorm.Model
	Currency string  `json:"currency" gorm:"uniqueIndex"` // Делаем валюту уникальной
	Rate     float64 `json:"rate"`
}

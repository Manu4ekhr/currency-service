package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ABSAPIURL      string
	BankAPIURL     string
	UpdateInterval time.Duration
	DatabaseDSN    string
}

func GetConfig() Config {
	updateInterval, err := strconv.ParseInt(strings.TrimSpace(os.Getenv("UPDATE_INTERVAL")), 10, 64)
	if err != nil || updateInterval <= 0 {
		updateInterval = 600 // 10 минут по умолчанию
	}

	return Config{
		ABSAPIURL:      strings.TrimSpace(getEnv("ABS_API_URL", "https://default-abs-api.com")),
		BankAPIURL:     strings.TrimSpace(getEnv("BANK_API_URL", "https://default-bank-api.com")),
		UpdateInterval: time.Duration(updateInterval) * time.Second,
		DatabaseDSN:    strings.TrimSpace(getEnv("DATABASE_DSN", "currency.db")), // SQLite по умолчанию
	}
}

func getEnv(key, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}

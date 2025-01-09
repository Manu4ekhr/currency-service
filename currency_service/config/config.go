// package config
package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ABSAPIURL      string
	BankAPIURL     string
	UpdateInterval time.Duration
}

func GetConfig() Config {
	updateInterval, err := strconv.Atoi(getEnv("UPDATE_INTERVAL", "600")) // 600 секунд = 10 минут
	if err != nil {
		updateInterval = 600
	}

	return Config{
		ABSAPIURL:      getEnv("ABS_API_URL", "https://.............................."),
		BankAPIURL:     getEnv("BANK_API_URL", "https://............................."),
		UpdateInterval: time.Duration(updateInterval) * time.Second,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

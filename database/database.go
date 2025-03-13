package database

import (
	"errors"
	"fmt"
	"log"

	"currency-service/models"

	"gorm.io/driver/sqlite" // Используем SQLite (можно заменить на PostgreSQL)
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Автоматическое создание таблиц
	err = DB.AutoMigrate(&models.CurrencyRate{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database successfully connected and migrated")
}

func SaveRates(rates []models.CurrencyRate) error {
	for _, rate := range rates {
		var existingRate models.CurrencyRate
		if err := DB.Where("currency = ?", rate.Currency).First(&existingRate).Error; err == nil {
			// Валюта найдена → обновляем курс
			existingRate.Rate = rate.Rate
			if err := DB.Save(&existingRate).Error; err != nil {
				return fmt.Errorf("failed to update rate for %s: %w", rate.Currency, err)
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// Валюта не найдена → создаем новую запись
			if err := DB.Create(&rate).Error; err != nil {
				return fmt.Errorf("failed to insert new rate for %s: %w", rate.Currency, err)
			}
		} else {
			return fmt.Errorf("failed to check existing rate for %s: %w", rate.Currency, err)
		}
	}
	return nil
}

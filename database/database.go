package database

import (
	"currency-service/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func InitDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&entity.CurrencyOperation{})
	if err != nil {
		log.Fatalf("Не удалось выполнить миграцию базы данных: %v", err)
		return nil, err
	}

	log.Println("База данных успешно подключена и мигрирована")
	return db, nil
}

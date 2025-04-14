package repository

import (
	"currency-service/entity"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RatesRepository interface {
	SaveOrGetOperation(op entity.CurrencyOperation) (*entity.CurrencyOperation, error)
	UpdateStatus(extID string, oldStatus entity.OperationStatus, newStatus entity.OperationStatus) error
	GetRatesByStatus(status entity.OperationStatus) ([]entity.CurrencyOperation, error)
}

type ratesRepo struct {
	DB *gorm.DB
}

func NewRatesRepository(db *gorm.DB) RatesRepository {
	return &ratesRepo{DB: db}
}

// Сохраняем новую операцию или возвращаем уже существующую
func (r *ratesRepo) SaveOrGetOperation(op entity.CurrencyOperation) (*entity.CurrencyOperation, error) {
	var existing entity.CurrencyOperation

	err := r.DB.Where("ext_id = ?", op.ExtID).First(&existing).Error
	if err == nil {
		// Найдена существующая операция — возвращаем её
		return &existing, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("ошибка при поиске операции: %w", err)
	}

	// Не найдена — создаём новую
	if err := r.DB.Create(&op).Error; err != nil {
		return nil, fmt.Errorf("не удалось сохранить новую операцию: %w", err)
	}
	return &op, nil
}

// Получить все операции по статусу
func (r *ratesRepo) GetRatesByStatus(status entity.OperationStatus) ([]entity.CurrencyOperation, error) {
	var ops []entity.CurrencyOperation
	err := r.DB.Where("status = ?", status).Find(&ops).Error
	if err != nil {
		return nil, fmt.Errorf("не удалось получить операции со статусом %s: %w", status, err)
	}
	return ops, nil
}

// Обновляем статус операции по ext_id с проверкой текущего статуса и блокировкой SELECT FOR UPDATE
func (r *ratesRepo) UpdateStatus(extID string, oldStatus, newStatus entity.OperationStatus) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var op entity.CurrencyOperation
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("ext_id = ?", extID).First(&op).Error
		if err != nil {
			return err
		}

		if op.Status != oldStatus {
			return fmt.Errorf("статус не совпадает: ожидался %s, в базе %s", oldStatus, op.Status)
		}

		op.Status = newStatus
		return tx.Save(&op).Error
	})
}

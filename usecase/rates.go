package usecase

import (
	"currency-service/adapter"
	"currency-service/entity"
	"currency-service/repository"
	"fmt"
	"log"
)

type RatesUsecase struct {
	absAdapter    adapter.Adapter
	repo          repository.RatesRepository
	websiteSender adapter.WebsiteSender
}

func NewRatesUsecase(abs adapter.Adapter, repo repository.RatesRepository, website adapter.WebsiteSender) *RatesUsecase {
	return &RatesUsecase{
		absAdapter:    abs,
		repo:          repo,
		websiteSender: website,
	}
}

func (u *RatesUsecase) FetchSaveAndSendRates() error {
	log.Println("Получаем курсы валют из АБС...")
	data, err := u.absAdapter.Fetcher()
	if err != nil {
		return err
	}
	log.Printf("Получено %d курсов валют", len(data.CurrRates))

	var opsToSend []entity.CurrencyOperation

	for _, r := range data.CurrRates {
		operation := entity.CurrencyOperation{
			From:     "ABS",
			To:       "Website",
			ExtID:    fmt.Sprintf("abs-%s", r.Cur),
			Type:     entity.TypeNBT, // Можно заменить, если логика будет расширяться
			Status:   entity.StatusPending,
			Currency: r.Cur,
			Rate:     r.Sell,
		}

		// Сохраняем или получаем существующую
		savedOp, err := u.repo.SaveOrGetOperation(operation)
		if err != nil {
			log.Printf("Ошибка сохранения операции: %v", err)
			continue
		}

		// Пропускаем, если уже отправлено
		if savedOp.Status == entity.StatusSended {
			log.Printf("Валюта %s уже отправлена, пропускаем", savedOp.Currency)
			continue
		}

		opsToSend = append(opsToSend, *savedOp)
	}

	// Отправляем на сайт
	if len(opsToSend) == 0 {
		log.Println("Нет новых курсов для отправки.")
		return nil
	}

	log.Println("Отправляем курсы на сайт банка...")
	if err := u.websiteSender.SendRates(opsToSend); err != nil {
		return err
	}
	log.Println("Курсы успешно отправлены на сайт")

	// Обновляем статус на "sended"
	for _, op := range opsToSend {
		_ = u.repo.UpdateStatus(op.ExtID, op.Status, entity.StatusSended)

	}
	log.Println("Статусы обновлены на 'sended'")

	return nil
}

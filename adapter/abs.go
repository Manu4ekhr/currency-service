package adapter

import (
	"currency-service/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type adapter struct {
	client http.Client
	URL    string
}

type Adapter interface {
	Fetcher() (*models.ExchangeRate, error)
}

func NewAdapter(client http.Client, url string) Adapter {
	return &adapter{
		client: client,
		URL:    url,
	}
}

// Сборщик реализует адаптер
func (a *adapter) Fetcher() (*models.ExchangeRate, error) {
	// 1. Делаем GET-запрос
	resp, err := a.client.Get(a.URL + "/currency_rates")
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// 2. Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %w", err)
	}

	// 3. Проверяем, что код ответа 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус: %d — %s", resp.StatusCode, string(body))
	}

	// 4. Распарсим JSON в структуру models.ExchangeRate
	var result models.ExchangeRate
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка при разборе JSON: %w", err)
	}

	// 5. Вернём результат (а не сохраняем в базу — usecase этим займётся)
	return &result, nil
}

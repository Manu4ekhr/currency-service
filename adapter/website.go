package adapter

import (
	"bytes"
	"currency-service/entity"
	"currency-service/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// Интерфейс для отправки курсов на сайт
type WebsiteSender interface {
	SendRates(rate []entity.CurrencyOperation) error
}

// Реализация адаптера
type websiteAdapter struct {
	client   http.Client
	URL      string
	AuthURL  string
	Email    string
	Password string
}

// Конструктор адаптера
func NewWebsiteSender(client http.Client, url, authURL, email, password string) WebsiteSender {
	return &websiteAdapter{
		client:   client,
		URL:      url,
		AuthURL:  authURL,
		Email:    email,
		Password: password,
	}
}

// Метод отправки курсов на сайт
func (w *websiteAdapter) SendRates(rates []entity.CurrencyOperation) error {
	// Получаем токен авторизации
	token, err := w.getAuthToken()
	if err != nil {
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	// Преобразуем []CurrencyRate в []RateType
	payload := convertToWebsiteFormat(rates)
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка кодирования JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, w.URL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус ответа от сайта: %d", resp.StatusCode)
	}
	return nil
}

// Получение access token (авторизация)
func (w *websiteAdapter) getAuthToken() (string, error) {
	type AuthRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type AuthResponse struct {
		Token string `json:"token"`
	}
	//authUrl := fmt.Sprintf("%s%s", url, "/login") // Очень важно! Учитывать это изменение! (изменение в .env)
	url := fmt.Sprintf("%s/login", w.AuthURL)
	body := AuthRequest{
		Email:    w.Email,
		Password: w.Password,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("ошибка кодирования JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неожиданный статус ответа: %d", resp.StatusCode)
	}

	var res AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return res.Token, nil
}

// Вспомогательная функция: преобразуем []entity.CurrencyOperation в []RateType
func convertToWebsiteFormat(rates []entity.CurrencyOperation) []models.RateType {
	var rateList []models.Rate
	for _, r := range rates {
		rateList = append(rateList, models.Rate{
			Currency: r.Currency,
			Buy:      fmt.Sprintf("%.4f", r.Rate),
			Sale:     fmt.Sprintf("%.4f", r.Rate),
		})
	}
	return []models.RateType{
		{
			RateNameTJ:  "Интиқолҳо",
			RateNameRU:  "Переводы",
			RateNameUSD: "Translations",
			Rates:       rateList,
			Key:         "translations",
		},
	}
}

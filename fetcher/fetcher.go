package fetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"currency-service/database"
	"currency-service/models"
)

var client = &http.Client{Timeout: 10 * time.Second}

func FetchRates(ctx context.Context, absAPIURL string) ([]models.CurrencyRate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, absAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from ABS API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rates []models.CurrencyRate
	if err := json.Unmarshal(body, &rates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(rates) == 0 {
		return nil, errors.New("no currency rates found in response")
	}

	// Сохраняем в базу данных
	if err := database.SaveRates(rates); err != nil {
		return nil, fmt.Errorf("failed to save rates to database: %w", err)
	}

	return rates, nil
}

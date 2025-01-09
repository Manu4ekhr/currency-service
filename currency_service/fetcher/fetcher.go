// package fetcher
package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"currency-service/models"
)

func FetchRates(ctx context.Context, absAPIURL string) ([]models.CurrencyRate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, absAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from ABS API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var rates []models.CurrencyRate
	if err := json.Unmarshal(body, &rates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return rates, nil
}

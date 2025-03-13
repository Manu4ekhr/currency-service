package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"currency-service/models"
)

var client = &http.Client{Timeout: 10 * time.Second}

func SendRates(ctx context.Context, bankAPIURL string, rates []models.CurrencyRate) error {
	if len(rates) == 0 {
		return errors.New("no currency rates to send")
	}

	payload, err := json.Marshal(rates)
	if err != nil {
		return fmt.Errorf("failed to marshal rates: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bankAPIURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send data to bank API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

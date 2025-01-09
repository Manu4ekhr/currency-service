// package sender
package sender 

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"currency-service/models"
)

func SendRates(ctx context.Context, bankAPIURL string, rates []models.CurrencyRate) error {
	payload, err := json.Marshal(rates)
	if err != nil {
		return fmt.Errorf("failed to marshal rates: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bankAPIURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send data to bank API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	return nil
}

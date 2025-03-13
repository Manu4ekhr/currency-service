package main

import (
	"context"
	"currency-service/config"
	"currency-service/database"
	"currency-service/fetcher"
	"currency-service/sender"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.GetConfig()
	database.InitDatabase(cfg.DatabaseDSN) // Подключаем БД

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(cfg.UpdateInterval)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Currency rate service started")

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Fetching currency rates...")
				rates, err := fetcher.FetchRates(ctx, cfg.ABSAPIURL)
				if err != nil {
					log.Printf("Error fetching rates: %v", err)
					continue
				}

				log.Println("Sending currency rates to bank API...")
				if err := sender.SendRates(ctx, cfg.BankAPIURL, rates); err != nil {
					log.Printf("Error sending rates: %v", err)
				} else {
					log.Println("Rates successfully updated on the bank website")
				}
			case <-ctx.Done():
				log.Println("Shutting down service...")
				return
			}
		}
	}()

	<-sigs
	log.Println("Received shutdown signal")
	cancel()
	time.Sleep(2 * time.Second) // Даем время завершить текущие операции
	log.Println("Service stopped")
}

package main

import (
	"context"
	"currency-service/adapter"
	"currency-service/config"
	"currency-service/database"
	"currency-service/job"
	"currency-service/repository"
	"currency-service/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения")
	}

	cfg := config.GetConfig()

	// Инициализация базы данных
	db, err := database.InitDatabase(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Ошибка при инициализации БД: %v", err)
	}

	// Создаем адаптеры и зависимости
	absAdapter := adapter.NewAdapter(http.Client{}, cfg.ABSAPIURL)
	websiteSender := adapter.NewWebsiteSender(http.Client{}, cfg.BankAPIURL, cfg.AuthURL, cfg.LoginEmail, cfg.LoginPassword)
	ratesRepo := repository.NewRatesRepository(db)
	useCase := usecase.NewRatesUsecase(absAdapter, ratesRepo, websiteSender)

	// Запуск контекста и планировщика
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler, err := job.NewScheduler()
	if err != nil {
		log.Fatalf("Ошибка при создании планировщика: %v", err)
	}

	err = scheduler.Register(cfg.UpdateInterval, func() {
		if err := useCase.FetchSaveAndSendRates(); err != nil {
			log.Printf("Ошибка при выполнении задачи: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Ошибка при регистрации задачи: %v", err)
	}

	scheduler.Start()
	log.Println("Сервис курсов валют запущен с использованием gocron")

	// Ожидаем сигнал завершения
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("Получен сигнал завершения (Ctrl+C или остановка)")

	scheduler.Shutdown()
	cancel()
	time.Sleep(2 * time.Second)
	log.Println("Сервис остановлен")
}

# Currency Service

Сервис для автоматического получения курсов валют из АБС ЦФТ и публикации их на сайт банка.

## Функционал:
Получает курс валют с АБС ЦФТ.  
Сохраняет курсы в базу данных (SQLite/PostgreSQL).  
Отправляет актуальные курсы на сайт банка.  
Логирует все процессы.

## Установка и запуск:

### 1. Клонируем репозиторий:
```sh
git clone https://github.com/Manu4ekhr/currency-service.git 
cd currency-service

## API:

1️. Получить список курсов валют

**GET** `/rates`

Описание:

Возвращает актуальный список валют и их курсов.

Пример запроса:

```sh
curl -X GET http://localhost:7540/rates

Пример ответа (200 OK):

[
  { "currency": "USD", "rate": 11.25 },
  { "currency": "EUR", "rate": 12.40 }
]

2. Обновить курс валют

POST /update

Описание:
Принимает список валют и обновляет их в системе.

Пример запроса:

curl -X POST http://localhost:7540/update \
  -H "Content-Type: application/json" \
  -d '[{"currency": "USD", "rate": 11.30}, {"currency": "EUR", "rate": 12.50}]'

Пример ответа (200 OK):

{ "message": "Rates updated successfully" }
🔹 Ошибки (400 Bad Request):

{ "error": "Invalid data format" }
package models

// Отправляем на сайт (формат получения данных на сайте)
type Rate struct {
	Currency string `json:"title"`
	Buy      string `json:"value_buy"`
	Sale     string `json:"value_sale"`
}

type RateType struct {
	RateNameTJ  string `json:"title_tg"`
	RateNameRU  string `json:"title_ru"`
	RateNameUSD string `json:"title_en"`
	Rates       []Rate `json:"data"`
	Key         string `json:"translations"`
}

// Получаем из АБС
type ExchangeRate struct {
	Base       string  `json:"base"`
	LastUpdate string  `json:"last_update"`
	CurrRates  []Rates `json:"rates"`
}

type Rates struct {
	UpdateID         int64   `json:"update_id"`
	Cur              string  `json:"cur"`
	Buy              float64 `json:"buy"`
	Sell             float64 `json:"sell"`
	Mode             string  `json:"mode"`
	ExchangeRateType string  `json:"exchange_rate_type"`
}

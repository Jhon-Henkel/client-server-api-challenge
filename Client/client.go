package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ExchangeRate struct {
	USD float64 `json:"bid"`
}

func main() {
	exchangeRate, err := requestDollarExchangeRate()
	if err != nil {
		panic(err)
	}
	println(exchangeRate.USD)
}

func requestDollarExchangeRate() (*ExchangeRate, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3000*time.Microsecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var exchangeRate ExchangeRate
	err = json.NewDecoder(res.Body).Decode(&exchangeRate)
	if err != nil {
		return nil, err
	}
	return &exchangeRate, nil
}
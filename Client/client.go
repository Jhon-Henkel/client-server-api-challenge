package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type ExchangeRate struct {
	USD string `json:"bid"`
}

func main() {
	exchangeRate, err := requestDollarExchangeRate()
	if err != nil {
		log.Println(err)
		return
	}
	println(exchangeRate.USD)
	saveResult(exchangeRate.USD)
}

func requestDollarExchangeRate() (*ExchangeRate, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var exchangeRate ExchangeRate
	err = json.NewDecoder(res.Body).Decode(&exchangeRate)
	if err != nil {
		return nil, err
	}
	return &exchangeRate, nil
}

func saveResult(result string) {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println(err)
		return
	}
	file.WriteString("Dólar: " + result + "\n")
	file.Close()
}

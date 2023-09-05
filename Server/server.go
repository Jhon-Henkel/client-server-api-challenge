package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ExchangeRate struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", dollarExchangeRateHandler)
	http.ListenAndServe(":8080", nil)
}

func dollarExchangeRateHandler(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/cotacao" {
		log.Println("Not found")
		response.WriteHeader(http.StatusNotFound)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	exchangeRate, err := dollarExchangeRateRequest()
	if err != nil {
		log.Println(err.Error())
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(response).Encode(exchangeRate.USDBRL)
}

func dollarExchangeRateRequest() (*ExchangeRate, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

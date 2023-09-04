package main

import (
	"encoding/json"
	"io"
	"net/http"
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
		response.WriteHeader(http.StatusNotFound)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	exchangeRate, err := DollarExchangeRateRequest()
	if err != nil {
		http.Error(response, "Internal server error", http.StatusInternalServerError)
	}
	json.NewEncoder(response).Encode(exchangeRate.USDBRL)
}

func DollarExchangeRateRequest() (*ExchangeRate, error) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var exchangeRate ExchangeRate
	err = json.Unmarshal(res, &exchangeRate)
	if err != nil {
		return nil, err
	}
	return &exchangeRate, nil
}

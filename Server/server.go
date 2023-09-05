package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExchangeRate struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type ExchangeRateDB struct {
	ID  int `gorm:"primaryKey"`
	bid string
	gorm.Model
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
	err = saveExchangeRate(exchangeRate)
	if err != nil {
		log.Println(err)
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

func saveExchangeRate(exchangeRate *ExchangeRate) error {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&ExchangeRateDB{})
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	dbError := db.WithContext(ctx).Create(&ExchangeRateDB{bid: exchangeRate.USDBRL.Bid})
	if dbError.Error != nil {
		return dbError.Error
	}
	return nil
}

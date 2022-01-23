package stockapi

import (
	"context"
	"stockx-backend/db/models"
	"strconv"
	"time"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

type SymbolInfo struct {
	CurrentStock finnhub.Quote           `json:"stock"`
	Profile      finnhub.CompanyProfile2 `json:"profile"`
	News         []finnhub.CompanyNews   `json:"news"`
	Hold         models.HoldStocks       `json:"hold"`
	Credits      float32                 `json:"credits"` //users available credits for trading
	BuyingPower  float32                 `json:"buyingPower"`
}

func GetCompanyInfo(symbol string) (SymbolInfo, error) {
	finnhubClient := finnhub.NewAPIClient(FinnhubConfiguration).DefaultApi

	readyChan := make(chan int)

	var err error

	var currentStockInfo finnhub.Quote

	go func(s string, c chan int) {
		currentStockInfo, _, err = finnhubClient.Quote(context.Background()).Symbol(s).Execute()
		if err != nil {
			c <- 1
		} else {
			c <- 0
		}
	}(symbol, readyChan)

	var companyProfile finnhub.CompanyProfile2

	go func(s string, c chan int) {
		companyProfile, _, err = finnhubClient.CompanyProfile2(context.Background()).Symbol(s).Execute()
		if err != nil {
			c <- 1
		} else {
			c <- 0
		}
	}(symbol, readyChan)

	var companyNews []finnhub.CompanyNews

	dt := time.Now()
	nowDate := strconv.Itoa(dt.Year()) + "-" + dt.Month().String() + "-" + strconv.Itoa(dt.Day())
	from := dt.AddDate(0, -1, 0)
	fromDate := strconv.Itoa(from.Year()) + "-" + from.Month().String() + "-" + strconv.Itoa(from.Day())

	go func(s string, c chan int) {
		companyNews, _, err = finnhubClient.CompanyNews(context.Background()).Symbol(s).From(fromDate).To(nowDate).Execute()
		if err != nil {
			c <- 1
		} else {
			c <- 0
		}
	}(symbol, readyChan)

	for i := 0; i < 3; i++ {
		if <-readyChan != 0 {
			return SymbolInfo{}, err
		}
	}

	return SymbolInfo{CurrentStock: currentStockInfo, Profile: companyProfile, News: companyNews}, nil
}

func GetQuoteForSymbol(symbol string) (finnhub.Quote, error) {
	finnhubClient := finnhub.NewAPIClient(FinnhubConfiguration).DefaultApi

	currentStockInfo, _, err := finnhubClient.Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return finnhub.Quote{}, err
	}

	return currentStockInfo, nil
}

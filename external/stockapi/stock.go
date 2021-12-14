package stockapi

import (
	"context"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func GetStockSymbols(q string) ([]finnhub.StockSymbol, error) {
	finnhubClient := finnhub.NewAPIClient(FinnhubConfiguration).DefaultApi

	res, _, err := finnhubClient.StockSymbols(context.Background()).Exchange("US").Mic("XNAS").Currency("USD").SecurityType("Common Stock").Execute()
	if err != nil {
		return []finnhub.StockSymbol{}, err
	}

	return res, nil
}

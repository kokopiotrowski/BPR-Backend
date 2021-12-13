package stock

import (
	"stockx-backend/db"
	"stockx-backend/db/models"
	"stockx-backend/external/stockapi"
	"stockx-backend/reserr"
)

func GetSymbolInfo(email, symbol string) (stockapi.SymbolInfo, error) {
	symbolInfo, err := stockapi.GetCompanyInfo(symbol)
	if err != nil {
		return stockapi.SymbolInfo{}, reserr.Internal("Symbol info error", err, "Could not retrieve symbol information")
	}

	symbolInfo.Hold = models.HoldStocks{
		HoldLong:  []models.HoldLong{},
		HoldShort: []models.HoldShort{},
	}

	trades, err := db.GetTradesFromTableForUser(email)
	if err != nil {
		return stockapi.SymbolInfo{}, reserr.Internal("Symbol info error", err, "Could not retrieve users information for this symbol")
	}

	for _, t := range trades.HoldLong {
		if t.Symbol == symbol {
			symbolInfo.Hold.HoldLong = append(symbolInfo.Hold.HoldLong, t)
		}
	}

	for _, t := range trades.HoldShort {
		if t.Symbol == symbol {
			symbolInfo.Hold.HoldShort = append(symbolInfo.Hold.HoldShort, t)
		}
	}

	return symbolInfo, nil
}

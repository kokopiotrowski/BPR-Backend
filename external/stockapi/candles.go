package stockapi

import (
	"context"
	"errors"
	"stockx-backend/reserr"
	"time"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

type Candles struct {
	Ohlcv [][]interface{} `json:"ohlcv"` // time (int64), o, h, l, c, v
}

func GetStockCandles(symbol string) (Candles, error) {
	finnhubClient := finnhub.NewAPIClient(FinnhubConfiguration).DefaultApi

	res, _, err := finnhubClient.StockCandles(context.Background()).Symbol(symbol).Resolution("D").From(time.Now().AddDate(-1, 0, 0).Unix()).To(time.Now().Unix()).Execute()
	if err != nil {
		return Candles{}, err
	}

	if *res.S != "ok" {
		return Candles{}, reserr.Internal("error", errors.New("response from external api is empty"), "Failed to retrieve chart data for this stock.\nSearch for more popular ones.")
	}

	candles := Candles{
		Ohlcv: make([][]interface{}, len(*res.T)),
	}

	if *res.S == "ok" {
		for i := 0; i < len(*res.T); i++ {
			candles.Ohlcv[i] = []interface{}{
				(*res.T)[i] * 1000,
				(*res.O)[i],
				(*res.H)[i],
				(*res.L)[i],
				(*res.C)[i],
				(*res.V)[i]}
		}
	}

	return candles, nil
}

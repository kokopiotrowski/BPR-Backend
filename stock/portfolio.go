package stock

import (
	"errors"
	"stockx-backend/db"
	"stockx-backend/db/models"
	"stockx-backend/external/stockapi"
	"stockx-backend/reserr"
)

func GetPortfolio(email string) (models.Trades, error) {
	trade, err := db.GetTradesFromTableForUser(email)
	if err != nil {
		return models.Trades{}, err
	}

	readyChan := make(chan int)

	trade.AccountValue = trade.Credits

	for i := 0; i < len(trade.HoldLong); i++ {
		go func(index int) {
			q, err := stockapi.GetQuoteForSymbol(trade.HoldLong[index].Symbol)
			if err != nil {
				readyChan <- 1
				return
			}

			trade.HoldLong[index].CurrentPrice = q.GetC()
			trade.HoldLong[index].Gain = (trade.HoldLong[index].CurrentPrice - trade.HoldLong[index].Price) * float32(trade.HoldLong[index].Amount)
			trade.HoldLong[index].GainP = (trade.HoldLong[index].Gain / (trade.HoldLong[index].Price * float32(trade.HoldLong[index].Amount))) * 100

			trade.AccountValue += trade.HoldLong[index].CurrentPrice * float32(trade.HoldLong[index].Amount)
			readyChan <- 0
		}(i)
	}

	for i := 0; i < len(trade.HoldLong); i++ {
		if <-readyChan == 1 {
			return models.Trades{}, reserr.Internal("External api fail", errors.New("failed to retrieve data from external api"), "Failed to retrieve portfolio from server")
		}
	}

	for i := 0; i < len(trade.HoldShort); i++ {
		go func(index int) {
			q, err := stockapi.GetQuoteForSymbol(trade.HoldShort[index].Symbol)
			if err != nil {
				readyChan <- 1
				return
			}

			trade.HoldShort[index].CurrentPrice = q.GetC()
			trade.HoldShort[index].Gain = (trade.HoldShort[index].Price - trade.HoldShort[index].CurrentPrice) * float32(trade.HoldShort[index].Amount)
			trade.HoldShort[index].GainP = (trade.HoldShort[index].Gain / (trade.HoldShort[index].Price * float32(trade.HoldShort[index].Amount))) * 100

			trade.AccountValue -= trade.HoldShort[index].CurrentPrice * float32(trade.HoldShort[index].Amount)

			readyChan <- 0
		}(i)
	}

	for i := 0; i < len(trade.HoldShort); i++ {
		if <-readyChan == 1 {
			return models.Trades{}, reserr.Internal("External api fail", errors.New("failed to retrieve data from external api"), "Failed to retrieve portfolio from server")
		}
	}

	return trade, nil
}

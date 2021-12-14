package stock

import (
	"errors"
	"stockx-backend/db"
	"stockx-backend/db/models"
	"stockx-backend/external/stockapi"
	"stockx-backend/reserr"
	"time"
)

func BuyStock(email string, item models.BoughtStock) error {
	if item.Amount < 1 {
		return reserr.BadRequest("info", errors.New("buying 0 stocks"), "How do you think that would work? :)")
	}

	trades, err := db.GetTradesFromTableForUser(email)
	if err != nil {
		return err
	}

	stock, err := stockapi.GetQuoteForSymbol(item.Symbol)
	if err != nil {
		return reserr.Internal("error", err, "Failed to get information about the stock")
	}

	if trades.Credits < *stock.C*float32(item.Amount) {
		return reserr.BadRequest("info", errors.New("not enough credits to buyt stocks"), "Not enough credits to buy stocks")
	}

	item.Price = *stock.C
	trades.Credits -= item.Price * float32(item.Amount)

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return err
	}

	dt := time.Now().In(loc)

	item.Date = dt.Unix() * 1000
	trades.BoughtStocks = append(trades.BoughtStocks, item)

	exists := false

	for i := 0; i < len(trades.HoldLong); i++ {
		if trades.HoldLong[i].Symbol == item.Symbol {
			currentBuyPrice := ((trades.HoldLong[i].Price * float32(trades.HoldLong[i].Amount)) + (item.Price * float32(item.Amount))) / (float32(trades.HoldLong[i].Amount) + float32(item.Amount))
			trades.HoldLong[i].Price = currentBuyPrice
			trades.HoldLong[i].Amount += item.Amount
			exists = true

			break
		}
	}

	if !exists {
		trades.HoldLong = append(trades.HoldLong, models.HoldLong{
			Symbol: item.Symbol,
			Price:  item.Price,
			Amount: item.Amount,
		})
	}

	err = db.PutTradesInTheTable(email, trades)
	if err != nil {
		return reserr.Internal("error", err, "Failed to update trades. Purchase has not been made")
	}

	return nil
}

func ShortStock(email string, item models.ShortStock) error {
	if item.Amount < 1 {
		return reserr.BadRequest("info", errors.New("sell 0 stocks"), "How do you think that would work? :)")
	}

	trades, err := db.GetTradesFromTableForUser(email)
	if err != nil {
		return err
	}

	stock, err := stockapi.GetQuoteForSymbol(item.Symbol)
	if err != nil {
		return reserr.Internal("error", err, "Failed to get information about the stock")
	}

	item.Price = *stock.C
	trades.Credits += item.Price * float32(item.Amount)

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return err
	}

	dt := time.Now().In(loc)

	item.Date = dt.Unix() * 1000
	trades.ShortStocks = append(trades.ShortStocks, item)

	exists := false

	for i := 0; i < len(trades.HoldShort); i++ {
		if trades.HoldShort[i].Symbol == item.Symbol {
			currentBuyPrice := ((trades.HoldShort[i].Price * float32(trades.HoldShort[i].Amount)) + (item.Price * float32(item.Amount))) / (float32(trades.HoldShort[i].Amount) + float32(item.Amount))
			trades.HoldShort[i].Price = currentBuyPrice
			trades.HoldShort[i].Amount += item.Amount
			exists = true

			break
		}
	}

	if !exists {
		trades.HoldShort = append(trades.HoldShort, models.HoldShort{
			Symbol: item.Symbol,
			Price:  item.Price,
			Amount: item.Amount,
		})
	}

	err = db.PutTradesInTheTable(email, trades)

	return err
}

func SellStock(email string, item models.SoldStock) error {
	if item.Amount < 1 {
		return reserr.BadRequest("info", errors.New("sell 0 stocks"), "How do you think that would work? :)")
	}

	trades, err := db.GetTradesFromTableForUser(email)
	if err != nil {
		return err
	}

	stock, err := stockapi.GetQuoteForSymbol(item.Symbol)
	if err != nil {
		return reserr.Internal("error", err, "Failed to get information about the stock")
	}

	amount := item.Amount

	for i := len(trades.HoldLong) - 1; i >= 0; i-- {
		if amount == 0 {
			break
		}

		if trades.HoldLong[i].Symbol == item.Symbol {
			for amount > 0 && trades.HoldLong[i].Amount > 0 {
				amount--
				trades.HoldLong[i].Amount--
			}

			if trades.HoldLong[i].Amount == 0 {
				trades.HoldLong = append(trades.HoldLong[:i], trades.HoldLong[i+1:]...)
				continue
			}
		}
	}

	if amount > 0 {
		return reserr.BadRequest("info", errors.New("not enough stocks held by user to sell"), "You don't own that many stocks of "+item.Symbol+" to sell! Firstly, buy some more :)")
	}

	item.Price = *stock.C
	trades.Credits += item.Price * float32(item.Amount)

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return err
	}

	dt := time.Now().In(loc)

	item.Date = dt.Unix() * 1000
	trades.SoldStocks = append(trades.SoldStocks, item)

	err = db.PutTradesInTheTable(email, trades)

	return err
}

func BuyToCover(email string, item models.BoughtToCover) error {
	if item.Amount < 1 {
		return reserr.BadRequest("info", errors.New("info"), "How do you think that would work?")
	}

	trades, err := db.GetTradesFromTableForUser(email)
	if err != nil {
		return err
	}

	stock, err := stockapi.GetQuoteForSymbol(item.Symbol)
	if err != nil {
		return reserr.Internal("error", err, "Failed to get information about the stock")
	}

	amount := item.Amount

	for i := len(trades.HoldShort) - 1; i >= 0; i-- {
		if amount == 0 {
			break
		}

		if trades.HoldShort[i].Symbol == item.Symbol {
			for amount > 0 && trades.HoldShort[i].Amount > 0 {
				amount--
				trades.HoldShort[i].Amount--
			}

			if trades.HoldShort[i].Amount == 0 {
				trades.HoldShort = append(trades.HoldShort[:i], trades.HoldShort[i+1:]...)
				continue
			}
		}
	}

	if amount > 0 {
		return reserr.BadRequest("info", errors.New("not that many stocks to cover"), "You don't have that many stocks of "+item.Symbol+" to cover")
	}

	item.Price = *stock.C
	trades.Credits -= item.Price * float32(item.Amount)

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return err
	}

	dt := time.Now().In(loc)

	item.Date = dt.Unix() * 1000
	trades.BoughtToCover = append(trades.BoughtToCover, item)

	err = db.PutTradesInTheTable(email, trades)

	return err
}

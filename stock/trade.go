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
		return err
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

	for _, hl := range trades.HoldLong {
		if hl.Symbol == item.Symbol {
			currentBuyPrice := ((hl.Price * float32(hl.Amount)) + (item.Price * float32(item.Amount))) / (float32(hl.Amount) + float32(item.Amount))
			hl.Price = currentBuyPrice
			hl.Amount += item.Amount
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

	return err
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
		return err
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

	for _, hs := range trades.HoldShort {
		if hs.Symbol == item.Symbol {
			currentBuyPrice := ((hs.Price * float32(hs.Amount)) + (item.Price * float32(item.Amount))) / (float32(hs.Amount) + float32(item.Amount))
			hs.Price = currentBuyPrice
			hs.Amount += item.Amount
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
		return err
	}

	amount := item.Amount

	for i := len(trades.HoldLong) - 1; i >= 0; i-- {
		if trades.HoldLong[i].Symbol == item.Symbol {
			for amount > 0 && trades.HoldLong[i].Amount > 0 {
				amount--
				trades.HoldLong[i].Amount--
			}

			if amount == 0 {
				break
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
		return err
	}

	amount := item.Amount

	for i := len(trades.HoldShort) - 1; i >= 0; i-- {
		if trades.HoldShort[i].Symbol == item.Symbol {
			for amount > 0 && trades.HoldShort[i].Amount > 0 {
				amount--
				trades.HoldShort[i].Amount--
			}

			if amount == 0 {
				break
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
	trades.Credits += item.Price * float32(item.Amount)

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

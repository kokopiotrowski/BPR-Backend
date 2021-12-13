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
	trades, err := db.GetTradesFromTableForUser(email)
	if err != nil {
		return err
	}

	stock, err := stockapi.GetQuoteForSymbol(item.Symbol)
	if err != nil {
		return err
	}

	if trades.Credits < *stock.C*float32(item.Amount) {
		return reserr.BadRequest("not enough credits to buy stocks", nil, "not enough credits to buy stocks")
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
	trades.HoldLong = append(trades.HoldLong, models.HoldLong{
		Symbol: item.Symbol,
		Price:  item.Price,
		Amount: item.Amount,
		Date:   item.Date,
	})

	err = db.PutTradesInTheTable(email, trades)

	return err
}

func ShortStock(email string, item models.ShortStock) error {
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

	item.Date = dt.Unix()
	trades.ShortStocks = append(trades.ShortStocks, item)
	trades.HoldShort = append(trades.HoldShort, models.HoldShort{
		Symbol: item.Symbol,
		Price:  item.Price,
		Amount: item.Amount,
		Date:   item.Date,
	})

	err = db.PutTradesInTheTable(email, trades)

	return err
}

func SellStock(email string, item models.SoldStock) error {
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
		return reserr.Information("", errors.New(""), "You don't own that many stocks of "+item.Symbol+" to sell! Firstly, buy some more :)")
	}

	item.Price = *stock.C
	trades.Credits += item.Price * float32(item.Amount)

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return err
	}

	dt := time.Now().In(loc)

	item.Date = dt.Unix()
	trades.SoldStocks = append(trades.SoldStocks, item)

	err = db.PutTradesInTheTable(email, trades)

	return err
}

func BuyToCover(email string, item models.BoughtToCover) error {
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
		return reserr.Information("failed to cover stocks - you don't have that many stocks to cover", errors.New("failed to cover stocks - you don't have that many stocks to cover"), "failed to cover stocks - you don't have that many stocks of "+item.Symbol+" to cover")
	}

	item.Price = *stock.C
	trades.Credits += item.Price * float32(item.Amount)

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return err
	}

	dt := time.Now().In(loc)

	item.Date = dt.Unix()
	trades.BoughtToCover = append(trades.BoughtToCover, item)

	err = db.PutTradesInTheTable(email, trades)

	return err
}

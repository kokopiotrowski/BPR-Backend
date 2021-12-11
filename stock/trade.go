package stock

import (
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
		return reserr.BadRequest("not enough credits to buy stocks", nil, "")
	}

	item.Price = *stock.C
	trades.Credits -= item.Price * float32(item.Amount)

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return err
	}

	dt := time.Now().In(loc)

	item.Date = dt.Unix()
	trades.BoughtStocks = append(trades.BoughtStocks, item)

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

	for i := len(trades.BoughtStocks) - 1; i >= 0; i-- {
		if trades.BoughtStocks[i].Symbol == item.Symbol {
			for amount > 0 && trades.BoughtStocks[i].Amount > 0 {
				amount--
				trades.BoughtStocks[i].Amount--
			}

			if amount == 0 {
				break
			}

			if trades.BoughtStocks[i].Amount == 0 {
				trades.BoughtStocks = append(trades.BoughtStocks[:i], trades.BoughtStocks[i+1:]...)
				continue
			}
		}
	}

	if amount > 0 {
		return reserr.BadRequest("you don't own that many stocks to sell", nil, "")
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

	for i := len(trades.ShortStocks) - 1; i >= 0; i-- {
		if trades.ShortStocks[i].Symbol == item.Symbol {
			for amount > 0 && trades.ShortStocks[i].Amount > 0 {
				amount--
				trades.ShortStocks[i].Amount--
			}

			if amount == 0 {
				break
			}

			if trades.ShortStocks[i].Amount == 0 {
				trades.ShortStocks = append(trades.ShortStocks[:i], trades.ShortStocks[i+1:]...)
				continue
			}
		}
	}

	if amount > 0 {
		return reserr.BadRequest("failed to cover stocks - you don't have that many stocks to cover", nil, "")
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

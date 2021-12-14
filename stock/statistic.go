package stock

import (
	"fmt"
	"stockx-backend/db"
	"stockx-backend/db/models"
	"stockx-backend/reserr"
	"time"
)

type ReturnStatTable struct {
	Name string    `json:"name"`
	X    []string  `json:"x"`
	Y    []float32 `json:"y"`
}

func GetStatisticsForUser(email string) ([]ReturnStatTable, error) {
	stats, err := db.GetStatisticsFromTableForUser(email)
	if err != nil {
		return []ReturnStatTable{}, reserr.Internal("error", err, "Failed to retrieve statistics")
	}

	accountValueTable := ReturnStatTable{
		Name: "Account Value",
		X:    []string{},
		Y:    []float32{},
	}

	for _, acval := range stats.AccountValue {
		accountValueTable.X = append(accountValueTable.X, acval.Date)
		accountValueTable.Y = append(accountValueTable.Y, acval.AcountValue)
	}

	creditsTable := ReturnStatTable{
		Name: "Credits",
		X:    []string{},
		Y:    []float32{},
	}

	for _, accred := range stats.Credits {
		creditsTable.X = append(creditsTable.X, accred.Date)
		creditsTable.Y = append(creditsTable.Y, accred.Credits)
	}

	amountOwnedStockTable := ReturnStatTable{
		Name: "Amount of owned stocks",
		X:    []string{},
		Y:    []float32{},
	}

	for _, ownedStockAm := range stats.OwnedStocksAmount {
		amountOwnedStockTable.X = append(amountOwnedStockTable.X, ownedStockAm.Date)
		amountOwnedStockTable.Y = append(amountOwnedStockTable.Y, float32(ownedStockAm.Amount))
	}

	return []ReturnStatTable{
		accountValueTable,
		creditsTable,
		amountOwnedStockTable,
	}, nil
}

func TrackStatistics() {
	fmt.Println("Collecting statistics...")

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return
	}

	dt := time.Now().In(loc)

	currentDate := dt.Format("02-01-2006")

	listOfUsers, err := db.GetListOfRegisteredUsers()
	if err != nil {
		return
	}

	for _, email := range listOfUsers.Users {
		trades, err := GetPortfolio(email)
		for err != nil {
			time.Sleep(1 * time.Second)

			trades, err = GetPortfolio(email)
		}

		statistics, err := db.GetStatisticsFromTableForUser(email)
		for err != nil {
			time.Sleep(1 * time.Second)

			statistics, err = db.GetStatisticsFromTableForUser(email)
		}

		statistics.AccountValue = append(statistics.AccountValue, models.AccountValueStatus{
			Date:        currentDate,
			AcountValue: trades.AccountValue,
		})

		statistics.Credits = append(statistics.Credits, models.CreditsStatus{
			Date:    currentDate,
			Credits: trades.Credits,
		})

		amountOfOwnedStocks := int64(0)

		for _, hl := range trades.HoldLong {
			amountOfOwnedStocks += hl.Amount
		}

		for _, hs := range trades.HoldShort {
			amountOfOwnedStocks += hs.Amount
		}

		statistics.OwnedStocksAmount = append(statistics.OwnedStocksAmount, models.OwnedStocksAmountStatus{
			Date:   currentDate,
			Amount: amountOfOwnedStocks,
		})

		statistics.HoldStocks = append(statistics.HoldStocks, models.HoldStocksStatus{
			Date:   currentDate,
			HLong:  trades.HoldLong,
			HShort: trades.HoldShort,
		})

		err = db.PutStatisticsInTheTable(email, statistics)
		if err != nil {
			return
		}
	}

	fmt.Println("Finished collecting statistics.")
}

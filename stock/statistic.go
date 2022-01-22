package stock

import (
	"fmt"
	"log"
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

	totalGainLossTable := ReturnStatTable{
		Name: "Total gain/loss",
		X:    []string{},
		Y:    []float32{},
	}

	for _, gain := range stats.TotalGain {
		totalGainLossTable.X = append(amountOwnedStockTable.X, gain.Date)
		totalGainLossTable.Y = append(amountOwnedStockTable.Y, gain.Gain)
	}

	return []ReturnStatTable{
		accountValueTable,
		creditsTable,
		amountOwnedStockTable,
		totalGainLossTable,
	}, nil
}

func TrackStatistics() {
	log.Println("Collecting statistics...")

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
		portfolio, err := GetPortfolio(email)
		for err != nil {
			time.Sleep(30 * time.Second)

			portfolio, err = GetPortfolio(email)
		}

		statistics, err := db.GetStatisticsFromTableForUser(email)
		for err != nil {
			time.Sleep(1 * time.Second)

			statistics, err = db.GetStatisticsFromTableForUser(email)
		}

		statistics.AccountValue = append(statistics.AccountValue, models.AccountValueStatus{
			Date:        currentDate,
			AcountValue: portfolio.AccountValue,
		})

		statistics.Credits = append(statistics.Credits, models.CreditsStatus{
			Date:    currentDate,
			Credits: portfolio.Credits,
		})

		amountOfOwnedStocks := int64(0)

		for _, hl := range portfolio.HoldLong {
			amountOfOwnedStocks += hl.Amount
		}

		for _, hs := range portfolio.HoldShort {
			amountOfOwnedStocks += hs.Amount
		}

		statistics.OwnedStocksAmount = append(statistics.OwnedStocksAmount, models.OwnedStocksAmountStatus{
			Date:   currentDate,
			Amount: amountOfOwnedStocks,
		})

		statistics.HoldStocks = append(statistics.HoldStocks, models.HoldStocksStatus{
			Date:   currentDate,
			HLong:  portfolio.HoldLong,
			HShort: portfolio.HoldShort,
		})

		totalGain := float32(0)

		for _, hs := range portfolio.HoldShort {
			totalGain += hs.Gain
		}

		for _, hl := range portfolio.HoldLong {
			totalGain += hl.Gain
		}

		statistics.TotalGain = append(statistics.TotalGain, models.TotalGain{
			Date: currentDate,
			Gain: totalGain,
		})

		err = db.PutStatisticsInTheTable(email, statistics)
		if err != nil {
			return
		}
	}

	fmt.Println("Finished collecting statistics.")
}

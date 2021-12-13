package stock

import (
	"fmt"
	"stockx-backend/db"
	"stockx-backend/db/models"
	"time"

	"github.com/robfig/cron/v3"
)

func StartTrackingStatistics() error {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}

	cronRunner := cron.New(cron.WithLocation(loc))

	_, err = cronRunner.AddFunc("0 16 * * 1-5", collectInformation)
	if err != nil {
		return err
	}

	cronRunner.Start()

	return nil
}

func collectInformation() {
	fmt.Println("Collecting statistics...")

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return
	}

	dt := time.Now().In(loc)

	currentDate := dt.Format("01-02-2006")

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

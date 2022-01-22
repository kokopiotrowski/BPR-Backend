package jobs

import (
	"stockx-backend/external/stockapi"
	"stockx-backend/stock"
	"stockx-backend/users"
	"time"

	"github.com/robfig/cron/v3"
)

func Start() error {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}

	cronRunner := cron.New(cron.WithLocation(loc))

	_, err = cronRunner.AddFunc("0 16 * * 1-5", stock.TrackStatistics) // every weekday at 4pm
	if err != nil {
		return err
	}

	_, err = cronRunner.AddFunc("0 17 * * 1-5", users.TrackRanking) // every weekday at 5pm
	if err != nil {
		return err
	}

	_, err = cronRunner.AddFunc("35 9 * * 1-5", stockapi.ReloadStockInfo) // 5 minutes after opening stock market reload information about each symbol
	if err != nil {
		return err
	}

	cronRunner.Start()

	return nil
}

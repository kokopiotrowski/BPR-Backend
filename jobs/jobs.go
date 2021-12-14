package jobs

import (
	"stockx-backend/stock"
	"stockx-backend/users"
	"time"

	"github.com/robfig/cron/v3"
)

func StartSteve() error {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}

	cronRunner := cron.New(cron.WithLocation(loc))

	_, err = cronRunner.AddFunc("0 16 * * 1-5", stock.TrackStatistics)
	if err != nil {
		return err
	}

	_, err = cronRunner.AddFunc("0 17 * * 1-5", users.TrackRanking)
	if err != nil {
		return err
	}

	cronRunner.Start()

	return nil
}

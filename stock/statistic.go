package stock

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func StartTrackingStatistics() error {
	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return err
	}

	cronRunner := cron.New(cron.WithLocation(loc))

	_, err = cronRunner.AddFunc("@every 1m", collectInformation)
	if err != nil {
		return err
	}

	cronRunner.Start()

	return nil
}

func collectInformation() {
	fmt.Printf("Collecting information for statistics...")
}

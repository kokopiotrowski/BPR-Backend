package users

import (
	"fmt"
	"stockx-backend/db"
	"stockx-backend/db/models"
	"time"
)

func TrackRanking() {
	fmt.Println("Generating ranking...")

	listOfUsers, err := db.GetListOfRegisteredUsers()
	for err != nil {
		time.Sleep(1 * time.Second)

		listOfUsers, err = db.GetListOfRegisteredUsers()
	}

	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return
	}

	dt := time.Now().In(loc)

	currentDate := dt.Format("02-01-2006")
	yesterday := dt.AddDate(0, 0, -1).Format("02-01-2006")

	ranking := models.Rankings{
		Region: "Europe/Copenhagen",
		Date:   currentDate,
	}

	for _, email := range listOfUsers.Users {
		var statistics models.Statistics

		addedToStats := false

		statistics, err = db.GetStatisticsFromTableForUser(email)
		for err != nil {
			time.Sleep(1 * time.Second)

			statistics, err = db.GetStatisticsFromTableForUser(email)
		}

		if statistics.AccountValue[len(statistics.AccountValue)-1].Date == yesterday {
			usersAccountValue := models.UserAccountValue{
				Username:     statistics.Email,
				AccountValue: statistics.AccountValue[len(statistics.AccountValue)-1].AcountValue,
			}

			for i, acv := range ranking.AccountValueRating {
				if acv.AccountValue <= statistics.AccountValue[len(statistics.AccountValue)-1].AcountValue {
					ranking.AccountValueRating = append(ranking.AccountValueRating[:i+1], ranking.AccountValueRating[i:]...)
					ranking.AccountValueRating[i] = usersAccountValue
					addedToStats = true
				}
			}

			if !addedToStats {
				ranking.AccountValueRating = append(ranking.AccountValueRating, usersAccountValue)
			}
		}
	}

	for _, email := range listOfUsers.Users {
		var trades models.Trades

		trades, err = db.GetTradesFromTableForUser(email)
		for err != nil {
			time.Sleep(1 * time.Second)

			trades, err = db.GetTradesFromTableForUser(email)
		}

		for _, holdL := range trades.HoldLong {
			for i, srank := range ranking.MostPopularStock {
				if srank.Symbol == holdL.Symbol {
					srank.Amount += holdL.Amount
					break
				}

				if i == len(ranking.MostPopularStock)-1 {
					ranking.MostPopularStock = append(ranking.MostPopularStock, models.MostPopularStock{
						Symbol: holdL.Symbol,
						Amount: holdL.Amount,
					})

					break
				}
			}
		}
	}

	ranking.MostPopularStock = sortPopularSymbols(ranking.MostPopularStock)

	err = db.PutRankingsInTheTable(ranking)

	for err != nil {
		time.Sleep(1 * time.Second)

		err = db.PutRankingsInTheTable(ranking)
	}

	fmt.Println("Finished ranking generation.")
}

func sortPopularSymbols(stocks []models.MostPopularStock) []models.MostPopularStock {
	for i := 0; i < len(stocks)-1; i++ {
		if stocks[i].Amount > stocks[i+1].Amount {
			stocks[i], stocks[i+1] = stocks[i+1], stocks[i]
		}
	}

	return stocks
}

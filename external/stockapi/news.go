package stockapi

import (
	"context"
	"stockx-backend/reserr"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func GetNews() ([]finnhub.MarketNews, error) {
	finnhubClient := finnhub.NewAPIClient(FinnhubConfiguration).DefaultApi

	news, response, err := finnhubClient.MarketNews(context.Background()).Category("forex").Execute()
	if err != nil {
		return []finnhub.MarketNews{}, reserr.Internal("failed to retrieve news", err, "News from external api are either not available or the api is down")
	}
	defer response.Body.Close()

	return news, nil
}

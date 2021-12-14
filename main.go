package main

import (
	"log"
	"net/http"
	"stockx-backend/conf"
	"stockx-backend/db"
	"stockx-backend/external/stockapi"
	"stockx-backend/jobs"
	"stockx-backend/network"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/rs/cors"
)

func main() {
	flags, err := conf.ParseFlags()
	if err != nil {
		panic(1)
	}

	config, err := conf.ReadConfig()
	if err != nil {
		panic(1)
	}

	stockapi.FinnhubConfiguration = finnhub.NewConfiguration()
	stockapi.FinnhubConfiguration.AddDefaultHeader("X-Finnhub-Token", config.StockAPI.Key)

	db.InitDB(flags.IsProduction)

	if flags.IsProduction {
		log.Printf("PRODUCTION Server started on port %v", config.Server.ProdPort)
	} else {
		log.Printf("DEV Server started on port %v", config.Server.DevPort)
	}

	router := network.NewRouter()
	handler := cors.AllowAll().Handler(router)

	go stockapi.StartListening(config.StockAPI.Key)

	err = jobs.StartSteve()
	if err != nil {
		panic(1)
	}

	if flags.IsProduction {
		log.Fatal(http.ListenAndServe(config.Server.ProdPort, handler))
	} else {
		log.Fatal(http.ListenAndServe(config.Server.DevPort, handler))
	}
}

package main

import (
	"errors"
	"log"
	"net/http"
	"os"
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

	if conf.FlagConf.IsLoggingOn {
		var f *os.File

		f, err = os.OpenFile("logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	}

	stockapi.FinnhubConfiguration = finnhub.NewConfiguration()
	stockapi.FinnhubConfiguration.AddDefaultHeader("X-Finnhub-Token", config.StockAPI.Key)

	db.ListAvailableTables(flags.IsProduction)

	if flags.IsProduction {
		log.Printf("PRODUCTION Server started on port %v", config.Server.ProdPort)
	} else {
		log.Printf("DEV Server started on port %v", config.Server.DevPort)
	}

	router := network.NewRouter()
	handler := cors.AllowAll().Handler(router)

	go stockapi.StartListening(config.StockAPI.Key)

	err = jobs.Start()
	if err != nil {
		panic(1)
	}

	if flags.IsProduction {
		log.Fatal(http.ListenAndServe(config.Server.ProdPort, handler))
	} else {
		log.Fatal(http.ListenAndServe(config.Server.DevPort, handler))
	}

	defer rec()
}

func rec() {
	var err error

	if r := recover(); r != nil {
		log.Println("Recovered in f", r)
		// find out exactly what the error was and set err
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			// Fallback err (per specs, error strings should be lowercase w/o punctuation
			err = errors.New("unknown panic")
		}
	}

	log.Print(err)
}

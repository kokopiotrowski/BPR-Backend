package main

import (
	"log"
	"net/http"
	"stockx-backend/conf"
	"stockx-backend/db"
	"stockx-backend/network"

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

	db.InitDB(flags.IsProduction)

	if flags.IsProduction {
		log.Printf("PRODUCTION Server started on port %v", config.Server.ProdPort)
	} else {
		log.Printf("DEV Server started on port %v", config.Server.DevPort)
	}

	router := network.NewRouter()
	handler := cors.AllowAll().Handler(router)

	if flags.IsProduction {
		log.Fatal(http.ListenAndServe(config.Server.ProdPort, handler))
	} else {
		log.Fatal(http.ListenAndServe(config.Server.DevPort, handler))
	}
}

package main

import (
	"log"
	"net/http"
	"stockx-backend/conf"
	"stockx-backend/network"
)

const (
	ProductionPort = ":80"
	DevPort        = ":3000"
)

func main() {
	co, err := conf.Init()
	if err != nil {
		panic(1)
	}

	log.Printf("Server started")

	router := network.NewRouter()
	if co.IsProduction {
		log.Fatal(http.ListenAndServe(ProductionPort, router))
	} else {
		log.Fatal(http.ListenAndServe(DevPort, router))
	}
}

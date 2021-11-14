package main

import (
	"log"
	"net/http"

	"./network"
)

func main() {
	log.Printf("Server started")

	router := network.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}

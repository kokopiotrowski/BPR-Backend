package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
	fmt.Print("Server running")

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Siema, tu Konrad :)"))
	if err != nil {
		panic(err)
	}
}

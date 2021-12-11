package api

import (
	"net/http"
	"stockx-backend/auth"
	"stockx-backend/external/stockapi"
	"stockx-backend/util"
)

func GetStockSymbols(w http.ResponseWriter, r *http.Request) {
	var username string
	if auth.CheckIfAuthorized(w, r, &username) {
		var query string
		if util.DecodeFormValueAsString(w, r, "q", &query) {
			symbols, err := stockapi.GetStockSymbols(query)
			if err != nil {
				util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to retrieve stock symbols from stockapi", err)
				return
			}

			util.RespondWithJSON(w, r, http.StatusOK, symbols, nil)
		}
	}
}

func GetStockCandles(w http.ResponseWriter, r *http.Request) {
	var username string
	if auth.CheckIfAuthorized(w, r, &username) {
		var symbol string
		if util.DecodeFormValueAsString(w, r, "symbol", &symbol) {
			candleData, err := stockapi.GetStockCandles(symbol)
			if err != nil {
				util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to retrieve stock candle data from stockapi", err)
				return
			}

			util.RespondWithJSON(w, r, http.StatusOK, candleData, nil)
		}
	}
}

func GetCurrentStockPrice(w http.ResponseWriter, r *http.Request) {

}

func GetCompanyInfo(w http.ResponseWriter, r *http.Request) {
	var username string
	if auth.CheckIfAuthorized(w, r, &username) {
		var symbol string
		if util.DecodeFormValueAsString(w, r, "symbol", &symbol) {
			companyInfo, err := stockapi.GetCompanyInfo(symbol)
			if err != nil {
				util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to retrieve company information from stockapi", err)
				return
			}

			util.RespondWithJSON(w, r, http.StatusOK, companyInfo, nil)
		}
	}
}

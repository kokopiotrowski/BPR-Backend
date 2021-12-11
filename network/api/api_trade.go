package api

import (
	"net/http"
	"stockx-backend/auth"
	"stockx-backend/db/models"
	"stockx-backend/stock"
	"stockx-backend/util"
)

func BuyStockLong(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		var bs models.BoughtStock
		if util.DecodeBodyAsJSON(w, r, &bs) {
			err := stock.BuyStock(email, bs)
			if err != nil {
				util.RespondWithJSON(w, r, http.StatusInternalServerError, "", err)
				return
			}

			util.RespondWithJSON(w, r, http.StatusOK, "", nil)
		}
	}
}

func SellStockLong(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		var ss models.SoldStock
		if util.DecodeBodyAsJSON(w, r, &ss) {
			err := stock.SellStock(email, ss)
			if err != nil {
				util.RespondWithJSON(w, r, http.StatusInternalServerError, "", err)
				return
			}

			util.RespondWithJSON(w, r, http.StatusOK, "", nil)
		}
	}
}

func BuyStockShort(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		var ss models.ShortStock
		if util.DecodeBodyAsJSON(w, r, &ss) {
			err := stock.ShortStock(email, ss)
			if err != nil {
				util.RespondWithJSON(w, r, http.StatusInternalServerError, "", err)
				return
			}

			util.RespondWithJSON(w, r, http.StatusOK, "", nil)
		}
	}
}

func BuyToCover(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		var btc models.BoughtToCover
		if util.DecodeBodyAsJSON(w, r, &btc) {
			err := stock.BuyToCover(email, btc)
			if err != nil {
				util.RespondWithJSON(w, r, http.StatusInternalServerError, "", err)
				return
			}

			util.RespondWithJSON(w, r, http.StatusOK, "", nil)
		}
	}
}

package api

import (
	"net/http"

	"stockx-backend/auth"
	"stockx-backend/util"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials auth.Credentials
	if util.DecodeBodyAsJSON(w, r, &credentials) {
		token, err := auth.LogIn(credentials)
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusUnauthorized, "NOT STONKS", err)
			return
		}

		util.RespondWithJSON(w, r, http.StatusOK, `{"Authorization":"`+token.Token+`"}`, nil)

		return
	}

	util.RespondWithJSON(w, r, http.StatusInternalServerError, "Received login body wasn't parsed", nil)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerBody auth.Register
	if util.DecodeBodyAsJSON(w, r, &registerBody) {
		_, err := auth.RegisterUser(registerBody)
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, "NOT STONKS", err)
			return
		}

		util.RespondWithJSON(w, r, http.StatusOK, "User registered - check email", nil)

		return
	}

	util.RespondWithJSON(w, r, http.StatusInternalServerError, "Received register body wasn't parsed", nil)
}

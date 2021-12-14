package api

import (
	"fmt"
	"net/http"
	"stockx-backend/auth"
	"stockx-backend/stock"
	"stockx-backend/users"
	"stockx-backend/util"
)

func UserSettingsGet(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		user, err := users.GetUser(email)

		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to retrieve user", err)
			return
		}

		util.RespondWithJSON(w, r, http.StatusOK, user, nil)
	}
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		err := users.DeleteUser(email)
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to delete user", err)
			return
		}

		util.RespondWithJSON(w, r, http.StatusOK, "User deleted", nil)
	}
}

func UserChangePassword(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		var userChangePassword users.UpdatedPassword
		if util.DecodeBodyAsJSON(w, r, &userChangePassword) {
			err := users.ChangePassword(email, userChangePassword)
			if err != nil {
				util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to change user's password", err)
				return
			}

			util.RespondWithJSON(w, r, http.StatusOK, "User password changed", nil)
		}
	}
}

func GetUserPortfolio(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		portfolio, err := stock.GetPortfolio(email)
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to retrieve user's portfolio", err)
			return
		}

		util.RespondWithJSON(w, r, http.StatusOK, portfolio, nil)
	}
}

func GetUserStatistics(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		stats, err := stock.GetStatisticsForUser(email)
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to retrieve user's statistics", err)
			return
		}

		util.RespondWithJSON(w, r, http.StatusOK, stats, nil)
	}
}

func ConfirmAccount(w http.ResponseWriter, r *http.Request) {
	var token string
	if util.DecodeFormValueAsString(w, r, "t", &token) {
		page, err := users.ConfirmUserAccount(token)
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to retrieve user's statistics", err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, page)
	}
}

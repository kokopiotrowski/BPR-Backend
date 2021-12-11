package api

import (
	"net/http"
	"stockx-backend/auth"
	"stockx-backend/external/stockapi"
	"stockx-backend/util"
)

func GetNews(w http.ResponseWriter, r *http.Request) {
	var username string
	if auth.CheckIfAuthorized(w, r, &username) {
		news, err := stockapi.GetNews()
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to get news", err)
			return
		}

		util.RespondWithJSON(w, r, http.StatusOK, news, nil)
	}
}

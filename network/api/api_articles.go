package api

import (
	"net/http"
	"stockx-backend/external/stockapi"
	"stockx-backend/util"
)

func GetNews(w http.ResponseWriter, r *http.Request) {
	// err := auth.TokenValid(r)
	// if err != nil {
	// 	util.RespondWithJSON(w, r, http.StatusUnauthorized, "Unauthenticated", err)
	// 	return
	// }
	news, err := stockapi.GetNews()
	if err != nil {
		util.RespondWithJSON(w, r, http.StatusInternalServerError, "failed to get news", err)
		return
	}

	util.RespondWithJSON(w, r, http.StatusOK, news, nil)
}

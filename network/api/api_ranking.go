package api

import (
	"net/http"
	"stockx-backend/auth"
	"stockx-backend/users"
	"stockx-backend/util"
)

func GetRanking(w http.ResponseWriter, r *http.Request) {
	var email string
	if auth.CheckIfAuthorized(w, r, &email) {
		ranking, err := users.GetRanking()
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, nil, err)
			return
		}

		util.RespondWithJSON(w, r, http.StatusOK, ranking, nil)
	}
}

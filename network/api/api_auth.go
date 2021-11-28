package api

import (
	"net/http"

	"stockx-backend/util"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJSON(w, r, http.StatusOK, "User login", nil)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJSON(w, r, http.StatusOK, "User register", nil)
}

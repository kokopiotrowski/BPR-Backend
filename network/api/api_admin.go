package api

import (
	"io/ioutil"
	"net/http"
	"stockx-backend/auth"
	"stockx-backend/users"
	"stockx-backend/util"
)

func GetListOfRegisteredUsers(w http.ResponseWriter, r *http.Request) {
	if admin, err := auth.CheckIfAdmin(r); admin && err == nil {
		listOfRegisteredUsers, err := users.GetListOfUserData()
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, nil, err)
			return
		}

		util.RespondWithJSON(w, r, http.StatusOK, listOfRegisteredUsers, nil)

		return
	}

	util.RespondWithJSON(w, r, http.StatusForbidden, nil, nil)
}

func DeleteRegisteredUser(w http.ResponseWriter, r *http.Request) {
	if admin, err := auth.CheckIfAdmin(r); admin && err == nil {
		var email string
		if util.DecodeFormValueAsString(w, r, "email", &email) {
			err := users.DeleteUser(email)
			if err != nil {
				util.RespondWithJSON(w, r, http.StatusInternalServerError, nil, err)
				return
			}

			util.RespondWithJSON(w, r, http.StatusOK, nil, nil)

			return
		}
	}

	util.RespondWithJSON(w, r, http.StatusForbidden, nil, nil)
}

func GetLogs(w http.ResponseWriter, r *http.Request) {
	if admin, err := auth.CheckIfAdmin(r); admin && err == nil {
		fileBytes, err := ioutil.ReadFile("logfile.log")
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")

		_, err = w.Write(fileBytes)
		if err != nil {
			util.RespondWithJSON(w, r, http.StatusInternalServerError, nil, err)
			return
		}

		return
	}

	util.RespondWithJSON(w, r, http.StatusForbidden, nil, nil)
}

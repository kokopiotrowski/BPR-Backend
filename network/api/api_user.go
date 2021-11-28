package api

import (
	"net/http"

	"../../util"
)

func UserSettingsGet(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJSON(w, r, http.StatusOK, "User settings", nil)
}

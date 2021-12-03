package network

import (
	"net/http"
	"stockx-backend/logger"
	"stockx-backend/network/api"
	"stockx-backend/util"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJSON(w, r, http.StatusOK, "INDEX", nil)
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"",
		Index,
	},

	Route{
		"LoginUser",
		strings.ToUpper("Post"),
		"/auth/login",
		api.LoginUser,
	},

	Route{
		"RegisterUser",
		strings.ToUpper("Post"),
		"/auth/register",
		api.RegisterUser,
	},

	Route{
		"UserSettingsGet",
		strings.ToUpper("Get"),
		"/user/settings",
		api.UserSettingsGet,
	},

	Route{
		"GetNews",
		strings.ToUpper("Get"),
		"/news",
		api.GetNews,
	},
}

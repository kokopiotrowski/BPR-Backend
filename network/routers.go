package network

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"../logger"
	"./api"
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
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/k0k0piotrowski/StockX/1.0.0/",
		Index,
	},

	Route{
		"LoginUser",
		strings.ToUpper("Post"),
		"/k0k0piotrowski/StockX/1.0.0/auth/login",
		api.LoginUser,
	},

	Route{
		"RegisterUser",
		strings.ToUpper("Post"),
		"/k0k0piotrowski/StockX/1.0.0/auth/register",
		api.RegisterUser,
	},

	Route{
		"UserSettingsGet",
		strings.ToUpper("Get"),
		"/k0k0piotrowski/StockX/1.0.0/user/settings",
		api.UserSettingsGet,
	},
}

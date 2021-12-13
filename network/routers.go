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
		"/index",
		Index,
	},

	Route{
		"Websocket",
		strings.ToUpper("Get"),
		"/ws",
		api.WsEndpoint,
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
		"UserSettingsDelete",
		strings.ToUpper("Delete"),
		"/user/settings/delete",
		api.UserDelete,
	},

	Route{
		"UserUpdatePassword",
		strings.ToUpper("Post"),
		"/user/settings/password",
		api.UserChangePassword,
	},

	Route{
		"UserGetPortfolio",
		strings.ToUpper("Get"),
		"/user/portfolio",
		api.GetUserPortfolio,
	},

	Route{
		"GetStockSymbols",
		strings.ToUpper("Get"),
		"/stock/symbols", // form values - q(string)
		api.GetStockSymbols,
	},

	Route{
		"GetCandlesOfStocks",
		strings.ToUpper("Get"),
		"/stock/candles", // form values - symbol(string)
		api.GetStockCandles,
	},

	Route{
		"GetSymbolInfo",
		strings.ToUpper("Get"),
		"/stock/company/info", // form values - symbol(string)
		api.GetSymbolInfo,
	},

	Route{
		"GetCurrentStockPrice",
		strings.ToUpper("Get"),
		"/stock/current", // form values - symbol(string)
		api.GetCurrentStockPrice,
	},

	Route{
		"BuyStockLong",
		strings.ToUpper("Post"),
		"/stock/long/buy",
		api.BuyStockLong,
	},

	Route{
		"SellStockLong",
		strings.ToUpper("Post"),
		"/stock/long/sell",
		api.SellStockLong,
	},

	Route{
		"BuyStockShort",
		strings.ToUpper("Post"),
		"/stock/short/sell",
		api.BuyStockShort,
	},
	Route{
		"BuyStockLong",
		strings.ToUpper("Post"),
		"/stock/short/cover",
		api.BuyToCover,
	},

	Route{
		"GetNews",
		strings.ToUpper("Get"),
		"/news",
		api.GetNews,
	},
}

package routes

import (
	"net/http"

	"../jwtauth"
	"../util"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter()
	for _, route := range composeAllRoutes() {
		var handler http.Handler

		handler = route.HandlerFunc
		if route.Pattern == "/signin" {
			handler = util.Logger(handler, route.Name)
		}

		if route.Pattern != "/signin" {
			handler = util.Logger(jwtauth.AuthMiddleware(handler), route.Name)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler).
			Queries(route.Queries...)

	}

	return router
}

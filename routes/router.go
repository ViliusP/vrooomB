package routes

import (
	"net/http"

	"../util"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter()
	for _, route := range composeAllRoutes() {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = util.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler).
			Queries(route.Queries...)

	}

	return router
}

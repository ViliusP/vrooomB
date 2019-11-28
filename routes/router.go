package routes

import (
	"net/http"

	"strconv"

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
			handler = util.Logger(jwtauth.AuthMiddleware(VarsCheckMiddleware(handler)), route.Name)
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

func VarsCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range mux.Vars(r) {
			if parsedV, err := strconv.Atoi(v); err != nil || parsedV < 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})

}

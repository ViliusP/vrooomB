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
		var fakeHandler http.Handler
		var handler http.Handler
		fakeRoute := Route{
			"CORS",
			http.MethodOptions,
			route.Pattern,
			corsHandler,
			[]string{},
		}
		fakeHandler = fakeRoute.HandlerFunc
		handler = route.HandlerFunc
		if route.Pattern == "/signin" {
			handler = util.Logger(AllowCORSHeader(handler), route.Name)
			fakeHandler = util.Logger(AllowCORSHeader(fakeHandler), fakeRoute.Name)
		}

		if route.Pattern != "/signin" {
			handler = util.Logger(jwtauth.AuthMiddleware(VarsCheckMiddleware(AllowCORSHeader(handler))), route.Name)
			fakeHandler = util.Logger(jwtauth.AuthMiddleware(VarsCheckMiddleware(AllowCORSHeader(fakeHandler))), fakeRoute.Name)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler).
			Queries(route.Queries...)

		router.
			Methods(fakeRoute.Method).
			Path(fakeRoute.Pattern).
			Name(fakeRoute.Name + " " + route.Name).
			Handler(fakeHandler).
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

func AllowCORSHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Allow-Origin")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		next.ServeHTTP(w, r)
	})
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

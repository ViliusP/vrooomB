package util

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%8s\t%25s\t%20s\t%10s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

package logger

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		// var email string
		// if auth.CheckIfAuthorized(w, r, &email) {
		// 	log.Printf(
		// 		"%s %s %s %s %s",
		// 		r.Method,
		// 		email,
		// 		r.RequestURI,
		// 		name,
		// 		time.Since(start),
		// 	)
		// } else

		if r.RequestURI != "/index" {
			log.Printf(
				"%s %s %s %s",
				r.Method,
				r.RequestURI,
				name,
				time.Since(start),
			)
		}
	})
}

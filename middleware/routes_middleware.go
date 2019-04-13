package middleware

import (
	"log"
	"net/http"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

//WithContentType -
func WithContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Printf("Setting content-type to json")
		next.ServeHTTP(w, r)
	}
}

//WithCORS -
func WithCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		log.Printf("Setting CORS to allow *")
		next.ServeHTTP(w, r)
	}
}

//WithLogging -
func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Logged connection from %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

func withTracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tracing request for %s", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

// ChainMiddleware composes a list of middleware functions such that they are
// called in a stack, i.e. a(b(c()))
func ChainMiddleware(ms ...middleware) middleware {
	return func(route http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			currentMiddleware := route // route is the  main request handler
			for _, m := range ms {
				currentMiddleware = m(currentMiddleware)
			}

			currentMiddleware(w, r)
		}
	}
}

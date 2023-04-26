package server

import "net/http"

const (
	authHeader = "X-API-Key"
)

// authenticate is a middleware that checks for a valid API key in the request.
func authenticate(keys map[string]struct{}, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get(authHeader)
		if len(key) == 0 {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}
		if _, ok := keys[key]; !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

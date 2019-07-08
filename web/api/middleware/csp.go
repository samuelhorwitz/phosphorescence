package middleware

import (
	"fmt"
	"net/http"
)

func CSP(phosphorOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'none';base-uri 'none';form-action 'none';frame-ancestors %s;block-all-mixed-content;navigate-to 'self' https://accounts.spotify.com;", phosphorOrigin))
			next.ServeHTTP(w, r)
		})
	}
}

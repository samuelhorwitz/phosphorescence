package middleware

import (
	"errors"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"net/http"
)

func SpotifyLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := r.Context().Value(SessionContextKey).(*session.Session)
		if !ok {
			common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
			return
		}
		lmtErr := tollbooth.LimitByKeys(spotifyLimiter, []string{sess.SpotifyID})
		if lmtErr != nil {
			spotifyLimiter.ExecOnLimitReached(w, r)
			common.Fail(w, fmt.Errorf("Phosphorescence Spotify API rate limiting hit: %s", lmtErr.Message), lmtErr.StatusCode)
			return
		}
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"errors"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"net/http"
)

func SpotifyLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spotifyID, ok := r.Context().Value(SpotifyIDContextKey).(string)
		if !ok {
			common.Fail(w, errors.New("No Spotify ID on request context"), http.StatusInternalServerError)
			return
		}
		lmtErr := tollbooth.LimitByKeys(spotifyLimiter, []string{spotifyID})
		if lmtErr != nil {
			spotifyLimiter.ExecOnLimitReached(w, r)
			common.Fail(w, fmt.Errorf("Phosphorescence Spotify API rate limiting hit: %s", lmtErr.Message), lmtErr.StatusCode)
			return
		}
		next.ServeHTTP(w, r)
	})
}

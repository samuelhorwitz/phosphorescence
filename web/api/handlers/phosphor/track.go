package phosphor

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

func GetTracks(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	getTracks(w, r, sess.SpotifyCountry, chi.URLParam(r, "trackIDs"))
}

func GetTracksUnauthenticated(w http.ResponseWriter, r *http.Request) {
	getTracks(w, r, chi.URLParam(r, "region"), chi.URLParam(r, "trackIDs"))
}

func getTracks(w http.ResponseWriter, r *http.Request, region, trackIDsStr string) {
	trackIDs := strings.Split(trackIDsStr, ",")
	tracks, err := models.GetTracks(r.Context(), region, trackIDs)
	if err != nil {
		code := http.StatusInternalServerError
		if _, ok := err.(models.TrackNotFoundInRegionError); ok {
			code = http.StatusNotFound
		}
		common.Fail(w, fmt.Errorf("Could not get track: %s", err), code)
		return
	}
	common.JSON(w, map[string]interface{}{"tracks": tracks})
}

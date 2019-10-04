package phosphor

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/handlers"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

func GetAlbum(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	getAlbum(w, r, sess.SpotifyCountry, chi.URLParam(r, "albumID"))
}

func GetAlbumUnauthenticated(w http.ResponseWriter, r *http.Request) {
	getAlbum(w, r, chi.URLParam(r, "region"), chi.URLParam(r, "albumID"))
}

func getAlbum(w http.ResponseWriter, r *http.Request, region, albumID string) {
	if albumID == "" {
		common.Fail(w, errors.New("Must include album ID"), http.StatusBadRequest)
		return
	}
	album, err := models.GetAlbumTracks(r.Context(), region, albumID)
	if err != nil {
		code := http.StatusInternalServerError
		if httpErr, ok := err.(handlers.HTTPError); ok {
			code = httpErr.Code
		}
		common.Fail(w, fmt.Errorf("Could not get album: %s", err), code)
		return
	}
	common.JSON(w, map[string]interface{}{"album": album})
}

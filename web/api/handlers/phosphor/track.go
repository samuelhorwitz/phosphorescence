package phosphor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/handlers"
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

func GetTrackPreviews(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	getTrackPreviews(w, r, sess.SpotifyCountry, chi.URLParam(r, "trackIDs"))
}

func GetTrackPreviewsUnauthenticated(w http.ResponseWriter, r *http.Request) {
	getTrackPreviews(w, r, chi.URLParam(r, "region"), chi.URLParam(r, "trackIDs"))
}

func getTracks(w http.ResponseWriter, r *http.Request, region, trackIDsStr string) {
	tracks, err := getTracksFromModel(r.Context(), region, trackIDsStr)
	if err != nil {
		code := http.StatusInternalServerError
		if httpErr, ok := err.(handlers.HTTPError); ok {
			code = httpErr.Code
		}
		common.Fail(w, err, code)
		return
	}
	common.JSON(w, map[string]interface{}{"tracks": tracks})
}

func getTrackPreviews(w http.ResponseWriter, r *http.Request, region, trackIDsStr string) {
	tracks, err := getTracksFromModel(r.Context(), region, trackIDsStr)
	if err != nil {
		code := http.StatusInternalServerError
		if httpErr, ok := err.(handlers.HTTPError); ok {
			code = httpErr.Code
		}
		common.Fail(w, fmt.Errorf("Could not get track previews: %s", err), code)
		return
	}
	trackPreviews := make(map[string]string)
	for _, track := range tracks {
		if track.Track.PreviewURL != "" {
			trackPreviews[track.ID] = track.Track.PreviewURL
		}
	}
	common.JSON(w, map[string]interface{}{"trackPreviews": trackPreviews})
}

func getTracksFromModel(ctx context.Context, region, trackIDsStr string) ([]*models.SpotifyTrackEnvelope, error) {
	if trackIDsStr == "" {
		return nil, handlers.NewHTTPError(errors.New("No track IDs specified"), http.StatusBadRequest)
	}
	trackIDs := strings.Split(trackIDsStr, ",")
	tracks, err := models.GetTracks(ctx, region, trackIDs)
	if err != nil {
		code := http.StatusInternalServerError
		if _, ok := err.(models.TrackNotFoundInRegionError); ok {
			code = http.StatusNotFound
		}
		return nil, handlers.NewHTTPError(fmt.Errorf("Could not get tracks: %s", err), code)
	}
	return tracks, nil
}

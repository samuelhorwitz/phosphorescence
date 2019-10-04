package phosphor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/handlers"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

var errNoTracks = errors.New("Must include at least one track")

func GetPlaylist(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	getPlaylist(w, r, sess.SpotifyCountry, chi.URLParam(r, "playlistID"))
}

func GetPlaylistUnauthenticated(w http.ResponseWriter, r *http.Request) {
	getPlaylist(w, r, chi.URLParam(r, "region"), chi.URLParam(r, "playlistID"))
}

func getPlaylist(w http.ResponseWriter, r *http.Request, region, playlistID string) {
	if playlistID == "" {
		common.Fail(w, errors.New("Must include playlist ID"), http.StatusBadRequest)
		return
	}
	playlist, err := models.GetPlaylist(r.Context(), region, playlistID)
	if err != nil {
		code := http.StatusInternalServerError
		if httpErr, ok := err.(handlers.HTTPError); ok {
			code = httpErr.Code
		}
		common.Fail(w, fmt.Errorf("Could not get playlist: %s", err), code)
		return
	}
	common.JSON(w, map[string]interface{}{"playlist": playlist})
}

func CreatePrivatePlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, err := createPlaylist(r)
	if err != nil {
		code := http.StatusInternalServerError
		if httpErr, ok := err.(handlers.HTTPError); ok {
			code = httpErr.Code
		}
		common.Fail(w, fmt.Errorf("Could not create playlist: %s", err), code)
		return
	}
	common.JSON(w, map[string]interface{}{"playlist": playlistID})
}

func MakePlaylistOfficial(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "playlistID")
	if playlistID == "" {
		common.Fail(w, errors.New("Must include playlist ID"), http.StatusBadRequest)
		return
	}
	err := models.MakePlaylistPublic(r.Context(), playlistID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not make playlist public: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"playlist": playlistID})
}

func createPlaylist(r *http.Request) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", handlers.NewHTTPError(fmt.Errorf("Could not read request body: %s", err), http.StatusBadRequest)
	}
	var requestBody struct {
		Image            string `json:"image"`
		UTCOffsetMinutes int    `json:"utcOffsetMinutes"`
		Tracks           []struct {
			Name string `json:"name"`
			URI  string `json:"uri"`
		} `json:"tracks"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return "", handlers.NewHTTPError(fmt.Errorf("Could not parse request body: %s", err), http.StatusBadRequest)
	}
	if len(requestBody.Tracks) < 1 {
		return "", handlers.NewHTTPError(errNoTracks, http.StatusBadRequest)
	}
	var trackURIs []string
	for _, track := range requestBody.Tracks {
		trackURIs = append(trackURIs, track.URI)
	}
	playlistID, err := models.CreatePlaylist(r.Context(), requestBody.Tracks[0].Name, requestBody.Image, requestBody.UTCOffsetMinutes, trackURIs)
	if err != nil {
		return "", fmt.Errorf("Failed to create playlist: %s", err)
	}
	return playlistID, nil
}

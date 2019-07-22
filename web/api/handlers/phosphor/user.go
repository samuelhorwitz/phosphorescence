package phosphor

import (
	"errors"
	"fmt"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"net/http"
	"time"
)

type User struct {
	SpotifyID string `json:"spotifyId"`
	Name      string `json:"name"`
	Country   string `json:"country"`
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	spotifyID, ok := r.Context().Value(middleware.SpotifyIDContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify ID on request context"), http.StatusInternalServerError)
		return
	}
	name, ok := r.Context().Value(middleware.SpotifyNameContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify name on request context"), http.StatusInternalServerError)
		return
	}
	country, ok := r.Context().Value(middleware.SpotifyCountryContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify country on request context"), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"user": User{SpotifyID: spotifyID, Name: name, Country: country}})
}

func ListCurrentUserScripts(w http.ResponseWriter, r *http.Request) {
	spotifyID, ok := r.Context().Value(middleware.SpotifyIDContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify ID on request context"), http.StatusInternalServerError)
		return
	}
	count, _ := r.Context().Value(middleware.PageCountContextKey).(uint64)
	if !ok {
		common.Fail(w, errors.New("No page count on request context"), http.StatusInternalServerError)
		return
	}
	from, _ := r.Context().Value(middleware.PageCursorContextKey).(time.Time)
	if !ok {
		common.Fail(w, errors.New("No page cursor on request context"), http.StatusInternalServerError)
		return
	}
	scripts, err := models.GetScriptsBySpotifyUserID(spotifyID, count, from, true)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get scripts for user: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"scripts": scripts})
}

func ListSpotifyDevices(w http.ResponseWriter, r *http.Request) {
	spotifyToken, ok := r.Context().Value(middleware.SpotifyTokenContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify token on request context"), http.StatusInternalServerError)
		return
	}
	body, err := getDevices(r.Context(), spotifyToken)
	if err != nil {
		common.Fail(w, fmt.Errorf("Unable to get devices: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSONRaw(w, body)
}

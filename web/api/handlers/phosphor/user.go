package phosphor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"io/ioutil"
	"net/http"
	"strings"
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
	req, err := http.NewRequestWithContext(r.Context(), "GET", "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not build Spotify devices request: %s", err), http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not make Spotify devices request: %s", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		common.Fail(w, fmt.Errorf("Spotify devices request responded with %d", res.StatusCode), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read Spotify devices response: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSONRaw(w, body)
}

func TransferPlayback(w http.ResponseWriter, r *http.Request) {
	spotifyToken, ok := r.Context().Value(middleware.SpotifyTokenContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify token on request context"), http.StatusInternalServerError)
		return
	}
	deviceID := chi.URLParam(r, "deviceID")
	if deviceID == "" {
		common.Fail(w, errors.New("Must specify device ID"), http.StatusBadRequest)
		return
	}
	playState := strings.ToLower(r.URL.Query().Get("playState"))
	body, err := json.Marshal(struct {
		DeviceIDs []string `json:"device_ids"`
		Play      bool     `json:"play"`
	}{
		DeviceIDs: []string{deviceID},
		Play:      playState == "play",
	})
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not build Spotify device transfer request body: %s", err), http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), "PUT", "https://api.spotify.com/v1/me/player", bytes.NewBuffer(body))
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not build Spotify device transfer request: %s", err), http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not make Spotify device transfer request: %s", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		common.Fail(w, fmt.Errorf("Spotify device transfer request responded with %d", res.StatusCode), http.StatusInternalServerError)
		return
	}
	if playState == "pause" {
		body, err := json.Marshal(struct {
			DeviceID string `json:"device_id"`
		}{
			DeviceID: deviceID,
		})
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify pause request body: %s", err), http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequestWithContext(r.Context(), "PUT", "https://api.spotify.com/v1/me/player/pause", bytes.NewBuffer(body))
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify pause request: %s", err), http.StatusInternalServerError)
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not make Spotify pause request: %s", err), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		if !(res.StatusCode == http.StatusNoContent || res.StatusCode == http.StatusNotFound) {
			common.Fail(w, fmt.Errorf("Spotify pause request responded with %d", res.StatusCode), http.StatusInternalServerError)
			return
		}
	}
	common.JSON(w, map[string]interface{}{"deviceTransfer": true})
}

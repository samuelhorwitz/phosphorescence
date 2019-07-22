package phosphor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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
	var devicesBody json.RawMessage
	err = common.ExponentialBackoff(250*time.Millisecond, 100*time.Millisecond, func(escape func(error)) bool {
		devicesBody, err = getDevices(r.Context(), spotifyToken)
		if err != nil {
			escape(err)
			return false
		}
		var deviceResponse struct {
			Devices []struct {
				ID       string `json:"id"`
				IsActive bool   `json:"is_active"`
			} `json:"devices"`
		}
		err = json.Unmarshal(devicesBody, &deviceResponse)
		if err != nil {
			escape(err)
			return false
		}
		for _, device := range deviceResponse.Devices {
			if device.ID == deviceID {
				if device.IsActive {
					return true
				} else {
					break
				}
			}
		}
		return false
	})
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get devices: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSONRaw(w, devicesBody)
}

func getDevices(ctx context.Context, spotifyToken string) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		return json.RawMessage{}, fmt.Errorf("Could not build Spotify devices request: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return json.RawMessage{}, fmt.Errorf("Could not make Spotify devices request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return json.RawMessage{}, fmt.Errorf("Spotify devices request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return json.RawMessage{}, fmt.Errorf("Could not read Spotify devices response: %s", err)
	}
	return json.RawMessage(body), nil
}

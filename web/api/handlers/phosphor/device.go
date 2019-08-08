package phosphor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"golang.org/x/oauth2"
)

func TransferPlayback(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
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
	sess.SpotifyToken.SetAuthHeader(req)
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
		sess.SpotifyToken.SetAuthHeader(req)
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
	time.Sleep(250 * time.Millisecond)
	common.ExponentialBackoff(100*time.Millisecond, func() bool {
		devicesBody, err := getDevices(r.Context(), sess.SpotifyToken)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not get devices: %s", err), http.StatusInternalServerError)
			return true
		}
		var deviceResponse struct {
			Devices []struct {
				ID       string `json:"id"`
				IsActive bool   `json:"is_active"`
			} `json:"devices"`
		}
		err = json.Unmarshal(devicesBody, &deviceResponse)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not parse devices response: %s", err), http.StatusInternalServerError)
			return true
		}
		deviceStillExists := false
		for _, device := range deviceResponse.Devices {
			if device.ID == deviceID {
				deviceStillExists = true
				if device.IsActive {
					common.JSONRaw(w, devicesBody)
					return true
				}
				break
			}
		}
		if !deviceStillExists {
			common.FailWithRawJSON(w, fmt.Errorf("Device %s does not exist", deviceID), devicesBody, http.StatusNotFound)
			return true
		}
		return false
	})
}

func getDevices(ctx context.Context, token *oauth2.Token) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		return json.RawMessage{}, fmt.Errorf("Could not build Spotify devices request: %s", err)
	}
	token.SetAuthHeader(req)
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

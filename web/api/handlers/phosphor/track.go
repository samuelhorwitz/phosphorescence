package phosphor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"io/ioutil"
	"net/http"
)

func GetTrackData(w http.ResponseWriter, r *http.Request) {
	trackID := chi.URLParam(r, "trackID")
	spotifyToken, ok := r.Context().Value(middleware.SpotifyTokenContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify token on request context"), http.StatusInternalServerError)
		return
	}
	var rawTrackData json.RawMessage
	var featuresData json.RawMessage
	{
		req, err := http.NewRequestWithContext(common.HandlerTimeoutCancelContext(r), "GET", fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID), nil)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify track request: %s", err), http.StatusInternalServerError)
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not make Spotify track request: %s", err), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			common.Fail(w, fmt.Errorf("Spotify track request responded with %d", res.StatusCode), http.StatusForbidden)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not read Spotify track response: %s", err), http.StatusInternalServerError)
			return
		}
		rawTrackData = json.RawMessage(body)
	}
	{
		req, err := http.NewRequestWithContext(common.HandlerTimeoutCancelContext(r), "GET", fmt.Sprintf("https://api.spotify.com/v1/audio-features?ids=%s", trackID), nil)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify track audio feature request: %s", err), http.StatusInternalServerError)
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not make Spotify track audio feature request: %s", err), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			common.Fail(w, fmt.Errorf("Spotify track audio feature request responded with %d", res.StatusCode), http.StatusForbidden)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not read Spotify track audio feature response: %s", err), http.StatusInternalServerError)
			return
		}
		var parsedBody struct {
			AudioFeatures []json.RawMessage `json:"audio_features"`
		}
		err = json.Unmarshal(body, &parsedBody)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not parse Spotify track audio feature response: %s", err), http.StatusInternalServerError)
			return
		}
		if len(parsedBody.AudioFeatures) == 0 {
			common.Fail(w, fmt.Errorf("No Spotify track audio features on response: %s", err), http.StatusInternalServerError)
			return
		}
		featuresData = parsedBody.AudioFeatures[0]
	}
	common.JSON(w, map[string]interface{}{"track": struct {
		Track    json.RawMessage `json:"track"`
		Features json.RawMessage `json:"features"`
	}{
		Track:    rawTrackData,
		Features: featuresData,
	}})
}

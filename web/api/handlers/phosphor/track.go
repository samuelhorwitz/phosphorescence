package phosphor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/tracks"
	"io/ioutil"
	"net/http"
)

func GetTrackData(w http.ResponseWriter, r *http.Request) {
	country, ok := r.Context().Value(middleware.SpotifyCountryContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify country on request context"), http.StatusInternalServerError)
		return
	}
	trackID := chi.URLParam(r, "trackID")
	// First let's check if the track is in our huge cache of tracks
	// and return it if it is, no Spotify request needed.
	track, ok := tracks.GetTrack(trackID)
	if ok {
		canPlay, err := checkIfTrackPlayableInRegion(country, track.Track)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not check if track playable in region: %s", err), http.StatusInternalServerError)
			return
		}
		if !canPlay {
			common.Fail(w, fmt.Errorf("Track not playable in region %s", country), http.StatusNotFound)
			return
		}
		common.JSON(w, map[string]interface{}{"track": track})
		return
	}
	// We don't have that track cached, so let's reach out to Spotify
	spotifyToken, ok := r.Context().Value(middleware.SpotifyTokenContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify token on request context"), http.StatusInternalServerError)
		return
	}
	var rawTrackData json.RawMessage
	var featuresData json.RawMessage
	{
		req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID), nil)
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
			common.Fail(w, fmt.Errorf("Spotify track request responded with %d", res.StatusCode), http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not read Spotify track response: %s", err), http.StatusInternalServerError)
			return
		}
		rawTrackData = json.RawMessage(body)
	}
	canPlay, err := checkIfTrackPlayableInRegion(country, rawTrackData)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not check if track playable in region: %s", err), http.StatusInternalServerError)
		return
	}
	if !canPlay {
		common.Fail(w, fmt.Errorf("Track not playable in region %s", country), http.StatusNotFound)
		return
	}
	{
		req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("https://api.spotify.com/v1/audio-features?ids=%s", trackID), nil)
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
			common.Fail(w, fmt.Errorf("Spotify track audio feature request responded with %d", res.StatusCode), http.StatusInternalServerError)
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
			common.Fail(w, fmt.Errorf("No Spotify track audio features on response: %s", err), http.StatusNotFound)
			return
		}
		featuresData = parsedBody.AudioFeatures[0]
		if string(featuresData) == "null" {
			common.Fail(w, errors.New("Null audio features for track"), http.StatusNotFound)
			return
		}
	}
	common.JSON(w, map[string]interface{}{"track": struct {
		Track    json.RawMessage `json:"track"`
		Features json.RawMessage `json:"features"`
	}{
		Track:    rawTrackData,
		Features: featuresData,
	}})
}

func checkIfTrackPlayableInRegion(country string, track json.RawMessage) (bool, error) {
	var parsedBody struct {
		AvailableMarkets []string `json:"available_markets"`
	}
	err := json.Unmarshal(track, &parsedBody)
	if err != nil {
		return false, fmt.Errorf("Could not parse Spotify track response: %s", err)
	}
	for _, market := range parsedBody.AvailableMarkets {
		if market == country {
			return true, nil
		}
	}
	return false, nil
}

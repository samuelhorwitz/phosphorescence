package phosphor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"github.com/samuelhorwitz/phosphorescence/api/tracks"
)

// This struct is used instead json.RawMessage directly because we
// want to filter out keys that we don't have on tracks when they
// are pulled by the job. The job uses playlist track getting filters
// built into Spotify's endpoint and this struct mimics the result
// and implicitly discards the extra track data when remarshaled.
type spotifyTrack struct {
	AvailableMarkets []string        `json:"available_markets"`
	DurationMS       json.RawMessage `json:"duration_ms"`
	ExternalURLs     json.RawMessage `json:"external_urls"`
	ID               json.RawMessage `json:"id"`
	Name             json.RawMessage `json:"name"`
	Popularity       json.RawMessage `json:"popularity"`
	URI              json.RawMessage `json:"uri"`
	Album            struct {
		ExternalURLs json.RawMessage `json:"external_urls"`
		ID           json.RawMessage `json:"id"`
		Images       json.RawMessage `json:"images"`
		Name         json.RawMessage `json:"name"`
		Artists      json.RawMessage `json:"artists"`
	} `json:"album"`
	Artists json.RawMessage `json:"artists"`
}

func GetTrackData(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	trackID := chi.URLParam(r, "trackID")
	// First let's check if the track is in our huge cache of tracks
	// and return it if it is, no Spotify request needed.
	track, ok := tracks.GetTrack(trackID)
	var trackData spotifyTrack
	if ok {
		err := json.Unmarshal(track.Track, &trackData)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not parse Spotify track: %s", err), http.StatusInternalServerError)
			return
		}
		canPlay := checkIfTrackPlayableInRegion(sess.SpotifyCountry, trackData)
		if !canPlay {
			common.Fail(w, fmt.Errorf("Track not playable in region %s", sess.SpotifyCountry), http.StatusNotFound)
			return
		}
		common.JSON(w, map[string]interface{}{"track": track})
		return
	}
	// We don't have that track cached, so let's reach out to Spotify
	var featuresData json.RawMessage
	{
		req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID), nil)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify track request: %s", err), http.StatusInternalServerError)
			return
		}
		sess.SpotifyToken.SetAuthHeader(req)
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
		err = json.Unmarshal(body, &trackData)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not parse Spotify track response: %s", err), http.StatusInternalServerError)
			return
		}
	}
	canPlay := checkIfTrackPlayableInRegion(sess.SpotifyCountry, trackData)
	if !canPlay {
		common.Fail(w, fmt.Errorf("Track not playable in region %s", sess.SpotifyCountry), http.StatusNotFound)
		return
	}
	{
		req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("https://api.spotify.com/v1/audio-features?ids=%s", trackID), nil)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify track audio feature request: %s", err), http.StatusInternalServerError)
			return
		}
		sess.SpotifyToken.SetAuthHeader(req)
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
		Track    spotifyTrack    `json:"track"`
		Features json.RawMessage `json:"features"`
	}{
		Track:    trackData,
		Features: featuresData,
	}})
}

func checkIfTrackPlayableInRegion(country string, track spotifyTrack) bool {
	for _, market := range track.AvailableMarkets {
		if market == country {
			return true
		}
	}
	return false
}

package phosphor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	AvailableMarkets []string        `json:"available_markets,omitempty"`
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
	Artists    json.RawMessage `json:"artists"`
	IsPlayable bool            `json:"is_playable,omitempty"`
	LinkedFrom json.RawMessage `json:"linked_from,omitempty"`
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
	track, err := getTrackFromJSON(sess, trackID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get track from JSON, refusing to continue: %s", err), http.StatusInternalServerError)
		return
	}
	if track != nil {
		common.JSON(w, map[string]interface{}{"track": track})
		return
	}
	// We don't have that track cached (or it's region locked), so let's reach out to Spotify
	trackData, err := getTrackFromSpotify(r, sess, trackID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get track from Spotify: %s", err), http.StatusInternalServerError)
		return
	}
	if !trackData.IsPlayable {
		common.Fail(w, fmt.Errorf("Track not playable in region %s", sess.SpotifyCountry), http.StatusNotFound)
		return
	}
	// Finally lets get the audio features
	audioFeatures, err := getAudioFeatures(r, sess, trackID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get audio features for track: %s", err), http.StatusInternalServerError)
		return
	}
	featuresData, err := extractAudioFeatures(audioFeatures)
	if err != nil {
		common.Fail(w, fmt.Errorf("No audio features for track: %s", err), http.StatusNotFound)
		return
	}
	common.JSON(w, map[string]interface{}{"track": struct {
		Track    *spotifyTrack   `json:"track"`
		Features json.RawMessage `json:"features"`
	}{
		Track:    trackData,
		Features: featuresData,
	}})
}

func getTrackFromJSON(sess *session.Session, trackID string) (*tracks.TrackData, error) {
	track, ok := tracks.GetTrack(trackID)
	if !ok {
		return nil, nil
	}
	var trackData spotifyTrack
	err := json.Unmarshal(track.Track, &trackData)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify track: %s", err)
	}
	canPlay := checkIfTrackPlayableInRegion(sess.SpotifyCountry, &trackData)
	if !canPlay {
		if !isProduction {
			log.Println("Track found in JSON, but not playable in region, will look for linked track")
		}
		return nil, nil
	}
	return &track, nil
}

func getTrackFromSpotify(r *http.Request, sess *session.Session, trackID string) (*spotifyTrack, error) {
	req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("https://api.spotify.com/v1/tracks/%s?market=from_token", trackID), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not build Spotify track request: %s", err)
	}
	sess.SpotifyToken.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not make Spotify track request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Spotify track request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read Spotify track response: %s", err)
	}
	var trackData spotifyTrack
	err = json.Unmarshal(body, &trackData)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify track response: %s", err)
	}
	return &trackData, nil
}

func getAudioFeatures(r *http.Request, sess *session.Session, trackID string) ([]json.RawMessage, error) {
	req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("https://api.spotify.com/v1/audio-features?ids=%s", trackID), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not build Spotify track audio feature request: %s", err)
	}
	sess.SpotifyToken.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not make Spotify track audio feature request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Spotify track audio feature request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read Spotify track audio feature response: %s", err)
	}
	var parsedBody struct {
		AudioFeatures []json.RawMessage `json:"audio_features"`
	}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify track audio feature response: %s", err)
	}
	return parsedBody.AudioFeatures, nil
}

func extractAudioFeatures(audioFeatures []json.RawMessage) (json.RawMessage, error) {
	if len(audioFeatures) == 0 {
		return nil, errors.New("No Spotify track audio features on response")
	}
	featuresData := audioFeatures[0]
	if string(featuresData) == "null" {
		return nil, errors.New("Null audio features for track")
	}
	return featuresData, nil
}

func checkIfTrackPlayableInRegion(country string, track *spotifyTrack) bool {
	for _, market := range track.AvailableMarkets {
		if market == country {
			return true
		}
	}
	return false
}

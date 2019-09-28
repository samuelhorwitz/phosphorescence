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
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"github.com/samuelhorwitz/phosphorescence/api/tracks"
)

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
	// SpotifyTrackEnvelope
	common.JSON(w, map[string]interface{}{"track": models.SpotifyTrackEnvelope{
		ID:       trackID,
		Track:    trackData,
		Features: audioFeatures,
	}})
}

func getTrackFromJSON(sess *session.Session, trackID string) (*models.SpotifyTrackEnvelope, error) {
	track, ok := tracks.GetTrack(trackID)
	if !ok {
		return nil, nil
	}
	canPlay := checkIfTrackPlayableInRegion(sess.SpotifyCountry, track.Track)
	if !canPlay {
		if !isProduction {
			log.Println("Track found in JSON, but not playable in region, will look for linked track")
		}
		return nil, nil
	}
	return track, nil
}

func getTrackFromSpotify(r *http.Request, sess *session.Session, trackID string) (*models.SpotifyTrack, error) {
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
	var trackData models.SpotifyTrack
	err = json.Unmarshal(body, &trackData)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify track response: %s", err)
	}
	trackData.AvailableMarkets = nil
	trackData.Album.Images = findBestImage(trackData.Album.Images)
	return &trackData, nil
}

func getAudioFeatures(r *http.Request, sess *session.Session, trackID string) (*models.SpotifyFeatures, error) {
	req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("https://api.spotify.com/v1/audio-features/%s", trackID), nil)
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
	var spotifyFeatures models.SpotifyFeatures
	err = json.Unmarshal(body, &spotifyFeatures)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify track audio feature response: %s", err)
	}
	return &spotifyFeatures, nil
}

func checkIfTrackPlayableInRegion(country string, track *models.SpotifyTrack) bool {
	for _, market := range track.AvailableMarkets {
		if market == country {
			return true
		}
	}
	return false
}

func findBestImage(images []models.SpotifyImage) []models.SpotifyImage {
	var bestSize int
	var bestImage models.SpotifyImage
	for _, img := range images {
		size := img.Width * img.Height
		if size > bestSize {
			bestImage = img
			bestSize = size
		}
	}
	bestImage.Width = 0
	bestImage.Height = 0
	return []models.SpotifyImage{bestImage}
}

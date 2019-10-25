package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/samuelhorwitz/phosphorescence/api/cache"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"golang.org/x/oauth2"
)

type SpotifyTrackEnvelope struct {
	ID       string           `json:"id"`
	Track    *SpotifyTrack    `json:"track"`
	Features *SpotifyFeatures `json:"features"`
}

type SpotifyTrack struct {
	ID               string              `json:"id,omitempty"`
	Album            SpotifyAlbum        `json:"album"`
	Artists          []SpotifyArtist     `json:"artists"`
	Name             string              `json:"name"`
	Popularity       int                 `json:"popularity"`
	AvailableMarkets []string            `json:"available_markets,omitempty"`
	IsPlayable       bool                `json:"is_playable"`
	LinkedFrom       *SpotifyLinkedTrack `json:"linked_from,omitempty"`
	IsLocal          bool                `json:"is_local,omitempty"`
}

type SpotifyLinkedTrack struct {
	ID string `json:"id"`
}

type SpotifyArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SpotifyFeatures struct {
	ID                  string  `json:"id,omitempty"`
	Danceability        float64 `json:"danceability"`
	Energy              float64 `json:"energy"`
	Key                 int     `json:"key"`
	Loudness            float64 `json:"loudness"`
	Mode                int     `json:"mode"`
	Speechiness         float64 `json:"speechiness"`
	Acousticness        float64 `json:"acousticness"`
	Instrumentalness    float64 `json:"instrumentalness"`
	Liveness            float64 `json:"liveness"`
	Valence             float64 `json:"valence"`
	Tempo               float64 `json:"tempo"`
	DurationMillseconds int     `json:"duration_ms"`
	TimeSignature       int     `json:"time_signature"`
}

const maxTracksPerSpotifyRequest = 50
const maxTrackFeaturessPerSpotifyRequest = 100

type TrackNotFoundInRegionError struct {
	region string
}

func (e TrackNotFoundInRegionError) Error() string {
	return fmt.Sprintf("Track not playable in region %s", e.region)
}

func (t SpotifyTrack) MarshalJSON() ([]byte, error) {
	var tmp struct {
		Album      SpotifyAlbum        `json:"album"`
		Artists    []SpotifyArtist     `json:"artists"`
		Name       string              `json:"name"`
		Popularity int                 `json:"popularity"`
		IsPlayable bool                `json:"is_playable"`
		LinkedFrom *SpotifyLinkedTrack `json:"linked_from,omitempty"`
	}
	tmp.Album = t.Album
	tmp.Artists = t.Artists
	tmp.Name = t.Name
	tmp.Popularity = t.Popularity
	tmp.IsPlayable = t.IsPlayable
	tmp.LinkedFrom = t.LinkedFrom
	return json.Marshal(&tmp)
}

func (f SpotifyFeatures) MarshalJSON() ([]byte, error) {
	var tmp struct {
		Danceability        float64 `json:"danceability"`
		Energy              float64 `json:"energy"`
		Key                 int     `json:"key"`
		Loudness            float64 `json:"loudness"`
		Mode                int     `json:"mode"`
		Speechiness         float64 `json:"speechiness"`
		Acousticness        float64 `json:"acousticness"`
		Instrumentalness    float64 `json:"instrumentalness"`
		Liveness            float64 `json:"liveness"`
		Valence             float64 `json:"valence"`
		Tempo               float64 `json:"tempo"`
		DurationMillseconds int     `json:"duration_ms"`
		TimeSignature       int     `json:"time_signature"`
	}
	tmp.Danceability = f.Danceability
	tmp.Energy = f.Energy
	tmp.Key = f.Key
	tmp.Loudness = f.Loudness
	tmp.Mode = f.Mode
	tmp.Speechiness = f.Speechiness
	tmp.Acousticness = f.Acousticness
	tmp.Instrumentalness = f.Instrumentalness
	tmp.Liveness = f.Liveness
	tmp.Valence = f.Valence
	tmp.Tempo = f.Tempo
	tmp.DurationMillseconds = f.DurationMillseconds
	tmp.TimeSignature = f.TimeSignature
	return json.Marshal(&tmp)
}

func (e SpotifyTrackEnvelope) OriginalID() string {
	if e.Track != nil && e.Track.LinkedFrom != nil {
		return e.Track.LinkedFrom.ID
	}
	return e.ID
}

func GetTrack(ctx context.Context, region, trackID string) (*SpotifyTrackEnvelope, error) {
	phosphorescenceToken, err := spotifyclient.GetAppToken()
	if err != nil {
		return nil, fmt.Errorf("Could not get Spotify application token: %s", err)
	}
	// First, let's see if we have this track in the cache for the region
	cachedTrack := getTrackFromCache(region, trackID)
	if cachedTrack != nil {
		if !cachedTrack.Track.IsPlayable {
			return nil, TrackNotFoundInRegionError{region}
		}
		return cachedTrack, nil
	}
	// We don't have that track cached for that region, so let's reach out to Spotify
	trackData, err := getTrackFromSpotify(ctx, phosphorescenceToken, region, trackID)
	if err != nil {
		return nil, fmt.Errorf("Could not get track from Spotify: %s", err)
	}
	envelope := SpotifyTrackEnvelope{
		ID:    trackData.ID,
		Track: trackData,
	}
	if !trackData.IsPlayable {
		setTrackInCache(region, trackID, &envelope)
		return nil, TrackNotFoundInRegionError{region}
	}
	// Finally lets get the audio features
	audioFeatures, err := getAudioFeatures(ctx, phosphorescenceToken, trackID)
	if err != nil {
		return nil, fmt.Errorf("Could not get audio features for track: %s", err)
	}
	envelope.Features = audioFeatures
	setTrackInCache(region, trackID, &envelope)
	return &envelope, nil
}

func GetTracks(ctx context.Context, region string, trackIDs []string) ([]*SpotifyTrackEnvelope, error) {
	phosphorescenceToken, err := spotifyclient.GetAppToken()
	if err != nil {
		return nil, fmt.Errorf("Could not get Spotify application token: %s", err)
	}
	if len(trackIDs) > maxTracksPerRequest {
		return nil, ErrTooManyTracks
	}
	return getTracks(ctx, phosphorescenceToken, region, trackIDs)
}

func getTracks(ctx context.Context, token *oauth2.Token, region string, trackIDs []string) ([]*SpotifyTrackEnvelope, error) {
	var missingFromCache []string
	missingFromCacheMap := make(map[string]int)
	cachedTracks := getTracksFromCache(region, trackIDs)
	tracksMap := make(map[string]*SpotifyTrackEnvelope)
	for i, trackID := range trackIDs {
		cachedTrack, ok := cachedTracks[trackID]
		if ok && cachedTrack != nil && cachedTrack.Track != nil && cachedTrack.Track.IsPlayable && cachedTrack.Features != nil {
			tracksMap[trackID] = &SpotifyTrackEnvelope{
				ID:       cachedTrack.ID,
				Track:    cachedTrack.Track,
				Features: cachedTrack.Features,
			}
		} else {
			missingFromCache = append(missingFromCache, trackID)
			missingFromCacheMap[trackID] = i
		}
	}
	if len(missingFromCache) > 0 {
		missingTracks, err := getTracksFromSpotify(ctx, token, region, missingFromCache)
		if err != nil {
			return nil, fmt.Errorf("Could not get missing tracks: %s", err)
		}
		for _, unenclosedTrack := range missingTracks {
			track := unenclosedTrack
			envelope := SpotifyTrackEnvelope{
				ID:    track.ID,
				Track: &track,
			}
			tracksMap[envelope.OriginalID()] = &envelope
		}
	}
	var tracks []*SpotifyTrackEnvelope
	for _, trackID := range trackIDs {
		tracks = append(tracks, tracksMap[trackID])
	}
	var err error
	tracks, err = populateAudioFeatures(ctx, token, region, tracks)
	if err != nil {
		return nil, fmt.Errorf("Could not get missing audio features: %s", err)
	}
	return tracks, nil
}

func populateAudioFeatures(ctx context.Context, token *oauth2.Token, region string, tracks []*SpotifyTrackEnvelope) ([]*SpotifyTrackEnvelope, error) {
	var missingFromCache []string
	missingFromCacheMap := make(map[string]int)
	var trackIDs []string
	for _, track := range tracks {
		trackIDs = append(trackIDs, track.OriginalID())
	}
	cachedTracks := getTracksFromCache(region, trackIDs)
	for i, track := range tracks {
		trackID := track.OriginalID()
		cachedTrack, ok := cachedTracks[trackID]
		if ok && cachedTrack != nil && cachedTrack.Track != nil && cachedTrack.Track.IsPlayable && cachedTrack.Features != nil {
			track.Features = cachedTrack.Features
		} else {
			missingFromCache = append(missingFromCache, trackID)
			missingFromCacheMap[trackID] = i
		}
	}
	if len(missingFromCache) > 0 {
		audioFeaturesFromSpotify, err := getManyAudioFeatures(ctx, token, missingFromCache)
		if err != nil {
			return nil, fmt.Errorf("Could not get missing audio features: %s", err)
		}
		for _, features := range audioFeaturesFromSpotify {
			featuresClosure := features
			id := featuresClosure.ID
			index := missingFromCacheMap[id]
			tracks[index].Features = &featuresClosure
			setTrackInCache(region, id, tracks[index])
		}
	}
	return tracks, nil
}

func getTrackFromSpotify(ctx context.Context, token *oauth2.Token, region, trackID string) (*SpotifyTrack, error) {
	tracks, err := getTracksFromSpotify(ctx, token, region, []string{trackID})
	if err != nil {
		return nil, fmt.Errorf("Could not get track: %s", err)
	}
	firstTrack := tracks[0]
	return &firstTrack, nil
}

func getTracksFromSpotify(ctx context.Context, token *oauth2.Token, region string, allTrackIDs []string) ([]SpotifyTrack, error) {
	trackIDPages := pageTrackIDs(allTrackIDs, maxTracksPerSpotifyRequest)
	var tracks []SpotifyTrack
	for _, trackIDs := range trackIDPages {
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.spotify.com/v1/tracks?ids=%s&market=%s", strings.Join(trackIDs, ","), region), nil)
		if err != nil {
			return nil, fmt.Errorf("Could not build Spotify track request: %s", err)
		}
		token.SetAuthHeader(req)
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
		var trackData struct {
			Tracks []SpotifyTrack `json:"tracks"`
		}
		err = json.Unmarshal(body, &trackData)
		if err != nil {
			return nil, fmt.Errorf("Could not parse Spotify track response: %s", err)
		}
		for _, track := range trackData.Tracks {
			track.AvailableMarkets = nil
			track.Album.Images = findBestImage(track.Album.Images)
		}
		tracks = append(tracks, trackData.Tracks...)
	}
	return tracks, nil
}

func getAudioFeatures(ctx context.Context, token *oauth2.Token, trackID string) (*SpotifyFeatures, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.spotify.com/v1/audio-features/%s", trackID), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not build Spotify track audio feature request: %s", err)
	}
	token.SetAuthHeader(req)
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
	var spotifyFeatures SpotifyFeatures
	err = json.Unmarshal(body, &spotifyFeatures)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify track audio feature response: %s", err)
	}
	return &spotifyFeatures, nil
}

func getManyAudioFeatures(ctx context.Context, token *oauth2.Token, allTrackIDs []string) ([]SpotifyFeatures, error) {
	trackIDPages := pageTrackIDs(allTrackIDs, maxTrackFeaturessPerSpotifyRequest)
	var features []SpotifyFeatures
	for _, trackIDs := range trackIDPages {
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.spotify.com/v1/audio-features?ids=%s", strings.Join(trackIDs, ",")), nil)
		if err != nil {
			return nil, fmt.Errorf("Could not build Spotify track audio feature request: %s", err)
		}
		token.SetAuthHeader(req)
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
		var response struct {
			AudioFeatures []SpotifyFeatures `json:"audio_features"`
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("Could not parse Spotify track audio feature response: %s", err)
		}
		features = append(features, response.AudioFeatures...)
	}
	return features, nil
}

func findBestImage(images []SpotifyImage) []SpotifyImage {
	// from the docs: "The cover art for the album in various sizes, widest first."
	// https://developer.spotify.com/documentation/web-api/reference/object-model/#album-object-simplified
	if len(images) == 0 {
		return images
	}
	bestImage := images[0]
	bestImage.Width = 0
	bestImage.Height = 0
	return []SpotifyImage{bestImage}
}

func getTrackFromCache(region, trackID string) *SpotifyTrackEnvelope {
	cachedTrack, ok := cache.GetTrack(region, trackID)
	if !ok {
		return nil
	}
	envelope := SpotifyTrackEnvelope{
		ID: cachedTrack.ID,
	}
	var track SpotifyTrack
	var features SpotifyFeatures
	var err error
	err = json.Unmarshal([]byte(cachedTrack.Track), &track)
	if err != nil {
		if !isProduction {
			log.Printf("Could not unmarshal cached track: %s", err)
		}
		return nil
	}
	envelope.Track = &track
	if cachedTrack.Features != "" {
		err = json.Unmarshal([]byte(cachedTrack.Features), &features)
		if err != nil {
			if !isProduction {
				log.Printf("Could not unmarshal cached track features: %s", err)
			}
			return nil
		}
		envelope.Features = &features
	}
	return &envelope
}

func getTracksFromCache(region string, trackIDs []string) map[string]*SpotifyTrackEnvelope {
	envelopes := make(map[string]*SpotifyTrackEnvelope)
	cachedTracks := cache.GetTracks(region, trackIDs)
	if cachedTracks == nil {
		return envelopes
	}
	for id, cachedTrack := range cachedTracks {
		envelope := SpotifyTrackEnvelope{
			ID: cachedTrack.ID,
		}
		var track SpotifyTrack
		var features SpotifyFeatures
		var err error
		err = json.Unmarshal([]byte(cachedTrack.Track), &track)
		if err != nil {
			if !isProduction {
				log.Printf("Could not unmarshal cached track: %s", err)
			}
			continue
		}
		envelope.Track = &track
		if cachedTrack.Features != "" {
			err = json.Unmarshal([]byte(cachedTrack.Features), &features)
			if err != nil {
				if !isProduction {
					log.Printf("Could not unmarshal cached track features: %s", err)
				}
				continue
			}
			envelope.Features = &features
		}
		envelopes[id] = &envelope
	}
	return envelopes
}

func setTrackInCache(region, trackID string, envelope *SpotifyTrackEnvelope) bool {
	var cachedTrack cache.CachedTrack
	trackJSON, err := json.Marshal(envelope.Track)
	if err != nil {
		if !isProduction {
			log.Printf("Could not marshal track: %s", err)
		}
		return false
	}
	featuresJSON, err := json.Marshal(envelope.Features)
	if err != nil {
		if !isProduction {
			log.Printf("Could not marshal track features: %s", err)
		}
		return false
	}
	cachedTrack.ID = envelope.ID
	cachedTrack.Track = string(trackJSON)
	cachedTrack.Features = string(featuresJSON)
	cache.SetTrack(region, trackID, cachedTrack)
	return true
}

func dedupeTrackIDs(trackIDs []string) (dedupedTrackIDs []string) {
	trackIDDedupeMap := make(map[string]bool)
	for _, trackID := range trackIDs {
		if trackIDDedupeMap[trackID] {
			continue
		}
		dedupedTrackIDs = append(dedupedTrackIDs, trackID)
		trackIDDedupeMap[trackID] = true
	}
	return dedupedTrackIDs
}

func pageTrackIDs(trackIDs []string, chunkSize uint) (pagedTrackIDs [][]string) {
	dedupedTrackIDs := dedupeTrackIDs(trackIDs)
	totalTracks := uint(len(dedupedTrackIDs))
	for i := uint(0); i < totalTracks; i += chunkSize {
		end := i + chunkSize
		if end > totalTracks {
			end = totalTracks
		}
		pagedTrackIDs = append(pagedTrackIDs, dedupedTrackIDs[i:end])
	}
	return pagedTrackIDs
}

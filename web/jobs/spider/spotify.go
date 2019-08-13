package spider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var spotifyClient = &http.Client{
	Timeout: 60 * time.Second,
}

var spotifyToken string

type TrackEnvelope struct {
	Track    json.RawMessage `json:"track"`
	Features json.RawMessage `json:"features"`
}

func initializeToken(spotifyClientID, spotifySecret string) error {
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(getSpotifyTokenRequestBody(spotifyClientID, spotifySecret)))
	if err != nil {
		return fmt.Errorf("Could not build Spotify token request: %s", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := spotifyClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Spotify token request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Spotify token request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Could not read Spotify token response: %s", err)
	}
	var parsedBody struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return fmt.Errorf("Could not parse Spotify token response: %s", err)
	}
	spotifyToken = parsedBody.AccessToken
	return nil
}

func getSpotifyTokenRequestBody(spotifyClientID, spotifySecret string) string {
	body := url.Values{}
	body.Set("client_id", spotifyClientID)
	body.Set("client_secret", spotifySecret)
	body.Set("grant_type", "client_credentials")
	return body.Encode()
}

func getTracksFromPlaylist(playlistID string) (map[string]*TrackEnvelope, error) {
	nextURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?limit=100&fields=next,items(track(available_markets,duration_ms,external_urls,id,name,popularity,album(external_urls,id,images,name,artists),artists))", playlistID)
	tracks := make(map[string]*TrackEnvelope)
	for nextURL != "" {
		log.Printf("Handling %s...", nextURL)
		req, err := http.NewRequest("GET", nextURL, nil)
		if err != nil {
			return tracks, fmt.Errorf("Could not build Spotify playlist request: %s", err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
		res, err := spotifyClient.Do(req)
		if err != nil {
			return tracks, fmt.Errorf("Could not make Spotify playlist request: %s", err)
		}
		defer res.Body.Close()
		if res.StatusCode == http.StatusTooManyRequests {
			log.Println("Spotify asked us too back off")
			retryAfterSeconds, err := strconv.Atoi(res.Header.Get("Retry-After"))
			if err != nil {
				return nil, fmt.Errorf("Could not parse retry after header: %s", err)
			}
			log.Printf("Waiting for %d seconds...", retryAfterSeconds)
			time.Sleep(time.Duration(retryAfterSeconds) * time.Second)
			return getTrackFeatures(tracks)
		}
		if res.StatusCode != http.StatusOK {
			return tracks, fmt.Errorf("Spotify playlist request responded with %d", res.StatusCode)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return tracks, fmt.Errorf("Could not read Spotify playlist response: %s", err)
		}
		var parsedBody struct {
			Next  string
			Items []struct {
				Track json.RawMessage `json:"track"`
			} `json:"items"`
		}
		err = json.Unmarshal(body, &parsedBody)
		if err != nil {
			return tracks, fmt.Errorf("Could not parse Spotify playlist response: %s", err)
		}
		for _, item := range parsedBody.Items {
			var track struct {
				ID string `json:"id"`
			}
			err = json.Unmarshal(item.Track, &track)
			if err != nil {
				return tracks, fmt.Errorf("Could not parse Spotify track: %s", err)
			}
			tracks[track.ID] = &TrackEnvelope{Track: item.Track}
		}
		nextURL = parsedBody.Next
	}
	return tracks, nil
}

func getTrackFeaturesInBatches(allTracks map[string]*TrackEnvelope) (map[string]*TrackEnvelope, error) {
	var trackIDs []string
	for id := range allTracks {
		trackIDs = append(trackIDs, id)
	}
	newTracks := make(map[string]*TrackEnvelope)
	for i := 0; i < len(trackIDs); i += 100 {
		end := i + 100
		if len(trackIDs) < end {
			end = len(trackIDs)
		}
		ids := trackIDs[i:end]
		log.Printf("Getting feature batch %d of %d, ids %s to %s...", (i/100)+1, int64(math.Ceil(float64(len(trackIDs))/100.0)), ids[0], ids[len(ids)-1])
		batchOfTracks := make(map[string]*TrackEnvelope)
		for _, id := range ids {
			batchOfTracks[id] = allTracks[id]
		}
		var err error
		tracks, err := getTrackFeatures(batchOfTracks)
		if err != nil {
			return nil, fmt.Errorf("Could not get batch of track features: %s", err)
		}
		for id, track := range tracks {
			newTracks[id] = track
		}
	}
	return newTracks, nil
}

func getTrackFeatures(tracks map[string]*TrackEnvelope) (map[string]*TrackEnvelope, error) {
	var ids []string
	for id := range tracks {
		ids = append(ids, id)
	}
	idsParam := strings.Join(ids, ",")
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/audio-features?ids=%s", idsParam), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not build Spotify track features request: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
	res, err := spotifyClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not make Spotify track features request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusTooManyRequests {
		log.Println("Spotify asked us too back off")
		retryAfterSeconds, err := strconv.Atoi(res.Header.Get("Retry-After"))
		if err != nil {
			return nil, fmt.Errorf("Could not parse retry after header: %s", err)
		}
		log.Printf("Waiting for %d seconds...", retryAfterSeconds)
		time.Sleep(time.Duration(retryAfterSeconds) * time.Second)
		return getTrackFeatures(tracks)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Spotify track features request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read Spotify track features response: %s", err)
	}
	var parsedBody struct {
		AudioFeatures []json.RawMessage `json:"audio_features"`
	}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify track features response: %s", err)
	}
	featuresByID := make(map[string]json.RawMessage)
	for _, trackFeatures := range parsedBody.AudioFeatures {
		var parsedBody struct {
			ID string `json:"id"`
		}
		err = json.Unmarshal(trackFeatures, &parsedBody)
		if err != nil {
			return nil, fmt.Errorf("Could not parse Spotify track feature response: %s", err)
		}
		featuresByID[parsedBody.ID] = trackFeatures
	}
	for id, track := range tracks {
		if track == nil {
			delete(tracks, id)
			continue
		}
		if trackFeatures, ok := featuresByID[id]; ok {
			track.Features = trackFeatures
		} else {
			delete(tracks, id)
		}
	}
	return tracks, nil
}

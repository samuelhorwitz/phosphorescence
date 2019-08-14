package phosphor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

type User struct {
	SpotifyID     string `json:"spotifyId"`
	Name          string `json:"name"`
	Country       string `json:"country"`
	Authenticated bool   `json:"authenticated"`
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	common.JSON(w, map[string]interface{}{
		"user": User{
			SpotifyID:     sess.SpotifyID,
			Name:          sess.SpotifyName,
			Country:       sess.SpotifyCountry,
			Authenticated: sess.Authenticated,
		},
	})
}

func GetCurrentlyPlaying(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), "GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not build Spotify currently playing request: %s", err), http.StatusInternalServerError)
		return
	}
	sess.SpotifyToken.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not make Spotify currently playing request: %s", err), http.StatusInternalServerError)
		return
	}
	fetchedAt := time.Now()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		common.Fail(w, fmt.Errorf("Spotify currently playing request responded with %d", res.StatusCode), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read Spotify currently playing response: %s", err), http.StatusInternalServerError)
		return
	}
	var parsedBody struct {
		IsPlaying            bool            `json:"is_playing"`
		Track                json.RawMessage `json:"item"`
		CurrentlyPlayingType string          `json:"currently_playing_type"`
		Timestamp            int             `json:"timestamp"`
		ProgressMilliseconds int             `json:"progress_ms"`
	}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not parse Spotify currently playing response: %s", err), http.StatusInternalServerError)
		return
	}
	if parsedBody.CurrentlyPlayingType != "track" {
		common.Fail(w, fmt.Errorf("Spotify currently playing is not a track: %s", err), http.StatusNotFound)
		return
	}
	common.JSON(w, map[string]interface{}{
		"isPlaying":        parsedBody.IsPlaying,
		"progress":         parsedBody.ProgressMilliseconds,
		"fetchedAtSpotify": parsedBody.Timestamp,
		"fetchedAt":        fetchedAt.UnixNano() / int64(time.Millisecond),
		"track":            parsedBody.Track,
	})
}

func ListCurrentUserScripts(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.AuthenticatedSessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
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
	scripts, err := models.GetScriptsBySpotifyUserID(sess.SpotifyID, count, from, true)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get scripts for user: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"scripts": scripts})
}

func ListSpotifyDevices(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	body, err := getDevices(r.Context(), sess.SpotifyToken)
	if err != nil {
		common.Fail(w, fmt.Errorf("Unable to get devices: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSONRaw(w, body)
}

func CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read request body: %s", err), http.StatusInternalServerError)
		return
	}
	var requestBody struct {
		UTCOffsetMinutes int `json:"utcOffsetMinutes"`
		Tracks           []struct {
			Name string `json:"name"`
			URI  string `json:"uri"`
		} `json:"tracks"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not parse request body: %s", err), http.StatusInternalServerError)
		return
	}
	if len(requestBody.Tracks) < 1 {
		common.Fail(w, errors.New("Must include at least one track"), http.StatusBadRequest)
		return
	}
	var createdPlaylist struct {
		ID string `json:"id"`
	}
	{
		var createPlaylistBody struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		now := time.Now().UTC().Add(time.Duration(requestBody.UTCOffsetMinutes) * time.Minute)
		createPlaylistBody.Name = requestBody.Tracks[0].Name
		createPlaylistBody.Description = fmt.Sprintf("Created by Phosphorescence on %s at %s. Visit https://phosphor.me to create more trance playlists!", now.Format("Monday, January _2"), now.Format("3:04 PM"))
		createPlaylistBodyJSON, err := json.Marshal(createPlaylistBody)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not marshal create playlist request body: %s", err), http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequestWithContext(r.Context(), "POST", fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", sess.SpotifyID), bytes.NewBuffer(createPlaylistBodyJSON))
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify create playlist request: %s", err), http.StatusInternalServerError)
			return
		}
		sess.SpotifyToken.SetAuthHeader(req)
		req.Header.Set("Content-Type", "application/json")
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not make Spotify create playlist request: %s", err), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
			common.Fail(w, fmt.Errorf("Spotify create playlist request responded with %d", res.StatusCode), http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not read Spotify create playlist response: %s", err), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &createdPlaylist)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not parse Spotify create playlist response: %s", err), http.StatusInternalServerError)
			return
		}
	}
	{
		var addTracksBody struct {
			URIs []string `json:"uris"`
		}
		for _, track := range requestBody.Tracks {
			addTracksBody.URIs = append(addTracksBody.URIs, track.URI)
		}
		addTracksBodyJSON, err := json.Marshal(addTracksBody)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not marshal add tracks request body: %s", err), http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequestWithContext(r.Context(), "POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", createdPlaylist.ID), bytes.NewBuffer(addTracksBodyJSON))
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify add tracks request: %s", err), http.StatusInternalServerError)
			return
		}
		sess.SpotifyToken.SetAuthHeader(req)
		req.Header.Set("Content-Type", "application/json")
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not make Spotify add tracks request: %s", err), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			common.Fail(w, fmt.Errorf("Spotify add tracks request responded with %d", res.StatusCode), http.StatusInternalServerError)
			return
		}
	}
	{
		req, err := http.NewRequestWithContext(r.Context(), "PUT", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/images", createdPlaylist.ID), bytes.NewBuffer([]byte(playlistImageBase64)))
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify change image request: %s", err), http.StatusInternalServerError)
			return
		}
		sess.SpotifyToken.SetAuthHeader(req)
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not make Spotify change image request: %s", err), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusAccepted {
			common.Fail(w, fmt.Errorf("Spotify change image request responded with %d", res.StatusCode), http.StatusInternalServerError)
			return
		}
	}
	common.JSON(w, map[string]interface{}{"playlist": createdPlaylist.ID})
}

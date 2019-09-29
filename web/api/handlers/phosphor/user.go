package phosphor

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	SpotifyID         string `json:"spotifyId"`
	Name              string `json:"name"`
	Country           string `json:"country"`
	Product           string `json:"product"`
	Authenticated     bool   `json:"authenticated"`
	GoogleAnalyticsID string `json:"gaId,omitempty"`
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	sess, err := session.UpdateSessionDetailsFromSpotify(r, sess)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not update session details: %s", err), http.StatusInternalServerError)
		return
	}
	googleAnalyticsIDSum := hmac.New(sha256.New, googleAnalyticsSecret)
	googleAnalyticsIDSum.Write([]byte(sess.SpotifyID))
	googleAnalyticsID := googleAnalyticsIDSum.Sum(nil)
	common.JSON(w, map[string]interface{}{
		"user": User{
			SpotifyID:         sess.SpotifyID,
			Name:              sess.SpotifyName,
			Country:           sess.SpotifyCountry,
			Product:           sess.SpotifyProduct,
			Authenticated:     sess.Authenticated,
			GoogleAnalyticsID: hex.EncodeToString(googleAnalyticsID),
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

func CreateAndFollowPlaylist(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	playlistID, err := createPlaylist(r)
	if err != nil {
		code := http.StatusInternalServerError
		if err == errNoTracks {
			code = http.StatusBadRequest
		}
		common.Fail(w, fmt.Errorf("Failed to create playlist: %s", err), code)
		return
	}
	err = followPlaylist(r, sess.SpotifyToken, playlistID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Failed to follow playlist: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"playlist": playlistID})
}

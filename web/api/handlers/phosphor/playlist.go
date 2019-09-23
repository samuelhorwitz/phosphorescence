package phosphor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"golang.org/x/oauth2"
)

var errNoTracks = errors.New("Must include at least one track")

func CreatePrivatePlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, err := createPlaylist(r)
	if err != nil {
		code := http.StatusInternalServerError
		if err == errNoTracks {
			code = http.StatusBadRequest
		}
		common.Fail(w, fmt.Errorf("Failed to create playlist: %s", err), code)
		return
	}
	common.JSON(w, map[string]interface{}{"playlist": playlistID})
}

func MakeOfficialPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "playlistID")
	phosphorescenceToken, err := getPhosphorescenceToken()
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get Spotify application user token: %s", err), http.StatusInternalServerError)
		return
	}
	err = followPlaylist(r, phosphorescenceToken, playlistID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Failed to follow playlist: %s", err), http.StatusInternalServerError)
		return
	}
	var updatePlaylistBody struct {
		Public bool `json:"public"`
	}
	updatePlaylistBody.Public = true
	updatePlaylistBodyJSON, err := json.Marshal(updatePlaylistBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not marshal update playlist request body: %s", err), http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), "PUT", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playlistID), bytes.NewBuffer(updatePlaylistBodyJSON))
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not build Spotify update playlist request: %s", err), http.StatusInternalServerError)
		return
	}
	phosphorescenceToken.SetAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not make Spotify update playlist request: %s", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
		common.Fail(w, fmt.Errorf("Spotify update playlist request responded with %d", res.StatusCode), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"playlist": playlistID})
}

func createPlaylist(r *http.Request) (string, error) {
	phosphorescenceToken, err := getPhosphorescenceToken()
	if err != nil {
		return "", fmt.Errorf("Could not get Spotify application user token: %s", err)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read request body: %s", err)
	}
	var requestBody struct {
		Image            string `json:"image"`
		UTCOffsetMinutes int    `json:"utcOffsetMinutes"`
		Tracks           []struct {
			Name string `json:"name"`
			URI  string `json:"uri"`
		} `json:"tracks"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return "", fmt.Errorf("Could not parse request body: %s", err)
	}
	if len(requestBody.Tracks) < 1 {
		return "", errNoTracks
	}
	var createdPlaylist struct {
		ID string `json:"id"`
	}
	{
		var createPlaylistBody struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Public      bool   `json:"public"`
		}
		createPlaylistBody.Name = getPlaylistName(requestBody.Tracks[0].Name)
		createPlaylistBody.Description = getPlaylistDescription(requestBody.UTCOffsetMinutes)
		createPlaylistBody.Public = false
		createPlaylistBodyJSON, err := json.Marshal(createPlaylistBody)
		if err != nil {
			return "", fmt.Errorf("Could not marshal create playlist request body: %s", err)
		}
		req, err := http.NewRequestWithContext(r.Context(), "POST", fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", phosphorescenceSpotifyID), bytes.NewBuffer(createPlaylistBodyJSON))
		if err != nil {
			return "", fmt.Errorf("Could not build Spotify create playlist request: %s", err)
		}
		phosphorescenceToken.SetAuthHeader(req)
		req.Header.Set("Content-Type", "application/json")
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("Could not make Spotify create playlist request: %s", err)
		}
		defer res.Body.Close()
		if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
			return "", fmt.Errorf("Spotify create playlist request responded with %d", res.StatusCode)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("Could not read Spotify create playlist response: %s", err)
		}
		err = json.Unmarshal(body, &createdPlaylist)
		if err != nil {
			return "", fmt.Errorf("Could not parse Spotify create playlist response: %s", err)
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
			return "", fmt.Errorf("Could not marshal add tracks request body: %s", err)
		}
		req, err := http.NewRequestWithContext(r.Context(), "POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", createdPlaylist.ID), bytes.NewBuffer(addTracksBodyJSON))
		if err != nil {
			return "", fmt.Errorf("Could not build Spotify add tracks request: %s", err)
		}
		phosphorescenceToken.SetAuthHeader(req)
		req.Header.Set("Content-Type", "application/json")
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("Could not make Spotify add tracks request: %s", err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			return "", fmt.Errorf("Spotify add tracks request responded with %d", res.StatusCode)
		}
	}
	{
		var buf *bytes.Buffer
		if requestBody.Image != "" {
			buf = bytes.NewBuffer([]byte(requestBody.Image))
		} else {
			buf = bytes.NewBuffer([]byte(playlistImageBase64))
		}
		req, err := http.NewRequestWithContext(r.Context(), "PUT", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/images", createdPlaylist.ID), buf)
		if err != nil {
			return "", fmt.Errorf("Could not build Spotify change image request: %s", err)
		}
		phosphorescenceToken.SetAuthHeader(req)
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("Could not make Spotify change image request: %s", err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusAccepted {
			return "", fmt.Errorf("Spotify change image request responded with %d", res.StatusCode)
		}
	}
	// We "delete" the playlist to avoid clutter on the account. Deleting is really just unfollowing.
	{
		req, err := http.NewRequestWithContext(r.Context(), "DELETE", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/followers", createdPlaylist.ID), nil)
		if err != nil {
			return "", fmt.Errorf("Could not build Spotify unfollow playlist request: %s", err)
		}
		phosphorescenceToken.SetAuthHeader(req)
		req.Header.Set("Content-Type", "application/json")
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("Could not make Spotify unfollow playlist request: %s", err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return "", fmt.Errorf("Spotify unfollow playlist request responded with %d", res.StatusCode)
		}
	}
	return createdPlaylist.ID, nil
}

func followPlaylist(r *http.Request, token *oauth2.Token, playlistID string) error {
	req, err := http.NewRequestWithContext(r.Context(), "PUT", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/followers", playlistID), nil)
	if err != nil {
		return fmt.Errorf("Could not build Spotify follow request: %s", err)
	}
	token.SetAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Spotify follow request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Spotify follow request responded with %d", res.StatusCode)
	}
	return nil
}

func getPlaylistName(firstTrackName string) string {
	return fmt.Sprintf("phosphor.me | %s", firstTrackName)
}

func getPlaylistDescription(utcOffsetMinutes int) string {
	now := time.Now().UTC().Add(time.Duration(utcOffsetMinutes) * time.Minute)
	return fmt.Sprintf("Created by Phosphorescence on %s at %s. Visit phosphor.me to create more trance playlists!", now.Format("Monday, January _2"), now.Format("3:04 PM"))
}

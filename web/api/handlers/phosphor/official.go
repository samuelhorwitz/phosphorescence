package phosphor

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

type playlist struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tracks      struct {
		Items []struct {
			Track struct {
				URI string `json:"uri"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

func MakeOfficialPlaylist(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	playlistID := chi.URLParam(r, "playlistID")
	playlist, err := getPlaylist(r, sess, playlistID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get playlist: %s", err), http.StatusInternalServerError)
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
		createPlaylistBody.Name = playlist.Name
		createPlaylistBody.Description = playlist.Description
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
		for _, track := range playlist.Tracks.Items {
			addTracksBody.URIs = append(addTracksBody.URIs, track.Track.URI)
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
	var base64Image string
	{
		res, err := safeHTTPClient.Get(playlist.Images[0].URL)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not get playlist image: %s", err), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		image, err := ioutil.ReadAll(res.Body)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not read playlist image: %s", err), http.StatusInternalServerError)
			return
		}
		base64Image = base64.StdEncoding.EncodeToString(image)
	}
	{
		buf := bytes.NewBuffer([]byte(base64Image))
		req, err := http.NewRequestWithContext(r.Context(), "PUT", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/images", createdPlaylist.ID), buf)
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

func getPlaylist(r *http.Request, sess *session.Session, playlistID string) (*playlist, error) {
	req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playlistID), nil)
	if err != nil {
		return nil, fmt.Errorf("Could not build Spotify playlist request: %s", err)
	}
	sess.SpotifyToken.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not make Spotify playlist request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Spotify playlist request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read Spotify playlist response: %s", err)
	}
	var playlistData playlist
	err = json.Unmarshal(body, &playlistData)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify playlist response: %s", err)
	}
	return &playlistData, nil
}

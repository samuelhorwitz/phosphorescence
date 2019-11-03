package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"golang.org/x/oauth2"
)

type SpotifyPlaylist struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Owner       SpotifyUser                `json:"owner"`
	Images      []SpotifyImage             `json:"images"`
	Tracks      SpotifyPlaylistTrackPaging `json:"tracks"`
}

type SpotifyPlaylistTrackPaging struct {
	Items []SpotifyPlaylistTrack `json:"items"`
	Total int                    `json:"total"`
	Next  string                 `json:"next"`
}

type SpotifyPlaylistTrack struct {
	IsLocal bool         `json:"is_local"`
	Track   SpotifyTrack `json:"track"`
}

type Playlist struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Owner       SpotifyUser             `json:"owner"`
	Images      []SpotifyImage          `json:"images"`
	Tracks      []*SpotifyTrackEnvelope `json:"tracks"`
}

const maxTracksPerRequest = 500

var ErrTooManyTracks = fmt.Errorf("Too many tracks (max %d)", maxTracksPerRequest)

func GetPlaylist(ctx context.Context, region, playlistID string) (*Playlist, error) {
	phosphorescenceToken, err := spotifyclient.GetAppToken()
	if err != nil {
		return nil, fmt.Errorf("Could not get Spotify application token: %s", err)
	}
	playlist, err := getPlaylist(ctx, phosphorescenceToken, region, playlistID)
	if err != nil {
		return nil, fmt.Errorf("Could not get playlist: %s", err)
	}
	return playlist, nil
}

func GetSimplePlaylist(ctx context.Context, playlistID string) (*Playlist, error) {
	phosphorescenceToken, err := spotifyclient.GetAppToken()
	if err != nil {
		return nil, fmt.Errorf("Could not get Spotify application token: %s", err)
	}
	playlist, err := getPlaylist(ctx, phosphorescenceToken, "", playlistID)
	if err != nil {
		return nil, fmt.Errorf("Could not get playlist: %s", err)
	}
	return playlist, nil
}

func CreatePlaylist(ctx context.Context, firstTrackName string, base64Image string, utcOffsetMinutes int, trackURIs []string) (string, error) {
	phosphorescenceToken, err := spotifyclient.GetAppUserToken()
	if err != nil {
		return "", fmt.Errorf("Could not get Spotify application user token: %s", err)
	}
	createdPlaylistID, err := createPlaylist(ctx, phosphorescenceToken, firstTrackName, utcOffsetMinutes)
	if err != nil {
		return "", fmt.Errorf("Could not create playlist: %s", err)
	}
	err = addTracksToPlaylist(ctx, phosphorescenceToken, createdPlaylistID, trackURIs)
	if err != nil {
		return "", fmt.Errorf("Could not add tracks to playlist: %s", err)
	}
	err = setPlaylistImage(ctx, phosphorescenceToken, createdPlaylistID, base64Image)
	if err != nil {
		return "", fmt.Errorf("Could not set playlist image: %s", err)
	}
	err = unfollowPlaylist(ctx, phosphorescenceToken, createdPlaylistID)
	if err != nil {
		return "", fmt.Errorf("Could not unfollow playlist: %s", err)
	}
	return createdPlaylistID, nil
}

func MakePlaylistPublic(ctx context.Context, playlistID string) error {
	phosphorescenceToken, err := spotifyclient.GetAppUserToken()
	if err != nil {
		return fmt.Errorf("Could not get Spotify application user token: %s", err)
	}
	err = followPlaylist(ctx, phosphorescenceToken, playlistID)
	if err != nil {
		return fmt.Errorf("Failed to follow playlist: %s", err)
	}
	err = makePlaylistPublic(ctx, phosphorescenceToken, playlistID)
	if err != nil {
		return fmt.Errorf("Failed to make playlist public: %s", err)
	}
	return nil
}

func FollowPlaylist(ctx context.Context, sess *session.Session, playlistID string) error {
	return followPlaylist(ctx, sess.SpotifyToken, playlistID)
}

func getPlaylist(ctx context.Context, token *oauth2.Token, region, playlistID string) (*Playlist, error) {
	spotifyPlaylist, err := getSpotifyPlaylist(ctx, token, region, playlistID)
	if err != nil {
		return nil, fmt.Errorf("Could not get Spotify playlist: %s", err)
	}
	playlist := Playlist{
		ID:          spotifyPlaylist.ID,
		Name:        spotifyPlaylist.Name,
		Description: spotifyPlaylist.Description,
		Owner:       spotifyPlaylist.Owner,
		Images:      spotifyPlaylist.Images,
	}
	if spotifyPlaylist.Tracks.Total > maxTracksPerRequest {
		return nil, ErrTooManyTracks
	}
	playlist.Tracks, err = getSpotifyPlaylistTracks(ctx, token, region, spotifyPlaylist)
	if err != nil {
		return nil, fmt.Errorf("Could not get track data for playlist tracks: %s", err)
	}
	playlist.Tracks, err = populateAudioFeatures(ctx, token, region, playlist.Tracks)
	if err != nil {
		return nil, fmt.Errorf("Could not get audio features: %s", err)
	}
	return &playlist, nil
}

func getSpotifyPlaylist(ctx context.Context, token *oauth2.Token, region, playlistID string) (*SpotifyPlaylist, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playlistID)
	if region != "" {
		url = fmt.Sprintf("%s?market=%s", url, region)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not build Spotify playlist request: %s", err)
	}
	token.SetAuthHeader(req)
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
	var playlistData SpotifyPlaylist
	err = json.Unmarshal(body, &playlistData)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify playlist response: %s", err)
	}
	return &playlistData, nil
}

func getSpotifyPlaylistTracks(ctx context.Context, token *oauth2.Token, region string, spotifyPlaylist *SpotifyPlaylist) (trackData []*SpotifyTrackEnvelope, err error) {
	trackPage := spotifyPlaylist.Tracks
	for true {
		for _, playlistTrack := range trackPage.Items {
			if playlistTrack.IsLocal {
				continue
			}
			track := playlistTrack.Track
			if region != "" && !track.IsPlayable {
				continue
			}
			trackData = append(trackData, &SpotifyTrackEnvelope{
				ID:    track.ID,
				Track: &track,
			})
		}
		if trackPage.Next != "" {
			trackPage, err = getSpotifyPlaylistTrackPage(ctx, token, trackPage.Next)
			if err != nil {
				return nil, fmt.Errorf("Could not get next playlist track page: %s", err)
			}
		} else {
			break
		}
	}
	return trackData, nil
}

func getSpotifyPlaylistTrackPage(ctx context.Context, token *oauth2.Token, pageURL string) (SpotifyPlaylistTrackPaging, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
	if err != nil {
		return SpotifyPlaylistTrackPaging{}, fmt.Errorf("Could not build Spotify playlist page request: %s", err)
	}
	token.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return SpotifyPlaylistTrackPaging{}, fmt.Errorf("Could not make Spotify playlist page request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return SpotifyPlaylistTrackPaging{}, fmt.Errorf("Spotify playlist page request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SpotifyPlaylistTrackPaging{}, fmt.Errorf("Could not read Spotify playlist page response: %s", err)
	}
	var playlistPage SpotifyPlaylistTrackPaging
	err = json.Unmarshal(body, &playlistPage)
	if err != nil {
		return SpotifyPlaylistTrackPaging{}, fmt.Errorf("Could not parse Spotify playlist response: %s", err)
	}
	return playlistPage, nil
}

func createPlaylist(ctx context.Context, token *oauth2.Token, firstTrackName string, utcOffsetMinutes int) (string, error) {
	var createPlaylistBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Public      bool   `json:"public"`
	}
	createPlaylistBody.Name = createPlaylistName(firstTrackName)
	createPlaylistBody.Description = createPlaylistDescription(utcOffsetMinutes)
	createPlaylistBody.Public = false
	createPlaylistBodyJSON, err := json.Marshal(createPlaylistBody)
	if err != nil {
		return "", fmt.Errorf("Could not marshal create playlist request body: %s", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", spotifyclient.AppUserSpotifyID()), bytes.NewBuffer(createPlaylistBodyJSON))
	if err != nil {
		return "", fmt.Errorf("Could not build Spotify create playlist request: %s", err)
	}
	token.SetAuthHeader(req)
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
	var createdPlaylist SpotifyPlaylist
	err = json.Unmarshal(body, &createdPlaylist)
	if err != nil {
		return "", fmt.Errorf("Could not parse Spotify create playlist response: %s", err)
	}
	return createdPlaylist.ID, nil
}

func addTracksToPlaylist(ctx context.Context, token *oauth2.Token, playlistID string, trackURIs []string) error {
	var addTracksBody struct {
		URIs []string `json:"uris"`
	}
	for _, trackURI := range trackURIs {
		addTracksBody.URIs = append(addTracksBody.URIs, trackURI)
	}
	addTracksBodyJSON, err := json.Marshal(addTracksBody)
	if err != nil {
		return fmt.Errorf("Could not marshal add tracks request body: %s", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID), bytes.NewBuffer(addTracksBodyJSON))
	if err != nil {
		return fmt.Errorf("Could not build Spotify add tracks request: %s", err)
	}
	token.SetAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Spotify add tracks request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("Spotify add tracks request responded with %d", res.StatusCode)
	}
	return nil
}

func setPlaylistImage(ctx context.Context, token *oauth2.Token, playlistID string, base64Image string) error {
	var buf *bytes.Buffer
	if base64Image != "" {
		buf = bytes.NewBuffer([]byte(base64Image))
	} else {
		buf = bytes.NewBuffer([]byte(playlistImageBase64))
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/images", playlistID), buf)
	if err != nil {
		return fmt.Errorf("Could not build Spotify change image request: %s", err)
	}
	token.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Spotify change image request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Spotify change image request responded with %d", res.StatusCode)
	}
	return nil
}

func followPlaylist(ctx context.Context, token *oauth2.Token, playlistID string) error {
	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/followers", playlistID), nil)
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

func unfollowPlaylist(ctx context.Context, token *oauth2.Token, playlistID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/followers", playlistID), nil)
	if err != nil {
		return fmt.Errorf("Could not build Spotify unfollow playlist request: %s", err)
	}
	token.SetAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Spotify unfollow playlist request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Spotify unfollow playlist request responded with %d", res.StatusCode)
	}
	return nil
}

func makePlaylistPublic(ctx context.Context, token *oauth2.Token, playlistID string) error {
	var updatePlaylistBody struct {
		Public bool `json:"public"`
	}
	updatePlaylistBody.Public = true
	updatePlaylistBodyJSON, err := json.Marshal(updatePlaylistBody)
	if err != nil {
		return fmt.Errorf("Could not marshal update playlist request body: %s", err)
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playlistID), bytes.NewBuffer(updatePlaylistBodyJSON))
	if err != nil {
		return fmt.Errorf("Could not build Spotify update playlist request: %s", err)
	}
	token.SetAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Spotify update playlist request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Spotify update playlist request responded with %d", res.StatusCode)
	}
	return nil
}

func createPlaylistName(firstTrackName string) string {
	return fmt.Sprintf("phosphor.me | %s", firstTrackName)
}

func createPlaylistDescription(utcOffsetMinutes int) string {
	now := time.Now().UTC().Add(time.Duration(utcOffsetMinutes) * time.Minute)
	return fmt.Sprintf("Created by Phosphorescence on %s at %s. Visit phosphor.me to create more trance playlists!", now.Format("Monday, January _2"), now.Format("3:04 PM"))
}

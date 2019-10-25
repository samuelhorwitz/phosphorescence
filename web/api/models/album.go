package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"golang.org/x/oauth2"
)

type SpotifyAlbum struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Artists []SpotifyArtist `json:"artists"`
	Images  []SpotifyImage  `json:"images"`
}

type SpotifyAlbumTrackPaging struct {
	Next  string         `json:"next"`
	Total int            `json:"total"`
	Items []SpotifyTrack `json:"items"`
}

func GetAlbumTracks(ctx context.Context, region, albumID string) ([]*SpotifyTrackEnvelope, error) {
	phosphorescenceToken, err := spotifyclient.GetAppToken()
	if err != nil {
		return nil, fmt.Errorf("Could not get Spotify application token: %s", err)
	}
	tracks, err := getAlbumTracks(ctx, phosphorescenceToken, region, albumID)
	if err != nil {
		return nil, fmt.Errorf("Could not get album: %s", err)
	}
	return tracks, nil
}

func getAlbumTracks(ctx context.Context, token *oauth2.Token, region, albumID string) ([]*SpotifyTrackEnvelope, error) {
	spotifyAlbumTracks, err := getSpotifyAlbumTracks(ctx, token, albumID)
	if err != nil {
		return nil, fmt.Errorf("Could not get Spotify album: %s", err)
	}
	var trackIDs []string
	for _, track := range spotifyAlbumTracks {
		trackIDs = append(trackIDs, track.ID)
	}
	tracks, err := getTracks(ctx, token, region, trackIDs)
	if err != nil {
		return nil, fmt.Errorf("Could not get tracks: %s", err)
	}
	var cleanedTracks []*SpotifyTrackEnvelope
	for _, unenclosedTrack := range tracks {
		track := unenclosedTrack
		if !track.Track.IsPlayable {
			continue
		}
		cleanedTracks = append(cleanedTracks, track)
	}
	return cleanedTracks, nil
}

func getSpotifyAlbumTracks(ctx context.Context, token *oauth2.Token, albumID string) ([]SpotifyTrack, error) {
	firstPageURL := fmt.Sprintf("https://api.spotify.com/v1/albums/%s/tracks?limit=50", albumID)
	currentPage, err := getSpotifyAlbumTracksPage(ctx, token, firstPageURL)
	if err != nil {
		return nil, fmt.Errorf("Could not get Spotify album tracks first page: %s", err)
	}
	if currentPage.Total > maxTracksPerRequest {
		return nil, ErrTooManyTracks
	}
	var allTracks []SpotifyTrack
	allTracks = append(allTracks, currentPage.Items...)
	for currentPage.Next != "" {
		currentPage, err = getSpotifyAlbumTracksPage(ctx, token, currentPage.Next)
		if err != nil {
			return nil, fmt.Errorf("Could not get Spotify album tracks next page %s: %s", currentPage.Next, err)
		}
		allTracks = append(allTracks, currentPage.Items...)
	}
	return allTracks, nil
}

func getSpotifyAlbumTracksPage(ctx context.Context, token *oauth2.Token, pageURL string) (*SpotifyAlbumTrackPaging, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not build Spotify album request: %s", err)
	}
	token.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not make Spotify album request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Spotify album request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read Spotify album response: %s", err)
	}
	var trackPage SpotifyAlbumTrackPaging
	err = json.Unmarshal(body, &trackPage)
	if err != nil {
		return nil, fmt.Errorf("Could not parse Spotify album response: %s", err)
	}
	return &trackPage, nil
}

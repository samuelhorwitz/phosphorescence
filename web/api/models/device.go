package models

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

type SpotifyDevices struct {
	Devices []SpotifyDevice `json:"devices"`
}

type SpotifyDevice struct {
	ID               string `json:"id"`
	IsActive         bool   `json:"is_active"`
	IsPrivateSession bool   `json:"is_private_session"`
	IsRestricted     bool   `json:"is_restricted"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	VolumePercent    int    `json:"volume_percent"`
}

type Playback struct {
	IsPlaying        bool         `json:"isPlaying"`
	Progress         int          `json:"progress"`
	FetchedAtSpotify int          `json:"fetchedAtSpotify"`
	FetchedAt        int64        `json:"fetchedAt"`
	Track            SpotifyTrack `json:"track"`
}

type PlayState int

const (
	PlayStateUndefined PlayState = iota
	PlayStatePause
	PlayStatePlay
)

var ErrCurrentlyPlayingNotTrack = errors.New("Spotify currently playing is not a track")
var ErrLocalTrack = errors.New("Spotify currently playing is a local track")

func GetDevices(ctx context.Context, sess *session.Session) (SpotifyDevices, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		return SpotifyDevices{}, fmt.Errorf("Could not build Spotify devices request: %s", err)
	}
	sess.SpotifyToken.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return SpotifyDevices{}, fmt.Errorf("Could not make Spotify devices request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return SpotifyDevices{}, fmt.Errorf("Spotify devices request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SpotifyDevices{}, fmt.Errorf("Could not read Spotify devices response: %s", err)
	}
	var devices SpotifyDevices
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return SpotifyDevices{}, fmt.Errorf("Could not parse Spotify devices response: %s", err)
	}
	return devices, nil
}

func TransferPlayback(ctx context.Context, sess *session.Session, deviceID string, playState PlayState) error {
	body, err := json.Marshal(struct {
		DeviceIDs []string `json:"device_ids"`
		Play      bool     `json:"play"`
	}{
		DeviceIDs: []string{deviceID},
		Play:      playState == PlayStatePlay,
	})
	if err != nil {
		return fmt.Errorf("Could not build Spotify device transfer request body: %s", err)
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", "https://api.spotify.com/v1/me/player", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Could not build Spotify device transfer request: %s", err)
	}
	sess.SpotifyToken.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Spotify device transfer request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Spotify device transfer request responded with %d", res.StatusCode)
	}
	if playState == PlayStatePause {
		err = Pause(ctx, sess, deviceID)
		if err != nil {
			return fmt.Errorf("Could not pause after transfering playback: %s", err)
		}
	}
	return nil
}

func Pause(ctx context.Context, sess *session.Session, deviceID string) error {
	body, err := json.Marshal(struct {
		DeviceID string `json:"device_id"`
	}{
		DeviceID: deviceID,
	})
	if err != nil {
		return fmt.Errorf("Could not build Spotify pause request body: %s", err)
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", "https://api.spotify.com/v1/me/player/pause", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Could not build Spotify pause request body: %s", err)
	}
	sess.SpotifyToken.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Spotify pause request: %s", err)
	}
	defer res.Body.Close()
	if !(res.StatusCode == http.StatusNoContent || res.StatusCode == http.StatusNotFound) {
		return fmt.Errorf("Spotify pause request responded with %d", res.StatusCode)
	}
	return nil
}

func GetCurrentPlayback(ctx context.Context, sess *session.Session) (Playback, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		return Playback{}, fmt.Errorf("Could not build Spotify currently playing request: %s", err)
	}
	sess.SpotifyToken.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return Playback{}, fmt.Errorf("Could not make Spotify currently playing request: %s", err)
	}
	fetchedAt := time.Now()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Playback{}, fmt.Errorf("Spotify currently playing request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Playback{}, fmt.Errorf("Could not read Spotify currently playing response: %s", err)
	}
	var parsedBody struct {
		IsPlaying            bool         `json:"is_playing"`
		Track                SpotifyTrack `json:"item"`
		CurrentlyPlayingType string       `json:"currently_playing_type"`
		Timestamp            int          `json:"timestamp"`
		ProgressMilliseconds int          `json:"progress_ms"`
	}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return Playback{}, fmt.Errorf("Could not parse Spotify currently playing response: %s", err)
	}
	if parsedBody.CurrentlyPlayingType != "track" {
		return Playback{}, ErrCurrentlyPlayingNotTrack
	}
	if parsedBody.Track.IsLocal {
		return Playback{}, ErrLocalTrack
	}
	return Playback{
		IsPlaying:        parsedBody.IsPlaying,
		Progress:         parsedBody.ProgressMilliseconds,
		FetchedAtSpotify: parsedBody.Timestamp,
		FetchedAt:        fetchedAt.UnixNano() / int64(time.Millisecond),
		Track:            parsedBody.Track,
	}, nil
}

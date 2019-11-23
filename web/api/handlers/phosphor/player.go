package phosphor

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/handlers"
	"github.com/samuelhorwitz/phosphorescence/api/models"
)

type playerPlaylist struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Owner       models.SpotifyUser    `json:"owner"`
	Images      []models.SpotifyImage `json:"images"`
	Tracks      []playerPlaylistTrack `json:"tracks"`
}

type playerPlaylistTrack struct {
	ID                  string                 `json:"id"`
	Artists             []models.SpotifyArtist `json:"artists"`
	Album               playerPlaylistAlbum    `json:"album"`
	Name                string                 `json:"name"`
	PreviewURL          string                 `json:"preview_url"`
	DurationMillseconds int                    `json:"duration_ms"`
}

type playerPlaylistAlbum struct {
	ID     string                `json:"id"`
	Name   string                `json:"name"`
	Images []models.SpotifyImage `json:"images"`
}

func GetPlayerPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "playlistID")
	if playlistID == "" {
		common.Fail(w, errors.New("Must include playlist ID"), http.StatusBadRequest)
		return
	}
	playlist, err := models.GetSimplePlaylist(r.Context(), playlistID)
	if err != nil {
		code := http.StatusInternalServerError
		if httpErr, ok := err.(handlers.HTTPError); ok {
			code = httpErr.Code
		}
		common.Fail(w, fmt.Errorf("Could not get playlist: %s", err), code)
		return
	}
	var responsePlaylist playerPlaylist
	responsePlaylist.ID = playlist.ID
	responsePlaylist.Name = playlist.Name
	responsePlaylist.Description = playlist.Description
	responsePlaylist.Owner = playlist.Owner
	responsePlaylist.Images = playlist.Images
	for _, track := range playlist.Tracks {
		responsePlaylist.Tracks = append(responsePlaylist.Tracks, playerPlaylistTrack{
			ID:      track.ID,
			Artists: track.Track.Artists,
			Album: playerPlaylistAlbum{
				ID:     track.Track.Album.ID,
				Name:   track.Track.Album.Name,
				Images: track.Track.Album.Images,
			},
			Name:                track.Track.Name,
			PreviewURL:          track.Track.PreviewURL,
			DurationMillseconds: track.Track.DurationMillseconds,
		})
	}
	common.JSON(w, map[string]interface{}{"playlist": responsePlaylist})
}

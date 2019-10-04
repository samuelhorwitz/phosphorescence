package models

type SpotifyImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type SpotifyUser struct {
	ID string `json:"id"`
}

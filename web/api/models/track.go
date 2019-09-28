package models

type SpotifyTrackEnvelope struct {
	ID       string           `json:"id"`
	Track    *SpotifyTrack    `json:"track"`
	Features *SpotifyFeatures `json:"features"`
}

type SpotifyTrack struct {
	Album            SpotifyAlbum        `json:"album"`
	Artists          []SpotifyArtist     `json:"artists"`
	Name             string              `json:"name"`
	Popularity       int                 `json:"popularity"`
	AvailableMarkets []string            `json:"available_markets,omitempty"`
	IsPlayable       bool                `json:"is_playable,omitempty"`
	LinkedFrom       *SpotifyLinkedTrack `json:"linked_from,omitempty"`
}

type SpotifyLinkedTrack struct {
	ID string `json:"id"`
}

type SpotifyAlbum struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Artists []SpotifyArtist `json:"artists"`
	Images  []SpotifyImage  `json:"images"`
}

type SpotifyArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SpotifyImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type SpotifyFeatures struct {
	Danceability        float64 `json:"danceability"`
	Energy              float64 `json:"energy"`
	Key                 int     `json:"key"`
	Loudness            float64 `json:"loudness"`
	Mode                int     `json:"mode"`
	Speechiness         float64 `json:"speechiness"`
	Acousticness        float64 `json:"acousticness"`
	Instrumentalness    float64 `json:"instrumentalness"`
	Liveness            float64 `json:"liveness"`
	Valence             float64 `json:"valence"`
	Tempo               float64 `json:"tempo"`
	DurationMillseconds int     `json:"duration_ms"`
	TimeSignature       int     `json:"time_signature"`
}

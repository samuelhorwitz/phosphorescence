package spider

type TrackEnvelope struct {
	ID       string   `json:"id"`
	Track    track    `json:"track"`
	Features features `json:"features"`
}

type track struct {
	ID               string   `json:"id,omitempty"`
	Album            album    `json:"album"`
	Artists          []artist `json:"artists"`
	Name             string   `json:"name"`
	Popularity       int      `json:"popularity"`
	AvailableMarkets []string `json:"available_markets,omitempty"`
}

type album struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []artist `json:"artists"`
	Images  []image  `json:"images"`
}

type artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type image struct {
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type features struct {
	ID                  string  `json:"id,omitempty"`
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

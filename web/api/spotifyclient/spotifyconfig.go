package spotifyclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var (
	spotifyConfig   *oauth2.Config
	baseHTTPTimeout time.Duration
)

type Config struct {
	SpotifyClientID string
	SpotifySecret   string
	APIOrigin       string
	BaseHTTPTimeout time.Duration
}

func Initialize(cfg *Config) {
	spotifyConfig = &oauth2.Config{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifySecret,
		Scopes: []string{
			"streaming",
			"user-read-birthdate",
			"user-read-email",
			"user-read-private",
			"user-read-playback-state",
			"user-read-recently-played",
		},
		Endpoint:    spotify.Endpoint,
		RedirectURL: fmt.Sprintf("%s/spotify/authorize/redirect", cfg.APIOrigin),
	}
	baseHTTPTimeout = cfg.BaseHTTPTimeout
}

func AuthCodeURL(state string) string {
	return spotifyConfig.AuthCodeURL(state)
}

func TokenExchange(code string) (*oauth2.Token, error) {
	return spotifyConfig.Exchange(newSpotifyHTTPClientContext(), code)
}

func GetToken(token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := spotifyConfig.TokenSource(newSpotifyHTTPClientContext(), token).Token()
	if err != nil {
		return nil, fmt.Errorf("Could not get access token: %s", err)
	}
	return newToken, nil
}

func newSpotifyHTTPClientContext() context.Context {
	return context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
		Timeout: baseHTTPTimeout,
	})
}

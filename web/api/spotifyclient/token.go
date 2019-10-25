package spotifyclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/spotify"
)

var (
	spotifyUserConfig    *oauth2.Config
	spotifyAppUserConfig *oauth2.Config
	spotifyAppConfig     *clientcredentials.Config
	appUserToken         *oauth2.Token
	tokenHTTPClient      *http.Client
	appUserSpotifyID     string
)

type Config struct {
	SpotifyClientID             string
	SpotifySecret               string
	APIOrigin                   string
	BaseHTTPTimeout             time.Duration
	PhosphorescenceSpotifyID    string
	PhosphorescenceRefreshToken string
}

func Initialize(cfg *Config) {
	spotifyUserConfig = &oauth2.Config{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifySecret,
		Scopes:       userScopes(),
		Endpoint:     spotify.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/spotify/authorize/redirect", cfg.APIOrigin),
	}
	spotifyAppUserConfig = &oauth2.Config{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifySecret,
		Scopes:       appUserScopes(),
		Endpoint:     spotify.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/spotify/authorize/redirect", cfg.APIOrigin),
	}
	spotifyAppConfig = &clientcredentials.Config{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifySecret,
		TokenURL:     spotify.Endpoint.TokenURL,
	}
	tokenHTTPClient = &http.Client{
		Timeout: cfg.BaseHTTPTimeout,
	}
	appUserSpotifyID = cfg.PhosphorescenceSpotifyID
	appUserToken = &oauth2.Token{
		RefreshToken: cfg.PhosphorescenceRefreshToken,
	}
	if _, err := GetAppUserToken(); err != nil {
		log.Fatalf("Could not get Phosphorescence user token from refresh token: %s", err)
		return
	}
}

func AuthCodeURL(state string) string {
	return spotifyUserConfig.AuthCodeURL(state)
}

func TokenExchange(code string) (*oauth2.Token, error) {
	return spotifyUserConfig.Exchange(newSpotifyHTTPClientContext(), code)
}

// GetToken is for managing end user Spotify sessions
func GetToken(token *oauth2.Token) (*oauth2.Token, error) {
	return getToken(spotifyUserConfig, token)
}

func getToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := config.TokenSource(newSpotifyHTTPClientContext(), token).Token()
	if err != nil {
		return nil, fmt.Errorf("Could not get user token: %s", err)
	}
	return newToken, nil
}

// GetAppUserToken is for managing the Phosphorescence user, necessary for things
// like creating playlists under the application's name (not really, but from a
// user perspective it appears this way due to the Phosphorescence user branding).
func GetAppUserToken() (*oauth2.Token, error) {
	var err error
	appUserToken, err = getToken(spotifyAppUserConfig, appUserToken)
	if err != nil {
		return nil, fmt.Errorf("Could not get app user token: %s", err)
	}
	return appUserToken, nil
}

// GetAppToken is for pure server-to-server application stuff which doesn't need
// to be tied to a user.
func GetAppToken() (*oauth2.Token, error) {
	token, err := spotifyAppConfig.Token(newSpotifyHTTPClientContext())
	if err != nil {
		return nil, fmt.Errorf("Could not get app token: %s", err)
	}
	return token, nil
}

func AppUserSpotifyID() string {
	return appUserSpotifyID
}

func userScopes() []string {
	return []string{
		"streaming",
		"user-read-email",
		"user-read-private",
		"user-read-playback-state",
		"user-read-recently-played",
		"user-modify-playback-state",
		"playlist-modify-public",
	}
}

// this exists solely for future reference if we need
// a new token for the Phosphorescence user to create
// private playlists for itself
func appUserScopes() []string {
	return []string{
		"playlist-modify-public",
		"playlist-read-private",
		"playlist-modify-private",
		"ugc-image-upload",
	}
}

func newSpotifyHTTPClientContext() context.Context {
	return context.WithValue(context.Background(), oauth2.HTTPClient, tokenHTTPClient)
}

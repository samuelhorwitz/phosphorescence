package session

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	sessionExpiration      = 3 * time.Hour
	permanentExpiration    = 10 * 365 * 24 * time.Hour
	authRedirectExpiration = 2 * time.Minute
	magicLinkExpiration    = 10 * time.Minute
)

type Session struct {
	ID string
	// A session is not authenticated just because it has Spotify details and a token
	// Authentication means the user has been authenticated for our platform and the
	// OAuth2 authorization for Spotify does not count.
	Authenticated  bool   `redis:"authenticated"`
	Permanent      bool   `redis:"permanent"`
	SpotifyID      string `redis:"spotify_id"`
	SpotifyName    string `redis:"spotify_name"`
	SpotifyCountry string `redis:"spotify_country"`
	SpotifyToken   *oauth2.Token
}

type rawSession struct {
	Session
	SpotifyAccessToken  string `redis:"spotify_access_token"`
	SpotifyRefreshToken string `redis:"spotify_refresh_token"`
	SpotifyTokenExpiry  int64  `redis:"spotify_token_expiry"`
}

func (r rawSession) session(id string) (_ *Session, isRefreshed bool, err error) {
	s := r.Session
	s.SpotifyToken, isRefreshed, err = r.token()
	if err != nil {
		return nil, false, fmt.Errorf("Could not rehydrate and refresh session token: %s", err)
	}
	s.ID = id
	return &s, isRefreshed, nil
}

func (r rawSession) token() (*oauth2.Token, bool, error) {
	return refreshIfNeeded(&oauth2.Token{
		AccessToken:  r.SpotifyAccessToken,
		RefreshToken: r.SpotifyRefreshToken,
		Expiry:       time.Unix(r.SpotifyTokenExpiry, 0),
	})
}

func (s Session) GetEmail(r *http.Request) (string, error) {
	user, err := getUser(r, s.SpotifyToken)
	if err != nil {
		return "", fmt.Errorf("Could not get user email from Spotify: %s", err)
	}
	return user.Email, nil
}

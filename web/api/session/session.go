package session

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"time"
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

type user struct {
	ID      string `json:"id"`
	Name    string `json:"display_name"`
	Country string `json:"country"`
}

const (
	sessionExpiration   = 3 * time.Hour
	refreshExpiration   = 4 * time.Hour
	permanentExpiration = 10 * 365 * 24 * time.Hour
)

const (
	sessionCookieName = "sid"
	refreshCookieName = "ref"
)

var (
	cookieDomain string
)

type Config struct {
	CookieDomain string
}

func Initialize(cfg *Config) {
	cookieDomain = cfg.CookieDomain
}

func Create(w http.ResponseWriter, r *http.Request, token *oauth2.Token, permanent bool) error {
	_, _, err := createSession(w, r, token, permanent)
	if err != nil {
		return fmt.Errorf("Could not create session: %s", err)
	}
	return nil
}

func Find(w http.ResponseWriter, r *http.Request, sessionID string, refreshID string) (*Session, error) {
	if sessionID != "" {
		sess, err := liveSession(sessionID, refreshID)
		if err == nil {
			return sess, nil
		}
	}
	if refreshID != "" {
		return reviveSession(w, r, refreshID)
	}
	clearCookies(w)
	return nil, errors.New("No session or refresh IDs")
}

func Destroy(w http.ResponseWriter, sess *Session, refreshID string) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	redisConn.Send("MULTI")
	redisConn.Send("DEL", getSessionKey(sess.ID))
	redisConn.Send("DEL", getRefreshKey(refreshID))
	redisConn.Do("EXEC") // Ignore Redis failure, delete cookies no matter what
	clearCookies(w)
}

func createSession(w http.ResponseWriter, r *http.Request, token *oauth2.Token, permanent bool) (_ string, _ string, err error) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	sessionIDBytes := make([]byte, 32)
	_, err = rand.Read(sessionIDBytes)
	if err != nil {
		return "", "", fmt.Errorf("Could not create session id: %s", err)
	}
	sessionID := hex.EncodeToString(sessionIDBytes)
	refreshIDBytes := make([]byte, 32)
	_, err = rand.Read(refreshIDBytes)
	if err != nil {
		return "", "", fmt.Errorf("Could not create refresh id: %s", err)
	}
	refreshID := hex.EncodeToString(refreshIDBytes)
	user, err := getUser(r, token)
	if err != nil {
		return "", "", fmt.Errorf("Could not get user data from Spotify: %s", err)
	}
	sessionKey := getSessionKey(sessionID)
	refreshKey := getRefreshKey(refreshID)
	redisConn.Send("MULTI")
	redisConn.Send("HMSET", sessionKey,
		"spotify_id", user.ID,
		"spotify_name", user.Name,
		"spotify_country", user.Country,
		"spotify_access_token", token.AccessToken,
		"spotify_refresh_token", token.RefreshToken,
		"spotify_token_expiry", token.Expiry.Unix(),
		"permanent", permanent,
	)
	redisConn.Send("EXPIRE", sessionKey, uint64(sessionExpiration.Seconds()))
	redisConn.Send("HMSET", refreshKey,
		"spotify_access_token", token.AccessToken,
		"spotify_refresh_token", token.RefreshToken,
		"spotify_token_expiry", token.Expiry.Unix(),
		"permanent", permanent,
	)
	if !permanent {
		redisConn.Send("EXPIRE", refreshKey, uint64(refreshExpiration.Seconds()))
	}
	_, err = redisConn.Do("EXEC")
	if err != nil {
		return "", "", fmt.Errorf("Could not save session: %s", err)
	}
	setCookies(w, sessionID, refreshID, permanent)
	return sessionID, refreshID, nil
}

func liveSession(sessionID string, refreshID string) (*Session, error) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	sessionKey := getSessionKey(sessionID)
	refreshKey := getRefreshKey(refreshID)
	var rawSess rawSession
	sessionData, err := redis.Values(redisConn.Do("HGETALL", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("Could not find live session: %s", err)
	}
	err = redis.ScanStruct(sessionData, &rawSess)
	if err != nil {
		return nil, fmt.Errorf("Could not parse live session: %s", err)
	}
	sess, isRefreshed, err := rawSess.session(sessionID)
	if err != nil {
		return nil, fmt.Errorf("Could not rehydrate live session: %s", err)
	}
	if isRefreshed {
		redisConn.Send("MULTI")
		redisConn.Send("HMSET", sessionKey,
			"spotify_access_token", sess.SpotifyToken.AccessToken,
			"spotify_refresh_token", sess.SpotifyToken.RefreshToken,
			"spotify_token_expiry", sess.SpotifyToken.Expiry.Unix(),
		)
		if refreshID != "" {
			redisConn.Send("HMSET", refreshKey,
				"spotify_access_token", sess.SpotifyToken.AccessToken,
				"spotify_refresh_token", sess.SpotifyToken.RefreshToken,
				"spotify_token_expiry", sess.SpotifyToken.Expiry.Unix(),
			)
		}
		_, err = redisConn.Do("EXEC")
		if err != nil {
			return nil, fmt.Errorf("Could not update live session: %s", err)
		}
	}
	return sess, nil
}

func reviveSession(w http.ResponseWriter, r *http.Request, refreshID string) (*Session, error) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	refreshKey := getRefreshKey(refreshID)
	refreshData, err := redis.Values(redisConn.Do("HGETALL", refreshKey))
	if err != nil {
		return nil, fmt.Errorf("No refresh data: %s", err)
	}
	_, err = redisConn.Do("DEL", refreshKey)
	if err != nil {
		return nil, fmt.Errorf("Could not delete refresh data: %s", err)
	}
	var refresh rawSession
	err = redis.ScanStruct(refreshData, &refresh)
	if err != nil {
		return nil, fmt.Errorf("Could not scan refresh data: %s", err)
	}
	token, _, err := refresh.token()
	if err != nil {
		return nil, fmt.Errorf("Could not get token from refresh data: %s", err)
	}
	sessionID, newRefreshID, err := createSession(w, r, token, refresh.Permanent)
	if err != nil {
		return nil, fmt.Errorf("Could not revive session: %s", err)
	}
	return liveSession(sessionID, newRefreshID)
}

func setCookies(w http.ResponseWriter, sessionID string, refreshID string, permanent bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	refreshCookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    refreshID,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	if permanent {
		refreshCookie.Expires = time.Now().Add(permanentExpiration)
		refreshCookie.MaxAge = int(permanentExpiration.Seconds())
	}
	http.SetCookie(w, refreshCookie)
}

func clearCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func getSessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

func getRefreshKey(refreshID string) string {
	return fmt.Sprintf("refresh:%s", refreshID)
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

func refreshIfNeeded(token *oauth2.Token) (*oauth2.Token, bool, error) {
	newToken, err := spotifyclient.GetToken(token)
	if err != nil {
		return nil, false, fmt.Errorf("Could not refresh session token: %s", err)
	}
	return newToken, newToken.AccessToken != token.AccessToken, nil
}

func getUser(r *http.Request, token *oauth2.Token) (user, error) {
	req, err := http.NewRequestWithContext(r.Context(), "GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return user{}, fmt.Errorf("Could not build Spotify profile request: %s", err)
	}
	token.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return user{}, fmt.Errorf("Could not make Spotify profile request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return user{}, fmt.Errorf("Spotify profile request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return user{}, fmt.Errorf("Could not read Spotify profile response: %s", err)
	}
	var parsedBody user
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return user{}, fmt.Errorf("Could not parse Spotify profile response: %s", err)
	}
	return parsedBody, nil
}

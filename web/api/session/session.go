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
	Email   string `json:"email"`
}

const (
	sessionExpiration      = 3 * time.Hour
	permanentExpiration    = 10 * 365 * 24 * time.Hour
	authRedirectExpiration = 2 * time.Minute
	magicLinkExpiration    = 10 * time.Minute
)

const (
	SessionCookieName = "sid"
	RefreshCookieName = "ref"
)

const (
	fixedSessionPrefix = "fixedsession"
	sessionPointerKey  = "session_pointer"
	refreshPointerKey  = "refresh_pointer"
)

var (
	cookieDomain string
	isProduction bool
)

type Config struct {
	CookieDomain string
	IsProduction bool
}

func Initialize(cfg *Config) {
	cookieDomain = cfg.CookieDomain
	isProduction = cfg.IsProduction
	initializeReaper()
}

func Create(w http.ResponseWriter, r *http.Request, token *oauth2.Token, permanent bool) error {
	_, err := createSession(w, r, token, permanent)
	if err != nil {
		return fmt.Errorf("Could not create session: %s", err)
	}
	return nil
}

func Find(w http.ResponseWriter, r *http.Request, sessionID string, refreshID string) (*Session, error) {
	if sessionID != "" {
		sess, err := liveSession(sessionID)
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

func Destroy(w http.ResponseWriter, sessionID string, refreshID string) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	sessionKey := getSessionPointerKey(sessionID)
	fixedSessionKey, _ := redis.String(redisConn.Do("GET", sessionKey))
	redisConn.Send("MULTI")
	redisConn.Send("DEL", sessionKey)
	redisConn.Send("DEL", getRefreshKey(refreshID))
	redisConn.Send("DEL", fixedSessionKey)
	redisConn.Do("EXEC") // Ignore Redis failure, delete cookies no matter what
	clearCookies(w)
}

func CreateMagicLink(sess *Session) (string, error) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	sessionKey := getSessionPointerKey(sess.ID)
	fixedSessionKey, err := redis.String(redisConn.Do("GET", sessionKey))
	if err != nil {
		return "", fmt.Errorf("Could not get session from pointer: %s", err)
	}
	magicLinkBytes := make([]byte, 32)
	_, err = rand.Read(magicLinkBytes)
	if err != nil {
		return "", fmt.Errorf("Could not create magic link: %s", err)
	}
	magicLink := hex.EncodeToString(magicLinkBytes)
	magicLinkKey := getMagicLinkKey(magicLink)
	_, err = redisConn.Do("SETEX", magicLinkKey, uint64((magicLinkExpiration).Seconds()), fixedSessionKey)
	if err != nil {
		return "", fmt.Errorf("Could not add magic link to Redis: %s", err)
	}
	return magicLink, nil
}

func Upgrade(w http.ResponseWriter, sess *Session, magicLink string) error {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	magicLinkKey := getMagicLinkKey(magicLink)
	var sessionID string
	if sess != nil {
		sessionID = sess.ID
	}
	handleSessionKey := getSessionPointerKey(sessionID)
	handlerFixedSessionKey, _ := redis.String(redisConn.Do("GET", handleSessionKey))
	initiatorFixedSessionKey, err := redis.String(redisConn.Do("GET", magicLinkKey))
	if err != nil {
		return fmt.Errorf("Could not find magic link in Redis: %s", err)
	}
	_, err = redisConn.Do("DEL", magicLinkKey)
	if err != nil {
		return fmt.Errorf("Could not delete magic link from Redis: %s", err)
	}
	err = authenticateSession(initiatorFixedSessionKey)
	if err != nil {
		return fmt.Errorf("Could not authenticate initiator session: %s", err)
	}
	if handlerFixedSessionKey == initiatorFixedSessionKey {
		return nil
	}
	err = cloneSession(w, initiatorFixedSessionKey)
	if err != nil {
		return fmt.Errorf("Could not clone initiator session for handler: %s", err)
	}
	return nil
}

func CreateAuthRedirect(permanent bool) (string, error) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	stateBytes := make([]byte, 32)
	_, err := rand.Read(stateBytes)
	if err != nil {
		return "", fmt.Errorf("Could not create state token: %s", err)
	}
	state := hex.EncodeToString(stateBytes)
	stateKey := getStateKey(state)
	_, err = redisConn.Do("SETEX", stateKey, uint64((authRedirectExpiration).Seconds()), permanent)
	if err != nil {
		return "", fmt.Errorf("Could not add state token to Redis: %s", err)
	}
	return state, nil
}

func CheckAuthRedirect(state string) (bool, error) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	stateKey := getStateKey(state)
	isPermanent, err := redis.Bool(redisConn.Do("GET", stateKey))
	if err != nil {
		return false, fmt.Errorf("Could not find state in Redis: %s", err)
	}
	_, err = redisConn.Do("DEL", stateKey)
	if err != nil {
		return false, fmt.Errorf("Could not delete state in Redis: %s", err)
	}
	return isPermanent, nil
}

func createSession(w http.ResponseWriter, r *http.Request, token *oauth2.Token, permanent bool) (_ string, err error) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	fixedSessionIDBytes := make([]byte, 32)
	_, err = rand.Read(fixedSessionIDBytes)
	if err != nil {
		return "", fmt.Errorf("Could not create fixed session id: %s", err)
	}
	fixedSessionID := hex.EncodeToString(fixedSessionIDBytes)
	user, err := getUser(r, token)
	if err != nil {
		return "", fmt.Errorf("Could not get user data from Spotify: %s", err)
	}
	fixedSessionKey := getSessionKey(user.ID, fixedSessionID)
	_, err = redisConn.Do("HMSET", fixedSessionKey,
		"spotify_id", user.ID,
		"spotify_name", user.Name,
		"spotify_country", user.Country,
		"spotify_access_token", token.AccessToken,
		"spotify_refresh_token", token.RefreshToken,
		"spotify_token_expiry", token.Expiry.Unix(),
		"permanent", permanent,
	)
	if err != nil {
		return "", fmt.Errorf("Could not save session: %s", err)
	}
	return createSessionPointer(w, fixedSessionKey, permanent)
}

func createSessionPointer(w http.ResponseWriter, fixedSessionKey string, permanent bool) (_ string, err error) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	sessionIDBytes := make([]byte, 32)
	_, err = rand.Read(sessionIDBytes)
	if err != nil {
		return "", fmt.Errorf("Could not create session id: %s", err)
	}
	sessionID := hex.EncodeToString(sessionIDBytes)
	refreshIDBytes := make([]byte, 32)
	_, err = rand.Read(refreshIDBytes)
	if err != nil {
		return "", fmt.Errorf("Could not create refresh id: %s", err)
	}
	refreshID := hex.EncodeToString(refreshIDBytes)
	sessionKey := getSessionPointerKey(sessionID)
	refreshKey := getRefreshKey(refreshID)
	redisConn.Send("MULTI")
	redisConn.Send("HMGET", fixedSessionKey, sessionPointerKey, refreshPointerKey)
	redisConn.Send("SETEX", sessionKey, uint64(sessionExpiration.Seconds()), fixedSessionKey)
	redisConn.Send("SET", refreshKey, fixedSessionKey)
	if !permanent {
		redisConn.Send("EXPIRE", refreshKey, uint64((sessionExpiration + (1 * time.Hour)).Seconds()))
	}
	redisConn.Send("HMSET", fixedSessionKey, sessionPointerKey, sessionKey, refreshPointerKey, refreshKey)
	res, err := redis.Values(redisConn.Do("EXEC"))
	if err != nil {
		return "", fmt.Errorf("Could not save session pointers: %s", err)
	}
	pointers, err := redis.Values(res[0], nil)
	if err != nil {
		return "", fmt.Errorf("Could not get old session pointers: %s", err)
	}
	_, err = redisConn.Do("DEL", pointers...)
	if err != nil {
		return "", fmt.Errorf("Could not delete old session pointers: %s", err)
	}
	setCookies(w, sessionID, refreshID, permanent)
	return sessionID, nil
}

func cloneSession(w http.ResponseWriter, fixedSessionKeyToClone string) error {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	spotifyUserID, err := redis.String(redisConn.Do("HGET", fixedSessionKeyToClone, "spotify_id"))
	if err != nil {
		return fmt.Errorf("Could not get Spotify user ID: %s", err)
	}
	newFixedSessionIDBytes := make([]byte, 32)
	_, err = rand.Read(newFixedSessionIDBytes)
	if err != nil {
		return fmt.Errorf("Could not create new fixed session id: %s", err)
	}
	newFixedSessionID := hex.EncodeToString(newFixedSessionIDBytes)
	newFixedSessionKey := getSessionKey(spotifyUserID, newFixedSessionID)
	oldFixedSession, err := redis.String(redisConn.Do("DUMP", fixedSessionKeyToClone))
	if err != nil {
		return fmt.Errorf("Could not dump old fixed session: %s", err)
	}
	redisConn.Send("MULTI")
	redisConn.Send("RESTORE", newFixedSessionKey, 0, oldFixedSession)
	redisConn.Send("HDEL", newFixedSessionKey, sessionPointerKey)
	redisConn.Send("HDEL", newFixedSessionKey, refreshPointerKey)
	_, err = redisConn.Do("EXEC")
	if err != nil {
		return fmt.Errorf("Could not restore old fixed session to new key: %s", err)
	}
	_, err = createSessionPointer(w, newFixedSessionKey, true)
	if err != nil {
		return fmt.Errorf("Could not create session pointer: %s", err)
	}
	return nil
}

func authenticateSession(fixedSessionKey string) error {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	_, err := redisConn.Do("HSET", fixedSessionKey, "authenticated", true)
	if err != nil {
		return fmt.Errorf("Could not authenticate session: %s", err)
	}
	sessionPointerID, _ := redis.String(redisConn.Do("HGET", fixedSessionKey, sessionPointerKey))
	_, err = redisConn.Do("DEL", sessionPointerID)
	if err != nil {
		return fmt.Errorf("Could not delete sessions: %s", err)
	}
	return nil
}

func liveSession(sessionID string) (*Session, error) {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	sessionKey := getSessionPointerKey(sessionID)
	fixedSessionKey, err := redis.String(redisConn.Do("GET", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("Could not get session from pointer: %s", err)
	}
	var rawSess rawSession
	sessionData, err := redis.Values(redisConn.Do("HGETALL", fixedSessionKey))
	if err != nil {
		return nil, fmt.Errorf("Could not get session data: %s", err)
	}
	err = redis.ScanStruct(sessionData, &rawSess)
	if err != nil {
		return nil, fmt.Errorf("Could not parse fixed session: %s", err)
	}
	sess, isRefreshed, err := rawSess.session(sessionID)
	if err != nil {
		return nil, fmt.Errorf("Could not rehydrate live session: %s", err)
	}
	if isRefreshed {
		_, err = redisConn.Do("HMSET", fixedSessionKey,
			"spotify_access_token", sess.SpotifyToken.AccessToken,
			"spotify_refresh_token", sess.SpotifyToken.RefreshToken,
			"spotify_token_expiry", sess.SpotifyToken.Expiry.Unix(),
		)
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
	fixedSessionKey, err := redis.String(redisConn.Do("GET", refreshKey))
	if err != nil {
		return nil, fmt.Errorf("Could not get session from refresh pointer: %s", err)
	}
	_, err = redisConn.Do("DEL", refreshKey)
	if err != nil {
		return nil, fmt.Errorf("Could not delete refresh data: %s", err)
	}
	sessionData, err := redis.Values(redisConn.Do("HGETALL", fixedSessionKey))
	if err != nil {
		return nil, fmt.Errorf("Could not get session data: %s", err)
	}
	var rawSess rawSession
	err = redis.ScanStruct(sessionData, &rawSess)
	if err != nil {
		return nil, fmt.Errorf("Could not parse fixed session: %s", err)
	}
	sessionID, err := createSessionPointer(w, fixedSessionKey, rawSess.Permanent)
	if err != nil {
		return nil, fmt.Errorf("Could not create session pointers: %s", err)
	}
	return liveSession(sessionID)
}

func setCookies(w http.ResponseWriter, sessionID string, refreshID string, permanent bool) {
	sessionCookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	if isProduction {
		sessionCookie.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(w, sessionCookie)
	refreshCookie := &http.Cookie{
		Name:     RefreshCookieName,
		Value:    refreshID,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	if permanent {
		refreshCookie.Expires = time.Now().Add(permanentExpiration)
		refreshCookie.MaxAge = int(permanentExpiration.Seconds())
	}
	if isProduction {
		refreshCookie.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(w, refreshCookie)
}

func clearCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func getSessionKey(spotifyUserID, fixedSessionID string) string {
	return fmt.Sprintf("%s:%s:%s", fixedSessionPrefix, spotifyUserID, fixedSessionID)
}

func getSessionPointerKey(sessionID string) string {
	return fmt.Sprintf("sessionptr:%s", sessionID)
}

func getRefreshKey(refreshID string) string {
	return fmt.Sprintf("refreshptr:%s", refreshID)
}

func getMagicLinkKey(magicLink string) string {
	return fmt.Sprintf("magiclink:%s", magicLink)
}

func getStateKey(state string) string {
	return fmt.Sprintf("state:%s", state)
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
	var parsedUser user
	err = json.Unmarshal(body, &parsedUser)
	if err != nil {
		return user{}, fmt.Errorf("Could not parse Spotify profile response: %s", err)
	}
	return parsedUser, nil
}

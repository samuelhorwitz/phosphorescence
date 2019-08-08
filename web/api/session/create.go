package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"golang.org/x/oauth2"
)

const (
	fixedSessionPrefix = "fixedsession"
	sessionPointerKey  = "session_pointer"
	refreshPointerKey  = "refresh_pointer"
)

func Create(w http.ResponseWriter, r *http.Request, token *oauth2.Token, permanent bool) error {
	_, err := createSession(w, r, token, permanent)
	if err != nil {
		return fmt.Errorf("Could not create session: %s", err)
	}
	return nil
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

func getSessionKey(spotifyUserID, fixedSessionID string) string {
	return fmt.Sprintf("%s:%s:%s", fixedSessionPrefix, spotifyUserID, fixedSessionID)
}

func getSessionPointerKey(sessionID string) string {
	return fmt.Sprintf("sessionptr:%s", sessionID)
}

func getRefreshKey(refreshID string) string {
	return fmt.Sprintf("refreshptr:%s", refreshID)
}

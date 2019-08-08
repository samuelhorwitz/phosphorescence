package session

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/samuelhorwitz/phosphorescence/api/common"
)

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

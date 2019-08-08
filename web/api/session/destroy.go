package session

import (
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/samuelhorwitz/phosphorescence/api/common"
)

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

func DestroyAllByUser(w http.ResponseWriter, sess *Session) error {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	sessionID := sess.ID
	sessionKey := getSessionPointerKey(sessionID)
	currentFixedSessionKey, err := redis.String(redisConn.Do("GET", sessionKey))
	if err != nil {
		return fmt.Errorf("Could not get current fixed session key: %s", err)
	}
	spotifyUserID := sess.SpotifyID
	iter := 0
	for {
		scanResult, err := redis.Values(redisConn.Do("SCAN", iter, "MATCH", fmt.Sprintf("%s:%s:*", fixedSessionPrefix, spotifyUserID), "COUNT", 1000))
		if err != nil {
			return fmt.Errorf("Could not scan Redis for user sessions: %s", err)
		}
		iter, err = redis.Int(scanResult[0], nil)
		if err != nil {
			return fmt.Errorf("Could not get new iterator: %s", err)
		}
		fixedSessionKeys, err := redis.Strings(scanResult[1], nil)
		if err != nil {
			return fmt.Errorf("Could not get batch of sessions: %s", err)
		}
		for _, fixedSessionKey := range fixedSessionKeys {
			if fixedSessionKey == currentFixedSessionKey {
				continue
			}
			err = deleteFixedSessionAndAssociated(redisConn, fixedSessionKey)
			if err != nil {
				return fmt.Errorf("Could not delete fixed session: %s", err)
			}
		}
		if iter == 0 {
			break
		}
	}
	err = deleteFixedSessionAndAssociated(redisConn, currentFixedSessionKey)
	if err != nil {
		return fmt.Errorf("Could not delete current fixed session: %s", err)
	}
	clearCookies(w)
	return nil
}

func deleteFixedSessionAndAssociated(redisConn redis.Conn, fixedSessionKey string) error {
	pointers, err := redis.Values(redisConn.Do("HMGET", fixedSessionKey, sessionPointerKey, refreshPointerKey))
	if err != nil {
		return fmt.Errorf("Could not get session pointers: %s", err)
	}
	var toDelete []interface{}
	toDelete = append(toDelete, fixedSessionKey)
	toDelete = append(toDelete, pointers...)
	_, err = redisConn.Do("DEL", toDelete...)
	if err != nil {
		return fmt.Errorf("Could not delete session and associated pointers: %s", err)
	}
	return nil
}

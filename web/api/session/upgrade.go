package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/samuelhorwitz/phosphorescence/api/common"
)

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

func getMagicLinkKey(magicLink string) string {
	return fmt.Sprintf("magiclink:%s", magicLink)
}

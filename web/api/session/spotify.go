package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/samuelhorwitz/phosphorescence/api/common"
)

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

func getStateKey(state string) string {
	return fmt.Sprintf("state:%s", state)
}

package session

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"log"
	"time"
)

func initializeReaper() {
	go func() {
		iter := 0
		for {
			time.Sleep(3 * time.Second)
			iter = reap(iter)
		}
	}()
}

func reap(iter int) int {
	redisConn := common.RedisPool.Get()
	defer redisConn.Close()
	scanResult, err := redis.Values(redisConn.Do("SCAN", iter, "MATCH", fmt.Sprintf("%s:*", fixedSessionPrefix)))
	if err != nil {
		if !isProduction {
			log.Printf("Could not scan: %s", err)
		}
		return iter
	}
	newIter, err := redis.Int(scanResult[0], nil)
	if err != nil {
		if !isProduction {
			log.Printf("Could not get new iterator: %s", err)
		}
		return iter
	}
	sessions, err := redis.Strings(scanResult[1], nil)
	if err != nil {
		if !isProduction {
			log.Printf("Could not get sessions: %s", err)
		}
		return iter
	}
	for _, session := range sessions {
		pointers, err := redis.Values(redisConn.Do("HMGET", session, sessionPointerKey, refreshPointerKey))
		if err != nil {
			if !isProduction {
				log.Printf("Could not get session pointers: %s", err)
			}
			continue
		}
		pointees, err := redis.Strings(redisConn.Do("MGET", pointers...))
		if err != nil {
			if !isProduction {
				log.Printf("Could not lookup session pointers: %s", err)
			}
			continue
		}
		refCount := len(pointees)
		for _, pointee := range pointees {
			if pointee == "" {
				refCount--
			}
		}
		if refCount == 0 {
			_, err = redisConn.Do("DEL", session)
			if err != nil {
				if !isProduction {
					log.Printf("Could not delete 0 ref count session: %s", err)
				}
				continue
			}
		}
	}
	return newIter
}

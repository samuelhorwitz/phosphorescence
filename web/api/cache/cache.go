package cache

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

const twelveHoursInSeconds uint64 = 60 * 60 * 12
const idField = "id"
const trackField = "track"
const featuresField = "features"

var (
	isProduction bool
	redisPool    *redis.Pool
)

type CachedTrack struct {
	ID       string
	Track    string
	Features string
}

type CachedPlaylist string

type Config struct {
	IsProduction bool
	RedisHost    string
}

func Initialize(cfg *Config) {
	isProduction = cfg.IsProduction
	redisPool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.RedisHost)
			if err != nil {
				log.Fatalf("Could not connect to Redis: %s", err)
			}
			return c, err
		},
	}
}

func GetTrack(region, id string) (CachedTrack, bool) {
	redisConn := redisPool.Get()
	defer redisConn.Close()
	trackJSONs, err := redis.StringMap(redisConn.Do("HGETALL", getTrackKey(region, id)))
	if err != nil {
		if !isProduction {
			log.Printf("cache miss from lookup error: %s", err)
		}
		return CachedTrack{}, false
	}
	return CachedTrack{
		ID:       trackJSONs[idField],
		Track:    trackJSONs[trackField],
		Features: trackJSONs[featuresField],
	}, true
}

func GetTracks(region string, ids []string) map[string]CachedTrack {
	redisConn := redisPool.Get()
	defer redisConn.Close()
	cachedTracks := make(map[string]CachedTrack)
	for _, id := range ids {
		redisConn.Send("HGETALL", getTrackKey(region, id))
	}
	redisConn.Flush()
	for _, id := range ids {
		trackJSONs, err := redis.StringMap(redisConn.Receive())
		if err != nil {
			if !isProduction {
				log.Printf("cache miss from lookup error: %s", err)
			}
			continue
		}
		cachedTracks[id] = CachedTrack{
			ID:       trackJSONs[idField],
			Track:    trackJSONs[trackField],
			Features: trackJSONs[featuresField],
		}
	}
	return cachedTracks
}

func SetTrack(region, id string, trackEnvelope CachedTrack) bool {
	redisConn := redisPool.Get()
	defer redisConn.Close()
	trackKey := getTrackKey(region, id)
	_, err := redisConn.Do("HMSET", trackKey, idField, trackEnvelope.ID, trackField, trackEnvelope.Track, featuresField, trackEnvelope.Features)
	if err != nil {
		if !isProduction {
			log.Printf("cache warm failed from write error: %s", err)
		}
		return false
	}
	_, err = redisConn.Do("EXPIRE", trackKey, twelveHoursInSeconds)
	if err != nil {
		if !isProduction {
			log.Printf("cache warm failed from expiration error: %s", err)
		}
		return false
	}
	return true
}

func getTrackKey(region, id string) string {
	return fmt.Sprintf("track:%s:%s", region, id)
}

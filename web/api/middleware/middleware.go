package middleware

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gomodule/redigo/redis"
	"log"
)

type contextKey string

var (
	phosphorLimiter *limiter.Limiter
	spotifyLimiter  *limiter.Limiter
	redisPool       *redis.Pool
)

type Config struct {
	RateLimitPerSecond int
	RedisHost          string
}

func Initialize(cfg *Config) {
	phosphorLimiter = tollbooth.NewLimiter(float64(cfg.RateLimitPerSecond), nil)
	spotifyLimiter = tollbooth.NewLimiter(4, nil)
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

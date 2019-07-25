package middleware

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

type contextKey string

var (
	phosphorLimiter *limiter.Limiter
	spotifyLimiter  *limiter.Limiter
	ipLimiter       *limiter.Limiter
)

type Config struct {
	RateLimitPerSecond int
}

func Initialize(cfg *Config) {
	phosphorLimiter = tollbooth.NewLimiter(float64(cfg.RateLimitPerSecond), nil)
	spotifyLimiter = tollbooth.NewLimiter(4, nil)
	ipLimiter = tollbooth.NewLimiter(float64(cfg.RateLimitPerSecond), nil)
}

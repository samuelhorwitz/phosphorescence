package middleware

import (
	"log"
	"net/http"
	"time"

	"net/url"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

type contextKey string

var (
	phosphorLimiter  *limiter.Limiter
	spotifyLimiter   *limiter.Limiter
	ipLimiter        *limiter.Limiter
	googleHTTPClient *http.Client
	recaptchaSecret  string
	phosphorHost     string
	twitterHost      string
	isProduction     bool
)

type Config struct {
	RateLimitPerSecond int
	RecaptchaSecret    string
	PhosphorOrigin     string
	TwitterOrigin      string
	IsProduction       bool
}

func Initialize(cfg *Config) {
	phosphorLimiter = tollbooth.NewLimiter(float64(cfg.RateLimitPerSecond), nil)
	spotifyLimiter = tollbooth.NewLimiter(10, nil)
	ipLimiter = tollbooth.NewLimiter(float64(cfg.RateLimitPerSecond), nil)
	googleHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	recaptchaSecret = cfg.RecaptchaSecret
	phosphorOriginParsed, err := url.Parse(cfg.PhosphorOrigin)
	if err != nil {
		log.Fatalf("Could not parse Phosphorescence origin: %s", err)
		return
	}
	phosphorHost = phosphorOriginParsed.Hostname()
	twitterOriginParsed, err := url.Parse(cfg.TwitterOrigin)
	if err != nil {
		log.Fatalf("Could not parse Phosphorescence Twitter origin: %s", err)
		return
	}
	twitterHost = twitterOriginParsed.Hostname()
	isProduction = cfg.IsProduction
}

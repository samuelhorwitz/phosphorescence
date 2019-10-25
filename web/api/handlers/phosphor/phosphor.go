package phosphor

import (
	"net/http"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

var (
	phosphorOrigin string
	isProduction   bool
	noHTML         *bluemonday.Policy
	safeHTTPClient *http.Client
)

type Config struct {
	PhosphorOrigin string
	IsProduction   bool
}

func Initialize(cfg *Config) {
	phosphorOrigin = cfg.PhosphorOrigin
	isProduction = cfg.IsProduction
	safeHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	noHTML = bluemonday.StrictPolicy()
}

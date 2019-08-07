package phosphor

import (
	"net/http"
	"time"
)

var (
	phosphorOrigin string
	isProduction   bool
	mailgunAPIKey  string
	mailgunClient  *http.Client
)

type Config struct {
	PhosphorOrigin string
	IsProduction   bool
	MailgunAPIKey  string
}

func Initialize(cfg *Config) {
	phosphorOrigin = cfg.PhosphorOrigin
	isProduction = cfg.IsProduction
	mailgunAPIKey = cfg.MailgunAPIKey
	mailgunClient = &http.Client{
		Timeout: 10 * time.Second,
	}
}

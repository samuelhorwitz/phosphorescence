package main

import (
	"time"
)

type config struct {
	isProduction                         bool
	phosphorOrigin                       string
	spotifyClientID                      string
	spotifySecret                        string
	spacesID                             string
	spacesSecret                         string
	spacesTracksEndpoint                 string
	spacesTracksRegion                   string
	spacesScriptsEndpoint                string
	spacesScriptsRegion                  string
	postgresConnectionString             string
	postgresMaxOpenConnections           int
	postgresMaxIdleConnections           int
	postgresMaxConnectionLifetimeMinutes int
	readTimeout                          time.Duration
	writeTimeout                         time.Duration
	idleTimeout                          time.Duration
	rateLimitPerSecond                   int
	redisHost                            string
}

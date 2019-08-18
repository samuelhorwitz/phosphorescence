package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/handlers/phosphor"
	"github.com/samuelhorwitz/phosphorescence/api/handlers/spotify"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"github.com/samuelhorwitz/phosphorescence/api/tracks"
)

func main() {
	isProduction := os.Getenv("ENV") == "production"
	if !isProduction {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Could not load .env file: %s", err)
			return
		}
	}
	pgMaxOpen, err := strconv.Atoi(os.Getenv("PG_MAX_OPEN_CONNS"))
	if err != nil {
		log.Fatalf("Could not parse PG max open conns: %s", err)
		return
	}
	pgMaxIdle, err := strconv.Atoi(os.Getenv("PG_MAX_IDLE_CONNS"))
	if err != nil {
		log.Fatalf("Could not parse PG max idle conns: %s", err)
		return
	}
	pgLifetime, err := strconv.Atoi(os.Getenv("PG_MAX_CONN_LIFETIME_MINUTES"))
	if err != nil {
		log.Fatalf("Could not parse PG conn lifetime: %s", err)
		return
	}
	rateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT_PER_SECOND"))
	if err != nil {
		log.Fatalf("Could not parse rate limit: %s", err)
		return
	}
	cfg := &config{
		isProduction:                         isProduction,
		phosphorOrigin:                       os.Getenv("PHOSPHOR_ORIGIN"),
		apiOrigin:                            os.Getenv("API_ORIGIN"),
		cookieDomain:                         os.Getenv("COOKIE_DOMAIN"),
		spotifyClientID:                      os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifySecret:                        os.Getenv("SPOTIFY_SECRET"),
		spacesID:                             os.Getenv("SPACES_ID"),
		spacesSecret:                         os.Getenv("SPACES_SECRET"),
		spacesTracksEndpoint:                 os.Getenv("SPACES_TRACKS_ENDPOINT"),
		spacesTracksRegion:                   os.Getenv("SPACES_TRACKS_REGION"),
		spacesScriptsEndpoint:                os.Getenv("SPACES_SCRIPTS_ENDPOINT"),
		spacesScriptsRegion:                  os.Getenv("SPACES_SCRIPTS_REGION"),
		postgresConnectionString:             os.Getenv("PG_CONNECTION_STRING"),
		postgresMaxOpenConnections:           pgMaxOpen,
		postgresMaxIdleConnections:           pgMaxIdle,
		postgresMaxConnectionLifetimeMinutes: pgLifetime,
		readTimeout:                          5 * time.Second,
		writeTimeout:                         10 * time.Second,
		idleTimeout:                          120 * time.Second,
		handlerTimeout:                       5 * time.Second,
		rateLimitPerSecond:                   rateLimit,
		redisHost:                            os.Getenv("REDIS_HOST"),
		mailgunAPIKey:                        os.Getenv("MAILGUN_API_KEY"),
	}
	migrate(cfg)
	initialize(cfg)
	run(cfg)
}

func initialize(cfg *config) {
	log.Println("Initializing...")
	rand.Seed(time.Now().UnixNano())
	log.Println("Randomness initialized")
	common.Initialize(&common.Config{
		IsProduction:   cfg.isProduction,
		SpotifyTimeout: cfg.handlerTimeout,
		RedisHost:      cfg.redisHost,
	})
	log.Println("Common initialized")
	middleware.Initialize(&middleware.Config{
		RateLimitPerSecond: cfg.rateLimitPerSecond,
	})
	log.Println("Middleware initialized")
	spotify.Initialize(&spotify.Config{
		IsProduction:   cfg.isProduction,
		PhosphorOrigin: cfg.phosphorOrigin,
		SpacesID:       cfg.spacesID,
		SpacesSecret:   cfg.spacesSecret,
		SpacesEndpoint: cfg.spacesTracksEndpoint,
		SpacesRegion:   cfg.spacesTracksRegion,
	})
	log.Println("Spotify handlers initialized")
	phosphor.Initialize(&phosphor.Config{
		IsProduction:   cfg.isProduction,
		PhosphorOrigin: cfg.phosphorOrigin,
		MailgunAPIKey:  cfg.mailgunAPIKey,
	})
	log.Println("Phosphor handlers initialized")
	spotifyclient.Initialize(&spotifyclient.Config{
		SpotifyClientID: cfg.spotifyClientID,
		SpotifySecret:   cfg.spotifySecret,
		APIOrigin:       cfg.apiOrigin,
		BaseHTTPTimeout: cfg.handlerTimeout,
	})
	log.Println("Spotify client initialized")
	session.Initialize(&session.Config{
		CookieDomain: cfg.cookieDomain,
		IsProduction: cfg.isProduction,
	})
	log.Println("Session handling initialized")
	models.Initialize(&models.Config{
		IsProduction:             cfg.isProduction,
		SpacesID:                 cfg.spacesID,
		SpacesSecret:             cfg.spacesSecret,
		PostgresConnectionString: cfg.postgresConnectionString,
		PostgresMaxOpen:          cfg.postgresMaxOpenConnections,
		PostgresMaxIdle:          cfg.postgresMaxIdleConnections,
		PostgreMaxLifetime:       cfg.postgresMaxConnectionLifetimeMinutes,
		SpacesScriptsEndpoint:    cfg.spacesScriptsEndpoint,
		SpacesScriptsRegion:      cfg.spacesScriptsRegion,
	})
	log.Println("Models initialized")
	tracks.Initialize(&tracks.Config{
		SpacesID:       cfg.spacesID,
		SpacesSecret:   cfg.spacesSecret,
		SpacesEndpoint: cfg.spacesTracksEndpoint,
		SpacesRegion:   cfg.spacesTracksRegion,
	})
	log.Println("Tracks initialized")
}

func run(cfg *config) {
	host := getHost(cfg)
	srv := &http.Server{
		Addr:         host,
		Handler:      initializeRoutes(cfg),
		ReadTimeout:  cfg.readTimeout,
		WriteTimeout: cfg.writeTimeout,
		IdleTimeout:  cfg.idleTimeout,
	}
	log.Printf("API listening on %s.", host)
	if cfg.isProduction {
		log.Fatal(srv.ListenAndServe())
	} else {
		log.Fatal(srv.ListenAndServeTLS("phosphor.localhost.crt", "phosphor.localhost.key"))
	}
}

func getHost(cfg *config) string {
	if cfg.isProduction {
		return ":80"
	} else {
		addr := os.Getenv("HOST")
		port := os.Getenv("PORT")
		if port == "" {
			port = "3002"
		}
		return fmt.Sprintf("%s:%s", addr, port)
	}
}

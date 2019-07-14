package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/handlers/spotify"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	cfg := &config{
		isProduction:                         isProduction,
		phosphorOrigin:                       os.Getenv("PHOSPHOR_ORIGIN"),
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
	}
	migrate(cfg)
	initialize(cfg)
	run(cfg)
}

func initialize(cfg *config) {
	common.Initialize(&common.Config{
		SpotifyTimeout: cfg.writeTimeout,
	})
	spotify.Initialize(&spotify.Config{
		IsProduction:    cfg.isProduction,
		SpotifyClientID: cfg.spotifyClientID,
		SpotifySecret:   cfg.spotifySecret,
		PhosphorOrigin:  cfg.phosphorOrigin,
		SpacesID:        cfg.spacesID,
		SpacesSecret:    cfg.spacesSecret,
		SpacesEndpoint:  cfg.spacesTracksEndpoint,
		SpacesRegion:    cfg.spacesTracksRegion,
	})
	models.Initialize(&models.Config{
		SpacesID:                 cfg.spacesID,
		SpacesSecret:             cfg.spacesSecret,
		PostgresConnectionString: cfg.postgresConnectionString,
		PostgresMaxOpen:          cfg.postgresMaxOpenConnections,
		PostgresMaxIdle:          cfg.postgresMaxIdleConnections,
		PostgreMaxLifetime:       cfg.postgresMaxConnectionLifetimeMinutes,
		SpacesScriptsEndpoint:    cfg.spacesScriptsEndpoint,
		SpacesScriptsRegion:      cfg.spacesScriptsRegion,
	})
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

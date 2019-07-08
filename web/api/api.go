package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/samuelhorwitz/phosphorescence/api/scripts"
	"github.com/samuelhorwitz/phosphorescence/api/spotify"
	"log"
	"net/http"
	"os"
)

func main() {
	isProduction := os.Getenv("ENV") == "production"
	if !isProduction {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Could not load .env file: %s", err)
			return
		}
	}
	cfg := &config{
		isProduction:             isProduction,
		phosphorOrigin:           os.Getenv("PHOSPHOR_ORIGIN"),
		spotifyClientID:          os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifySecret:            os.Getenv("SPOTIFY_SECRET"),
		apiOrigin:                os.Getenv("API_ORIGIN"),
		spacesID:                 os.Getenv("SPACES_ID"),
		spacesSecret:             os.Getenv("SPACES_SECRET"),
		spacesTracksEndpoint:     os.Getenv("SPACES_TRACKS_ENDPOINT"),
		spacesTracksRegion:       os.Getenv("SPACES_TRACKS_REGION"),
		spacesScriptsEndpoint:    os.Getenv("SPACES_SCRIPTS_ENDPOINT"),
		spacesScriptsRegion:      os.Getenv("SPACES_SCRIPTS_REGION"),
		postgresConnectionString: os.Getenv("PG_CONNECTION_STRING"),
	}
	migrate(cfg)
	initialize(cfg)
	run(cfg)
}

func initialize(cfg *config) {
	spotify.Initialize(&spotify.Config{
		IsProduction:    cfg.isProduction,
		SpotifyClientID: cfg.spotifyClientID,
		SpotifySecret:   cfg.spotifySecret,
		PhosphorOrigin:  cfg.phosphorOrigin,
		APIOrigin:       cfg.apiOrigin,
		SpacesID:        cfg.spacesID,
		SpacesSecret:    cfg.spacesSecret,
		SpacesEndpoint:  cfg.spacesTracksEndpoint,
		SpacesRegion:    cfg.spacesTracksRegion,
	})
	scripts.Initialize(&scripts.Config{
		SpacesID:                 cfg.spacesID,
		SpacesSecret:             cfg.spacesSecret,
		PostgresConnectionString: cfg.postgresConnectionString,
		SpacesEndpoint:           cfg.spacesScriptsEndpoint,
		SpacesRegion:             cfg.spacesScriptsRegion,
	})
}

func run(cfg *config) {
	host := getHost(cfg)
	log.Printf("API listening on %s.", host)
	if cfg.isProduction {
		log.Fatal(http.ListenAndServe(host, initializeRoutes(cfg)))
	} else {
		log.Fatal(http.ListenAndServeTLS(host, "phosphor.localhost.crt", "phosphor.localhost.key", initializeRoutes(cfg)))
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

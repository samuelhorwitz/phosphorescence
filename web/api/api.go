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
	initialize()
	run(&config{
		isProduction:   isProduction,
		phosphorOrigin: os.Getenv("PHOSPHOR_ORIGIN"),
	})
}

func initialize() {
	spotify.Initialize(&spotify.Config{
		IsProduction:    os.Getenv("ENV") == "production",
		SpotifyClientID: os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifySecret:   os.Getenv("SPOTIFY_SECRET"),
		PhosphorOrigin:  os.Getenv("PHOSPHOR_ORIGIN"),
		APIOrigin:       os.Getenv("API_ORIGIN"),
		SpacesID:        os.Getenv("SPACES_ID"),
		SpacesSecret:    os.Getenv("SPACES_SECRET"),
		SpacesEndpoint:  os.Getenv("SPACES_TRACKS_ENDPOINT"),
		SpacesRegion:    os.Getenv("SPACES_TRACKS_REGION"),
	})
	scripts.Initialize(&scripts.Config{
		SpacesID:                 os.Getenv("SPACES_ID"),
		SpacesSecret:             os.Getenv("SPACES_SECRET"),
		PostgresConnectionString: os.Getenv("PG_CONNECTION_STRING"),
		SpacesEndpoint:           os.Getenv("SPACES_SCRIPTS_ENDPOINT"),
		SpacesRegion:             os.Getenv("SPACES_SCRIPTS_REGION"),
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

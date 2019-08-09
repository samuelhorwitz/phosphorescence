package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/samuelhorwitz/phosphorescence/jobs/push"
	"github.com/samuelhorwitz/phosphorescence/jobs/spider"
)

func main() {
	var outFile string
	var testPush bool
	var fakeData bool
	flag.StringVar(&outFile, "out", "", "output file (opt)")
	flag.BoolVar(&testPush, "test", false, "test push (opt)")
	flag.BoolVar(&fakeData, "fake", false, "fake data (opt)")
	flag.Parse()
	isProduction := os.Getenv("ENV") == "production"
	if !isProduction {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Could not load .env file: %s", err)
			return
		}
	}
	cfg := &config{
		spotifyClientID: os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifySecret:   os.Getenv("SPOTIFY_SECRET"),
		spacesID:        os.Getenv("SPACES_ID"),
		spacesSecret:    os.Getenv("SPACES_SECRET"),
		spacesEndpoint:  os.Getenv("SPACES_TRACKS_ENDPOINT"),
		spacesRegion:    os.Getenv("SPACES_TRACKS_REGION"),
		outFile:         outFile,
		testPush:        testPush || fakeData,
		fakeData:        fakeData,
	}
	run(cfg)
}

func run(cfg *config) {
	var trackJSON []byte
	var err error
	if cfg.fakeData {
		trackJSON = []byte(`{"foo":"bar"}`)
	} else {
		tracks, err := spider.GetTracks(&spider.Config{
			SpotifyClientID: cfg.spotifyClientID,
			SpotifySecret:   cfg.spotifySecret,
		})
		if err != nil {
			log.Fatalf("Could not get tracks: %s", err)
		}
		trackJSON, err = json.Marshal(tracks)
		if err != nil {
			log.Fatalf("Could not marshal tracks JSON: %s", err)
		}
	}
	if cfg.outFile != "" {
		err = dumpTracksJSON(trackJSON, cfg.outFile)
		if err != nil {
			log.Fatalf("Could not dump tracks JSON: %s", err)
		}
	}
	filename := "tracks.json"
	if cfg.testPush {
		filename = "tracks-test.json"
	}
	err = push.PushTracks(&push.Config{
		SpacesID:       cfg.spacesID,
		SpacesSecret:   cfg.spacesSecret,
		SpacesEndpoint: cfg.spacesEndpoint,
		SpacesRegion:   cfg.spacesRegion,
		Key:            filename,
	}, trackJSON)
	if err != nil {
		log.Fatalf("Could not push tracks: %s", err)
	}
	log.Println("Success")
}

func dumpTracksJSON(trackJSON []byte, filename string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Could not get working directory: %s", err)
	}
	file, err := os.Create(filepath.Join(wd, filename))
	if err != nil {
		return fmt.Errorf("Could not create file: %s", err)
	}
	defer file.Close()
	_, err = io.WriteString(file, string(trackJSON))
	if err != nil {
		return fmt.Errorf("Could not write file: %s", err)
	}
	return file.Sync()
}

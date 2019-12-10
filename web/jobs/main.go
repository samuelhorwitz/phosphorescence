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
	var bucketUpdatesOnly bool
	flag.StringVar(&outFile, "out", "", "output file (opt)")
	flag.BoolVar(&testPush, "test", false, "test push (opt)")
	flag.BoolVar(&fakeData, "fake", false, "fake data (opt)")
	flag.BoolVar(&bucketUpdatesOnly, "bucket", false, "only make updates to bucket (opt)")
	flag.Parse()
	isProduction := os.Getenv("ENV") == "production"
	if !isProduction {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Could not load .env file: %s", err)
			return
		}
	}
	cfg := &config{
		spotifyClientID:   os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifySecret:     os.Getenv("SPOTIFY_SECRET"),
		spacesID:          os.Getenv("SPACES_ID"),
		spacesSecret:      os.Getenv("SPACES_SECRET"),
		spacesEndpoint:    os.Getenv("SPACES_TRACKS_ENDPOINT"),
		spacesRegion:      os.Getenv("SPACES_TRACKS_REGION"),
		outFile:           outFile,
		testPush:          testPush || fakeData,
		fakeData:          fakeData,
		bucketUpdatesOnly: bucketUpdatesOnly,
	}
	run(cfg)
}

func run(cfg *config) {
	var tracks map[string][]*spider.TrackEnvelope
	var allTracks []*spider.TrackEnvelope
	var err error
	pushCfg := &push.Config{
		SpacesID:       cfg.spacesID,
		SpacesSecret:   cfg.spacesSecret,
		SpacesEndpoint: cfg.spacesEndpoint,
		SpacesRegion:   cfg.spacesRegion,
	}
	if cfg.bucketUpdatesOnly {
		err = push.MakeBucketUpdates(pushCfg)
		if err != nil {
			log.Fatalf("Could not make bucket updates: %s", err)
		}
		log.Println("Success")
		return
	}
	if cfg.fakeData {
		tracks = map[string][]*spider.TrackEnvelope{"US": []*spider.TrackEnvelope{}}
		allTracks = []*spider.TrackEnvelope{}
	} else {
		allTracks, tracks, err = spider.GetTracks(&spider.Config{
			SpotifyClientID: cfg.spotifyClientID,
			SpotifySecret:   cfg.spotifySecret,
		})
		if err != nil {
			log.Fatalf("Could not get tracks: %s", err)
		}
	}
	if cfg.outFile != "" {
		err = dumpTracksJSON(tracks, cfg.outFile)
		if err != nil {
			log.Fatalf("Could not dump tracks JSON: %s", err)
		}
		return
	}
	pushCfg.Key = "tracks.{region}.json"
	if cfg.testPush {
		pushCfg.Key = "tracks-test.{region}.json"
	}
	err = push.PushTracks(pushCfg, allTracks, tracks)
	if err != nil {
		log.Fatalf("Could not push tracks: %s", err)
	}
	log.Println("Success")
}

func dumpTracksJSON(tracks interface{}, filename string) error {
	trackJSON, err := json.Marshal(tracks)
	if err != nil {
		log.Fatalf("Could not marshal tracks JSON: %s", err)
	}
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

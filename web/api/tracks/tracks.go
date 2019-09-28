package tracks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/models"
)

var trackListing struct {
	sync.RWMutex
	loaded bool
	tracks map[string]*models.SpotifyTrackEnvelope
}

func Initialize(cfg *Config) {
	s3Session, err := common.InitializeS3(&common.AWSConfig{
		Config: &aws.Config{
			Endpoint: aws.String(cfg.SpacesEndpoint),
			Region:   aws.String(cfg.SpacesRegion),
		},
		AccessKeyID:     cfg.SpacesID,
		SecretAccessKey: cfg.SpacesSecret,
	})
	if err != nil {
		log.Fatalf("Could not initialize Spaces connection: %s", err)
		return
	}
	go func(s3Session *session.Session, isProduction bool) {
		s3Service := s3.New(s3Session)
		for {
			if !isProduction {
				log.Println("Attempting to download config...")
			}
			tracksJSON, err := downloadConfig(s3Service)
			if err == nil {
				if !isProduction {
					log.Println("Config file downloaded, ready to parse")
				}
				var newTracks []*models.SpotifyTrackEnvelope
				err = json.Unmarshal(tracksJSON, &newTracks)
				trackListing.Lock()
				if err != nil {
					if !trackListing.loaded {
						log.Fatalf("Could not parse first JSON tracks, exiting: %s", err)
					} else {
						log.Println("Could not parse JSON tracks: %s", err)
					}
				} else {
					if !cfg.IsProduction {
						log.Println("New tracks listing loaded")
					}
					tracksMap := make(map[string]*models.SpotifyTrackEnvelope)
					for _, track := range newTracks {
						tracksMap[track.ID] = track
					}
					trackListing.tracks = tracksMap
					trackListing.loaded = true
				}
				trackListing.Unlock()
			} else {
				log.Println("Could not load new config, skipping: %s", err)
			}
			if !isProduction {
				log.Println("Config loader sleeping for 10 minutes...")
			}
			time.Sleep(10 * time.Minute)
		}
	}(s3Session, cfg.IsProduction)
}

func GetTrack(id string) (*models.SpotifyTrackEnvelope, bool) {
	trackListing.RLock()
	track, ok := trackListing.tracks[id]
	trackListing.RUnlock()
	return track, ok
}

func downloadConfig(s3Service *s3.S3) ([]byte, error) {
	res, err := s3Service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("phosphorescence-tracks"),
		Key:    aws.String("tracks.global.json"),
	})
	if err != nil {
		return nil, fmt.Errorf("Could not load track listing into memory: %s", err)
	}
	defer res.Body.Close()
	tracksJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read track listing body: %s", err)
	}
	return tracksJSON, nil
}

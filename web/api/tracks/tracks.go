package tracks

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"io/ioutil"
	"log"
)

type TrackData struct {
	Track    json.RawMessage `json:"track"`
	Features json.RawMessage `json:"features"`
}

var trackListing map[string]TrackData

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
	s3Service := s3.New(s3Session)
	res, err := s3Service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("phosphorescence"),
		Key:    aws.String("tracks.json"),
	})
	if err != nil {
		log.Fatalf("Could not load track listing into memory: %s", err)
		return
	}
	defer res.Body.Close()
	tracksJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Could not read track listing body: %s", err)
		return
	}
	err = json.Unmarshal(tracksJSON, &trackListing)
	if err != nil {
		log.Fatalf("Could not parse JSON tracks: %s", err)
		return
	}
}

func GetTrack(id string) (TrackData, bool) {
	track, ok := trackListing[id]
	return track, ok
}

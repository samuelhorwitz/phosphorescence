package spotify

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"log"
)

var (
	isProduction    bool
	spotifyClientID string
	spotifySecret   string
	apiOrigin       string
	phosphorOrigin  string
	s3Service       *s3.S3
)

type Config struct {
	IsProduction    bool
	SpotifyClientID string
	SpotifySecret   string
	APIOrigin       string
	PhosphorOrigin  string
	SpacesID        string
	SpacesSecret    string
	SpacesEndpoint  string
	SpacesRegion    string
}

func Initialize(cfg *Config) {
	isProduction = cfg.IsProduction
	spotifyClientID = cfg.SpotifyClientID
	spotifySecret = cfg.SpotifySecret
	apiOrigin = cfg.APIOrigin
	phosphorOrigin = cfg.PhosphorOrigin
	s3Session, err := common.InitializeS3(&common.AWSConfig{
		Config: &aws.Config{
			Endpoint: aws.String(cfg.SpacesEndpoint),
			Region:   aws.String(cfg.SpacesRegion),
		},
		AccessKeyID:     cfg.SpacesID,
		SecretAccessKey: cfg.SpacesSecret,
	})
	if err != nil {
		log.Fatalf("Could not Spaces connection: %s", err)
		return
	}
	s3Service = s3.New(s3Session)
}

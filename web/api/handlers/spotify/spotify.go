package spotify

import (
	"encoding/hex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"log"
)

var (
	isProduction    bool
	phosphorOrigin  string
	authStateSecret []byte
	s3Service       *s3.S3
)

type Config struct {
	IsProduction    bool
	PhosphorOrigin  string
	SpacesID        string
	SpacesSecret    string
	SpacesEndpoint  string
	SpacesRegion    string
	AuthStateSecret string
}

func Initialize(cfg *Config) {
	var err error
	authStateSecret, err = hex.DecodeString(cfg.AuthStateSecret)
	if err != nil {
		log.Fatalf("Could not parse secret hex: %s", err)
		return
	}
	isProduction = cfg.IsProduction
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
		log.Fatalf("Could not create Spaces connection: %s", err)
		return
	}
	s3Service = s3.New(s3Session)
}

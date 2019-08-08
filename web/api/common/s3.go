package common

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const S3NotFoundCode = "NotFound"

type AWSConfig struct {
	*aws.Config
	AccessKeyID     string
	SecretAccessKey string
}

func InitializeS3(cfg *AWSConfig) (*session.Session, error) {
	// From what I can tell, the S3 Go library doesn't allow you to just pass these values in
	// They have to be in the environment.
	// :'( why
	os.Setenv("AWS_ACCESS_KEY_ID", cfg.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", cfg.SecretAccessKey)
	return session.NewSession(cfg.Config)
}

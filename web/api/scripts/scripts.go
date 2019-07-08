package scripts

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/lib/pq"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/satori/go.uuid"
	"log"
)

var scriptsNamespace = uuid.NewV5(uuid.NewV5(uuid.NamespaceDNS, "phosphor.me"), "scripts")
var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var (
	s3Service  *s3.S3
	postgresDB *sql.DB
)

type Config struct {
	SpacesID                 string
	SpacesSecret             string
	SpacesEndpoint           string
	SpacesRegion             string
	PostgresConnectionString string
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
	s3Service = s3.New(s3Session)
	postgresDB, err = sql.Open("postgres", cfg.PostgresConnectionString)
	if err != nil {
		log.Fatalf("Could not initialize Postgres: %s", err)
		return
	}
}

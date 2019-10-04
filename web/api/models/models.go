package models

import (
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/lib/pq"
	"github.com/samuelhorwitz/phosphorescence/api/common"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var (
	s3Service             *s3.S3
	s3Uploader            *s3manager.Uploader
	postgresDB            *sql.DB
	playlistImageBase64   string
	isProduction          bool
	googleAnalyticsSecret []byte
)

type Config struct {
	IsProduction             bool
	SpacesID                 string
	SpacesSecret             string
	SpacesScriptsEndpoint    string
	SpacesScriptsRegion      string
	PostgresConnectionString string
	PostgresMaxOpen          int
	PostgresMaxIdle          int
	PostgreMaxLifetime       int
	GoogleAnalyticsSecret    string
}

func Initialize(cfg *Config) {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Could not get executable path: %s", err)
		return
	}
	exPath := filepath.Dir(ex)
	playlistImage, err := ioutil.ReadFile(filepath.Join(exPath, "assets", "playlist_small.jpg"))
	if err != nil {
		log.Fatalf("Could not open playlist image: %s", err)
		return
	}
	playlistImageBase64 = base64.StdEncoding.EncodeToString(playlistImage)
	isProduction = cfg.IsProduction
	if googleAnalyticsSecret, err = hex.DecodeString(cfg.GoogleAnalyticsSecret); err != nil {
		log.Fatalf("Could not get Google Analytics secret: %s", err)
		return
	}
	s3Session, err := common.InitializeS3(&common.AWSConfig{
		Config: &aws.Config{
			Endpoint: aws.String(cfg.SpacesScriptsEndpoint),
			Region:   aws.String(cfg.SpacesScriptsRegion),
		},
		AccessKeyID:     cfg.SpacesID,
		SecretAccessKey: cfg.SpacesSecret,
	})
	if err != nil {
		log.Fatalf("Could not initialize Spaces connection: %s", err)
		return
	}
	s3Service = s3.New(s3Session)
	s3Uploader = s3manager.NewUploader(s3Session)
	postgresDB, err = sql.Open("postgres", cfg.PostgresConnectionString)
	if err != nil {
		log.Fatalf("Could not initialize Postgres: %s", err)
		return
	}
	postgresDB.SetMaxOpenConns(cfg.PostgresMaxOpen)
	postgresDB.SetMaxIdleConns(cfg.PostgresMaxIdle)
	postgresDB.SetConnMaxLifetime(time.Duration(cfg.PostgreMaxLifetime) * time.Minute)
	initializeRefreshers(cfg.IsProduction)
}

func initializeRefreshers(isProduction bool) {
	go func() {
		logInDev := func(l string) {
			if !isProduction {
				log.Println(l)
			}
		}
		for {
			time.Sleep(5 * time.Minute)
			logInDev("Refreshing materialized views...")
			var err error
			_, err = postgresDB.Exec("refresh materialized view searchables")
			if err != nil {
				log.Printf("Could not refresh searchables: %s", err)
				continue
			}
			logInDev("Searchables refreshed")
			_, err = postgresDB.Exec("refresh materialized view searchable_lexemes")
			if err != nil {
				log.Printf("Could not refresh searchable_lexemes: %s", err)
				continue
			}
			logInDev("Searchable lexemes refreshed")
			_, err = postgresDB.Exec("refresh materialized view searchable_tag_lexemes")
			if err != nil {
				log.Printf("Could not refresh searchable_tag_lexemes: %s", err)
				continue
			}
			logInDev("Searchable tag lexemes refreshed")
		}
	}()
}

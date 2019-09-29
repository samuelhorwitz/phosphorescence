package push

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/brotli/go/cbrotli"
	"github.com/samuelhorwitz/phosphorescence/jobs/spider"
)

var s3Session *session.Session
var s3Uploader *s3manager.Uploader
var s3Service *s3.S3

func PushTracks(cfg *Config, allTracks []*spider.TrackEnvelope, trackRegions map[string][]*spider.TrackEnvelope) (err error) {
	log.Println("Preparing to push tracks...")
	os.Setenv("AWS_ACCESS_KEY_ID", cfg.SpacesID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", cfg.SpacesSecret)
	s3Session, err = session.NewSession(&aws.Config{
		Endpoint: aws.String(cfg.SpacesEndpoint),
		Region:   aws.String(cfg.SpacesRegion),
	})
	if err != nil {
		return fmt.Errorf("Could not initialize Spaces S3: %s", err)
	}
	s3Uploader = s3manager.NewUploader(s3Session)
	s3Service = s3.New(s3Session)
	err = versioning()
	if err != nil {
		return fmt.Errorf("Could not ensure versioning enabled: %s", err)
	}
	err = lifecycle()
	if err != nil {
		return fmt.Errorf("Could not ensure versioning lifecycle enabled: %s", err)
	}
	log.Printf("Pushing %d global tracks", len(allTracks))
	key := strings.Replace(cfg.Key, "{region}", "global", 1)
	globalTrackJSON, err := json.Marshal(allTracks)
	if err != nil {
		log.Fatalf("Could not marshal tracks JSON: %s", err)
	}
	err = uploadGlobalTrackList(key, globalTrackJSON)
	if err != nil {
		return fmt.Errorf("Could not upload tracks to Spaces S3: %s", err)
	}
	for _, track := range allTracks {
		track.Track.AvailableMarkets = nil
	}
	for region, tracks := range trackRegions {
		log.Printf("Pushing %d tracks for region %s", len(tracks), region)
		key := strings.Replace(cfg.Key, "{region}", strings.ToLower(region), 1)
		trackJSON, err := json.Marshal(tracks)
		if err != nil {
			log.Fatalf("Could not marshal tracks JSON: %s", err)
		}
		err = uploadTrackList(key, strings.ToUpper(region), trackJSON)
		if err != nil {
			return fmt.Errorf("Could not upload tracks to Spaces S3: %s", err)
		}
	}
	err = cors()
	if err != nil {
		return fmt.Errorf("Could not set CORS config: %s", err)
	}
	return nil
}

func uploadTrackList(key, region string, trackJSON []byte) (err error) {
	reader, writer := io.Pipe()
	go func() {
		zw := cbrotli.NewWriter(writer, cbrotli.WriterOptions{
			Quality: 11,
		})
		if _, err = zw.Write(trackJSON); err != nil {
			err = fmt.Errorf("Could not compress tracks JSON: %s", err)
			return
		}
		if err := zw.Close(); err != nil {
			err = fmt.Errorf("Could not close compression buffer: %s", err)
			return
		}
		writer.Close()
	}()
	if err != nil {
		return fmt.Errorf("Could not zip tracks JSON: %s", err)
	}
	_, err = s3Uploader.Upload(&s3manager.UploadInput{
		Bucket:          aws.String("phosphorescence-tracks"),
		ACL:             aws.String("private"),
		ContentType:     aws.String("application/json"),
		ContentEncoding: aws.String("br"),
		CacheControl:    aws.String("private"),
		Metadata: aws.StringMap(map[string]string{
			"uncompressed-length": fmt.Sprintf("%d", len(trackJSON)),
			"region":              region,
		}),
		Expires: aws.Time(time.Now().AddDate(0, 0, 1).Add(6 * time.Hour)),
		Key:     aws.String(key),
		Body:    reader,
	})
	if err != nil {
		return fmt.Errorf("Could not upload to Spaces: %s", err)
	}
	return nil
}

func uploadGlobalTrackList(key string, trackJSON []byte) (err error) {
	reader, writer := io.Pipe()
	go func() {
		zw, err := gzip.NewWriterLevel(writer, gzip.BestCompression)
		if err != nil {
			err = fmt.Errorf("Could not create new GZIP writer: %s", err)
			return
		}
		if _, err = zw.Write(trackJSON); err != nil {
			err = fmt.Errorf("Could not compress tracks JSON: %s", err)
			return
		}
		if err := zw.Close(); err != nil {
			err = fmt.Errorf("Could not close compression buffer: %s", err)
			return
		}
		writer.Close()
	}()
	if err != nil {
		return fmt.Errorf("Could not zip tracks JSON: %s", err)
	}
	_, err = s3Uploader.Upload(&s3manager.UploadInput{
		Bucket:          aws.String("phosphorescence-tracks"),
		ACL:             aws.String("private"),
		ContentType:     aws.String("application/json"),
		ContentEncoding: aws.String("gzip"),
		CacheControl:    aws.String("private"),
		Metadata: aws.StringMap(map[string]string{
			"uncompressed-length": fmt.Sprintf("%d", len(trackJSON)),
		}),
		Expires: aws.Time(time.Now().AddDate(0, 0, 1).Add(6 * time.Hour)),
		Key:     aws.String(key),
		Body:    reader,
	})
	if err != nil {
		return fmt.Errorf("Could not upload to Spaces: %s", err)
	}
	return nil
}

func versioning() error {
	_, err := s3Service.PutBucketVersioning(&s3.PutBucketVersioningInput{
		Bucket: aws.String("phosphorescence-tracks"),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String("Enabled"),
		},
	})
	return err
}

func lifecycle() error {
	_, err := s3Service.PutBucketLifecycleConfiguration(&s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String("phosphorescence-tracks"),
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: []*s3.LifecycleRule{
				{
					ID:     aws.String("Expire Old Track Metadata"),
					Status: aws.String("Enabled"),
					AbortIncompleteMultipartUpload: &s3.AbortIncompleteMultipartUpload{
						DaysAfterInitiation: aws.Int64(7),
					},
					Expiration: &s3.LifecycleExpiration{
						Days: aws.Int64(7),
					},
					Filter: &s3.LifecycleRuleFilter{
						Prefix: aws.String("tracks"),
					},
				},
			},
		},
	})
	return err
}

func cors() error {
	_, err := s3Service.PutBucketCors(&s3.PutBucketCorsInput{
		Bucket: aws.String("phosphorescence-tracks"),
		CORSConfiguration: &s3.CORSConfiguration{
			CORSRules: []*s3.CORSRule{&s3.CORSRule{
				AllowedOrigins: aws.StringSlice([]string{"https://phosphor.me", "https://phosphor.localhost:3000"}),
				AllowedMethods: aws.StringSlice([]string{"GET"}),
				ExposeHeaders: aws.StringSlice([]string{
					"cache-control",
					"content-length",
					"content-type",
					"expires",
					"last-modified",
					"x-amz-meta-uncompressed-length",
					"x-amz-meta-region",
				}),
			}},
		},
	})
	return err
}

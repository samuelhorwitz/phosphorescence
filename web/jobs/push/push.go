package push

import (
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
	"time"
)

var s3Session *session.Session
var s3Uploader *s3manager.Uploader
var s3Service *s3.S3

func PushTracks(cfg *Config, trackJSON []byte) (err error) {
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
	err = backupOldTrackList()
	if err != nil {
		return fmt.Errorf("Could not back up old tracks to Spaces S3: %s", err)
	}
	err = uploadTrackList(cfg.Key, trackJSON)
	if err != nil {
		return fmt.Errorf("Could not upload tracks to Spaces S3: %s", err)
	}
	return nil
}

func backupOldTrackList() error {
	_, err := s3Service.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String("phosphorescence"),
		CopySource: aws.String("/phosphorescence/tracks.json"),
		Key:        aws.String(fmt.Sprintf("old/tracks-%d.json", time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)))),
	})
	if err != nil {
		return fmt.Errorf("Could not back up old track list: %s", err)
	}
	return nil
}

func uploadTrackList(key string, trackJSON []byte) (err error) {
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
		Bucket:          aws.String("phosphorescence"),
		ACL:             aws.String("private"),
		ContentType:     aws.String("application/json"),
		ContentEncoding: aws.String("gzip"),
		Key:             aws.String(key),
		Body:            reader,
	})
	if err != nil {
		return fmt.Errorf("Could not upload to Spaces: %s", err)
	}
	return nil
}

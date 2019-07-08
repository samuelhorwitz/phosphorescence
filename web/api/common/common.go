package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"net/http"
	"os"
	"time"
)

const S3NotFoundCode = "NotFound"

var SpotifyClient = &http.Client{
	Timeout: time.Second * 2,
}

type AWSConfig struct {
	*aws.Config
	AccessKeyID     string
	SecretAccessKey string
}

func InitializeS3(cfg *AWSConfig) (*session.Session, error) {
	os.Setenv("AWS_ACCESS_KEY_ID", cfg.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", cfg.SecretAccessKey)
	return session.NewSession(cfg.Config)
}

func JSONRaw(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func JSON(w http.ResponseWriter, data interface{}) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		Fail(w, fmt.Errorf("Could not marshal JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	JSONRaw(w, jsonResponse)
}

func Fail(w http.ResponseWriter, err error, status int) {
	log.Printf("Request failed, returning %d: %s", status, err)
	http.Error(w, "error", status)
}

func RollbackAndFail(w http.ResponseWriter, tx *sql.Tx, err error, status int) {
	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		err = fmt.Errorf("Rollback failed: %s; Original Error: %s", rollbackErr, err)
	}
	Fail(w, err, status)
}

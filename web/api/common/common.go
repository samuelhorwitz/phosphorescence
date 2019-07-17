package common

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

var PhosphorUUIDV5Namespace = uuid.NewV5(uuid.NamespaceDNS, "phosphor.me")

var SpotifyClient *spotifyclient.SpotifyClient

type Config struct {
	SpotifyTimeout time.Duration
}

func Initialize(cfg *Config) {
	SpotifyClient = &spotifyclient.SpotifyClient{
		Timeout: cfg.SpotifyTimeout,
		Client: &http.Client{
			Timeout: cfg.SpotifyTimeout,
		},
	}
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
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"error":true}`, status)
}

func TryToRollback(tx *sql.Tx, err error) error {
	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		err = fmt.Errorf("Rollback failed: %s; Original Error: %s", rollbackErr, err)
	}
	return err
}

func ParseScriptVersion(versionStr string) time.Time {
	var version time.Time
	if versionStr != "" {
		if tentativeVersion, err := time.Parse(time.RFC3339, versionStr); err == nil {
			version = tentativeVersion
		}
	}
	return version
}

func HandlerTimeoutCancelContext(r *http.Request) context.Context {
	reqCtx, reqCancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-r.Context().Done():
			reqCancel()
		}
	}()
	return reqCtx
}

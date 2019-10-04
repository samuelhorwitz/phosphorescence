package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"github.com/satori/go.uuid"
)

const simpleErrorPayload = `{"error":true}`

var PhosphorUUIDV5Namespace = uuid.NewV5(uuid.NamespaceDNS, "phosphor.me")

var (
	SpotifyClient *spotifyclient.SpotifyClient
	isProduction  bool
	RedisPool     *redis.Pool
)

type Config struct {
	IsProduction   bool
	SpotifyTimeout time.Duration
	RedisHost      string
}

func Initialize(cfg *Config) {
	SpotifyClient = &spotifyclient.SpotifyClient{
		Timeout: cfg.SpotifyTimeout,
		Client: &http.Client{
			Timeout: cfg.SpotifyTimeout,
		},
	}
	isProduction = cfg.IsProduction
	RedisPool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.RedisHost)
			if err != nil {
				log.Fatalf("Could not connect to Redis: %s", err)
			}
			return c, err
		},
	}
}

func JSON(w http.ResponseWriter, data interface{}) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		Fail(w, fmt.Errorf("Could not marshal JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func Fail(w http.ResponseWriter, err error, status int) {
	fail(w, err, status, false, simpleErrorPayload)
}

func FailWithJSON(w http.ResponseWriter, err error, data interface{}, status int) {
	jsonResponse, err := json.Marshal(struct {
		Error bool        `json:"error"`
		Data  interface{} `json:"data"`
	}{
		Error: true,
		Data:  data,
	})
	if err != nil {
		Fail(w, fmt.Errorf("Could not marshal JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	fail(w, err, status, false, string(jsonResponse))
}

func FailAndLog(w http.ResponseWriter, err error, status int) {
	fail(w, err, status, true, simpleErrorPayload)
}

func fail(w http.ResponseWriter, err error, status int, forceLog bool, payload string) {
	if forceLog || !isProduction {
		log.Printf("Request failed, returning %d: %s", status, err)
	}
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, payload, status)
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

func ExponentialBackoff(baseDuration time.Duration, timeout time.Duration, fn func() bool) {
	done := make(chan struct{})
	go func() {
		count := 0
		for !fn() {
			count++
			time.Sleep(baseDuration * time.Duration(math.Pow(2.0, float64(rand.Intn(count+1)))-1))
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
	}
}

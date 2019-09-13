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
	fail(w, err, status, false)
}

func FailAndLog(w http.ResponseWriter, err error, status int) {
	fail(w, err, status, true)
}

func fail(w http.ResponseWriter, err error, status int, forceLog bool) {
	if forceLog || !isProduction {
		log.Printf("Request failed, returning %d: %s", status, err)
	}
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"error":true}`, status)
}

func FailWithRawJSON(w http.ResponseWriter, err error, body []byte, status int) {
	if !isProduction {
		log.Printf("Request failed, returning %d: %s", status, err)
	}
	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		Fail(w, fmt.Errorf("Could not unmarshal JSON: %s", err), http.StatusInternalServerError)
		return
	}
	parsedBody["error"] = true
	jsonResponse, err := json.Marshal(parsedBody)
	if err != nil {
		Fail(w, fmt.Errorf("Could not marshal JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, string(jsonResponse), status)
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

func ExponentialBackoff(baseDuration time.Duration, fn func() bool) {
	count := 0
	for !fn() {
		count++
		time.Sleep(baseDuration * time.Duration(math.Pow(2.0, float64(rand.Intn(count+1)))-1))
	}
}

package phosphor

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

var (
	phosphorOrigin      string
	isProduction        bool
	mailgunAPIKey       string
	mailgunClient       *http.Client
	playlistImageBase64 string
	noHTML              *bluemonday.Policy
	safeHTTPClient      *http.Client
)

type Config struct {
	PhosphorOrigin string
	IsProduction   bool
	MailgunAPIKey  string
}

func Initialize(cfg *Config) {
	phosphorOrigin = cfg.PhosphorOrigin
	isProduction = cfg.IsProduction
	mailgunAPIKey = cfg.MailgunAPIKey
	mailgunClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	safeHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
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
	noHTML = bluemonday.StrictPolicy()
}

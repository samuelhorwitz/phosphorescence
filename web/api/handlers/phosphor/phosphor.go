package phosphor

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"golang.org/x/oauth2"
)

var (
	phosphorOrigin           string
	isProduction             bool
	mailgunAPIKey            string
	mailgunClient            *http.Client
	playlistImageBase64      string
	noHTML                   *bluemonday.Policy
	safeHTTPClient           *http.Client
	phosphorescenceSpotifyID string
	phosphorescenceToken     *oauth2.Token
)

type Config struct {
	PhosphorOrigin              string
	IsProduction                bool
	MailgunAPIKey               string
	PhosphorescenceSpotifyID    string
	PhosphorescenceRefreshToken string
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
	phosphorescenceSpotifyID = cfg.PhosphorescenceSpotifyID
	phosphorescenceToken = &oauth2.Token{
		RefreshToken: cfg.PhosphorescenceRefreshToken,
	}
	if _, err = getPhosphorescenceToken(); err != nil {
		log.Fatalf("Could not get Phosphorescence user token from refresh token: %s", err)
		return
	}
}

func getPhosphorescenceToken() (*oauth2.Token, error) {
	var err error
	phosphorescenceToken, err = spotifyclient.GetToken(phosphorescenceToken)
	if err != nil {
		// TODO how should we handle this happening?
		return nil, fmt.Errorf("Could not get Phosphorescence token: %s", err)
	}
	return phosphorescenceToken, nil
}

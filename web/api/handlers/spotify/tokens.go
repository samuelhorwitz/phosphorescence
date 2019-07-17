package spotify

import (
	"encoding/json"
	"fmt"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Tokens(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequestWithContext(common.HandlerTimeoutCancelContext(r), "POST", "https://accounts.spotify.com/api/token", strings.NewReader(getSpotifyTokenRequestBody(r)))
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not build Spotify token request: %s", err), http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not make Spotify token request: %s", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		common.Fail(w, fmt.Errorf("Spotify token request responded with %d", res.StatusCode), res.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read Spotify token response: %s", err), http.StatusInternalServerError)
		return
	}
	var parsedBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not parse Spotify token response: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, struct {
		AccessToken  string `json:"access"`
		RefreshToken string `json:"refresh"`
		ExpiresIn    int    `json:"expires"`
	}{
		AccessToken:  parsedBody.AccessToken,
		RefreshToken: parsedBody.RefreshToken,
		ExpiresIn:    parsedBody.ExpiresIn,
	})
}

func getSpotifyTokenRequestBody(r *http.Request) string {
	body := url.Values{}
	body.Set("client_id", spotifyClientID)
	body.Set("client_secret", spotifySecret)
	if r.URL.Query().Get("type") == "refresh" {
		body.Set("grant_type", "refresh_token")
		body.Set("refresh_token", r.URL.Query().Get("code"))
	} else {
		body.Set("grant_type", "authorization_code")
		body.Set("code", r.URL.Query().Get("code"))
		body.Set("redirect_uri", fmt.Sprintf("%s/auth/login", phosphorOrigin))
	}
	return body.Encode()
}

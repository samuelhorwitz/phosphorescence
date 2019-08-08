package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/samuelhorwitz/phosphorescence/api/common"
	"golang.org/x/oauth2"
)

type user struct {
	ID      string `json:"id"`
	Name    string `json:"display_name"`
	Country string `json:"country"`
	Email   string `json:"email"`
}

func getUser(r *http.Request, token *oauth2.Token) (user, error) {
	req, err := http.NewRequestWithContext(r.Context(), "GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return user{}, fmt.Errorf("Could not build Spotify profile request: %s", err)
	}
	token.SetAuthHeader(req)
	res, err := common.SpotifyClient.Do(req)
	if err != nil {
		return user{}, fmt.Errorf("Could not make Spotify profile request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return user{}, fmt.Errorf("Spotify profile request responded with %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return user{}, fmt.Errorf("Could not read Spotify profile response: %s", err)
	}
	var parsedUser user
	err = json.Unmarshal(body, &parsedUser)
	if err != nil {
		return user{}, fmt.Errorf("Could not parse Spotify profile response: %s", err)
	}
	return parsedUser, nil
}

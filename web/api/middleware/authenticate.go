package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"io/ioutil"
	"net/http"
)

type contextKey string

var SpotifyIDContextKey = contextKey("spotifyID")

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var spotifyToken string
		spotifyTokenCookie, err := r.Cookie("spotify_access")
		if err != nil {
			common.Fail(w, errors.New("No token"), http.StatusForbidden)
			return
		}
		spotifyToken = spotifyTokenCookie.Value
		if spotifyToken == "" {
			common.Fail(w, errors.New("Empty token"), http.StatusForbidden)
			return
		}
		req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not build Spotify profile request: %s", err), http.StatusInternalServerError)
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
		res, err := common.SpotifyClient.Do(req)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not make Spotify profile request: %s", err), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			common.Fail(w, fmt.Errorf("Spotify profile request responded with %d", res.StatusCode), http.StatusForbidden)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not read Spotify profile response: %s", err), http.StatusInternalServerError)
			return
		}
		var parsedBody struct {
			ID string `json:"id"`
		}
		err = json.Unmarshal(body, &parsedBody)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not parse Spotify profile response: %s", err), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), SpotifyIDContextKey, parsedBody.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

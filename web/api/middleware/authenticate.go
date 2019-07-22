package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/gomodule/redigo/redis"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"io/ioutil"
	"net/http"
)

const SpotifyTokenContextKey = contextKey("spotifyToken")
const SpotifyIDContextKey = contextKey("spotifyID")
const SpotifyNameContextKey = contextKey("spotifyName")
const SpotifyCountryContextKey = contextKey("spotifyCountry")

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var spotifyToken string
		spotifyTokenCookie, err := r.Cookie("spotify_access")
		if err != nil {
			common.Fail(w, errors.New("No token"), http.StatusUnauthorized)
			return
		}
		spotifyToken = spotifyTokenCookie.Value
		if spotifyToken == "" {
			common.Fail(w, errors.New("Empty token"), http.StatusUnauthorized)
			return
		}
		authenticateOnSpotify := false
		redisConn := redisPool.Get()
		defer redisConn.Close()
		sessionData, err := redis.StringMap(redisConn.Do("HGETALL", getSessionKey(spotifyToken)))
		if err != nil {
			if err == redis.ErrNil {
				authenticateOnSpotify = true
			} else {
				common.Fail(w, fmt.Errorf("Could not access Redis: %s", err), http.StatusInternalServerError)
				return
			}
		} else if sessionData["id"] == "" || sessionData["country"] == "" {
			// We have to have country and ID for the user, name is not as important
			authenticateOnSpotify = true
		}
		var parsedBody struct {
			ID      string `json:"id"`
			Name    string `json:"display_name"`
			Country string `json:"country"`
		}
		if authenticateOnSpotify {
			req, err := http.NewRequestWithContext(r.Context(), "GET", "https://api.spotify.com/v1/me", nil)
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
				common.Fail(w, fmt.Errorf("Spotify profile request responded with %d", res.StatusCode), http.StatusInternalServerError)
				return
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				common.Fail(w, fmt.Errorf("Could not read Spotify profile response: %s", err), http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(body, &parsedBody)
			if err != nil {
				common.Fail(w, fmt.Errorf("Could not parse Spotify profile response: %s", err), http.StatusInternalServerError)
				return
			}
			_, err = redisConn.Do("HMSET", getSessionKey(spotifyToken),
				"id", parsedBody.ID,
				"name", parsedBody.Name,
				"country", parsedBody.Country)
			if err != nil {
				common.Fail(w, fmt.Errorf("Could not cache Spotify profile response: %s", err), http.StatusInternalServerError)
				return
			}
			_, err = redisConn.Do("EXPIRE", getSessionKey(spotifyToken), 60*60)
			if err != nil {
				common.Fail(w, fmt.Errorf("Could not set cache expiration for Spotify profile response: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			parsedBody.ID = sessionData["id"]
			parsedBody.Name = sessionData["name"]
			parsedBody.Country = sessionData["country"]
		}
		lmtErr := tollbooth.LimitByKeys(phosphorLimiter, []string{parsedBody.ID})
		if lmtErr != nil {
			phosphorLimiter.ExecOnLimitReached(w, r)
			common.Fail(w, fmt.Errorf("Phosphorescence API rate limiting hit: %s", lmtErr.Message), lmtErr.StatusCode)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, SpotifyTokenContextKey, spotifyToken)
		ctx = context.WithValue(ctx, SpotifyIDContextKey, parsedBody.ID)
		ctx = context.WithValue(ctx, SpotifyNameContextKey, parsedBody.Name)
		ctx = context.WithValue(ctx, SpotifyCountryContextKey, parsedBody.Country)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getSessionKey(spotifyToken string) string {
	return fmt.Sprintf("spotify:session:%s", spotifyToken)
}

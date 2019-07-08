package spotify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	query := url.Values{}
	query.Set("client_id", spotifyClientID)
	query.Set("response_type", "code")
	query.Set("scope", strings.Join([]string{
		"streaming",
		"user-read-birthdate",
		"user-read-email",
		"user-read-private",
		"user-read-playback-state",
		"user-read-recently-played",
	}, " "))
	query.Set("redirect_uri", fmt.Sprintf("%s/auth/login", phosphorOrigin))
	http.Redirect(w, r, fmt.Sprintf("https://accounts.spotify.com/authorize?%s", query.Encode()), http.StatusFound)
}

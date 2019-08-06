package spotify

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type oauthState struct {
	Token      string `json:"token"`
	RememberMe bool   `json:"rememberMe"`
}

func Authorize(w http.ResponseWriter, r *http.Request) {
	var isPermanent bool
	permanentQP := r.URL.Query().Get("permanent")
	if permanentQP != "" {
		var err error
		isPermanent, err = strconv.ParseBool(permanentQP)
		if err != nil {
			if !isProduction {
				log.Printf("Could not parse permanence query param: %s", err)
			}
			http.Redirect(w, r, fmt.Sprintf("%s/auth/failed", phosphorOrigin), http.StatusFound)
			return
		}
	}
	state, err := newState(isPermanent)
	if err != nil {
		if !isProduction {
			log.Printf("Could not generate state: %s", err)
		}
		http.Redirect(w, r, fmt.Sprintf("%s/auth/failed", phosphorOrigin), http.StatusFound)
		return
	}
	http.Redirect(w, r, spotifyclient.AuthCodeURL(state), http.StatusFound)
}

func AuthorizeRedirect(w http.ResponseWriter, r *http.Request) {
	var state oauthState
	err := json.Unmarshal([]byte(r.URL.Query().Get("state")), &state)
	if err != nil {
		if !isProduction {
			log.Printf("Could not parse state JSON: %s", err)
		}
		http.Redirect(w, r, fmt.Sprintf("%s/auth/failed", phosphorOrigin), http.StatusFound)
		return
	}
	if !validateStateToken(state.Token) {
		if !isProduction {
			log.Printf("Invalid state %s", state)
		}
		http.Redirect(w, r, fmt.Sprintf("%s/auth/failed", phosphorOrigin), http.StatusFound)
		return
	}
	token, err := spotifyclient.TokenExchange(r.URL.Query().Get("code"))
	if err != nil {
		if !isProduction {
			log.Printf("Could not authorize user: %s", err)
		}
		http.Redirect(w, r, fmt.Sprintf("%s/auth/failed", phosphorOrigin), http.StatusFound)
		return
	}
	session.Create(w, r, token, state.RememberMe)
	http.Redirect(w, r, phosphorOrigin, http.StatusFound)
}

func Token(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	common.JSON(w, map[string]interface{}{"token": sess.SpotifyToken.AccessToken})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		if !isProduction {
			log.Println("No session on request context")
		}
		http.Redirect(w, r, fmt.Sprintf("%s/auth/failed", phosphorOrigin), http.StatusFound)
		return
	}
	refreshID, _ := r.Context().Value(middleware.RefreshIDContextKey).(string)
	session.Destroy(w, sess, refreshID)
	http.Redirect(w, r, phosphorOrigin, http.StatusFound)
}

func newState(isPermanent bool) (string, error) {
	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		return "", fmt.Errorf("Could not get entropy: %s", err)
	}
	currentTime := time.Now().Unix()
	state := oauthState{
		Token:      buildStateToken(id, currentTime),
		RememberMe: isPermanent,
	}
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return "", fmt.Errorf("Could not marshal state: %s", err)
	}
	return string(stateJSON), nil
}

func validateStateToken(token string) bool {
	parts := strings.Split(token, ":")
	// Don't care about errors at all in this function, we
	// will work with zero values and ignore errors, and
	// expectedly fail to validate when a zero value occurs.
	unixTime, _ := strconv.ParseInt(parts[1], 10, 64)
	parsedTime := time.Unix(unixTime, 0)
	// Quick return should not be a security issue;
	// this only leaks that the time is invalid which
	// an attacker composing attacks would already know.
	// The hmac compare below is the part that must
	// be secure against timing attacks and it is by
	// using the Go standard library for safety.
	if parsedTime.Add(2 * time.Minute).Before(time.Now()) {
		return false
	}
	id, _ := hex.DecodeString(parts[0])
	expected, _ := hex.DecodeString(parts[2])
	actual := buildStateTokenHMAC(id, unixTime)
	return hmac.Equal(expected, actual)
}

func buildStateToken(id []byte, unixTime int64) string {
	stateTokenHMAC := buildStateTokenHMAC(id, unixTime)
	return fmt.Sprintf("%s:%d:%s", hex.EncodeToString(id), unixTime, hex.EncodeToString(stateTokenHMAC))
}

func buildStateTokenHMAC(id []byte, unixTime int64) []byte {
	timeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeBytes, uint64(unixTime))
	h := hmac.New(sha256.New, authStateSecret)
	h.Write(id)
	h.Write(timeBytes)
	return h.Sum(nil)
}

package spotify

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
)

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
	state, err := session.CreateAuthRedirect(isPermanent)
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
	isPermanent, err := session.CheckAuthRedirect(r.URL.Query().Get("state"))
	if err != nil {
		if !isProduction {
			log.Printf("Invalid state: %s", err)
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
	session.Create(w, r, token, isPermanent)
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

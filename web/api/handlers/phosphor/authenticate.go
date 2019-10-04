package phosphor

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/mail"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	utcOffsetMinutes, _ := strconv.ParseInt(r.URL.Query().Get("utcOffsetMinutes"), 10, 64)
	email, err := sess.GetEmail(r)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get user email: %s", err), http.StatusInternalServerError)
		return
	}
	magicLink, err := session.CreateMagicLink(sess)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not create magic link: %s", err), http.StatusInternalServerError)
		return
	}
	err = mail.SendAuthentication(sess.SpotifyName, email, magicLink, utcOffsetMinutes)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not email user: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"email": email})
}

func AuthenticateRedirect(w http.ResponseWriter, r *http.Request) {
	sess, _ := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	magicLink := chi.URLParam(r, "magicLink")
	err := session.Upgrade(w, sess, magicLink)
	if err != nil {
		if !isProduction {
			log.Printf("Could not authenticate session: %s", err)
		}
		http.Redirect(w, r, fmt.Sprintf("%s/auth/failedupgrade", phosphorOrigin), http.StatusFound)
		return
	}
	http.Redirect(w, r, phosphorOrigin, http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var sessionID string
	sess, _ := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	refreshID, _ := r.Context().Value(middleware.RefreshIDContextKey).(string)
	if sess != nil {
		sessionID = sess.ID
	}
	session.Destroy(w, sessionID, refreshID)
	http.Redirect(w, r, phosphorOrigin, http.StatusFound)
}

func LogoutEverywhere(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	err := session.DestroyAllByUser(w, sess)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not destroy all user sessions: %s", err), http.StatusUnauthorized)
		return
	}
	common.JSON(w, map[string]interface{}{"signedout": true})
}

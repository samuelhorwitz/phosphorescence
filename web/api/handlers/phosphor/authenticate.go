package phosphor

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	err = sendMail(sess.SpotifyName, email, magicLink, utcOffsetMinutes)
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

func sendMail(name string, email string, magicLink string, utcOffsetMinutes int64) error {
	form := url.Values{}
	form.Add("from", `"Samuel @ Phosphorescence" <noreply@phosphor.me>`)
	form.Add("to", fmt.Sprintf(`"%s" <%s>`, strings.Replace(name, `"`, "", -1), strings.Replace(email, ">", "", -1)))
	form.Add("subject", fmt.Sprintf("Phosphorescence Log In Magic Link (Attempt at %s)", time.Now().UTC().Add(time.Duration(utcOffsetMinutes)*time.Minute).Format("3:04 PM on Monday, January _2")))
	form.Add("template", "magiclinklogin")
	form.Add("h:X-Mailgun-Variables", fmt.Sprintf(`{"token": "%s"}`, magicLink))
	req, err := http.NewRequest("POST", "https://api.mailgun.net/v3/phosphor.me/messages", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("Could not build Mailgun profile request: %s", err)
	}
	req.SetBasicAuth("api", mailgunAPIKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := mailgunClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not make Mailgun send mail request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Mailgun request responded with %d", res.StatusCode)
	}
	return nil
}

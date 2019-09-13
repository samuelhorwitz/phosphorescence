package session

import (
	"net/http"
	"time"
)

const (
	SessionCookieName = "sid"
	RefreshCookieName = "ref"
)

func setCookies(w http.ResponseWriter, sessionID string, refreshID string, permanent bool) {
	sessionCookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	if isProduction {
		sessionCookie.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(w, sessionCookie)
	refreshCookie := &http.Cookie{
		Name:     RefreshCookieName,
		Value:    refreshID,
		Domain:   cookieDomain,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	if permanent {
		refreshCookie.Expires = time.Now().Add(permanentExpiration)
		refreshCookie.MaxAge = int(permanentExpiration.Seconds())
	}
	if isProduction {
		refreshCookie.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(w, refreshCookie)
}

func clearCookies(w http.ResponseWriter) {
	sessionCookie := &http.Cookie{
		Name:    SessionCookieName,
		Value:   "",
		Domain:  cookieDomain,
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}
	http.SetCookie(w, sessionCookie)
	refreshCookie := &http.Cookie{
		Name:    RefreshCookieName,
		Value:   "",
		Domain:  cookieDomain,
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}
	http.SetCookie(w, refreshCookie)
}

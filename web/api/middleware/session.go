package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"net/http"
)

const SessionContextKey = contextKey("session")
const AuthenticatedSessionContextKey = contextKey("authenticatedSession")
const RefreshIDContextKey = contextKey("refreshID")

func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := ""
		refreshID := ""
		sessionIDCookie, err := r.Cookie(session.SessionCookieName)
		if err == nil {
			sessionID = sessionIDCookie.Value
		}
		refreshIDCookie, err := r.Cookie(session.RefreshCookieName)
		if err == nil {
			refreshID = refreshIDCookie.Value
		}
		sess, err := session.Find(w, r, sessionID, refreshID)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not get session for ID: %s", err), http.StatusUnauthorized)
			return
		}
		lmtErr := tollbooth.LimitByKeys(phosphorLimiter, []string{sess.SpotifyID})
		if lmtErr != nil {
			phosphorLimiter.ExecOnLimitReached(w, r)
			common.Fail(w, fmt.Errorf("Phosphorescence API rate limiting hit: %s", lmtErr.Message), lmtErr.StatusCode)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, SessionContextKey, sess)
		ctx = context.WithValue(ctx, RefreshIDContextKey, refreshID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthenticatedSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := r.Context().Value(SessionContextKey).(*session.Session)
		if !ok {
			common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
			return
		}
		if !sess.Authenticated {
			common.Fail(w, errors.New("Session is not authenticated for Phosphorescence"), http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, AuthenticatedSessionContextKey, sess)
		ctx = context.WithValue(ctx, SessionContextKey, nil)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

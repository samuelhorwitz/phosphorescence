package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"
	"github.com/satori/go.uuid"
	"net/http"
)

const ScriptContextKey = contextKey("script")
const ScriptVersionContextKey = contextKey("scriptVersion")

func AuthorizeReadScript(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := r.Context().Value(AuthenticatedSessionContextKey).(*session.Session)
		if !ok {
			common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
			return
		}
		scriptID, err := uuid.FromString(chi.URLParam(r, "scriptID"))
		if err != nil {
			common.Fail(w, errors.New("Invalid script ID"), http.StatusBadRequest)
			return
		}
		script, ok, err := models.GetScriptWithAuthorizationCheck(sess.SpotifyID, scriptID)
		if err != nil {
			common.Fail(w, fmt.Errorf("Cannot check script authorization: %s", err), http.StatusInternalServerError)
			return
		}
		if !ok {
			common.Fail(w, errors.New("User does not have access to script or script does not exist"), http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), ScriptContextKey, script)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthorizePrivateScriptActions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := r.Context().Value(AuthenticatedSessionContextKey).(*session.Session)
		if !ok {
			common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
			return
		}
		script, ok := r.Context().Value(ScriptContextKey).(models.Script)
		if !ok {
			common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
			return
		}
		if script.AuthorSpotifyID.String != sess.SpotifyID {
			common.Fail(w, errors.New("User is not author of script"), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthorizeReadScriptVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := r.Context().Value(AuthenticatedSessionContextKey).(*session.Session)
		if !ok {
			common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
			return
		}
		script, ok := r.Context().Value(ScriptContextKey).(models.Script)
		if !ok {
			common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
			return
		}
		scriptVersionID := common.ParseScriptVersion(chi.URLParam(r, "scriptVersionID"))
		scriptVersion, ok, err := models.GetScriptVersionWithAuthorizationCheck(sess.SpotifyID, script.ID, scriptVersionID)
		if err != nil {
			common.Fail(w, fmt.Errorf("Cannot check script version authorization: %s", err), http.StatusInternalServerError)
			return
		}
		if !ok {
			common.Fail(w, errors.New("User does not have access to script version or script version does not exist"), http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), ScriptVersionContextKey, scriptVersion)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

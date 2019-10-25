package phosphor

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/handlers"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	sess, err := session.UpdateSessionDetailsFromSpotify(r, sess)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not update session details: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{
		"user": models.GetUser(sess),
	})
}

func GetCurrentlyPlaying(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	currentlyPlaying, err := models.GetCurrentPlayback(r.Context(), sess)
	if err != nil {
		code := http.StatusInternalServerError
		if err == models.ErrCurrentlyPlayingNotTrack {
			code = http.StatusNotFound
		}
		common.Fail(w, fmt.Errorf("Could not get current playback: %s", err), code)
		return
	}
	common.JSON(w, currentlyPlaying)
}

func ListCurrentUserScripts(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.AuthenticatedSessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	count, _ := r.Context().Value(middleware.PageCountContextKey).(uint64)
	if !ok {
		common.Fail(w, errors.New("No page count on request context"), http.StatusInternalServerError)
		return
	}
	from, _ := r.Context().Value(middleware.PageCursorContextKey).(time.Time)
	if !ok {
		common.Fail(w, errors.New("No page cursor on request context"), http.StatusInternalServerError)
		return
	}
	scripts, err := models.GetScriptsBySpotifyUserID(sess.SpotifyID, count, from, true)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get scripts for user: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"scripts": scripts})
}

func ListSpotifyDevices(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	devices, err := models.GetDevices(r.Context(), sess)
	if err != nil {
		common.Fail(w, fmt.Errorf("Unable to get devices: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, devices)
}

func CreateAndFollowPlaylist(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	playlistID, err := createPlaylist(r)
	if err != nil {
		code := http.StatusInternalServerError
		if httpErr, ok := err.(handlers.HTTPError); ok {
			code = httpErr.Code
		}
		common.Fail(w, fmt.Errorf("Failed to create playlist: %s", err), code)
		return
	}
	err = models.FollowPlaylist(r.Context(), sess, playlistID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Failed to follow playlist: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"playlist": playlistID})
}

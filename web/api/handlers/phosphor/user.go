package phosphor

import (
	"errors"
	"fmt"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"net/http"
	"time"
)

type User struct {
	SpotifyID string `json:"spotifyId"`
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	spotifyID, ok := r.Context().Value(middleware.SpotifyIDContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify ID on request context"), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"user": User{SpotifyID: spotifyID}})
}

func ListCurrentUserScripts(w http.ResponseWriter, r *http.Request) {
	spotifyID, ok := r.Context().Value(middleware.SpotifyIDContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify ID on request context"), http.StatusInternalServerError)
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
	scripts, err := models.GetScriptsBySpotifyUserID(spotifyID, count, from, true)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get scripts for user: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"scripts": scripts})
}

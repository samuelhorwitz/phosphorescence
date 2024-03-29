package phosphor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"

	"github.com/rivo/uniseg"
)

const maxScriptVersionPageSize float64 = 10

func GetScript(w http.ResponseWriter, r *http.Request) {
	script, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	mostRecent, ok, err := models.GetMostRecentPublishedScriptVersion(script.ID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get most recent script version: %s", err), http.StatusInternalServerError)
		return
	}
	if ok {
		script.MostRecent = &mostRecent
	}
	common.JSON(w, map[string]interface{}{"script": script})
}

func GetScriptVersion(w http.ResponseWriter, r *http.Request) {
	scriptVersion, ok := r.Context().Value(middleware.ScriptVersionContextKey).(models.ScriptVersion)
	if !ok {
		common.Fail(w, errors.New("No script version on request context"), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"scriptVersion": scriptVersion})
}

func GetScriptVersions(w http.ResponseWriter, r *http.Request) {
	getScriptVersions(w, r, true)
}

func GetPrivateScriptVersions(w http.ResponseWriter, r *http.Request) {
	getScriptVersions(w, r, false)
}

func getScriptVersions(w http.ResponseWriter, r *http.Request, limitToPublished bool) {
	script, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	count, ok := r.Context().Value(middleware.PageCountContextKey).(uint64)
	if !ok {
		common.Fail(w, errors.New("No page count on request context"), http.StatusInternalServerError)
		return
	}
	from, ok := r.Context().Value(middleware.PageCursorContextKey).(time.Time)
	if !ok {
		common.Fail(w, errors.New("No page cursor on request context"), http.StatusInternalServerError)
		return
	}
	versions, err := models.GetScriptVersions(script.ID, count, from, limitToPublished)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get script versions: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"scriptVersions": versions})
}

func ListPublicScripts(w http.ResponseWriter, r *http.Request) {
	count, ok := r.Context().Value(middleware.PageCountContextKey).(uint64)
	if !ok {
		common.Fail(w, errors.New("No page count on request context"), http.StatusInternalServerError)
		return
	}
	from, ok := r.Context().Value(middleware.PageCursorContextKey).(time.Time)
	if !ok {
		common.Fail(w, errors.New("No page cursor on request context"), http.StatusInternalServerError)
		return
	}
	scripts, err := models.GetNewScripts(count, from)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get new scripts: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"scripts": scripts})
}

func CreateScript(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.AuthenticatedSessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read request body: %s", err), http.StatusInternalServerError)
		return
	}
	var requestBody struct {
		Script      string `json:"script"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not parse request body: %s", err), http.StatusInternalServerError)
		return
	}
	if err := validateScript(requestBody.Name, requestBody.Description, requestBody.Script); err != nil {
		common.Fail(w, fmt.Errorf("Bad request: %s", err), http.StatusBadRequest)
		return
	}
	createDetails, err := models.CreateScript(sess.SpotifyID, requestBody.Name, requestBody.Description, requestBody.Script)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not create script: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"create": createDetails})
}

func UpdateScript(w http.ResponseWriter, r *http.Request) {
	existingScript, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read request body: %s", err), http.StatusInternalServerError)
		return
	}
	var requestBody struct {
		Script      string `json:"script"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Permissions string `json:"permissions"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not parse request body: %s", err), http.StatusInternalServerError)
		return
	}
	if requestBody.Permissions == "" {
		common.Fail(w, errors.New("Permissions must be populated"), http.StatusBadRequest)
		return
	}
	if err := validateScript(requestBody.Name, requestBody.Description, requestBody.Script); err != nil {
		common.Fail(w, fmt.Errorf("Bad request: %s", err), http.StatusBadRequest)
		return
	}
	updateDetails, err := models.UpdateScript(existingScript.ID, requestBody.Name, requestBody.Description, requestBody.Script, requestBody.Permissions, models.ScriptSaveTypeDraft)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not update script: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"update": updateDetails})
}

func PublishScript(w http.ResponseWriter, r *http.Request) {
	existingScript, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read request body: %s", err), http.StatusInternalServerError)
		return
	}
	var requestBody struct {
		Script      string `json:"script"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Permissions string `json:"permissions"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not parse request body: %s", err), http.StatusInternalServerError)
		return
	}
	if requestBody.Permissions == "" {
		common.Fail(w, errors.New("Permissions must be populated"), http.StatusBadRequest)
		return
	}
	if err := validateScript(requestBody.Name, requestBody.Description, requestBody.Script); err != nil {
		common.Fail(w, fmt.Errorf("Bad request: %s", err), http.StatusBadRequest)
		return
	}
	updateDetails, err := models.UpdateScript(existingScript.ID, requestBody.Name, requestBody.Script, requestBody.Description, requestBody.Permissions, models.ScriptSaveTypePublished)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not update script: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"update": updateDetails})
}

func DeleteScript(w http.ResponseWriter, r *http.Request) {
	existingScript, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	err := models.DeleteScript(existingScript.ID)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not delete script: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"delete": true})
}

func DeleteScriptVersion(w http.ResponseWriter, r *http.Request) {
	existingScript, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	scriptVersion, ok := r.Context().Value(middleware.ScriptVersionContextKey).(models.ScriptVersion)
	if !ok {
		common.Fail(w, errors.New("No script version on request context"), http.StatusInternalServerError)
		return
	}
	if scriptVersion.CreatedAt.IsZero() {
		common.Fail(w, errors.New("Invalid script version"), http.StatusBadRequest)
		return
	}
	err := models.DeleteScriptVersion(existingScript.ID, scriptVersion.CreatedAt)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not delete script version: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"delete": true})
}

func ForkScript(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.AuthenticatedSessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	existingScript, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	scriptVersionID := common.ParseScriptVersion(r.URL.Query().Get("version"))
	forkDetails, err := models.ForkScript(sess.SpotifyID, existingScript.ID, existingScript.Name.String, scriptVersionID, true)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not fork script: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"fork": forkDetails})
}

func ForkScriptVersion(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.AuthenticatedSessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	existingScript, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	scriptVersion, ok := r.Context().Value(middleware.ScriptVersionContextKey).(models.ScriptVersion)
	if !ok {
		common.Fail(w, errors.New("No script version on request context"), http.StatusInternalServerError)
		return
	}
	forkDetails, err := models.ForkScript(sess.SpotifyID, existingScript.ID, existingScript.Name.String, scriptVersion.CreatedAt, true)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not fork script: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"fork": forkDetails})
}

func DuplicateScript(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.AuthenticatedSessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	existingScript, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read request body: %s", err), http.StatusInternalServerError)
		return
	}
	var requestBody struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not parse request body: %s", err), http.StatusInternalServerError)
		return
	}
	err = validateScriptName(requestBody.Name)
	if err != nil {
		common.Fail(w, fmt.Errorf("Bad request: %s", err), http.StatusBadRequest)
		return
	}
	var name string
	if requestBody.Name != "" {
		name = requestBody.Name
	} else {
		name = existingScript.Name.String
	}
	scriptVersionID := common.ParseScriptVersion(r.URL.Query().Get("version"))
	forkDetails, err := models.ForkScript(sess.SpotifyID, existingScript.ID, name, scriptVersionID, false)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not duplicate script: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"fork": forkDetails})
}

func DuplicateScriptVersion(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.AuthenticatedSessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	existingScript, ok := r.Context().Value(middleware.ScriptContextKey).(models.Script)
	if !ok {
		common.Fail(w, errors.New("No script on request context"), http.StatusInternalServerError)
		return
	}
	scriptVersion, ok := r.Context().Value(middleware.ScriptVersionContextKey).(models.ScriptVersion)
	if !ok {
		common.Fail(w, errors.New("No script version on request context"), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read request body: %s", err), http.StatusInternalServerError)
		return
	}
	var requestBody struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not parse request body: %s", err), http.StatusInternalServerError)
		return
	}
	err = validateScriptName(requestBody.Name)
	if err != nil {
		common.Fail(w, fmt.Errorf("Bad request: %s", err), http.StatusBadRequest)
		return
	}
	var name string
	if requestBody.Name != "" {
		name = requestBody.Name
	} else {
		name = existingScript.Name.String
	}
	forkDetails, err := models.ForkScript(sess.SpotifyID, existingScript.ID, name, scriptVersion.CreatedAt, false)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not duplicate script: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"fork": forkDetails})
}

func validateScript(name, description, script string) error {
	if script == "" {
		return errors.New("Script cannot be empty")
	}
	if description != noHTML.Sanitize(description) {
		return errors.New("Description contained HTML")
	}
	if err := validateScriptName(name); err != nil {
		return err
	}
	if len([]byte(description)) > 1024 {
		return errors.New("Description is too long")
	}
	return nil
}

func validateScriptName(name string) error {
	if uniseg.GraphemeClusterCount(name) > 60 {
		return errors.New("Name is too long")
	}
	return nil
}

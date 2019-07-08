package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

var PhosphorUUIDV5Namespace = uuid.NewV5(uuid.NamespaceDNS, "phosphor.me")

var SpotifyClient = &http.Client{
	Timeout: time.Second * 2,
}

func JSONRaw(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func JSON(w http.ResponseWriter, data interface{}) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		Fail(w, fmt.Errorf("Could not marshal JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	JSONRaw(w, jsonResponse)
}

func Fail(w http.ResponseWriter, err error, status int) {
	log.Printf("Request failed, returning %d: %s", status, err)
	http.Error(w, "error", status)
}

func RollbackAndFail(w http.ResponseWriter, tx *sql.Tx, err error, status int) {
	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		err = fmt.Errorf("Rollback failed: %s; Original Error: %s", rollbackErr, err)
	}
	Fail(w, err, status)
}

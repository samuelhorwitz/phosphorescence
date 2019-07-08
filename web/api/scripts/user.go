package scripts

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

func UserScripts(w http.ResponseWriter, r *http.Request) {
	spotifyID, ok := r.Context().Value(middleware.SpotifyIDContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify ID on request context"), http.StatusForbidden)
		return
	}
	rows, err := psql.Select("scripts.id", "scripts.is_private").
		From("scripts").
		Join("users ON users.id = scripts.author_id").
		Where(sq.Eq{
			"users.spotify_id":   spotifyID,
			"scripts.deleted_at": nil,
		}).
		RunWith(postgresDB).Query()
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not get scripts from DB: %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type scannedRow struct {
		ID        uuid.UUID `json:"id"`
		IsPrivate bool      `json:"isPrivate"`
	}
	var respBody []scannedRow
	for rows.Next() {
		var id uuid.UUID
		var isPrivate bool
		err := rows.Scan(&id, &isPrivate)
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not scan row: %s", err), http.StatusInternalServerError)
			return
		}
		respBody = append(respBody, scannedRow{
			ID:        id,
			IsPrivate: isPrivate,
		})
	}
	err = rows.Err()
	if err != nil {
		common.Fail(w, fmt.Errorf("Error after scanning rows: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, respBody)
}

func SaveScript(w http.ResponseWriter, r *http.Request) {
	spotifyID, ok := r.Context().Value(middleware.SpotifyIDContextKey).(string)
	if !ok {
		common.Fail(w, errors.New("No Spotify ID on request context"), http.StatusForbidden)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not read request body: %s", err), http.StatusInternalServerError)
		return
	}
	var requestBody struct {
		Script string `json:"script"`
		Type   string `json:"type"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not parse request body: %s", err), http.StatusInternalServerError)
		return
	}
	if requestBody.Script == "" || requestBody.Type == "" {
		common.Fail(w, errors.New("We need a script and type on the body"), http.StatusBadRequest)
		return
	}
	scriptID := uuid.NewV5(scriptsNamespace, requestBody.Script)
	var scriptExistsOnS3 bool
	_, err = s3Service.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String("phosphorescence-scripts"),
		Key:    aws.String(scriptID.String()),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != common.S3NotFoundCode {
				common.Fail(w, fmt.Errorf("Could not reach Spaces: %s", aerr.Error()), http.StatusInternalServerError)
				return
			}
		} else {
			common.Fail(w, fmt.Errorf("Could not reach Spaces: %s", err), http.StatusInternalServerError)
			return
		}
	} else {
		scriptExistsOnS3 = true
	}
	if !scriptExistsOnS3 {
		var scriptGZIPBuffer bytes.Buffer
		zw := gzip.NewWriter(&scriptGZIPBuffer)
		_, err := zw.Write([]byte(requestBody.Script))
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not compress script: %s", err), http.StatusInternalServerError)
			return
		}
		if err := zw.Close(); err != nil {
			common.Fail(w, fmt.Errorf("Could not close compression buffer: %s", err), http.StatusInternalServerError)
			return
		}
		gzipReader, err := gzip.NewReader(&scriptGZIPBuffer)
		if err != nil {
			common.Fail(w, fmt.Errorf("Cannot read compressed buffer: %s", err), http.StatusInternalServerError)
			return
		}
		_, err = s3Service.PutObject(&s3.PutObjectInput{
			Bucket:          aws.String("phosphorescence-scripts"),
			ACL:             aws.String("public-read"),
			ContentType:     aws.String("application/javascript"),
			ContentEncoding: aws.String("gzip"),
			Key:             aws.String(scriptID.String()),
			Body:            aws.ReadSeekCloser(gzipReader),
		})
		if err != nil {
			common.Fail(w, fmt.Errorf("Could not upload to Spaces: %s", err), http.StatusInternalServerError)
			return
		}
	}
	tx, err := postgresDB.Begin()
	if err != nil {
		common.RollbackAndFail(w, tx, fmt.Errorf("Could not start transaction: %s", err), http.StatusInternalServerError)
		return
	}
	var userID uuid.UUID
	err = psql.Select("id").From("users").
		Where(sq.Eq{
			"spotify_id": spotifyID,
		}).RunWith(tx).QueryRow().Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			userID = uuid.NewV4()
			_, err = psql.Insert("users").Columns("id", "spotify_id").
				Values(userID, spotifyID).
				RunWith(tx).Exec()
			if err != nil {
				common.RollbackAndFail(w, tx, fmt.Errorf("Could not insert new user: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			common.RollbackAndFail(w, tx, fmt.Errorf("Could not query user: %s", err), http.StatusInternalServerError)
			return
		}
	}
	scriptInternalID := uuid.NewV4()
	_, err = psql.Insert("scripts").Columns("id", "author_id").
		Values(scriptInternalID, userID).
		RunWith(tx).Exec()
	if err != nil {
		common.RollbackAndFail(w, tx, fmt.Errorf("Could not insert new script: %s", err), http.StatusInternalServerError)
		return
	}
	_, err = psql.Insert("script_versions").Columns("script_id", "type", "file_id").
		Values(scriptInternalID, requestBody.Type, scriptID).
		RunWith(tx).Exec()
	if err != nil {
		common.RollbackAndFail(w, tx, fmt.Errorf("Could not insert new script version: %s", err), http.StatusInternalServerError)
		return
	}
	if err = tx.Commit(); err != nil {
		common.RollbackAndFail(w, tx, fmt.Errorf("Could not commit: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, struct {
		ID                string `json:"id"`
		MostRecentVersion string `json:"mostRecentVersion"`
	}{
		ID:                scriptInternalID.String(),
		MostRecentVersion: scriptID.String(),
	})
}

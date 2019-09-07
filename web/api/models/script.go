package models

import (
	"compress/gzip"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/satori/go.uuid"
	"io"
	"time"
)

var scriptsNamespace = uuid.NewV5(common.PhosphorUUIDV5Namespace, "scripts")

type Script struct {
	ID                      uuid.UUID      `json:"id"`
	Name                    nullString     `json:"name"`
	Description             nullString     `json:"description"`
	AuthorID                uuid.UUID      `json:"authorId"`
	AuthorSpotifyID         nullString     `json:"authorSpotifyId"`
	AuthorName              nullString     `json:"authorName"`
	ForkedFromScriptID      nullUUID       `json:"forkedFromScriptId"`
	ForkedFromScriptVersion nullTime       `json:"forkedFromScriptVersion"`
	IsPrivate               bool           `json:"isPrivate,omitempty"`
	MostRecent              *ScriptVersion `json:"mostRecent,omitempty"`
	CreatedAt               time.Time      `json:"createdAt"`
}

type ScriptVersion struct {
	CreatedAt time.Time      `json:"createdAt"`
	Type      ScriptSaveType `json:"type"`
	FileID    uuid.UUID      `json:"fileId"`
}

type CreateOrUpdateScriptResponse struct {
	ID          uuid.UUID `json:"id"`
	FileID      uuid.UUID `json:"fileId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions string    `json:"permissions"`
}

type ScriptSaveType string

const (
	ScriptSaveTypeDraft     ScriptSaveType = "draft"
	ScriptSaveTypePublished ScriptSaveType = "publish"
	ScriptSaveTypeFork      ScriptSaveType = "fork"
)

const (
	permissionsPublic  = "public"
	permissionsPrivate = "private"
)

func GetScriptWithAuthorizationCheck(spotifyUserID string, scriptID uuid.UUID) (Script, bool, error) {
	var script Script
	err := psql.Select(
		"scripts.id",
		"scripts.author_id",
		"scripts.name",
		"scripts.description",
		"users.spotify_id",
		"users.name",
		"scripts.forked_from_script_id",
		"scripts.forked_from_script_version_created_at",
		"scripts.is_private",
		"scripts.created_at").
		From("scripts_view as scripts").
		LeftJoin("users_view users on users.id = scripts.author_id").
		Where(sq.And{
			sq.Eq{"scripts.id": scriptID},
			sq.Or{
				sq.Eq{"users.spotify_id": spotifyUserID},
				sq.Eq{"scripts.is_private": false},
			},
		}).
		RunWith(postgresDB).QueryRow().Scan(
		&script.ID,
		&script.AuthorID,
		&script.Name,
		&script.Description,
		&script.AuthorSpotifyID,
		&script.AuthorName,
		&script.ForkedFromScriptID,
		&script.ForkedFromScriptVersion,
		&script.IsPrivate,
		&script.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Script{}, false, nil
		}
		return Script{}, false, fmt.Errorf("Could not query for script: %s", err)
	}
	return script, true, nil
}

func GetScriptVersionWithAuthorizationCheck(spotifyUserID string, scriptID uuid.UUID, scriptVersionID time.Time) (ScriptVersion, bool, error) {
	var scriptVersion ScriptVersion
	err := psql.Select("script_versions.created_at", "script_versions.type", "script_versions.file_id").
		From("script_versions_view as script_versions").
		Join("scripts_view scripts on scripts.id = script_versions.script_id").
		LeftJoin("users_view users on users.id = scripts.author_id").
		Where(sq.And{
			sq.Eq{"scripts.id": scriptID},
			sq.Or{
				sq.Eq{"users.spotify_id": spotifyUserID},
				sq.And{
					sq.Eq{"scripts.is_private": false},
					sq.Eq{"script_versions.type": ScriptSaveTypePublished},
				},
			},
		}).
		RunWith(postgresDB).QueryRow().Scan(&scriptVersion.CreatedAt, &scriptVersion.Type, &scriptVersion.FileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ScriptVersion{}, false, nil
		}
		return ScriptVersion{}, false, fmt.Errorf("Could not query for script version: %s", err)
	}
	return scriptVersion, true, nil
}

func GetScriptsBySpotifyUserID(spotifyUserID string, count uint64, from time.Time, includePrivate bool) (scripts []Script, err error) {
	where := sq.And{sq.Eq{"users.spotify_id": spotifyUserID}}
	if !from.IsZero() {
		where = append(where, sq.Lt{"scripts.created_at": from})
	}
	if !includePrivate {
		where = append(where, sq.Eq{"scripts.is_private": false})
	}
	rows, err := psql.Select(
		"scripts.id",
		"scripts.author_id",
		"scripts.name",
		"scripts.description",
		"users.spotify_id",
		"users.name",
		"scripts.forked_from_script_id",
		"scripts.forked_from_script_version_created_at",
		"scripts.is_private",
		"scripts.created_at").
		From("scripts_view as scripts").
		Join("users_view users on users.id = scripts.author_id").
		Where(where).
		OrderBy("scripts.created_at desc").
		Limit(count).
		RunWith(postgresDB).Query()
	if err != nil {
		return nil, fmt.Errorf("Could not get scripts from DB: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var script Script
		err := rows.Scan(
			&script.ID,
			&script.AuthorID,
			&script.Name,
			&script.Description,
			&script.AuthorSpotifyID,
			&script.AuthorName,
			&script.ForkedFromScriptID,
			&script.ForkedFromScriptVersion,
			&script.IsPrivate,
			&script.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Could not scan row: %s", err)
		}
		scripts = append(scripts, script)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Error after scanning rows: %s", err)
	}
	return scripts, nil
}

func GetMostRecentPublishedScriptVersion(scriptID uuid.UUID) (version ScriptVersion, ok bool, err error) {
	err = psql.Select("created_at", "type", "file_id").
		From("script_versions_view as script_versions").
		Where(sq.Eq{
			"script_id": scriptID,
			"type":      ScriptSaveTypePublished,
		}).
		OrderBy("created_at desc").
		Limit(1).
		RunWith(postgresDB).QueryRow().Scan(&version.CreatedAt, &version.Type, &version.FileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ScriptVersion{}, false, nil
		}
		return version, false, fmt.Errorf("Could not get most recent script version: %s", err)
	}
	return version, true, nil
}

func GetScriptVersions(scriptID uuid.UUID, count uint64, from time.Time, limitToPublished bool) (versions []ScriptVersion, err error) {
	where := sq.And{sq.Eq{"script_id": scriptID}}
	if !from.IsZero() {
		where = append(where, sq.Lt{"created_at": from})
	}
	if limitToPublished {
		where = append(where, sq.Eq{"type": ScriptSaveTypePublished})
	}
	sel := psql.Select("created_at", "type", "file_id").
		From("script_versions_view as script_versions").
		Where(where).
		OrderBy("created_at desc").
		Limit(count)
	rows, err := sel.RunWith(postgresDB).Query()
	if err != nil {
		return nil, fmt.Errorf("Could not get script versions: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var version ScriptVersion
		err := rows.Scan(&version.CreatedAt, &version.Type, &version.FileID)
		if err != nil {
			return nil, fmt.Errorf("Could not scan row: %s", err)
		}
		versions = append(versions, version)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Error after scanning rows: %s", err)
	}
	return versions, nil
}

func GetNewScripts(count uint64, from time.Time) (scripts []Script, err error) {
	where := sq.And{sq.Eq{"is_private": false}}
	if !from.IsZero() {
		where = append(where, sq.Lt{"scripts.created_at": from})
	}
	rows, err := psql.Select(
		"scripts.id",
		"scripts.author_id",
		"scripts.name",
		"scripts.description",
		"users.spotify_id",
		"users.name",
		"scripts.forked_from_script_id",
		"scripts.forked_from_script_version_created_at",
		"scripts.is_private",
		"scripts.created_at").
		From("scripts_view as scripts").
		LeftJoin("users_view users on users.id = scripts.author_id").
		Where(where).
		OrderBy("scripts.created_at desc").
		Limit(count).
		RunWith(postgresDB).Query()
	if err != nil {
		return nil, fmt.Errorf("Could not get scripts from DB: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var script Script
		err := rows.Scan(
			&script.ID,
			&script.AuthorID,
			&script.Name,
			&script.Description,
			&script.AuthorSpotifyID,
			&script.AuthorName,
			&script.ForkedFromScriptID,
			&script.ForkedFromScriptVersion,
			&script.IsPrivate,
			&script.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Could not scan row: %s", err)
		}
		scripts = append(scripts, script)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Error after scanning rows: %s", err)
	}
	return scripts, nil
}

func CreateScript(spotifyUserID, name, description, script string) (CreateOrUpdateScriptResponse, error) {
	scriptFileID, err := pushScriptToObjectStorage(script)
	if err != nil {
		return CreateOrUpdateScriptResponse{}, fmt.Errorf("Could not push script to Spaces: %s", err)
	}
	tx, err := postgresDB.Begin()
	if err != nil {
		return CreateOrUpdateScriptResponse{}, fmt.Errorf("Could not start transaction: %s", err)
	}
	userID, err := mapSpotifyIDToOurID(tx, spotifyUserID)
	if err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not get user ID from Spotify ID: %s", err))
	}
	scriptID := uuid.NewV4()
	_, err = psql.Insert("scripts").Columns("id", "author_id", "name", "description").
		Values(scriptID, userID, stringOrNull(name), stringOrNull(description)).
		RunWith(tx).Exec()
	if err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script: %s", err))
	}
	_, err = psql.Insert("script_versions").Columns("script_id", "type", "file_id").
		Values(scriptID, ScriptSaveTypeDraft, scriptFileID).
		RunWith(tx).Exec()
	if err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script version: %s", err))
	}
	if err = tx.Commit(); err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not commit: %s", err))
	}
	return CreateOrUpdateScriptResponse{
		ID:          scriptID,
		FileID:      scriptFileID,
		Name:        name,
		Permissions: permissionsPrivate,
	}, nil
}

func UpdateScript(scriptID uuid.UUID, name, description, script, permissions string, scriptSaveType ScriptSaveType) (CreateOrUpdateScriptResponse, error) {
	if name == "" && permissions == permissionsPublic {
		return CreateOrUpdateScriptResponse{}, errors.New("Public scripts must have a name")
	}
	res := CreateOrUpdateScriptResponse{}
	tx, err := postgresDB.Begin()
	if err != nil {
		return CreateOrUpdateScriptResponse{}, fmt.Errorf("Could not start transaction: %s", err)
	}
	if script != "" {
		scriptFileID, err := pushScriptToObjectStorage(script)
		if err != nil {
			return CreateOrUpdateScriptResponse{}, fmt.Errorf("Could not push script to Spaces: %s", err)
		}
		var mostRecentFileID uuid.UUID
		var mostRecentType ScriptSaveType
		err = psql.Select("file_id", "type").
			From("script_versions_view as script_versions").
			Where(sq.Eq{"script_id": scriptID}).
			OrderBy("created_at desc").
			Limit(1).
			RunWith(tx).QueryRow().Scan(&mostRecentFileID, &mostRecentType)
		if err != nil {
			return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not check most recent script version file ID: %s", err))
		}
		// Don't create a new version row if the code is the exact same, unless we are publishing and the previous save was not a publish
		if !uuid.Equal(mostRecentFileID, scriptFileID) || (mostRecentType != ScriptSaveTypePublished && scriptSaveType == ScriptSaveTypePublished) {
			_, err = psql.Insert("script_versions").Columns("script_id", "type", "file_id").
				Values(scriptID, scriptSaveType, scriptFileID).
				RunWith(tx).Exec()
			if err != nil {
				return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script version: %s", err))
			}
		}
		res.ID = scriptID
		res.FileID = scriptFileID
	}
	updateBuilder := psql.Update("scripts").Where(sq.Eq{"id": scriptID}).RunWith(tx)
	shouldUpdate := false
	if name != "" {
		updateBuilder.Set("name", name)
		res.Name = name
		shouldUpdate = true
	}
	if description != "" {
		updateBuilder.Set("description", description)
		res.Description = description
		shouldUpdate = true
	}
	if permissions != "" {
		var isPrivate bool
		switch permissions {
		case permissionsPublic:
			isPrivate = false
		case permissionsPrivate:
			isPrivate = true
		default:
			return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not update script permissions, invalid permissions %s", permissions))
		}
		updateBuilder.Set("isPrivate", isPrivate)
		res.Permissions = permissions
		shouldUpdate = true
	}
	if shouldUpdate {
		_, err := updateBuilder.Exec()
		if err != nil {
			return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not update script details: %s", err))
		}
	}
	if err = tx.Commit(); err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not commit: %s", err))
	}
	return res, nil
}

func DeleteScript(scriptID uuid.UUID) error {
	_, err := psql.Update("scripts").
		Set("deleted_at", sq.Expr("now()")).
		Where(sq.Eq{"id": scriptID}).
		RunWith(postgresDB).Exec()
	if err != nil {
		return fmt.Errorf("Could not mark script as deleted: %s", err)
	}
	return nil
}

func DeleteScriptVersion(scriptID uuid.UUID, scriptVersionID time.Time) error {
	_, err := psql.Update("script_versions").
		Set("deleted_at", sq.Expr("now()")).
		Where(sq.Eq{
			"script_id":  scriptID,
			"created_at": scriptVersionID,
		}).
		RunWith(postgresDB).Exec()
	if err != nil {
		return fmt.Errorf("Could not mark script version as deleted: %s", err)
	}
	return nil
}

func ForkScript(spotifyUserID string, toForkScriptID uuid.UUID, toForkScriptName string, toForkVersionID time.Time, onlyPublished bool) (CreateOrUpdateScriptResponse, error) {
	tx, err := postgresDB.Begin()
	if err != nil {
		return CreateOrUpdateScriptResponse{}, fmt.Errorf("Could not start transaction: %s", err)
	}
	userID, err := mapSpotifyIDToOurID(tx, spotifyUserID)
	if err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not get user ID from Spotify ID: %s", err))
	}
	where := sq.Eq{
		"script_id": toForkScriptID,
	}
	if onlyPublished {
		where["type"] = ScriptSaveTypePublished
	}
	sel := psql.Select("created_at", "type", "file_id").
		From("script_versions_view as script_versions").
		Limit(1)
	if !toForkVersionID.IsZero() {
		where["created_at"] = toForkVersionID
	} else {
		sel = sel.OrderBy("created_at desc")
	}
	var toForkVersion ScriptVersion
	err = sel.Where(where).
		RunWith(tx).QueryRow().Scan(&toForkVersion.CreatedAt, &toForkVersion.Type, &toForkVersion.FileID)
	if err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not get script version: %s", err))
	}
	scriptID := uuid.NewV4()
	_, err = psql.Insert("scripts").Columns("id", "author_id", "name", "forked_from_script_id", "forked_from_script_version_created_at").
		Values(scriptID, userID, stringOrNull(toForkScriptName), toForkScriptID, toForkVersion.CreatedAt).
		RunWith(tx).Exec()
	if err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script: %s", err))
	}
	_, err = psql.Insert("script_versions").Columns("script_id", "type", "file_id").
		Values(scriptID, ScriptSaveTypeFork, toForkVersion.FileID).
		RunWith(tx).Exec()
	if err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script version: %s", err))
	}
	if err = tx.Commit(); err != nil {
		return CreateOrUpdateScriptResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not commit: %s", err))
	}
	return CreateOrUpdateScriptResponse{
		ID:          scriptID,
		FileID:      toForkVersion.FileID,
		Name:        toForkScriptName,
		Permissions: permissionsPrivate,
	}, nil
}

func pushScriptToObjectStorage(script string) (scriptID uuid.UUID, err error) {
	scriptID = uuid.NewV5(scriptsNamespace, script)
	scriptInObjectStorage, err := isScriptIDInObjectStorage(scriptID)
	if err != nil {
		return scriptID, fmt.Errorf("Could access Spaces: %s", err)
	}
	if scriptInObjectStorage {
		return scriptID, nil
	}
	reader, writer := io.Pipe()
	go func() {
		zw := gzip.NewWriter(writer)
		if _, err = zw.Write([]byte(script)); err != nil {
			err = fmt.Errorf("Could not compress script: %s", err)
			return
		}
		if err := zw.Close(); err != nil {
			err = fmt.Errorf("Could not close compression buffer: %s", err)
			return
		}
		writer.Close()
	}()
	if err != nil {
		return scriptID, fmt.Errorf("Could not zip script: %s", err)
	}
	_, err = s3Uploader.Upload(&s3manager.UploadInput{
		Bucket:          aws.String("phosphorescence-scripts"),
		ACL:             aws.String("public-read"),
		ContentType:     aws.String("application/javascript"),
		ContentEncoding: aws.String("gzip"),
		CacheControl:    aws.String("public, max-age=31536000"),
		Key:             aws.String(scriptID.String()),
		Body:            reader,
	})
	if err != nil {
		return scriptID, fmt.Errorf("Could not upload to Spaces: %s", err)
	}
	return scriptID, nil
}

func isScriptInObjectStorage(script string) (bool, error) {
	return isScriptIDInObjectStorage(uuid.NewV5(scriptsNamespace, script))
}

func isScriptIDInObjectStorage(scriptID uuid.UUID) (bool, error) {
	_, err := s3Service.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String("phosphorescence-scripts"),
		Key:    aws.String(scriptID.String()),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == common.S3NotFoundCode {
				return false, nil
			}
			return false, fmt.Errorf("Could not reach Spaces: %s", aerr.Error())
		}
		return false, fmt.Errorf("Could not reach Spaces: %s", err)
	}
	return true, nil
}

func stringOrNull(str string) sql.NullString {
	var nullString sql.NullString
	if str == "" {
		nullString.Valid = false
	} else {
		nullString.Valid = true
		nullString.String = str
	}
	return nullString
}

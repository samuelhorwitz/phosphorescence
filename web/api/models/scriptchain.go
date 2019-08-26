package models

import (
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/satori/go.uuid"
	"time"
)

type ScriptChain struct {
	ID                           uuid.UUID           `json:"id"`
	Name                         nullString          `json:"name"`
	Description                  nullString          `json:"description"`
	AuthorID                     uuid.UUID           `json:"authorId"`
	AuthorSpotifyID              nullString          `json:"authorSpotifyId"`
	AuthorName                   nullString          `json:"authorName"`
	ForkedFromScriptChainID      nullUUID            `json:"forkedFromScriptChainId"`
	ForkedFromScriptChainVersion nullTime            `json:"forkedFromScriptChainVersion"`
	IsPrivate                    bool                `json:"isPrivate,omitempty"`
	MostRecent                   *ScriptChainVersion `json:"mostRecent,omitempty"`
	CreatedAt                    time.Time           `json:"createdAt"`
}

type ScriptChainVersion struct {
	CreatedAt time.Time      `json:"createdAt"`
	Type      ScriptSaveType `json:"type"`
	PrunerIDs []uuid.UUID    `json:"prunerIds"`
	SeederID  uuid.UUID      `json:"seederId"`
	BuilderID uuid.UUID      `json:"builderId"`
}

type CreateOrUpdateScriptChainResponse struct {
	ID          uuid.UUID   `json:"id"`
	PrunerIDs   []uuid.UUID `json:"prunerIds"`
	SeederID    uuid.UUID   `json:"seederId"`
	BuilderID   uuid.UUID   `json:"builderId"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Permissions string      `json:"permissions"`
}

func GetScriptChainWithAuthorizationCheck(spotifyUserID string, scriptChainID uuid.UUID) (ScriptChain, bool, error) {
	var scriptChain ScriptChain
	err := psql.Select(
		"script_chains.id",
		"script_chains.author_id",
		"script_chains.name",
		"script_chains.description",
		"users.spotify_id",
		"users.name",
		"script_chains.forked_from_script_chain_id",
		"script_chains.forked_from_script_chain_version_created_at",
		"script_chains.is_private",
		"script_chains.created_at").
		From("script_chains_view as script_chains").
		LeftJoin("users_view users on users.id = script_chains.author_id").
		Where(sq.And{
			sq.Eq{"script_chains.id": scriptChainID},
			sq.Or{
				sq.Eq{"users.spotify_id": spotifyUserID},
				sq.Eq{"script_chains.is_private": false},
			},
		}).
		RunWith(postgresDB).QueryRow().Scan(
		&scriptChain.ID,
		&scriptChain.AuthorID,
		&scriptChain.Name,
		&scriptChain.Description,
		&scriptChain.AuthorSpotifyID,
		&scriptChain.AuthorName,
		&scriptChain.ForkedFromScriptChainID,
		&scriptChain.ForkedFromScriptChainVersion,
		&scriptChain.IsPrivate,
		&scriptChain.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return ScriptChain{}, false, nil
		}
		return ScriptChain{}, false, fmt.Errorf("Could not query for script chain: %s", err)
	}
	return scriptChain, true, nil
}

func GetScriptChainVersionWithAuthorizationCheck(spotifyUserID string, scriptChainID uuid.UUID, scriptChainVersionID time.Time) (ScriptChainVersion, bool, error) {
	var scriptChainVersion ScriptChainVersion
	err := psql.Select("script_chain_versions.created_at", "script_chain_versions.type", "script_chain_versions.seeder_id", "script_chain_versions.builder_id", "script_chain_versions.pruner_ids").
		From("script_chain_versions_view as script_chain_versions").
		Join("script_chains_view script_chains on script_chains.id = script_chain_versions.script_id").
		LeftJoin("users_view users on users.id = script_chains.author_id").
		Where(sq.And{
			sq.Eq{"script_chains.id": scriptChainID},
			sq.Or{
				sq.Eq{"users.spotify_id": spotifyUserID},
				sq.And{
					sq.Eq{"script_chains.is_private": false},
				},
			},
		}).
		RunWith(postgresDB).QueryRow().Scan(&scriptChainVersion.CreatedAt, &scriptChainVersion.Type, &scriptChainVersion.SeederID, &scriptChainVersion.BuilderID, &scriptChainVersion.PrunerIDs)
	if err != nil {
		if err == sql.ErrNoRows {
			return ScriptChainVersion{}, false, nil
		}
		return ScriptChainVersion{}, false, fmt.Errorf("Could not query for script chain version: %s", err)
	}
	return scriptChainVersion, true, nil
}

func GetScriptChainsBySpotifyUserID(spotifyUserID string, count uint64, from time.Time, includePrivate bool) (scriptChains []ScriptChain, err error) {
	where := sq.And{sq.Eq{"users.spotify_id": spotifyUserID}}
	if !from.IsZero() {
		where = append(where, sq.Lt{"script_chains.created_at": from})
	}
	if !includePrivate {
		where = append(where, sq.Eq{"script_chains.is_private": false})
	}
	rows, err := psql.Select(
		"script_chains.id",
		"script_chains.author_id",
		"script_chains.name",
		"script_chains.description",
		"users.spotify_id",
		"users.name",
		"script_chains.forked_from_script_id",
		"script_chains.forked_from_script_version_created_at",
		"script_chains.is_private",
		"script_chains.created_at").
		From("script_chains_view as script_chains").
		Join("users_view users on users.id = script_chains.author_id").
		Where(where).
		OrderBy("script_chains.created_at desc").
		Limit(count).
		RunWith(postgresDB).Query()
	if err != nil {
		return nil, fmt.Errorf("Could not get script chains from DB: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var scriptChain ScriptChain
		err := rows.Scan(
			&scriptChain.ID,
			&scriptChain.AuthorID,
			&scriptChain.Name,
			&scriptChain.Description,
			&scriptChain.AuthorSpotifyID,
			&scriptChain.AuthorName,
			&scriptChain.ForkedFromScriptChainID,
			&scriptChain.ForkedFromScriptChainVersion,
			&scriptChain.IsPrivate,
			&scriptChain.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Could not scan row: %s", err)
		}
		scriptChains = append(scriptChains, scriptChain)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Error after scanning rows: %s", err)
	}
	return scriptChains, nil
}

func GetMostRecentScriptChainVersion(scriptChainID uuid.UUID) (version ScriptChainVersion, ok bool, err error) {
	err = psql.Select("created_at", "type", "seeder_id", "builder_id", "pruner_ids").
		From("script_chain_versions_view as script_chain_versions").
		Where(sq.Eq{
			"script_chain_id": scriptChainID,
		}).
		OrderBy("created_at desc").
		Limit(1).
		RunWith(postgresDB).QueryRow().Scan(&version.CreatedAt, &version.Type, &version.SeederID, &version.BuilderID, &version.PrunerIDs)
	if err != nil {
		if err == sql.ErrNoRows {
			return ScriptChainVersion{}, false, nil
		}
		return version, false, fmt.Errorf("Could not get most recent script chain version: %s", err)
	}
	return version, true, nil
}

func GetScriptChainVersions(scriptChainID uuid.UUID, count uint64, from time.Time) (versions []ScriptChainVersion, err error) {
	where := sq.And{sq.Eq{"script_chain_id": scriptChainID}}
	if !from.IsZero() {
		where = append(where, sq.Lt{"created_at": from})
	}
	sel := psql.Select("created_at", "type", "seeder_id", "builder_id", "pruner_ids").
		From("script_chain_versions_view as script_chain_versions").
		Where(where).
		OrderBy("created_at desc").
		Limit(count)
	rows, err := sel.RunWith(postgresDB).Query()
	if err != nil {
		return nil, fmt.Errorf("Could not get script chain versions: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var version ScriptChainVersion
		err := rows.Scan(&version.CreatedAt, &version.Type, &version.SeederID, &version.BuilderID, &version.PrunerIDs)
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

func GetNewScriptChains(count uint64, from time.Time) (scriptChains []ScriptChain, err error) {
	where := sq.And{sq.Eq{"is_private": false}}
	if !from.IsZero() {
		where = append(where, sq.Lt{"script_chains.created_at": from})
	}
	rows, err := psql.Select(
		"script_chains.id",
		"script_chains.author_id",
		"script_chains.name",
		"script_chains.description",
		"users.spotify_id",
		"users.name",
		"script_chains.forked_from_script_id",
		"script_chains.forked_from_script_version_created_at",
		"script_chains.is_private",
		"script_chains.created_at").
		From("script_chains_view as script_chains").
		LeftJoin("users_view users on users.id = script_chains.author_id").
		Where(where).
		OrderBy("script_chains.created_at desc").
		Limit(count).
		RunWith(postgresDB).Query()
	if err != nil {
		return nil, fmt.Errorf("Could not get script chains from DB: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var scriptChain ScriptChain
		err := rows.Scan(
			&scriptChain.ID,
			&scriptChain.AuthorID,
			&scriptChain.Name,
			&scriptChain.Description,
			&scriptChain.AuthorSpotifyID,
			&scriptChain.AuthorName,
			&scriptChain.ForkedFromScriptChainID,
			&scriptChain.ForkedFromScriptChainVersion,
			&scriptChain.IsPrivate,
			&scriptChain.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Could not scan row: %s", err)
		}
		scriptChains = append(scriptChains, scriptChain)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Error after scanning rows: %s", err)
	}
	return scriptChains, nil
}

func CreateScriptChain(spotifyUserID, name, description string, seederID, builderID uuid.UUID, prunerIDs []uuid.UUID) (CreateOrUpdateScriptChainResponse, error) {
	tx, err := postgresDB.Begin()
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, fmt.Errorf("Could not start transaction: %s", err)
	}
	userID, err := mapSpotifyIDToOurID(tx, spotifyUserID)
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not get user ID from Spotify ID: %s", err))
	}
	scriptID := uuid.NewV4()
	_, err = psql.Insert("script_chains").Columns("id", "author_id", "name", "description").
		Values(scriptID, userID, stringOrNull(name), stringOrNull(description)).
		RunWith(tx).Exec()
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script chain: %s", err))
	}
	_, err = psql.Insert("script_chain_versions").Columns("script_id", "seeder_id", "builder_id").
		Values(scriptID, seederID, builderID).
		RunWith(tx).Exec()
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script chain version: %s", err))
	}
	// TODO pruner
	if err = tx.Commit(); err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not commit: %s", err))
	}
	return CreateOrUpdateScriptChainResponse{
		ID:        scriptID,
		SeederID:  seederID,
		BuilderID: builderID,
		// TODO pruner IDs
		Name:        name,
		Permissions: permissionsPrivate,
	}, nil
}

func UpdateScriptChain(scriptChainID uuid.UUID, name, description string, seederID, builderID uuid.UUID, prunerIDs []uuid.UUID, permissions string) (CreateOrUpdateScriptChainResponse, error) {
	if name == "" && permissions == permissionsPublic {
		return CreateOrUpdateScriptChainResponse{}, errors.New("Public script chains must have a name")
	}
	res := CreateOrUpdateScriptChainResponse{}
	tx, err := postgresDB.Begin()
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, fmt.Errorf("Could not start transaction: %s", err)
	}
	if !uuid.Equal(seederID, uuid.Nil) || !uuid.Equal(builderID, uuid.Nil) || len(prunerIDs) > 0 {
		var mostRecentSeederID uuid.UUID
		var mostRecentBuilderID uuid.UUID
		var mostRecentPrunerIDs []uuid.UUID
		err = psql.Select("seeder_id", "builder_id", "pruner_ids").
			From("script_chain_versions_view as script_chain_versions").
			Where(sq.Eq{"script_chain_id": scriptChainID}).
			OrderBy("created_at desc").
			Limit(1).
			RunWith(tx).QueryRow().Scan(&mostRecentSeederID, &mostRecentBuilderID, &mostRecentPrunerIDs)
		if err != nil {
			return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not check most recent script chain version: %s", err))
		}
		// Don't create a new version row if the code is the exact same
		isSameAsOld := uuid.Equal(mostRecentSeederID, seederID) && uuid.Equal(mostRecentBuilderID, builderID)
		if isSameAsOld {
			if len(prunerIDs) != len(mostRecentPrunerIDs) {
				isSameAsOld = false
			} else {
				mostRecentPrunerIDsMap := make(map[string]bool)
				for _, mostRecentPrunerID := range mostRecentPrunerIDs {
					mostRecentPrunerIDsMap[mostRecentPrunerID.String()] = true
				}
				for _, prunerID := range prunerIDs {
					if !mostRecentPrunerIDsMap[prunerID.String()] {
						isSameAsOld = false
						break
					}
				}
			}
		}
		if !isSameAsOld {
			_, err = psql.Insert("script_chain_versions").Columns("script_id", "seeder_id", "builder_id").
				Values(scriptChainID, seederID, builderID).
				RunWith(tx).Exec()
			if err != nil {
				return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script chain version: %s", err))
			}
			// TODO pruner
		}
		res.ID = scriptChainID
		res.SeederID = seederID
		res.BuilderID = builderID
		// TODO pruner
	}
	updateBuilder := psql.Update("script_chains").Where(sq.Eq{"id": scriptChainID}).RunWith(tx)
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
			return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not update script chain permissions, invalid permissions %s", permissions))
		}
		updateBuilder.Set("isPrivate", isPrivate)
		res.Permissions = permissions
		shouldUpdate = true
	}
	if shouldUpdate {
		_, err := updateBuilder.Exec()
		if err != nil {
			return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not update script chain details: %s", err))
		}
	}
	if err = tx.Commit(); err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not commit: %s", err))
	}
	return res, nil
}

func DeleteScriptChain(scriptChainID uuid.UUID) error {
	_, err := psql.Update("script_chains").
		Set("deleted_at", sq.Expr("now()")).
		Where(sq.Eq{"id": scriptChainID}).
		RunWith(postgresDB).Exec()
	if err != nil {
		return fmt.Errorf("Could not mark script chain as deleted: %s", err)
	}
	return nil
}

func DeleteScriptChainVersion(scriptChainID uuid.UUID, scriptChainVersionID time.Time) error {
	_, err := psql.Update("script_chain_versions").
		Set("deleted_at", sq.Expr("now()")).
		Where(sq.Eq{
			"script_chain_id": scriptChainID,
			"created_at":      scriptChainVersionID,
		}).
		RunWith(postgresDB).Exec()
	if err != nil {
		return fmt.Errorf("Could not mark script chain version as deleted: %s", err)
	}
	return nil
}

func ForkScriptChain(spotifyUserID string, toForkScriptChainID uuid.UUID, toForkScriptChainName string, toForkVersionID time.Time) (CreateOrUpdateScriptChainResponse, error) {
	tx, err := postgresDB.Begin()
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, fmt.Errorf("Could not start transaction: %s", err)
	}
	userID, err := mapSpotifyIDToOurID(tx, spotifyUserID)
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not get user ID from Spotify ID: %s", err))
	}
	where := sq.Eq{
		"script_chain_id": toForkScriptChainID,
	}
	sel := psql.Select("created_at", "type", "seeder_id", "builder_id", "pruner_ids").
		From("script_chain_versions_view as script_chain_versions").
		Limit(1)
	if !toForkVersionID.IsZero() {
		where["created_at"] = toForkVersionID
	} else {
		sel = sel.OrderBy("created_at desc")
	}
	var toForkVersion ScriptChainVersion
	err = sel.Where(where).
		RunWith(tx).QueryRow().Scan(&toForkVersion.CreatedAt, &toForkVersion.Type, &toForkVersion.SeederID, &toForkVersion.BuilderID, &toForkVersion.PrunerIDs)
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not get script chain version: %s", err))
	}
	scriptChainID := uuid.NewV4()
	_, err = psql.Insert("script_chains").Columns("id", "author_id", "name", "forked_from_script_chain_id", "forked_from_script_chain_version_created_at").
		Values(scriptChainID, userID, stringOrNull(toForkScriptChainName), toForkScriptChainID, toForkVersion.CreatedAt).
		RunWith(tx).Exec()
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script chain: %s", err))
	}
	_, err = psql.Insert("script_chain_versions").Columns("script_id", "type", "seeder_id", "builder_id").
		Values(scriptChainID, ScriptSaveTypeFork, toForkVersion.SeederID, toForkVersion.BuilderID).
		RunWith(tx).Exec()
	if err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not insert new script chain version: %s", err))
	}
	// TODO pruner
	if err = tx.Commit(); err != nil {
		return CreateOrUpdateScriptChainResponse{}, common.TryToRollback(tx, fmt.Errorf("Could not commit: %s", err))
	}
	return CreateOrUpdateScriptChainResponse{
		ID:        scriptChainID,
		SeederID:  toForkVersion.SeederID,
		BuilderID: toForkVersion.BuilderID,
		// TODO pruner
		Name:        toForkScriptChainName,
		Permissions: permissionsPrivate,
	}, nil
}

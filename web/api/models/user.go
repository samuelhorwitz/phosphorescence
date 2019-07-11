package models

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/satori/go.uuid"
)

func mapSpotifyIDToOurID(tx *sql.Tx, spotifyUserID string) (uuid.UUID, error) {
	var userID uuid.UUID
	err := psql.Select("id").From("users_view").
		Where(sq.Eq{
			"spotify_id": spotifyUserID,
		}).RunWith(tx).QueryRow().Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			userID = uuid.NewV4()
			_, err = psql.Insert("users").Columns("id", "spotify_id").
				Values(userID, spotifyUserID).
				RunWith(tx).Exec()
			if err != nil {
				return uuid.Nil, fmt.Errorf("Could not insert new user: %s", err)
			}
			return userID, nil
		}
		return uuid.Nil, fmt.Errorf("Could not query user: %s", err)
	}
	return userID, nil
}

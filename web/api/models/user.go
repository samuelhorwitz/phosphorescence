package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/satori/go.uuid"

	"github.com/samuelhorwitz/phosphorescence/api/session"
)

type User struct {
	SpotifyID         string `json:"spotifyId"`
	Name              string `json:"name"`
	Country           string `json:"country"`
	Product           string `json:"product"`
	Authenticated     bool   `json:"authenticated"`
	GoogleAnalyticsID string `json:"gaId,omitempty"`
}

func GetUser(sess *session.Session) User {
	googleAnalyticsIDSum := hmac.New(sha256.New, googleAnalyticsSecret)
	googleAnalyticsIDSum.Write([]byte(sess.SpotifyID))
	googleAnalyticsID := googleAnalyticsIDSum.Sum(nil)
	return User{
		SpotifyID:         sess.SpotifyID,
		Name:              sess.SpotifyName,
		Country:           sess.SpotifyCountry,
		Product:           sess.SpotifyProduct,
		Authenticated:     sess.Authenticated,
		GoogleAnalyticsID: hex.EncodeToString(googleAnalyticsID),
	}
}

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

package models

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
	"time"
)

type nullUUID struct {
	uuid.NullUUID
}

type nullString struct {
	sql.NullString
}

type nullTime struct {
	pq.NullTime
}

func (u nullUUID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(u.UUID)
}

func (u *nullUUID) UnmarshalJSON(b []byte) error {
	var nonNullUUID uuid.UUID
	err := json.Unmarshal(b, &nonNullUUID)
	if err != nil {
		return err
	}
	u.UUID = nonNullUUID
	u.Valid = true
	return nil
}

func (s nullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (s *nullString) UnmarshalJSON(b []byte) error {
	var nonNullString string
	err := json.Unmarshal(b, &nonNullString)
	if err != nil {
		return err
	}
	s.String = nonNullString
	s.Valid = true
	return nil
}

func (t nullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}

func (t *nullTime) UnmarshalJSON(b []byte) error {
	var nonNullTime time.Time
	err := json.Unmarshal(b, &nonNullTime)
	if err != nil {
		return err
	}
	t.Time = nonNullTime
	t.Valid = true
	return nil
}

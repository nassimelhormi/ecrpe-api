package models

import (
	"time"
)

// Video struct
type Video struct {
	ID        int       `json:"id" db:"id"`
	UUID      string    `json:"uuid" db:"uuid"`
	Path      string    `json:"path" db:"path"`
	Duration  string    `json:"duration" db:"duration"`
	IsEncoded bool      `json:"is_encoded" db:"is_encoded"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	SessionID int       `json:"session_id" db:"session_id"`
}

package models

import (
	"time"
)

// Video struct
type Video struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Path      string    `json:"path"`
	Duration  string    `json:"duration"`
	IsEncoded bool      `json:"is_encoded"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

package models

import (
	"time"
)

// ClassPaper struct
type ClassPaper struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Path      string    `json:"path" db:"path"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	SessionID int       `json:"session_id" db:"session_id"`
}

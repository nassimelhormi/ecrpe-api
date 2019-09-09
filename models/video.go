package models

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Authorized bool
type Authorized bool

// MarshalGQL func
func (a Authorized) MarshalGQL(w io.Writer) {
	if a {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

// UnmarshalGQL func
func (a *Authorized) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		*a = strings.ToLower(v) == "true"
		return nil
	case bool:
		*a = Authorized(v)
		return nil
	default:
		return fmt.Errorf("%T is not a bool", v)
	}
}

// Video struct
type Video struct {
	ID        int       `json:"id" db:"id"`
	UUID      string    `json:"uuid" db:"uuid"`
	Path      string    `json:"path" db:"path"`
	Duration  string    `json:"duration" db:"duration"`
	IsEncoded bool      `json:"is_encoded" db:"is_encoded"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

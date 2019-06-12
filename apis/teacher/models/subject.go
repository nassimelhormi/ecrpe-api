package models

import (
	"github.com/francoispqt/gojay"
	"github.com/juleur/dash-exp/apis/utils"
)

// Subject struct
type Subject struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	Active    bool   `db:"active"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

// CreateSubject func
func (db DB) CreateSubject(s Subject) error {
	query := "INSERT INTO subjects(name, active, created_at) VALUES (?,?,?)"
	_, err := db.Queryx(query, s.Name, s.Active, utils.Datetime())
	return err
}

// GetSubject func
func (db DB) GetSubject(subjectID int64) (Subject, error) {
	s := Subject{}
	if err := db.Get(&s, "SELECT id, name FROM subjects WHERE id = ?", subjectID); err != nil {
		return Subject{}, err
	}
	return s, nil
}

// UnmarshalJSONObject implement gojay.UnmarshalerJSONObject
func (s *Subject) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "id":
		return dec.Int64(&s.ID)
	case "name":
		return dec.String(&s.Name)
	case "active":
		return dec.Bool(&s.Active)
	case "created_at":
		return dec.String(&s.CreatedAt)
	case "updated_at":
		return dec.String(&s.UpdatedAt)
	}
	return nil
}

// NKeys func
func (s *Subject) NKeys() int {
	return 5
}

// MarshalJSONObject implement MarshalerJSONObject
func (s *Subject) MarshalJSONObject(enc *gojay.Encoder) {
	enc.Int64KeyOmitEmpty("id", s.ID)
	enc.StringKeyOmitEmpty("name", s.Name)
	enc.BoolKeyOmitEmpty("active", s.Active)
	enc.StringKeyOmitEmpty("created_at", s.CreatedAt)
	enc.StringKeyOmitEmpty("updated_at", s.UpdatedAt)
}

// IsNil func
func (s *Subject) IsNil() bool {
	return s == nil
}

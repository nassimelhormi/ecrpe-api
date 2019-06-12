package models

import (
	"github.com/francoispqt/gojay"
	"github.com/juleur/dash-exp/apis/utils"
	"github.com/sirupsen/logrus"
)

// RefresherCourse struct
type RefresherCourse struct {
	ID         int64  `db:"id"`
	Year       string `db:"year"`
	IsFinished bool   `db:"is_finished"`
	CreatedAt  string `db:"created_at"`
	UpdatedAt  string `db:"updated_at"`
	SubjectID  int64  `db:"subject_id"`
}

// CreateRefresherCourse func
func (db DB) CreateRefresherCourse(rc RefresherCourse) error {
	firstQ := "INSERT INTO refresher_courses(year, created_at) VALUES (?, ?)"
	secondQ := "INSERT INTO subjects_refresher_courses(subject_id, refresher_course_id) VALUES (?,?)"
	res, err := db.Exec(firstQ, &rc.Year, utils.Datetime())
	if err != nil {
		return err
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	res, err = db.Exec(secondQ, rc.SubjectID, lastInsertID)
	return err
}

// GetRefresherCourse func
func (db DB) GetRefresherCourse(rcID int64) (RefresherCourse, error) {
	rc := RefresherCourse{}
	if err := db.Select(&rc, "SELECT id, year FROM refresher_courses WHERE id = ?", rcID); err != nil {
		logrus.Errorln(err)
		return RefresherCourse{}, err
	}
	return rc, nil
}

// UnmarshalJSONObject implement gojay.UnmarshalerJSONObject
func (rc *RefresherCourse) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "id":
		return dec.Int64(&rc.ID)
	case "year":
		return dec.String(&rc.Year)
	case "is_finished":
		return dec.Bool(&rc.IsFinished)
	case "created_at":
		return dec.String(&rc.CreatedAt)
	case "updated_at":
		return dec.String(&rc.UpdatedAt)
	case "subject_id":
		return dec.Int64(&rc.SubjectID)
	}
	return nil
}

// NKeys func
func (rc *RefresherCourse) NKeys() int {
	return 6
}

// MarshalJSONObject implement MarshalerJSONObject
func (rc *RefresherCourse) MarshalJSONObject(enc *gojay.Encoder) {
	enc.Int64KeyOmitEmpty("id", rc.ID)
	enc.StringKeyOmitEmpty("year", rc.Year)
	enc.BoolKey("is_finished", rc.IsFinished)
	enc.StringKeyOmitEmpty("created_at", rc.CreatedAt)
	enc.StringKeyOmitEmpty("updated_at", rc.UpdatedAt)
	enc.Int64KeyOmitEmpty("subject_id", rc.SubjectID)
}

// IsNil func
func (rc *RefresherCourse) IsNil() bool {
	return rc == nil
}

package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/juleur/dash-exp/apis/utils"
	"github.com/sirupsen/logrus"

	"github.com/francoispqt/gojay"
)

// Video struct
type Video struct {
	ID        int64  `db:"id"`
	UUID      string `db:"uuid"`
	Path      string `db:"path"`
	Duration  string `db:"duration"` // in minutes
	IsEncoded bool   `db:"is_encoded"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

// PublisherVideo updates video from a channel of videos
func PublisherVideo(tm *TeacherManager) {
	for v := range tm.PublishVideoCh {
		err := tm.DB.UpdateVideo(v)
		if err != nil {
			tm.Log.WithFields(logrus.Fields{"value": v}).Error(err)
		}
	}
}

// UpdateVideo updates video
func (db DB) UpdateVideo(v Video) error {
	query := "UPDATE videos SET path=?, duration=?, is_encoded=?, updated_at=? WHERE id = ?"
	_, err := db.Queryx(query, v.Path, v.Duration, v.IsEncoded, utils.Datetime(), v.ID)
	return err
}

// AddDuration converts seconds to minutes before affecting new value
func (v *Video) AddDuration(secondsDur []byte) {
	r := strings.Split(string(secondsDur), ".")
	s, err := strconv.Atoi(string(r[0]))
	if err != nil {
		v.Duration = "00:00"
		return
	}
	var format string
	if s < 59 {
		if s < 9 {
			format = fmt.Sprintf("00:0%d", s)
		} else {
			format = fmt.Sprintf("00:%d", s)
		}
	} else {
		mins := s / 60
		seconds := s % 60
		if seconds < 9 {
			format = fmt.Sprintf("%d:0%d", mins, seconds)
		} else {
			format = fmt.Sprintf("%d:%d", mins, seconds)
		}
	}
	v.Duration = format
}

// FileName func
func (v Video) FileName() string {
	return v.Path + "/" + v.UUID
}

// UnmarshalJSONObject implement gojay.UnmarshalerJSONObject
func (v *Video) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "id":
		return dec.Int64(&v.ID)
	case "uuid":
		return dec.String(&v.UUID)
	case "path":
		return dec.String(&v.Path)
	case "duration":
		return dec.String(&v.Duration)
	case "is_encoded":
		return dec.Bool(&v.IsEncoded)
	case "created_at":
		return dec.String(&v.CreatedAt)
	case "updated_at":
		return dec.String(&v.UpdatedAt)
	}
	return nil
}

// NKeys func
func (v *Video) NKeys() int {
	return 7
}

// MarshalJSONObject implement MarshalerJSONObject
func (v *Video) MarshalJSONObject(enc *gojay.Encoder) {
	enc.Int64KeyOmitEmpty("id", v.ID)
	enc.StringKeyOmitEmpty("uuid", v.UUID)
	enc.StringKeyOmitEmpty("path", v.Path)
	enc.StringKeyOmitEmpty("duration", v.Duration)
	enc.BoolKeyOmitEmpty("is_encoded", v.IsEncoded)
	enc.StringKeyOmitEmpty("created_at", v.CreatedAt)
	enc.StringKeyOmitEmpty("updated_at", v.UpdatedAt)
}

// IsNil func
func (v *Video) IsNil() bool {
	return v == nil
}

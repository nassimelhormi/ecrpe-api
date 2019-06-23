package models

import (
	"mime/multipart"
	"strconv"
	"strings"
	"unicode"

	"github.com/juleur/dash-exp/apis/utils"
	"github.com/sirupsen/logrus"

	"github.com/francoispqt/gojay"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Session struct
type Session struct {
	ID                  int64  `db:"id"`
	Title               string `db:"title"`
	Type                string `db:"type"`
	Description         string `db:"description"`
	Part                string `db:"part"`
	RecordedOn          string `db:"recorded_on"`
	CreatedAt           string `db:"created_at"`
	UpdatedAt           string `db:"updated_at"`
	SubjectID           int64
	SubjectName         string
	RefresherCourseID   int64 `db:"refresher_course_id"`
	RefresherCourseYear string
}

//CreateSessionVideo func
func (db DB) CreateSessionVideo(s Session) (Video, error) {
	firstQ := "INSERT INTO sessions(title, type, description, part, recorded_on, created_at, refresher_course_id) VALUES (?,?,?,?,?,?,?)"
	secondQ := "INSERT INTO videos(uuid, created_at) VALUES (?,?)"
	thriceQ := "INSERT INTO sessions_videos(session_id, video_id) VALUES (?,?)"
	session, err := db.Exec(firstQ, s.Title, s.Type, s.Description, s.Part, s.RecordedOn, utils.Datetime(), s.RefresherCourseID)
	if err != nil {
		return Video{}, err
	}
	sessionLastID, err := session.LastInsertId()
	if err != nil {
		return Video{}, err
	}
	video, err := db.Exec(secondQ, utils.HexKeyGenerator(22), utils.Datetime())
	if err != nil {
		return Video{}, err
	}
	videoLastID, err := video.LastInsertId()
	if err != nil {
		return Video{}, err
	}
	if _, err := db.Queryx(thriceQ, sessionLastID, videoLastID); err != nil {
		return Video{}, err
	}
	v := Video{}
	err = db.Get(&v, "SELECT id, uuid FROM videos WHERE id = ?", videoLastID)
	if err != nil {
		return Video{}, err
	}
	return v, nil
}

// MultipartFormToSession parses *Form.Value, assigns values to each other
// then return a Session struct
func MultipartFormToSession(s *Session, form *multipart.Form) error {
	for i, v := range form.Value {
		switch i {
		case "title":
			s.Title = v[0]
		case "type":
			s.Type = v[0]
		case "description":
			s.Description = v[0]
		case "part":
			s.Part = v[0]
		case "recorded_on":
			s.RecordedOn = v[0]
		case "subject_name":
			s.SubjectName = v[0]
		case "refresher_course_year":
			s.RefresherCourseYear = v[0]
		case "refresher_course_id":
			n, err := strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return err
			}
			s.RefresherCourseID = n
		}
	}
	return nil
}

// SubjectFolderName func
func (s Session) SubjectFolderName() string {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing 	marks
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	sLower := strings.ToLower(s.SubjectName)
	r, _, err := transform.String(t, sLower)
	if err != nil {
		logrus.Errorf("Transformer Error: %s\n", err)
		return sLower
	}
	return r
}

// UnmarshalJSONObject func
func (s *Session) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "id":
		return dec.Int64(&s.ID)
	case "title":
		return dec.String(&s.Title)
	case "type":
		return dec.String(&s.Type)
	case "description":
		return dec.String(&s.Description)
	case "part":
		return dec.String(&s.Part)
	case "recorded_on":
		return dec.String(&s.RecordedOn)
	case "created_at":
		return dec.String(&s.CreatedAt)
	case "updated_at":
		return dec.String(&s.UpdatedAt)
	case "subject_id":
		return dec.Int64(&s.SubjectID)
	case "subject_name":
		return dec.String(&s.SubjectName)
	case "refresher_course_id":
		return dec.Int64(&s.RefresherCourseID)
	case "refresher_course_year":
		return dec.String(&s.RefresherCourseYear)
	}
	return nil
}

// NKeys func
func (s *Session) NKeys() int {
	return 12
}

// MarshalJSONObject implement MarshalerJSONObject
func (s *Session) MarshalJSONObject(enc *gojay.Encoder) {
	enc.Int64KeyOmitEmpty("id", s.ID)
	enc.StringKeyOmitEmpty("title", s.Title)
	enc.StringKeyOmitEmpty("type", s.Type)
	enc.StringKeyOmitEmpty("description", s.Description)
	enc.StringKeyOmitEmpty("part", s.Part)
	enc.StringKeyOmitEmpty("recorded_on", s.RecordedOn)
	enc.StringKeyOmitEmpty("created_at", s.CreatedAt)
	enc.StringKeyOmitEmpty("updated_at", s.UpdatedAt)
	enc.Int64KeyOmitEmpty("subject_id", s.SubjectID)
	enc.StringKeyOmitEmpty("subject_year", s.SubjectName)
	enc.Int64KeyOmitEmpty("refresher_course_id", s.RefresherCourseID)
	enc.StringKeyOmitEmpty("refresher_course_year", s.RefresherCourseYear)
}

// IsNil func
func (s *Session) IsNil() bool {
	return s == nil
}

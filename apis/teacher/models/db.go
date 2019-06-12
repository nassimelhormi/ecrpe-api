package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Datastore interface
type Datastore interface {
	CreateRefresherCourse(rc RefresherCourse) error
	GetRefresherCourse(rcID int64) (RefresherCourse, error)
	CreateSessionVideo(s Session) (Video, error)
	CreateSubject(s Subject) error
	GetSubject(subjectID int64) (Subject, error)
	UpdateVideo(v Video) error
}

// DB struct
type DB struct {
	*sqlx.DB
}

// NewDB opens database
func NewDB() *DB {
	db, err := sqlx.Open("mysql", "chermak:pwd@tcp(127.0.0.1:7359)/ecrpe")
	if err != nil {
		logrus.Fatalln(err)
	}
	if err = db.Ping(); err != nil {
		logrus.Fatalln(err)
	}
	return &DB{db}
}

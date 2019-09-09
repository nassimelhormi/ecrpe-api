package models

import (
	"time"
)

// Session struct
type Session struct {
	ID                int       `json:"id" db:"id"`
	Title             string    `json:"title" db:"title"`
	Description       string    `json:"description" db:"description"`
	RecordedOn        time.Time `json:"recorded_on" db:"recorded_on"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	RefresherCourseID int       `json:"refresher_course_id" db:"refresher_course_id"`
}

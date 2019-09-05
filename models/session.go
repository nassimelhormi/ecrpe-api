package models

import (
	"time"
)

// Session struct
type Session struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	RecordedOn        time.Time `json:"recorded_on"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	RefresherCourseID int       `json:"refresher_course_id"`
}

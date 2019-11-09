package models

import "time"

// RefresherCourse struct
type RefresherCourse struct {
	ID         int       `json:"id"`
	Year       string    `json:"year"`
	IsFinished bool      `json:"is_finished"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

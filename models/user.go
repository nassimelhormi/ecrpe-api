package models

import "time"

// User struct
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsTeacher bool      `json:"is_teacher"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updateda_t"`
}

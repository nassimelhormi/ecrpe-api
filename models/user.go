package models

import "time"

// User struct
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	IsTeacher    bool      `json:"is_teacher" db:"is_teacher"`
	EncryptedPWD string    `db:"encrypted_pwd"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updateda_t" db:"updated_at"`
}

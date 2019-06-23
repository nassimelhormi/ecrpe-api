package models

import "time"

// User struct
type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phoneNumber"`
	Email       string    `json:"email"`
	CurrentRank int       `json:"currentRank"`
	IsTeacher   bool      `json:"isTeacher"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

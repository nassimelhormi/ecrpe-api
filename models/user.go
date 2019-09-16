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

type UserAuth struct {
	ID           int       `json:"id" db:"id"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	DeliveredAt  time.Time `json:"delivered_at" db:"delivered_at"`
	IsRevoked    bool      `json:"is_revoked" db:"is_revoked"`
	RevokedAt    time.Time `json:"revoked_at" db:"revoked_at"`
	OnLogin      bool      `json:"on_login" db:"on_login"`
	OnRefresh    bool      `json:"on_refresh" db:"on_refresh"`
	UserID       int       `json:"user_id" db:"user_id"`
	Username     string    `json:"username" db:"username"`
}

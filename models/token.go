package models

// Token struct
type Token struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

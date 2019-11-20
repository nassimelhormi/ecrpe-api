package utils

import "github.com/gbrlsnchs/jwt"

// CustomPayload struct
type CustomPayload struct {
	jwt.Payload
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
}

package models

import "time"

type RefreshToken struct {
	ID        int
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Tokens struct {
	Refresh string `json:"refresh_token"`
	JWT     string `json:"jwt_token"`
}

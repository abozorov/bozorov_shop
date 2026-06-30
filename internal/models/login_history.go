package models

import "time"

type LoginHistory struct {
	ID        int       `json:"-"`
	UserID    int       `json:"-"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

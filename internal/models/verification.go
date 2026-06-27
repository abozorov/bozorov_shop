package models

import "time"

type Verification struct {
	Code      int    `json:"otp_code"`
	Email     string `json:"email"`
	CreatedAt time.Time
}

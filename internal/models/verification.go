package models

type Verification struct {
	Code  int    `json:"otp_code"`
	Email string `json:"email"`
}

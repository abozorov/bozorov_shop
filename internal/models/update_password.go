package models

type UpdatePassword struct {
	UserID      int
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

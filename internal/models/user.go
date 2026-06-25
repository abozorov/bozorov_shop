package models

import (
	"strings"
	"time"
	"unicode/utf8"
)

const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Phone     string
	Password  string
	Role      string
	CreatedAt time.Time
	DeletedAt time.Time
}

func isEmail(email string) bool {
	hasDog := false
	k := 0
	for _, v := range email {
		if v == '@' {
			hasDog = true
			k++
		}
	}
	return hasDog && utf8.RuneCountInString(email) > 3 && k == 1
}

func isTajik(phone string) bool {
	phoneLen := utf8.RuneCountInString(phone)
	if phoneLen != 12 {
		return false
	}
	countryCode := "992"
	for i := 0; i < len(countryCode); i++ {
		if countryCode[i] != phone[i] {
			return false
		}
	}

	return true
}

func (u *User) Validate(create bool) bool {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	u.Phone = strings.TrimSpace(u.Phone)
	u.Password = strings.TrimSpace(u.Password)
	u.Role = strings.TrimSpace(u.Role)

	return (u.ID > 0 || create) &&
		u.Name != "" &&
		isEmail(u.Email) &&
		isTajik(u.Phone) &&
		(u.Password != "" || !create) &&
		u.Role != ""
}

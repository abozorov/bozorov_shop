package models

import (
	"strings"

	"github.com/abozorov/bozorov_shop/pkg/errs"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	if r.Name == "" || r.Email == "" || !isEmail(r.Email) || len(r.Password) < 8 {
		return errs.ErrBadRequestBody
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	if r.Email == "" || len(r.Password) < 8 {
		return errs.ErrBadRequestBody
	}
	return nil
}

package password

import (
	"errors"
	"fmt"

	"github.com/abozorov/bozorov_shop/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("passworrd.Hash: %w", err)
	}

	return string(hash), nil
}

func Compare(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("passworrd.Compare: %w, %w", errs.ErrIncorrectPassword, err)
		}
		return fmt.Errorf("passworrd.Compare: %w", err)
	}
	return nil
}

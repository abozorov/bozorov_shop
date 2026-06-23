package models

import (
	"strings"
	"time"
)

type Order struct {
	ID        int
	UserID    int
	Product   string
	Price     float64
	Status    string
	CreatedAt time.Timer
}

func (o *Order) Validate(create bool) bool {
	o.Product = strings.TrimSpace(o.Product)
	o.Status = strings.TrimSpace(o.Status)

	validStatuses := map[string]bool{
		"new":         true,
		"in_progress": true,
		"done":        true,
		"cancled":     true,
	}

	return (o.ID > 0 || create) &&
		o.UserID > 0 &&
		o.Product != "" &&
		o.Price > 0 &&
		validStatuses[o.Status]
}

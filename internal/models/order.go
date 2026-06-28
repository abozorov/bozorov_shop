package models

import (
	"strings"
	"time"
)

var (
	OrderStatusNew        = "new"
	OrderStatusInProgress = "in_progress"
	OrderStatusDone       = "done"
	OrderStatusCancle     = "cancled"
)

type Order struct {
	ID        int
	UserID    int
	Product   string
	Price     float64
	Status    string
	CreatedAt time.Time
}

func (o *Order) Validate(create bool) bool {
	o.Product = strings.TrimSpace(o.Product)
	o.Status = strings.TrimSpace(o.Status)

	validStatuses := map[string]bool{
		OrderStatusNew:        true,
		OrderStatusInProgress: true,
		OrderStatusDone:       true,
		OrderStatusCancle:     true,
	}

	return (o.ID > 0 || create) &&
		o.Product != "" &&
		o.Price > 0 &&
		(validStatuses[o.Status] || create)
}

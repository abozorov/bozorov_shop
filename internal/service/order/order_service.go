package orderservice

import (
	"context"
	"fmt"

	"github.com/abozorov/bozorov_shop/internal/models"
	orderrepo "github.com/abozorov/bozorov_shop/internal/repo/order"
	userrepo "github.com/abozorov/bozorov_shop/internal/repo/user"
	"github.com/abozorov/bozorov_shop/pkg/errs"
)

type OrderService struct {
	userR  *userrepo.UserRepo
	orderR *orderrepo.OrderRepo
}

func NewOrderService(userR *userrepo.UserRepo, orderR *orderrepo.OrderRepo) *OrderService {
	return &OrderService{
		userR:  userR,
		orderR: orderR,
	}
}

func (u *OrderService) Create(ctx context.Context, order models.Order) error {
	// validation
	if !order.Validate(true) {
		return fmt.Errorf("order_service.GetByID: %w", errs.ErrNotFound)
	}

	// creating
	err := u.orderR.Create(ctx, order)
	if err != nil {
		return fmt.Errorf("order_service.Create: %w", err)
	}

	return nil
}

func (u *OrderService) GetAll(ctx context.Context) ([]models.Order, error) {
	// get all orders
	allOrders, err := u.orderR.GetAll(ctx)
	if err != nil {
		return []models.Order{}, fmt.Errorf("order_service.GetAll: %w", err)
	}

	// get active orders
	activeOrders := make([]models.Order, 0, len(allOrders))
	for _, v := range allOrders {
		if v.DeletedAt.IsZero() {
			activeOrders = append(activeOrders, v)
		}
	}

	return activeOrders, nil
}

func (u *OrderService) GetByID(ctx context.Context, id int) (*models.Order, error) {
	// get all orders
	order, err := u.orderR.GetByID(ctx, id)
	if err != nil {
		return &models.Order{}, fmt.Errorf("order_service.GetByID: %w", err)
	}

	// get active orders
	if !order.DeletedAt.IsZero() {
		return &models.Order{}, fmt.Errorf("order_service.GetByID: %w", errs.ErrOrderNotFound)
	}
	return order, nil
}

func (u *OrderService) Update(ctx context.Context, order models.Order) error {
	// validation
	if !order.Validate(false) {
		return fmt.Errorf("order_service.Update: %w", errs.ErrBadRequestBody)
	}

	// updating
	err := u.orderR.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("order_service.Update: %w", err)
	}

	return nil
}

func (u *OrderService) DeleteOrder(ctx context.Context, id int) error {
	// delete order
	err := u.orderR.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order_service.DeleteByID: %w", err)
	}

	// delete order orders
	err = u.orderR.DeleteByOrderID(ctx, id)
	if err != nil {
		return fmt.Errorf("order_service.DeleteByID: %w", err)
	}

	return nil
}

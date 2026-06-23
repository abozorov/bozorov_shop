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
		return fmt.Errorf("order_service.Create: %w", errs.ErrBadRequestBody)
	}

	// check user from delete
	usr, err := u.userR.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("order_service.Create: %w", err)
	}
	if !usr.DeletedAt.IsZero() {
		return fmt.Errorf("order_service.Create: %w", errs.ErrUserNotFound)
	}

	// creating
	err = u.orderR.Create(ctx, order)
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
	return allOrders, nil
}

func (u *OrderService) GetByID(ctx context.Context, id int) (*models.Order, error) {
	// get all orders
	order, err := u.orderR.GetByID(ctx, id)
	if err != nil {
		return &models.Order{}, fmt.Errorf("order_service.GetByID: %w", err)
	}

	return order, nil
}

func (u *OrderService) Update(ctx context.Context, order models.Order) error {
	// validation
	if !order.Validate(false) {
		return fmt.Errorf("order_service.Update: %w", errs.ErrBadRequestBody)
	}

	// check user from delete
	usr, err := u.userR.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("order_service.Create: %w", err)
	}
	if !usr.DeletedAt.IsZero() {
		return fmt.Errorf("order_service.Create: %w", errs.ErrUserNotFound)
	}

	// updating
	err = u.orderR.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("order_service.Update: %w", err)
	}

	return nil
}

func (u *OrderService) CancleOrder(ctx context.Context, id int) error {
	// delete order
	err := u.orderR.CancleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order_service.CancleOrder: %w", err)
	}

	return nil
}

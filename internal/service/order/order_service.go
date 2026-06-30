package orderservice

import (
	"context"
	"fmt"

	"github.com/abozorov/bozorov_shop/internal/models"
	mycontext "github.com/abozorov/bozorov_shop/internal/my_context"
	"github.com/abozorov/bozorov_shop/internal/repo"
	"github.com/abozorov/bozorov_shop/pkg/errs"
)

type OrderService struct {
	userR  repo.UserRepo
	orderR repo.OrderRepo
}

func NewOrderService(userR repo.UserRepo, orderR repo.OrderRepo) *OrderService {
	return &OrderService{
		userR:  userR,
		orderR: orderR,
	}
}

func (o *OrderService) Create(ctx context.Context, order models.Order) error {
	// validation
	if !order.Validate(true) {
		return fmt.Errorf("order_service.Create: %w", errs.ErrBadRequestBody)
	}

	// check user from delete
	usr, err := o.userR.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("order_service.Create: %w", err)
	}
	if !usr.DeletedAt.IsZero() {
		return fmt.Errorf("order_service.Create: %w", errs.ErrUserNotFound)
	}

	// creating
	err = o.orderR.Create(ctx, order)
	if err != nil {
		return fmt.Errorf("order_service.Create: %w", err)
	}

	return nil
}

func (o *OrderService) GetAll(ctx context.Context) ([]models.Order, error) {
	// get all orders
	allOrders, err := o.orderR.GetAll(ctx)
	if err != nil {
		return []models.Order{}, fmt.Errorf("order_service.GetAll: %w", err)
	}

	// get active orders
	return allOrders, nil
}

func (o *OrderService) GetAllByUserID(ctx context.Context, userID int) ([]models.Order, error) {
	// get all orders
	allOrders, err := o.orderR.GetAllByUserID(ctx, userID)
	if err != nil {
		return []models.Order{}, fmt.Errorf("order_service.GetAllByUserID: %w", err)
	}

	// get active orders
	return allOrders, nil
}

func (o *OrderService) GetByID(ctx context.Context, id int) (*models.Order, error) {
	// get order
	order, err := o.orderR.GetByID(ctx, id)
	if err != nil {
		return &models.Order{}, fmt.Errorf("order_service.GetByID: %w", err)
	}

	// get userID
	userID, ok := ctx.Value(mycontext.UserIDKey).(int)
	if !ok {
		return &models.Order{}, fmt.Errorf("order_service.GetByID: %w", errs.ErrIncorrectLoginOrPassword)
	}
	if order.UserID != userID {
		return &models.Order{}, fmt.Errorf("order_service.GetByID: %w", errs.ErrBadRequest)
	}

	return order, nil
}

func (o *OrderService) Update(ctx context.Context, order models.Order) error {
	// validation
	if !order.Validate(false) {
		return fmt.Errorf("order_service.Update: %w", errs.ErrBadRequestBody)
	}

	// check user from delete
	usr, err := o.userR.GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("order_service.Update: %w", err)
	}
	if !usr.DeletedAt.IsZero() {
		return fmt.Errorf("order_service.Update: %w", errs.ErrUserNotFound)
	}

	// check user for order
	oldOrder, err := o.orderR.GetByID(ctx, order.ID)
	if err != nil {
		return fmt.Errorf("order_service.Update: %w", err)
	}
	if oldOrder.UserID != order.UserID ||
		(oldOrder.Status == models.OrderStatusNew && order.Status == models.OrderStatusCancle) {
		return fmt.Errorf("order_service.Update: %w", errs.ErrBadRequest)
	}

	// updating
	err = o.orderR.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("order_service.Update: %w", err)
	}

	return nil
}

func (o *OrderService) CancleOrder(ctx context.Context, id int) error {
	// get order
	order, err := o.orderR.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order_service.CancleOrder: %w", err)
	}

	// get userID
	userID, ok := ctx.Value(mycontext.UserIDKey).(int)
	if !ok {
		return fmt.Errorf("order_service.CancleOrder: %w", errs.ErrIncorrectLoginOrPassword)
	}

	// check user for orderr
	if order.UserID != userID || order.Status != models.OrderStatusNew {
		return fmt.Errorf("order_service.CancleOrder: %w", errs.ErrBadRequest)
	}

	// delete order
	err = o.orderR.CancleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order_service.CancleOrder: %w", err)
	}

	return nil
}

package userservice

import (
	"context"
	"fmt"

	"github.com/abozorov/bozorov_shop/internal/models"
	orderrepo "github.com/abozorov/bozorov_shop/internal/repo/order"
	userrepo "github.com/abozorov/bozorov_shop/internal/repo/user"
	"github.com/abozorov/bozorov_shop/pkg/errs"
)

type UserService struct {
	userR  *userrepo.UserRepo
	orderR *orderrepo.OrderRepo
}

func NewUserService(userR *userrepo.UserRepo, orderR *orderrepo.OrderRepo) *UserService {
	return &UserService{
		userR:  userR,
		orderR: orderR,
	}
}

func (u *UserService) Create(ctx context.Context, user models.User) error {
	// validation
	if !user.Validate(true) {
		return fmt.Errorf("user_service.GetByID: %w", errs.ErrUserNotFound)
	}

	// creating
	err := u.userR.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service.Create: %w", err)
	}

	return nil
}

func (u *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	// get all users
	allUsers, err := u.userR.GetAll(ctx)
	if err != nil {
		return []models.User{}, fmt.Errorf("user_service.GetAll: %w", err)
	}

	// get active users
	activeUsers := make([]models.User, 0, len(allUsers))
	for _, v := range allUsers {
		if v.DeletedAt.IsZero() {
			activeUsers = append(activeUsers, v)
		}
	}

	return activeUsers, nil
}

func (u *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	// get all users
	user, err := u.userR.GetByID(ctx, id)
	if err != nil {
		return &models.User{}, fmt.Errorf("user_service.GetByID: %w", err)
	}

	// get active users
	if !user.DeletedAt.IsZero() {
		return &models.User{}, fmt.Errorf("user_service.GetByID: %w", errs.ErrUserNotFound)
	}
	return user, nil
}

func (u *UserService) Update(ctx context.Context, user models.User) error {
	// validation
	if !user.Validate(false) {
		return fmt.Errorf("user_service.Update: %w", errs.ErrBadRequestBody)
	}

	// updating
	err := u.userR.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service.Update: %w", err)
	}

	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, id int) error {
	// delete user
	err := u.userR.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user_service.DeleteByID: %w", err)
	}

	// delete user orders
	err = u.orderR.DeleteByUserID(ctx, id)
	if err != nil {
		return fmt.Errorf("user_service.DeleteByID: %w", err)
	}

	return nil
}

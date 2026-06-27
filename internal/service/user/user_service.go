package userservice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/abozorov/bozorov_shop/internal/models"
	orderrepo "github.com/abozorov/bozorov_shop/internal/repo/order"
	userrepo "github.com/abozorov/bozorov_shop/internal/repo/user"
	"github.com/abozorov/bozorov_shop/pkg/errs"
	"github.com/abozorov/bozorov_shop/pkg/jwt"
	"github.com/abozorov/bozorov_shop/pkg/password"
)

type UserService struct {
	userR  *userrepo.UserRepo
	orderR *orderrepo.OrderRepo
	jwt    *jwt.JWTSecret
}

func NewUserService(userR *userrepo.UserRepo, orderR *orderrepo.OrderRepo, jwt *jwt.JWTSecret) *UserService {
	return &UserService{
		userR:  userR,
		orderR: orderR,
		jwt:    jwt,
	}
}

func (u *UserService) Register(ctx context.Context, request models.RegisterRequest) error {
	err := request.Validate()
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	exists, err := u.userR.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	if exists {
		return fmt.Errorf("user_service.Register: %w", errs.ErrUserAlreadyExists)
	}

	passwordHash, err := password.Hash(request.Password)
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: passwordHash,
		Role:     models.UserRole,
	}

	err = u.userR.Add(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	return nil
}

func (u *UserService) Login(ctx context.Context, request models.LoginRequest) (token string, err error) {
	err = request.Validate()
	if err != nil {
		return "", fmt.Errorf("user_service.Login: %w", err)
	}

	user, err := u.userR.GetByEmail(ctx, request.Email)
	if err != nil {
		return "", fmt.Errorf("user_service.Login: %w", err)
	}

	if !user.DeletedAt.IsZero() {
		return "", fmt.Errorf("user_service.Login: %w", errs.ErrUserNotFound)
	}

	err = password.Compare(user.Password, request.Password)
	if err != nil {
		return "", fmt.Errorf("user_service.Login: %w", err)
	}

	token, err = u.jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("user_service.Login: %w", err)
	}

	return token, nil
}

func (u *UserService) Create(ctx context.Context, user models.User) error {
	// validation
	if !user.Validate(true) {
		return fmt.Errorf("user_service.Create: %w", errs.ErrUserNotFound)
	}

	passwordHash, err := password.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = passwordHash

	// creating
	err = u.userR.Add(ctx, user)
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

	return allUsers, nil
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
	user.Name = strings.TrimSpace(user.Name)
	user.Phone = strings.TrimSpace(user.Phone)
	if user.Name == "" || user.Phone == "" {
		return fmt.Errorf("user_service.Update: %w", errs.ErrBadRequestBody)
	}

	// updating
	err := u.userR.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service.Update: %w", err)
	}

	return nil
}

func (u *UserService) UpdateUserRole(ctx context.Context, user models.User) error {
	// validation
	user.Role = strings.TrimSpace(user.Role)
	if user.Role == "" {
		return fmt.Errorf("user_service.UpdateUserRole: %w", errs.ErrBadRequestBody)
	}

	// updating
	err := u.userR.UpdateUserRole(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service.UpdateUserRole: %w", err)
	}

	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, id int) error {
	// delete user
	err := u.userR.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user_service.DeleteUser: %w", err)
	}

	// delete user orders
	err = u.orderR.DeleteByUserID(ctx, id)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("user_service.DeleteUser: %w", err)
	}

	return nil
}

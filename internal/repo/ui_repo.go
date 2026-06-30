package repo

import (
	"context"

	"github.com/abozorov/bozorov_shop/internal/models"
)

type (
	RefreshTokenRepo interface {
		Create(ctx context.Context, token models.RefreshToken) error
		GetByUserID(ctx context.Context, clientID int) (*models.RefreshToken, error)
		GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
		Update(ctx context.Context, token models.RefreshToken) error
		DeleteByUserID(ctx context.Context, userID int) error
		DeleteByToken(ctx context.Context, token string) error
		ExistByUserID(ctx context.Context, userID int) (bool, error)
		ExistByToken(ctx context.Context, token string) (bool, error)
	}

	OrderRepo interface {
		Create(ctx context.Context, order models.Order) error
		GetAll(ctx context.Context) ([]models.Order, error)
		GetAllByUserID(ctx context.Context, userID int) ([]models.Order, error)
		GetByID(ctx context.Context, id int) (*models.Order, error)
		Update(ctx context.Context, order models.Order) error
		CancleByID(ctx context.Context, id int) error
		DeleteByUserID(ctx context.Context, userID int) error
	}

	UserRepo interface {
		Add(ctx context.Context, user models.User) error
		ExistsByEmail(ctx context.Context, email string) (bool, error)
		Create(ctx context.Context, user models.User) error
		GetAll(ctx context.Context) ([]models.User, error)
		GetByID(ctx context.Context, id int) (*models.User, error)
		GetByEmail(ctx context.Context, email string) (*models.User, error)
		Update(ctx context.Context, user models.User) error
		UpdatePassword(ctx context.Context, user models.User) error
		UpdateUserRole(ctx context.Context, user models.User) error
		DeleteByID(ctx context.Context, id int) error
	}
)

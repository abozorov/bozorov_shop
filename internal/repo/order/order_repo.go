package orderrepo

import (
	"context"

	"github.com/abozorov/bozorov_shop/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (o *OrderRepo) Create(ctx context.Context, user models.Order) error {
	return nil

}

func (o *OrderRepo) GetAll(ctx context.Context) ([]models.Order, error) {
	return []models.Order{}, nil
}

func (o *OrderRepo) GetByID(ctx context.Context, id int) (*models.Order, error) {
	return &models.Order{}, nil

}

func (o *OrderRepo) Update(ctx context.Context, user models.Order) error {
	return nil

}

func (o *OrderRepo) DeleteByID(ctx context.Context, id int) error {
	return nil

}

func (o *OrderRepo) DeleteByUserID(ctx context.Context, userID int) error {
	return nil

}

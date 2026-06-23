package orderrepo

import (
	"context"
	"fmt"

	"github.com/abozorov/bozorov_shop/internal/models"
	"github.com/abozorov/bozorov_shop/pkg/errs"
	"github.com/jackc/pgx/v5/pgconn"
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

func execAnalysis(res pgconn.CommandTag, err error) error {
	if err != nil {
		return fmt.Errorf("order_repo.execAnalysis: %w", err)
	}
	if rows := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("order_repo.execAnalysis: %w", errs.ErrNotFound)
	}
	return nil
}

func (o *OrderRepo) Create(ctx context.Context, order models.Order) error {
	const query = `
		INSERT INTO orders(user_id, product, price) VALUES
		($1, $2, $3)
	`

	_, err := o.db.Exec(ctx, query,
		order.UserID,
		order.Product,
		order.Price,
	)
	if err != nil {
		return fmt.Errorf("order_repo.Create: %w", errs.PostgresToErrs(err))
	}
	return nil
}

func (o *OrderRepo) GetAll(ctx context.Context) ([]models.Order, error) {
	const query = `
		SELECT  id,
			user_id,
			product,
			price,
			status,
			created_at
		 FROM orders
	`
	rows, err := o.db.Query(ctx, query)
	if err != nil {
		return []models.Order{},
			fmt.Errorf("order_repo.GetAll: %w", errs.PostgresToErrs(err))
	}
	defer rows.Close()

	orders := make([]models.Order, 0, 100)
	order := models.Order{}

	for rows.Next() {
		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Product,
			&order.Price,
			&order.Status,
			&order.CreatedAt,
		)

		if err != nil {
			return []models.Order{},
				fmt.Errorf("order_repo.GetAll: %w", errs.PostgresToErrs(err))
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return orders, errs.ErrNotFound
	}
	return orders, nil
}

func (o *OrderRepo) GetByID(ctx context.Context, id int) (*models.Order, error) {
	const query = `
		SELECT  id,
			user_id,
			product,
			price,
			status,
			created_at
		 FROM orders
		 WHERE id = $1
	`
	row := o.db.QueryRow(ctx, query, id)
	order := models.Order{}

	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Product,
		&order.Price,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		return &models.Order{},
			fmt.Errorf("order_repo.GetByID: %w", errs.PostgresToErrs(err))
	}

	return &order, nil

}

func (o *OrderRepo) Update(ctx context.Context, order models.Order) error {
	const query = `
		UPDATE orders
		SET user_id = $2,
		product = $3,
		price = $4,
		status = $5
		WHERE id = $1
	`

	err := execAnalysis(o.db.Exec(ctx, query,
		order.ID,
		order.UserID,
		order.Product,
		order.Price,
		order.Status,
	))
	if err != nil {
		return fmt.Errorf("order_repo.Update: %w", errs.PostgresToErrs(err))
	}
	return nil

}

func (o *OrderRepo) CancleByID(ctx context.Context, id int) error {
	// update status with id
	const query = `
                UPDATE orders
                SET status = 'cancled'
                WHERE id=$1
	`

	err := execAnalysis(o.db.Exec(ctx, query, id))
	if err != nil {
		return fmt.Errorf("order_repo.DeleteByID: %w", errs.PostgresToErrs(err))
	}
	return nil
}

func (o *OrderRepo) DeleteByUserID(ctx context.Context, userID int) error {
	//
	const query = `
                UPDATE orders
                SET status = 'cancled'
                WHERE user_id=$1
	`

	err := execAnalysis(o.db.Exec(ctx, query, userID))
	if err != nil {
		return fmt.Errorf("order_repo.DeleteByUserID: %w", errs.PostgresToErrs(err))
	}
	return nil
}

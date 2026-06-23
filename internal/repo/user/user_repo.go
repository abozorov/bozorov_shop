package userrepo

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/abozorov/bozorov_shop/internal/models"
	"github.com/abozorov/bozorov_shop/pkg/errs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func execAnalysis(res pgconn.CommandTag, err error) error {
	if err != nil {
		return fmt.Errorf("user_repo.execAnalysis: %w", err)
	}
	if rows := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("user_repo.execAnalysis: %w", errs.ErrUserNotFound)
	}
	return nil
}

func (u *UserRepo) Create(ctx context.Context, user models.User) error {
	const query = `
		INSERT INTO users(name, email, phone, password_hash, role) VALUES
		($1, $2, $3, $4, $5)
	`

	passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(user.PasswordHash)))
	_, err := u.db.Exec(ctx, query,
		user.Name,
		user.Email,
		user.Phone,
		passwordHash,
		user.Role,
	)
	if err != nil {
		return fmt.Errorf("user_repo.Create: %w", errs.PostgresToErrs(err))
	}
	return nil
}

func (u *UserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	const query = `
		SELECT  id,
			name,
			email,
			phone,
			role,
			created_at,
			deleted_at
		 FROM users
	`
	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return []models.User{},
			fmt.Errorf("user_repo.GetAll: %w", errs.PostgresToErrs(err))
	}
	defer rows.Close()

	users := make([]models.User, 0, 100)
	user := models.User{}

	var deletedAt pgtype.Timestamptz
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Role,
			&user.CreatedAt,
			&deletedAt,
		)
		user.DeletedAt = deletedAt.Time
		if err != nil {
			return []models.User{},
				fmt.Errorf("user_repo.GetAll: %w", errs.PostgresToErrs(err))
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return users, errs.ErrNotFound
	}
	return users, nil
}

func (u *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	const query = `
		SELECT  id,
			name,
			email,
			phone,
			role,
			created_at,
			deleted_at
		 FROM users
		 WHERE id = $1
	`
	row := u.db.QueryRow(ctx, query, id)

	var (
		user      models.User
		deletedAt pgtype.Timestamptz
	)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&deletedAt,
	)
	if err != nil {
		return &models.User{},
			fmt.Errorf("user_repo.GetByID: %w", errs.PostgresToErrs(err))
	}
	user.DeletedAt = deletedAt.Time
	return &user, nil
}

func (u *UserRepo) Update(ctx context.Context, user models.User) error {
	const query = `
		UPDATE users
		SET name = $2,
		email = $3,
		phone = $4,
		role = $5
		WHERE id = $1 AND deleted_at IS null
	`

	err := execAnalysis(u.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Phone,
		user.Role,
	))

	if err != nil {
		return fmt.Errorf("user_repo.Update: %w", errs.PostgresToErrs(err))
	}
	return nil
}

func (u *UserRepo) DeleteByID(ctx context.Context, id int) error {
	// update with id
	const query = `
                UPDATE users
                SET deleted_at=current_timestamp
                WHERE id=$1 AND deleted_at IS null
	`
	err := execAnalysis(u.db.Exec(ctx, query, id))
	if err != nil {
		return fmt.Errorf("user_repo.DeleteByID: %w", errs.PostgresToErrs(err))
	}

	// write data
	return nil

}

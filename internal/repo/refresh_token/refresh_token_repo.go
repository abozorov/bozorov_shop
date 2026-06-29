package refreshtokenrepo

import (
	"context"
	"fmt"

	"github.com/abozorov/bozorov_shop/internal/models"
	"github.com/abozorov/bozorov_shop/pkg/errs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepo(db *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		db: db,
	}
}

func execAnalysis(res pgconn.CommandTag, err error) error {
	if err != nil {
		return fmt.Errorf("refresh_token_repo.execAnalysis: %w", err)
	}
	if rows := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("refresh_token_repo.execAnalysis: %w", errs.ErrNotFound)
	}
	return nil
}

func (r *RefreshTokenRepo) Create(ctx context.Context, token models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens(user_id, token_hash, expires_at, created_at) VALUES
		($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("refresh_token_repo.Create: %w", errs.PostgresToErrs(err))
	}
	return nil
}

func (r *RefreshTokenRepo) GetByUserID(ctx context.Context, clientID int) (*models.RefreshToken, error) {
	const query = `
		SELECT  id,
			user_id,
			token_hash,
			expires_at,
			created_at
		 FROM refresh_tokens
		 WHERE user_id = $1
	`
	row := r.db.QueryRow(ctx, query, clientID)
	token := models.RefreshToken{}

	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
	)
	if err != nil {
		return &models.RefreshToken{},
			fmt.Errorf("refresh_token_repo.GetByUserID: %w", errs.PostgresToErrs(err))
	}

	return &token, nil
}

func (r *RefreshTokenRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	const query = `
		SELECT  id,
			user_id,
			token_hash,
			expires_at,
			created_at
		 FROM refresh_tokens
		 WHERE token_hash = $1
	`
	row := r.db.QueryRow(ctx, query, tokenHash)
	token := models.RefreshToken{}

	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
	)
	if err != nil {
		return &models.RefreshToken{},
			fmt.Errorf("refresh_token_repo.GetByTokenHash: %w", errs.PostgresToErrs(err))
	}

	return &token, nil
}

func (r *RefreshTokenRepo) Update(ctx context.Context, token models.RefreshToken) error {
	const query = `
		UPDATE refresh_tokens
		SET token_hash = $2,
		expires_at = $3,
		created_at = $4
		WHERE user_id = $1
	`
	err := execAnalysis(
		r.db.Exec(ctx, query,
			token.UserID,
			token.TokenHash,
			token.ExpiresAt,
			token.CreatedAt,
		),
	)
	if err != nil {
		return fmt.Errorf("refresh_token_repo.Update: %w", errs.PostgresToErrs(err))
	}

	return nil
}

func (r *RefreshTokenRepo) DeleteByUserID(ctx context.Context, userID int) error {
	const query = `
		DELETE 
		FROM refresh_tokens
		WHERE user_id = $1
	`
	err := execAnalysis(r.db.Exec(ctx, query, userID))
	if err != nil {
		return fmt.Errorf("refresh_token_repo.Delete: %w", errs.PostgresToErrs(err))
	}

	return nil
}

func (r *RefreshTokenRepo) DeleteByToken(ctx context.Context, token string) error {
	const query = `
		DELETE 
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	err := execAnalysis(r.db.Exec(ctx, query, token))
	if err != nil {
		return fmt.Errorf("refresh_token_repo.Delete: %w", errs.PostgresToErrs(err))
	}

	return nil
}

func (r *RefreshTokenRepo) ExistByUserID(ctx context.Context, userID int) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM refresh_tokens
			WHERE user_id = $1
		);
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("user_repo.ExistByUserID: %w", errs.PostgresToErrs(err))
	}
	return exists, nil
}

func (r *RefreshTokenRepo) ExistByToken(ctx context.Context, token string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM refresh_tokens
			WHERE token_hash = $1
		);
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, token).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("user_repo.ExistByToken: %w", errs.PostgresToErrs(err))
	}
	return exists, nil
}

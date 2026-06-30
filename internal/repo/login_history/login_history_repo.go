package loginhistoryrepo

import (
	"context"
	"fmt"

	"github.com/abozorov/bozorov_shop/internal/models"
	"github.com/abozorov/bozorov_shop/pkg/errs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginHistoryRepo struct {
	db *pgxpool.Pool
}

func NewLoginHistoryRepo(db *pgxpool.Pool) *LoginHistoryRepo {
	return &LoginHistoryRepo{
		db: db,
	}
}

func execAnalysis(res pgconn.CommandTag, err error) error {
	if err != nil {
		return fmt.Errorf("login_history_repo.execAnalysis: %w", err)
	}
	if rows := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("login_history_repo.execAnalysis: %w", errs.ErrNotFound)
	}
	return nil
}

func (l *LoginHistoryRepo) Create(ctx context.Context, track models.LoginHistory) error {
	const query = `
		INSERT INTO login_history(user_id, ip, user_agent, created_at) VALUES
		($1, $2, $3, $4)
	`

	_, err := l.db.Exec(ctx, query,
		track.UserID,
		track.IP,
		track.UserAgent,
		track.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("login_history_repo.Create: %w", errs.PostgresToErrs(err))
	}
	return nil
}

func (l *LoginHistoryRepo) GetAll(ctx context.Context) ([]models.LoginHistory, error) {
	const query = `
		SELECT  id,
			user_id,
			ip,
			user_agent,
			created_at
		 FROM login_history
	`
	rows, err := l.db.Query(ctx, query)
	if err != nil {
		return []models.LoginHistory{},
			fmt.Errorf("login_history_repo.GetAll: %w", errs.PostgresToErrs(err))
	}
	defer rows.Close()

	history := make([]models.LoginHistory, 0, 100)
	track := models.LoginHistory{}

	for rows.Next() {
		err = rows.Scan(
			&track.ID,
			&track.UserID,
			&track.IP,
			&track.UserAgent,
			&track.CreatedAt,
		)

		if err != nil {
			return []models.LoginHistory{},
				fmt.Errorf("login_history_repo.GetAll: %w", errs.PostgresToErrs(err))
		}
		history = append(history, track)
	}

	return history, nil
}

func (l *LoginHistoryRepo) GetAllByUserID(ctx context.Context, userID int) ([]models.LoginHistory, error) {
	const query = `
		SELECT  id,
			user_id,
			ip,
			user_agent,
			created_at
		 FROM login_history
		 WHERE user_id=$1;
	`
	rows, err := l.db.Query(ctx, query, userID)
	if err != nil {
		return []models.LoginHistory{},
			fmt.Errorf("login_history_repo.GetAllByUserID: %w", errs.PostgresToErrs(err))
	}
	defer rows.Close()

	history := make([]models.LoginHistory, 0, 100)
	track := models.LoginHistory{}

	for rows.Next() {
		err = rows.Scan(
			&track.ID,
			&track.UserID,
			&track.IP,
			&track.UserAgent,
			&track.CreatedAt,
		)

		if err != nil {
			return []models.LoginHistory{},
				fmt.Errorf("login_history_repo.GetAllByUserID: %w", errs.PostgresToErrs(err))
		}
		history = append(history, track)
	}

	return history, nil
}

func (l *LoginHistoryRepo) GetByID(ctx context.Context, id int) (*models.LoginHistory, error) {
	const query = `
		SELECT  id,
			user_id,
			ip,
			user_agent,
			created_at
		 FROM login_history
		 WHERE id = $1
	`
	row := l.db.QueryRow(ctx, query, id)
	track := models.LoginHistory{}

	err := row.Scan(
		&track.ID,
		&track.UserID,
		&track.IP,
		&track.UserAgent,
		&track.CreatedAt,
	)
	if err != nil {
		return &models.LoginHistory{},
			fmt.Errorf("login_history_repo.GetByID: %w", errs.PostgresToErrs(err))
	}

	return &track, nil
}

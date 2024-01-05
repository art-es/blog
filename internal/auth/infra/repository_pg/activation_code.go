package repository_pg

import (
	"context"
	"database/sql"

	"github.com/art-es/blog/internal/common/repository/pg"
)

type activationCodeRepository struct {
	conn pg.Conn
}

func newActivationCodeRepository(conn pg.Conn) *activationCodeRepository {
	return &activationCodeRepository{conn: conn}
}

func (r *activationCodeRepository) Add(ctx context.Context, code string, userID int64) error {
	const query = `INSERT INTO activation_code (code, user_id) VALUES ($1, $2)`
	_, err := r.conn.ExecContext(ctx, query, code, userID)
	return err
}

func (r *activationCodeRepository) GetUserID(ctx context.Context, code string) (int64, error) {
	const query = `SELECT user_id FROM activation_code WHERE code=$1`
	var userID int64
	err := r.conn.QueryRowContext(ctx, query, code).
		Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return userID, err
}

func (r *activationCodeRepository) RemoveCodes(ctx context.Context, userID int64) error {
	const query = `DELETE FROM activation_code WHERE user_id=$1`
	_, err := r.conn.ExecContext(ctx, query, userID)
	return err
}

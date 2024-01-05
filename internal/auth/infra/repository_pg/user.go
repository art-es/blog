package repository_pg

import (
	"context"
	"database/sql"

	"github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/common/repository/pg"
)

type userRepository struct {
	conn pg.Conn
}

func newUserRepository(conn pg.Conn) *userRepository {
	return &userRepository{conn: conn}
}

func (r *userRepository) Activate(ctx context.Context, id int64) (bool, error) {
	const query = `UPDATE auth SET activate=TRUE WHERE id=$1`
	result, err := r.conn.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

func (r *userRepository) Get(ctx context.Context, id int64) (*domain.User, error) {
	const query = `SELECT id, name, password_hash FROM auth WHERE id=$1`
	user := &domain.User{}
	err := r.conn.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `SELECT id, name, password_hash FROM auth WHERE email=$1`
	user := &domain.User{}
	err := r.conn.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM auth WHERE email=$1)`
	var exists bool
	err := r.conn.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *userRepository) Exists(ctx context.Context, id int64) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM auth WHERE id=$1)`
	var exists bool
	err := r.conn.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	if user.ID == 0 {
		return r.insert(ctx, user)
	}
	return r.update(ctx, user)
}

func (r *userRepository) insert(ctx context.Context, user *domain.User) error {
	const query = `INSERT INTO auth (name, email, password_hash) 
		VALUES ($1, $2, $3) RETURNING id`
	return r.conn.QueryRowContext(ctx, query, user.Name, user.Email, user.PasswordHash).
		Scan(&user.ID)
}

func (r *userRepository) update(ctx context.Context, user *domain.User) error {
	const query = `UPDATE auth SET name=$2, email=$3, password_hash=$4 WHERE id=$1`
	_, err := r.conn.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash)
	return err
}

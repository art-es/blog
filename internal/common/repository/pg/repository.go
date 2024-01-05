package pg

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/art-es/blog/internal/common/repository"
)

// Conn - postgres database connection
type Conn interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Repository[B repository.TxBeginner[C], C repository.TxCommitter] struct {
	newTxComitter func(*Repository[B, C]) C
	// one of db and tx is not nil value
	db   *sql.DB
	tx   *sql.Tx
	conn Conn
}

func New[B repository.TxBeginner[C], C repository.TxCommitter](
	db *sql.DB,
	newTxComitter func(*Repository[B, C]) C,
) *Repository[B, C] {
	return &Repository[B, C]{
		newTxComitter: newTxComitter,
		db:            db,
		conn:          db,
	}
}

func (r *Repository[B, C]) BeginTx(ctx context.Context) (C, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		var c C
		return c, fmt.Errorf("[pg] tx starting error: %w", err)
	}

	txRepo := &Repository[B, C]{
		tx:   tx,
		conn: tx,
	}

	return r.newTxComitter(txRepo), nil
}

func (r *Repository[B, C]) Rollback() {
	if err := r.tx.Rollback(); err != nil {
		log.Printf("[ERROR] rollback error: %v\n", err)
	}
}

func (r *Repository[B, C]) Commit() error {
	return r.tx.Commit()
}

func (r *Repository[B, C]) Conn() Conn {
	return r.conn
}

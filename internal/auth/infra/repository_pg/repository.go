package repository_pg

import (
	"database/sql"

	"github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/common/repository/pg"
)

type Repository struct {
	*pg.Repository[domain.Repository, domain.TxCommitter]
}

func New(db *sql.DB) *Repository {
	return &Repository{
		Repository: pg.New[domain.Repository, domain.TxCommitter](db, newTxComitter),
	}
}

func newTxComitter(r *pg.Repository[domain.Repository, domain.TxCommitter]) domain.TxCommitter {
	return &Repository{
		Repository: r,
	}
}

func (r *Repository) User() domain.UserRepository {
	return newUserRepository(r.Conn())
}

func (r *Repository) ActivationCode() domain.ActivationCodeRepository {
	return newActivationCodeRepository(r.Conn())
}

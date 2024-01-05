//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock
package domain

import (
	"context"

	"github.com/art-es/blog/internal/common/repository"
)

type UserRepository interface {
	Activate(ctx context.Context, id int64) (bool, error)
	Get(ctx context.Context, id int64) (*User, error) // FIXME: dead code
	GetByEmail(ctx context.Context, email string) (*User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Exists(ctx context.Context, id int64) (bool, error)
	Save(ctx context.Context, user *User) error
}

type ActivationCodeRepository interface {
	Add(ctx context.Context, code string, userID int64) error
	GetUserID(ctx context.Context, code string) (int64, error)
	RemoveCodes(ctx context.Context, userID int64) error
}

type repositoryGetter interface {
	User() UserRepository
	ActivationCode() ActivationCodeRepository
}

type Repository interface {
	BeginTx(context.Context) (TxCommitter, error)
	// FIXME: generics not working with gomock
	// repository.TxBeginner[TxCommitter]
	repositoryGetter
}

type TxCommitter interface {
	repository.TxCommitter
	repositoryGetter
}

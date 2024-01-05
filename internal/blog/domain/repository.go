package domain

import (
	"context"

	"github.com/art-es/blog/internal/common/repository"
)

type ArticleRepository interface{}

type CategoryRepository interface{}

type repositoryGetter interface {
	Article() ArticleRepository
	Category() CategoryRepository
}

type Repository interface {
	BeginTx(context.Context) (TxCommitter, error)
	// repository.TxBeginner[TxCommitter]
	repositoryGetter
}

type TxCommitter interface {
	repository.TxCommitter
	repositoryGetter
}

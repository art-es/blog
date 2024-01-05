package repository

import "context"

type TxBeginner[T TxCommitter] interface {
	BeginTx(context.Context) (T, error)
}

type TxCommitter interface {
	Rollback()
	Commit() error
}

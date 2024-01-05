//go:generate mockgen -source=usecase_activate.go -destination=mock/usecase_activate.go -package=mock
package domain

import (
	"context"
	"fmt"

	"github.com/art-es/blog/internal/auth/dto"
)

type activator interface {
	Activate(ctx context.Context, code string, tx TxCommitter) error
}

type ActivateUsecase struct {
	repository Repository
	activator  activator
}

func NewActivateUsecase(
	repository Repository,
	activator activator,
) *ActivateUsecase {
	return &ActivateUsecase{
		repository: repository,
		activator:  activator,
	}
}

func (u *ActivateUsecase) Do(ctx context.Context, in *dto.ActivateIn) error {
	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("tx beginning error: %w", err)
	}

	if err = u.activator.Activate(ctx, in.Code, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("activation error: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx committing error: %w", err)
	}

	return nil
}

//go:generate mockgen -source=case_user_activate.go -destination=mock/case_user_activate.go -package=mock
package domain

import (
	"context"
	"fmt"

	"github.com/art-es/blog/internal/auth/dto"
)

type userActivator interface {
	Activate(ctx context.Context, activationCode string, tx TxCommitter) error
}

type UserActivateCase struct {
	repository Repository
	activator  userActivator
}

func NewUserActivateCase(
	repository Repository,
	activator userActivator,
) *UserActivateCase {
	return &UserActivateCase{
		repository: repository,
		activator:  activator,
	}
}

func (c *UserActivateCase) Use(ctx context.Context, in *dto.UserActivateIn) error {
	tx, err := c.repository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("tx beginning error: %w", err)
	}

	if err = c.activator.Activate(ctx, in.Code, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("activation error: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx committing error: %w", err)
	}

	return nil
}

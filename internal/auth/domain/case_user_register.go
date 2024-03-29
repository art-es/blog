//go:generate mockgen -source=case_user_register.go -destination=mock/case_user_register.go -package=mock
package domain

import (
	"context"
	"fmt"

	"github.com/art-es/blog/internal/auth/dto"
)

type passwordHashGenerator interface {
	Generate(password string) (string, error)
}

type activationCodeSender interface {
	SendCode(ctx context.Context, user *User, tx TxCommitter) error
}

type UserRegisterCase struct {
	repository            Repository
	passwordHashGenerator passwordHashGenerator
	activationCodeSender  activationCodeSender
}

func NewUserRegisterCase(
	repository Repository,
	passwordHashGenerator passwordHashGenerator,
	activationCodeSender activationCodeSender,
) *UserRegisterCase {
	return &UserRegisterCase{
		repository:            repository,
		passwordHashGenerator: passwordHashGenerator,
		activationCodeSender:  activationCodeSender,
	}
}

func (u *UserRegisterCase) Use(ctx context.Context, in *dto.UserRegisterIn) error {
	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("tx beginning error: %w", err)
	}

	if err = u.useInTx(ctx, in, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx committing error: %w", err)
	}

	return nil
}

func (u *UserRegisterCase) useInTx(ctx context.Context, in *dto.UserRegisterIn, tx TxCommitter) error {
	if err := validateEmail(ctx, in.Email, tx.User()); err != nil {
		return err
	}

	passwordHash, err := u.passwordHashGenerator.Generate(in.Password)
	if err != nil {
		return fmt.Errorf("password hash generation error: %w", err)
	}

	user := User{
		Name:         in.Name,
		Email:        in.Email,
		PasswordHash: passwordHash,
	}

	if err = tx.User().Save(ctx, &user); err != nil {
		return fmt.Errorf("auth saving error: %w", err)
	}

	if err = u.activationCodeSender.SendCode(ctx, &user, tx); err != nil {
		return fmt.Errorf("auth activation code sending error: %w", err)
	}

	return nil
}

func validateEmail(ctx context.Context, email string, userRepository UserRepository) error {
	exists, err := userRepository.EmailExists(ctx, email)
	if err != nil {
		return fmt.Errorf("email existence checking in repository error: %w", err)
	}
	if exists {
		return dto.ErrEmailIsBusy
	}
	return nil
}

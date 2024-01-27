package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/auth/domain/mock"
	"github.com/art-es/blog/internal/auth/dto"
)

func TestRegisterUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		repository            = mock.NewMockRepository(ctrl)
		passwordHashGenerator = mock.NewMockpasswordHashGenerator(ctrl)
		activationCodeSender  = mock.NewMockactivationCodeSender(ctrl)
	)

	var (
		ctx          = context.Background()
		name         = "dummyName"
		email        = "dummyEmail@example.com"
		password     = "dummyPassword1!"
		passwordHash = "dummyHashedPassword%$1"
		userID       = int64(1)
		in           = &dto.UserRegisterIn{
			Name:     name,
			Email:    email,
			Password: password,
		}
		noError = ""
	)

	for _, tt := range []struct {
		name   string
		setup  func()
		expErr string
	}{
		{
			name: "happy path",
			setup: func() {
				tx := mock.NewMockTxCommitter(ctrl)

				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(tx, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						r.EXPECT().
							EmailExists(gomock.Eq(ctx), gomock.Eq(email)).
							Return(false, nil)
						return r
					})

				passwordHashGenerator.EXPECT().
					Generate(gomock.Eq(password)).
					Return(passwordHash, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						user := &domain.User{
							ID:           0,
							Name:         name,
							Email:        email,
							PasswordHash: passwordHash,
						}
						r.EXPECT().
							Save(gomock.Eq(ctx), gomock.Eq(user)).
							DoAndReturn(func(ctx context.Context, user *domain.User) error {
								user.ID = userID
								return nil
							})
						return r
					})

				activationCodeSender.EXPECT().
					SendCode(
						gomock.Eq(ctx),
						gomock.Eq(&domain.User{
							ID:           userID,
							Name:         name,
							Email:        email,
							PasswordHash: passwordHash,
						}),
						gomock.Eq(tx),
					).
					Return(nil)

				tx.EXPECT().
					Commit().
					Return(nil)
			},
			expErr: noError,
		},
		{
			name: "error on beginning tx",
			setup: func() {
				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(nil, errors.New("dummy error"))
			},
			expErr: "tx beginning error: dummy error",
		},
		{
			name: "error on checking email existence",
			setup: func() {
				tx := mock.NewMockTxCommitter(ctrl)

				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(tx, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						r.EXPECT().
							EmailExists(gomock.Eq(ctx), gomock.Eq(email)).
							Return(false, errors.New("dummy error"))
						return r
					})

				tx.EXPECT().Rollback()
			},
			expErr: "email existence checking in repository error: dummy error",
		},
		{
			name: "email already exists",
			setup: func() {
				tx := mock.NewMockTxCommitter(ctrl)

				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(tx, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						r.EXPECT().
							EmailExists(gomock.Eq(ctx), gomock.Eq(email)).
							Return(true, nil)
						return r
					})

				tx.EXPECT().Rollback()
			},
			expErr: dto.ErrEmailIsBusy.Error(),
		},
		{
			name: "error on generating password hash",
			setup: func() {
				tx := mock.NewMockTxCommitter(ctrl)

				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(tx, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						r.EXPECT().
							EmailExists(gomock.Eq(ctx), gomock.Eq(email)).
							Return(false, nil)
						return r
					})

				passwordHashGenerator.EXPECT().
					Generate(gomock.Eq(password)).
					Return("", errors.New("dummy error"))

				tx.EXPECT().Rollback()
			},
			expErr: "password hash generation error: dummy error",
		},
		{
			name: "error on saving auth",
			setup: func() {
				tx := mock.NewMockTxCommitter(ctrl)

				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(tx, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						r.EXPECT().
							EmailExists(gomock.Eq(ctx), gomock.Eq(email)).
							Return(false, nil)
						return r
					})

				passwordHashGenerator.EXPECT().
					Generate(gomock.Eq(password)).
					Return(passwordHash, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						user := &domain.User{
							ID:           0,
							Name:         name,
							Email:        email,
							PasswordHash: passwordHash,
						}
						r.EXPECT().
							Save(gomock.Eq(ctx), gomock.Eq(user)).
							Return(errors.New("dummy error"))
						return r
					})

				tx.EXPECT().Rollback()
			},
			expErr: "auth saving error: dummy error",
		},
		{
			name: "error on sending activation code",
			setup: func() {
				tx := mock.NewMockTxCommitter(ctrl)

				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(tx, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						r.EXPECT().
							EmailExists(gomock.Eq(ctx), gomock.Eq(email)).
							Return(false, nil)
						return r
					})

				passwordHashGenerator.EXPECT().
					Generate(gomock.Eq(password)).
					Return(passwordHash, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						user := &domain.User{
							ID:           0,
							Name:         name,
							Email:        email,
							PasswordHash: passwordHash,
						}
						r.EXPECT().
							Save(gomock.Eq(ctx), gomock.Eq(user)).
							DoAndReturn(func(ctx context.Context, user *domain.User) error {
								user.ID = userID
								return nil
							})
						return r
					})

				activationCodeSender.EXPECT().
					SendCode(
						gomock.Eq(ctx),
						gomock.Eq(&domain.User{
							ID:           userID,
							Name:         name,
							Email:        email,
							PasswordHash: passwordHash,
						}),
						gomock.Eq(tx),
					).
					Return(errors.New("dummy error"))

				tx.EXPECT().Rollback()
			},
			expErr: "auth activation code sending error: dummy error",
		},
		{
			name: "error on tx committing",
			setup: func() {
				tx := mock.NewMockTxCommitter(ctrl)

				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(tx, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						r.EXPECT().
							EmailExists(gomock.Eq(ctx), gomock.Eq(email)).
							Return(false, nil)
						return r
					})

				passwordHashGenerator.EXPECT().
					Generate(gomock.Eq(password)).
					Return(passwordHash, nil)

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mock.NewMockUserRepository(ctrl)
						user := &domain.User{
							ID:           0,
							Name:         name,
							Email:        email,
							PasswordHash: passwordHash,
						}
						r.EXPECT().
							Save(gomock.Eq(ctx), gomock.Eq(user)).
							DoAndReturn(func(ctx context.Context, user *domain.User) error {
								user.ID = userID
								return nil
							})
						return r
					})

				activationCodeSender.EXPECT().
					SendCode(
						gomock.Eq(ctx),
						gomock.Eq(&domain.User{
							ID:           userID,
							Name:         name,
							Email:        email,
							PasswordHash: passwordHash,
						}),
						gomock.Eq(tx),
					).
					Return(nil)

				tx.EXPECT().
					Commit().
					Return(errors.New("dummy error"))
			},
			expErr: "tx committing error: dummy error",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			u := domain.NewUserRegisterCase(repository, passwordHashGenerator, activationCodeSender)
			err := u.Use(ctx, in)

			if tt.expErr == noError {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, tt.expErr)
		})
	}
}

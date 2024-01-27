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

func TestActivateUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		repository    = mock.NewMockRepository(ctrl)
		userActivator = mock.NewMockuserActivator(ctrl)
	)

	var (
		ctx     = context.Background()
		code    = "dummyCode"
		in      = &dto.UserActivateIn{Code: code}
		noError = ""
	)

	tests := []struct {
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

				userActivator.EXPECT().
					Activate(gomock.Eq(ctx), gomock.Eq(code), gomock.Eq(tx)).
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
			name: "error on activation",
			setup: func() {
				tx := mock.NewMockTxCommitter(ctrl)

				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(tx, nil)

				userActivator.EXPECT().
					Activate(gomock.Eq(ctx), gomock.Eq(code), gomock.Eq(tx)).
					Return(errors.New("dummy error"))

				tx.EXPECT().
					Rollback()
			},
			expErr: "activation error: dummy error",
		},
		{
			name: "error on tx committing",
			setup: func() {
				tx := mock.NewMockTxCommitter(ctrl)

				repository.EXPECT().
					BeginTx(gomock.Eq(ctx)).
					Return(tx, nil)

				userActivator.EXPECT().
					Activate(gomock.Eq(ctx), gomock.Eq(code), gomock.Eq(tx)).
					Return(nil)

				tx.EXPECT().
					Commit().
					Return(errors.New("dummy error"))
			},
			expErr: "tx committing error: dummy error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			u := domain.NewUserActivateCase(repository, userActivator)
			err := u.Use(ctx, in)

			if tt.expErr == noError {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expErr)
			}
		})
	}
}

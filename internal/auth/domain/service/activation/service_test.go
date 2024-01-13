package activation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/domain"
	mockdomain "github.com/art-es/blog/internal/auth/domain/mock"
	"github.com/art-es/blog/internal/auth/domain/service/activation"
	"github.com/art-es/blog/internal/auth/domain/service/activation/mock"
	"github.com/art-es/blog/internal/auth/dto"
	log_mock "github.com/art-es/blog/internal/common/log/mock"
	"github.com/art-es/blog/internal/common/testutil"
)

const noError = ""

func TestService_Activate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		logger = log_mock.NewMockLogger(ctrl)
		tx     = mockdomain.NewMockTxCommitter(ctrl)
	)

	var (
		ctx    = context.Background()
		userID = int64(1)
		code   = "dummyCode"
	)

	tests := []struct {
		name   string
		setup  func()
		expErr string
	}{
		{
			name: "happy path",
			setup: func() {
				tx.EXPECT().
					ActivationCode().
					DoAndReturn(func() domain.ActivationCodeRepository {
						r := mockdomain.NewMockActivationCodeRepository(ctrl)
						r.EXPECT().
							GetUserID(gomock.Eq(ctx), gomock.Eq(code)).
							Return(userID, nil)
						return r
					})

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mockdomain.NewMockUserRepository(ctrl)
						r.EXPECT().
							Activate(gomock.Eq(ctx), gomock.Eq(userID)).
							Return(true, nil)
						return r
					})

				tx.EXPECT().
					ActivationCode().
					DoAndReturn(func() domain.ActivationCodeRepository {
						r := mockdomain.NewMockActivationCodeRepository(ctrl)
						r.EXPECT().
							RemoveCodes(gomock.Eq(ctx), gomock.Eq(userID)).
							Return(nil)
						return r
					})
			},
			expErr: noError,
		},
		{
			name: "error on getting auth ID",
			setup: func() {
				tx.EXPECT().
					ActivationCode().
					DoAndReturn(func() domain.ActivationCodeRepository {
						r := mockdomain.NewMockActivationCodeRepository(ctrl)
						r.EXPECT().
							GetUserID(gomock.Eq(ctx), gomock.Eq(code)).
							Return(int64(0), errors.New("dummy error"))
						return r
					})
			},
			expErr: "auth ID getting from repository error: dummy error",
		},
		{
			name: "code not found",
			setup: func() {
				tx.EXPECT().
					ActivationCode().
					DoAndReturn(func() domain.ActivationCodeRepository {
						r := mockdomain.NewMockActivationCodeRepository(ctrl)
						r.EXPECT().
							GetUserID(gomock.Eq(ctx), gomock.Eq(code)).
							Return(int64(0), nil)
						return r
					})
			},
			expErr: dto.ErrActivationCodeNotFound.Error(),
		},
		{
			name: "error on activation",
			setup: func() {
				tx.EXPECT().
					ActivationCode().
					DoAndReturn(func() domain.ActivationCodeRepository {
						r := mockdomain.NewMockActivationCodeRepository(ctrl)
						r.EXPECT().
							GetUserID(gomock.Eq(ctx), gomock.Eq(code)).
							Return(userID, nil)
						return r
					})

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mockdomain.NewMockUserRepository(ctrl)
						r.EXPECT().
							Activate(gomock.Eq(ctx), gomock.Eq(userID)).
							Return(false, errors.New("dummy error"))
						return r
					})
			},
			expErr: "auth getting from repository error: dummy error",
		},
		{
			name: "auth not found",
			setup: func() {
				tx.EXPECT().
					ActivationCode().
					DoAndReturn(func() domain.ActivationCodeRepository {
						r := mockdomain.NewMockActivationCodeRepository(ctrl)
						r.EXPECT().
							GetUserID(gomock.Eq(ctx), gomock.Eq(code)).
							Return(userID, nil)
						return r
					})

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mockdomain.NewMockUserRepository(ctrl)
						r.EXPECT().
							Activate(gomock.Eq(ctx), gomock.Eq(userID)).
							Return(false, nil)
						return r
					})
			},
			expErr: dto.ErrUserNotFound.Error(),
		},
		{
			name: "error on removing codes",
			setup: func() {
				tx.EXPECT().
					ActivationCode().
					DoAndReturn(func() domain.ActivationCodeRepository {
						r := mockdomain.NewMockActivationCodeRepository(ctrl)
						r.EXPECT().
							GetUserID(gomock.Eq(ctx), gomock.Eq(code)).
							Return(userID, nil)
						return r
					})

				tx.EXPECT().
					User().
					DoAndReturn(func() domain.UserRepository {
						r := mockdomain.NewMockUserRepository(ctrl)
						r.EXPECT().
							Activate(gomock.Eq(ctx), gomock.Eq(userID)).
							Return(true, nil)
						return r
					})

				tx.EXPECT().
					ActivationCode().
					DoAndReturn(func() domain.ActivationCodeRepository {
						r := mockdomain.NewMockActivationCodeRepository(ctrl)
						r.EXPECT().
							RemoveCodes(gomock.Eq(ctx), gomock.Eq(userID)).
							Return(errors.New("dummy error"))
						return r
					})
			},
			expErr: "codes removing from repository error: dummy error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			s := activation.New(logger, nil)
			err := s.Activate(ctx, code, tx)

			if tt.expErr == noError {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, tt.expErr)
		})
	}
}

func TestService_SendCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		logger  = log_mock.NewMockLogger(ctrl)
		databus = mock.NewMockdatabus(ctrl)
		tx      = mockdomain.NewMockTxCommitter(ctrl)
	)

	var (
		ctx    = context.Background()
		userID = int64(1)
		email  = "dummyEmail@example.com"
		user   = &domain.User{ID: userID, Email: email}
	)

	tests := []struct {
		name   string
		user   *domain.User
		setup  func()
		expErr string
	}{
		{
			name: "happy path",
			user: user,
			setup: func() {
				repo := mockdomain.NewMockActivationCodeRepository(ctrl)
				repo.EXPECT().
					Add(gomock.Eq(ctx), testutil.IsUUID(), gomock.Eq(userID)).
					DoAndReturn(func(_ context.Context, code string, _ int64) error {
						msg := &dto.UserActivationEmailMessage{
							Email: email,
							Code:  code,
						}

						databus.EXPECT().
							ProduceActivationEmail(gomock.Eq(ctx), gomock.Eq(msg)).
							Return(nil)

						return nil
					})

				tx.EXPECT().ActivationCode().Return(repo)
			},
			expErr: noError,
		},
		{
			name: "error on adding new code",
			user: user,
			setup: func() {
				repo := mockdomain.NewMockActivationCodeRepository(ctrl)
				repo.EXPECT().
					Add(gomock.Eq(ctx), testutil.IsUUID(), gomock.Eq(userID)).
					Return(errors.New("dummy error"))

				tx.EXPECT().ActivationCode().Return(repo)
			},
			expErr: "activation code adding to repository error: dummy error",
		},
		{
			name: "error on producing email message",
			user: user,
			setup: func() {
				repo := mockdomain.NewMockActivationCodeRepository(ctrl)
				repo.EXPECT().
					Add(gomock.Eq(ctx), testutil.IsUUID(), gomock.Eq(userID)).
					DoAndReturn(func(_ context.Context, code string, _ int64) error {
						msg := &dto.UserActivationEmailMessage{
							Email: email,
							Code:  code,
						}

						databus.EXPECT().
							ProduceActivationEmail(gomock.Eq(ctx), gomock.Eq(msg)).
							Return(errors.New("dummy error"))

						return nil
					})

				tx.EXPECT().ActivationCode().Return(repo)
			},
			expErr: noError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			s := activation.New(logger, databus)
			err := s.SendCode(ctx, tt.user, tx)

			if tt.expErr == noError {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, tt.expErr)
		})
	}
}

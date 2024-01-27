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

func TestAuthenticateUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		userRepository    = mock.NewMockUserRepository(ctrl)
		passwordValidator = mock.NewMockpasswordValidator(ctrl)
		accessTokenIssuer = mock.NewMockaccessTokenIssuer(ctrl)
	)

	var (
		ctx      = context.Background()
		email    = "dummyEmail@example.com"
		password = "dummyPassword!%"
		user     = &domain.User{
			ID:           1,
			Email:        email,
			PasswordHash: "dummyPasswordHash",
		}
		accessTokenObject = &domain.AccessTokenObject{}
		accessToken       = "dummyAccessToken"
		in                = &dto.UserAuthenticateIn{
			Email:    email,
			Password: password,
		}
		noError = ""
	)

	tests := []struct {
		name   string
		setup  func()
		expOut *dto.UserAuthenticateOut
		expErr string
	}{
		{
			name: "happy path",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(gomock.Eq(ctx), gomock.Eq(email)).
					Return(user, nil)

				passwordValidator.EXPECT().
					Validate(gomock.Eq(password), gomock.Eq(user.PasswordHash)).
					Return(nil)

				accessTokenIssuer.EXPECT().
					NewObject(gomock.Eq(user.ID)).
					Return(accessTokenObject)

				accessTokenIssuer.EXPECT().
					Sign(gomock.Eq(accessTokenObject)).
					Return(accessToken, nil)
			},
			expOut: &dto.UserAuthenticateOut{AccessToken: accessToken},
			expErr: noError,
		},
		{
			name: "error on getting auth",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(gomock.Eq(ctx), gomock.Eq(email)).
					Return(nil, errors.New("dummy error"))
			},
			expErr: "auth getting by email error: dummy error",
		},
		{
			name: "auth not found",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(gomock.Eq(ctx), gomock.Eq(email)).
					Return(nil, nil)
			},
			expErr: dto.ErrUserNotFound.Error(),
		},
		{
			name: "error on validating password",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(gomock.Eq(ctx), gomock.Eq(email)).
					Return(user, nil)

				passwordValidator.EXPECT().
					Validate(gomock.Eq(password), gomock.Eq(user.PasswordHash)).
					Return(errors.New("dummy error"))
			},
			expErr: "password validate by hash error: dummy error",
		},
		{
			name: "error on creating access token",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(gomock.Eq(ctx), gomock.Eq(email)).
					Return(user, nil)

				passwordValidator.EXPECT().
					Validate(gomock.Eq(password), gomock.Eq(user.PasswordHash)).
					Return(nil)

				accessTokenIssuer.EXPECT().
					NewObject(gomock.Eq(user.ID)).
					Return(accessTokenObject)

				accessTokenIssuer.EXPECT().
					Sign(gomock.Eq(accessTokenObject)).
					Return("", errors.New("dummy error"))
			},
			expErr: "access token creation error: dummy error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			u := domain.NewUserAuthenticateCase(userRepository, passwordValidator, accessTokenIssuer)
			out, err := u.Use(ctx, in)

			assert.Equal(t, tt.expOut, out)

			if tt.expErr == noError {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expErr)
			}
		})
	}
}

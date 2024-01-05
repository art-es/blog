package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/auth/domain/mock"
	"github.com/art-es/blog/internal/auth/dto"
)

func TestRefreshTokenUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		userRepository       = mock.NewMockUserRepository(ctrl)
		accessTokenRefresher = mock.NewMockaccessTokenRefresher(ctrl)
	)

	var (
		ctx             = context.Background()
		userID          = int64(1)
		token           = "dummyAccessToken"
		refreshedToken  = "dummyRefreshedAccessToken"
		in              = &dto.RefreshTokenIn{AccessToken: token}
		now             = time.Now()
		refreshedObject = &domain.AccessTokenObject{ExpirationTime: now.Add(time.Minute), UserID: userID}
		objectFactory   = func() *domain.AccessTokenObject {
			return &domain.AccessTokenObject{ExpirationTime: now, UserID: userID}
		}
	)

	tests := []struct {
		name   string
		setup  func()
		expOut *dto.RefreshTokenOut
		expErr string
	}{
		{
			name: "happy path",
			setup: func() {
				object := objectFactory()

				accessTokenRefresher.EXPECT().
					Parse(gomock.Eq(token)).
					Return(object, nil)

				userRepository.EXPECT().
					Exists(gomock.Eq(ctx), gomock.Eq(userID)).
					Return(true, nil)

				accessTokenRefresher.EXPECT().
					Refresh(gomock.Eq(object)).
					DoAndReturn(func(o *domain.AccessTokenObject) {
						o.ExpirationTime = refreshedObject.ExpirationTime
					})

				accessTokenRefresher.EXPECT().
					Sign(gomock.Eq(refreshedObject)).
					Return(refreshedToken, nil)
			},
			expOut: &dto.RefreshTokenOut{AccessToken: refreshedToken},
			expErr: noError,
		},
		{
			name: "error on parsing access token",
			setup: func() {
				accessTokenRefresher.EXPECT().
					Parse(gomock.Eq(token)).
					Return(nil, errors.New("dummy error"))
			},
			expErr: "access token parsing error: dummy error",
		},
		{
			name: "error on getting auth",
			setup: func() {
				object := objectFactory()

				accessTokenRefresher.EXPECT().
					Parse(gomock.Eq(token)).
					Return(object, nil)

				userRepository.EXPECT().
					Exists(gomock.Eq(ctx), gomock.Eq(userID)).
					Return(false, errors.New("dummy error"))
			},
			expErr: "user checking existence in repository error: dummy error",
		},
		{
			name: "auth not found",
			setup: func() {
				object := objectFactory()

				accessTokenRefresher.EXPECT().
					Parse(gomock.Eq(token)).
					Return(object, nil)

				userRepository.EXPECT().
					Exists(gomock.Eq(ctx), gomock.Eq(userID)).
					Return(false, nil)
			},
			expErr: dto.ErrUserNotFound.Error(),
		},
		{
			name: "error on refreshing access token",
			setup: func() {
				object := objectFactory()

				accessTokenRefresher.EXPECT().
					Parse(gomock.Eq(token)).
					Return(object, nil)

				userRepository.EXPECT().
					Exists(gomock.Eq(ctx), gomock.Eq(userID)).
					Return(true, nil)

				accessTokenRefresher.EXPECT().
					Refresh(gomock.Eq(object)).
					DoAndReturn(func(o *domain.AccessTokenObject) {
						o.ExpirationTime = refreshedObject.ExpirationTime
					})

				accessTokenRefresher.EXPECT().
					Sign(gomock.Eq(refreshedObject)).
					Return("", errors.New("dummy error"))
			},
			expErr: "signing access token object error: dummy error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			u := domain.NewRefreshTokenUsecase(userRepository, accessTokenRefresher)
			out, err := u.Do(ctx, in)

			assert.Equal(t, tt.expOut, out)

			if tt.expErr == noError {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expErr)
			}
		})
	}

}

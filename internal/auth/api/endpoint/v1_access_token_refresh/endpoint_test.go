package v1_access_token_refresh

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/api/endpoint/v1_access_token_refresh/mock"
	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
	mock_api "github.com/art-es/blog/internal/common/api/mock"
)

func TestEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		accessTokenRefreshCase    = mock.NewMockaccessTokenRefreshCase(ctrl)
		serverErrorHandlerFactory = mock_api.NewMockServerErrorHandlerFactory(ctrl)
		serverErrorHandler        = mock_api.NewMockServerErrorHandler(ctrl)

		smellyAccessToken = "smelly access token"
		freshAccessToken  = "fresh access token"
		noError           = (error)(nil)
		dummyError        = errors.New("dummy error")

		expectedAccessTokenRefreshIn = &dto.AccessTokenRefreshIn{
			AccessToken: smellyAccessToken,
		}
		validAccessTokenRefreshOut = &dto.AccessTokenRefreshOut{
			AccessToken: freshAccessToken,
		}
		noAccessTokenRefreshOut = (*dto.AccessTokenRefreshOut)(nil)
	)

	serverErrorHandlerFactory.EXPECT().
		MakeHandler(gomock.Eq(method), gomock.Eq(path)).
		Return(serverErrorHandler).
		AnyTimes()

	tests := []struct {
		name    string
		setup   func()
		expCode int
		expBody string
	}{
		{
			name: "OK",
			setup: func() {
				accessTokenRefreshCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedAccessTokenRefreshIn)).
					Return(validAccessTokenRefreshOut, noError)
			},
			expCode: 200,
			expBody: `{"accessToken":"fresh access token"}`,
		},
		{
			name: "Unauthorized: invalid access token",
			setup: func() {
				accessTokenRefreshCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedAccessTokenRefreshIn)).
					Return(noAccessTokenRefreshOut, dto.ErrInvalidAccessToken)
			},
			expCode: 401,
			expBody: `{"message":"Please try to sign in again."}`,
		},
		{
			name: "Internal server error: unexpected error in use case",
			setup: func() {
				accessTokenRefreshCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedAccessTokenRefreshIn)).
					Return(noAccessTokenRefreshOut, dummyError)

				serverErrorHandler.EXPECT().
					Handle(gomock.Any(), gomock.Eq(dummyError)).
					DoAndReturn(func(ctx *gin.Context, _ error) {
						api.InternalServerErrorResponse(ctx)
					})
			},
			expCode: 500,
			expBody: `{"message":"Something went wrong, please try again later."}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, path, nil)
			r.Header.Set("X-Access-Token", "smelly access token")

			router := gin.New()
			Bind(router, accessTokenRefreshCase, serverErrorHandlerFactory)

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expCode, w.Code)
			assert.JSONEq(t, tt.expBody, w.Body.String())
		})
	}
}

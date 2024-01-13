package v1_token_refresh

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/api/handler/v1_token_refresh/mock"
	"github.com/art-es/blog/internal/auth/dto"
	api_mock "github.com/art-es/blog/internal/common/api/mock"
	"github.com/art-es/blog/internal/common/testutil"
)

func TestHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		usecase            = mock.NewMockusecase(ctrl)
		serverErrorHandler = api_mock.NewMockServerErrorHandler(ctrl)
		handler            = New(usecase, serverErrorHandler)
		router             = testutil.NewGinRouter(handler)
	)

	var (
		token             = "dummyAccessToken"
		refreshedToken    = "dummyRefreshedAccessToken"
		dummyResponseBody = `{"message":"dummy response"}`
		dummyError        = errors.New("dummy error")
		noError           = (error)(nil)
	)

	tests := []struct {
		name    string
		setup   func(r *http.Request)
		expCode int
		expBody string
	}{
		{
			name: "happy path",
			setup: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer dummyAccessToken")

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(&dto.RefreshTokenIn{AccessToken: token})).
					Return(&dto.RefreshTokenOut{AccessToken: refreshedToken}, noError)
			},
			expCode: 200,
			expBody: `{"ok":true,"accessToken":"dummyRefreshedAccessToken"}`,
		},
		{
			name: "token not specified",
			setup: func(r *http.Request) {
				r.Header.Del("Authorization")
			},
			expCode: 401,
			expBody: `{"ok":false}`,
		},
		{
			name: "invalid token",
			setup: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer dummyAccessToken")

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(&dto.RefreshTokenIn{AccessToken: token})).
					Return(nil, dto.ErrInvalidAccessToken)
			},
			expCode: 401,
			expBody: `{"ok":false}`,
		},
		{
			name: "server error",
			setup: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer dummyAccessToken")

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(&dto.RefreshTokenIn{AccessToken: token})).
					Return(nil, dummyError)

				serverErrorHandler.EXPECT().
					Handle(gomock.Eq(handler.Endpoint()), gomock.Any(), gomock.Eq(dummyError)).
					DoAndReturn(func(_ string, w http.ResponseWriter, _ error) {
						w.WriteHeader(500)
						w.Write([]byte(dummyResponseBody))
					})
			},
			expCode: 500,
			expBody: dummyResponseBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(handler.Method(), handler.Endpoint(), nil)

			if tt.setup != nil {
				tt.setup(r)
			}

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expCode, w.Code)
			assert.Equal(t, tt.expBody, w.Body.String())
		})
	}
}

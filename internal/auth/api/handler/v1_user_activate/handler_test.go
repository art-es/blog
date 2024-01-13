package v1_user_activate

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/api/handler/v1_user_activate/mock"
	"github.com/art-es/blog/internal/auth/dto"
	api_mock "github.com/art-es/blog/internal/common/api/mock"
	"github.com/art-es/blog/internal/common/testutil"
)

func TestHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		usecase            = mock.NewMockusecase(ctrl)
		validator          = api_mock.NewMockValidator(ctrl)
		serverErrorHandler = api_mock.NewMockServerErrorHandler(ctrl)
		handler            = New(usecase, validator, serverErrorHandler)
		router             = testutil.NewGinRouter(handler)
	)

	var (
		code              = "dummyCode"
		validationRule    = "required,uuid"
		dummyResponseBody = `{"message":"dummy response"}`
		dummyError        = errors.New("dummy error")
		noError           = (error)(nil)
	)

	tests := []struct {
		name    string
		path    string
		setup   func()
		expCode int
		expBody string
	}{
		{
			name: "happy path",
			path: "/v1/auth/user/activate/dummyCode",
			setup: func() {
				validator.EXPECT().
					Var(gomock.Eq(code), validationRule).
					Return(noError)

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(&dto.ActivateIn{Code: code})).
					Return(noError)
			},
			expCode: 200,
			expBody: `{"ok":true}`,
		},
		{
			name:    "empty code",
			path:    "/v1/auth/user/activate/",
			expCode: 404,
			expBody: "404 page not found",
		},
		{
			name: "validation error",
			path: "/v1/auth/user/activate/dummyCode",
			setup: func() {
				validator.EXPECT().
					Var(gomock.Eq(code), validationRule).
					Return(dummyError)
			},
			expCode: 400,
			expBody: `{"ok":false,"message":"dummy error"}`,
		},
		{
			name: "activation code not found",
			path: "/v1/auth/user/activate/dummyCode",
			setup: func() {
				validator.EXPECT().
					Var(gomock.Eq(code), validationRule).
					Return(noError)

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(&dto.ActivateIn{Code: code})).
					Return(dto.ErrActivationCodeNotFound)
			},
			expCode: 404,
			expBody: `{"ok":false}`,
		},
		{
			name: "user not found",
			path: "/v1/auth/user/activate/dummyCode",
			setup: func() {
				validator.EXPECT().
					Var(gomock.Eq(code), validationRule).
					Return(noError)

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(&dto.ActivateIn{Code: code})).
					Return(dto.ErrUserNotFound)
			},
			expCode: 404,
			expBody: `{"ok":false}`,
		},
		{
			name: "server error",
			path: "/v1/auth/user/activate/dummyCode",
			setup: func() {
				validator.EXPECT().
					Var(gomock.Eq(code), validationRule).
					Return(noError)

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(&dto.ActivateIn{Code: code})).
					Return(dummyError)

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
			if tt.setup != nil {
				tt.setup()
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(handler.Method(), tt.path, nil)
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expCode, w.Code)
			assert.Equal(t, tt.expBody, w.Body.String())
		})
	}
}

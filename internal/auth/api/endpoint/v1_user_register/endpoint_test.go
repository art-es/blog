package v1_user_register

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/api/endpoint/v1_user_register/mock"
	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
	mock_api "github.com/art-es/blog/internal/common/api/mock"
)

func TestEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		userRegisterCase          = mock.NewMockuserRegisterCase(ctrl)
		validator                 = mock_api.NewMockValidator(ctrl)
		serverErrorHandlerFactory = mock_api.NewMockServerErrorHandlerFactory(ctrl)
		serverErrorHandler        = mock_api.NewMockServerErrorHandler(ctrl)

		name       = "Ivan Ivanov"
		email      = "i.ivanov@example.com"
		password   = "Qwerty123!"
		noError    = (error)(nil)
		dummyError = errors.New("dummy error")

		expectedRequestInValidator = &request{
			Name:     name,
			Email:    email,
			Password: password,
		}
		expectedUserRegisterIn = &dto.UserRegisterIn{
			Name:     name,
			Email:    email,
			Password: password,
		}
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
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userRegisterCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserRegisterIn)).
					Return(noError)
			},
			expCode: 200,
			expBody: `{"message":"Please check your email to activate your account."}`,
		},
		{
			name: "Bad request: request validation failed",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(dummyError)
			},
			expCode: 400,
			expBody: `{"error":{"code":1001,"name":"Request validation failed"},"message":"dummy error"}`,
		},
		{
			name: "Bad request: busy email",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userRegisterCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserRegisterIn)).
					Return(dto.ErrEmailIsBusy)
			},
			expCode: 400,
			expBody: `{"error":{"code":2001,"name":"Busy email"},"message":"User with this email already exists."}`,
		},
		{
			name: "Internal server error: unexpected error in use case",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userRegisterCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserRegisterIn)).
					Return(dummyError)

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

			rBody := `{"name":"Ivan Ivanov","email":"i.ivanov@example.com","password":"Qwerty123!"}`
			r := httptest.NewRequest(method, path, io.NopCloser(bytes.NewBufferString(rBody)))
			w := httptest.NewRecorder()

			router := gin.New()
			Bind(router, userRegisterCase, validator, serverErrorHandlerFactory)

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expCode, w.Code)
			assert.JSONEq(t, tt.expBody, w.Body.String())
		})
	}
}

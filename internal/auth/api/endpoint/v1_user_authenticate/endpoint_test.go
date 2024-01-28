package v1_user_authenticate

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/api/endpoint/v1_user_authenticate/mock"
	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
	mock_api "github.com/art-es/blog/internal/common/api/mock"
)

func TestEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		userAuthenticateCase      = mock.NewMockuserAuthenticateCase(ctrl)
		validator                 = mock_api.NewMockValidator(ctrl)
		serverErrorHandlerFactory = mock_api.NewMockServerErrorHandlerFactory(ctrl)
		serverErrorHandler        = mock_api.NewMockServerErrorHandler(ctrl)

		email       = "i.ivanov@example.com"
		password    = "Qwerty123!"
		accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		noError     = (error)(nil)
		dummyError  = errors.New("dummy error")

		expectedRequestInValidator = &request{
			Email:    email,
			Password: password,
		}
		expectedUserAuthenticateIn = &dto.UserAuthenticateIn{
			Email:    email,
			Password: password,
		}
		validUserAuthenticateOut = &dto.UserAuthenticateOut{
			AccessToken: accessToken,
		}
		noUserAuthenticateOut = (*dto.UserAuthenticateOut)(nil)
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

				userAuthenticateCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserAuthenticateIn)).
					Return(validUserAuthenticateOut, noError)
			},
			expCode: 200,
			expBody: `{"accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}`,
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
			name: "Bad request: user not found",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userAuthenticateCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserAuthenticateIn)).
					Return(noUserAuthenticateOut, dto.ErrUserNotFound)
			},
			expCode: 400,
			expBody: `{"error":{"code":2003,"name":"Incorrect user credentials"},"message":"Email or password is incorrect."}`,
		},
		{
			name: "Bad request: incorrect password",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userAuthenticateCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserAuthenticateIn)).
					Return(noUserAuthenticateOut, dto.ErrIncorrectPassword)
			},
			expCode: 400,
			expBody: `{"error":{"code":2003,"name":"Incorrect user credentials"},"message":"Email or password is incorrect."}`,
		},
		{
			name: "Internal server error: unexpected error in use case",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userAuthenticateCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserAuthenticateIn)).
					Return(noUserAuthenticateOut, dummyError)

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

			rBody := `{"email":"i.ivanov@example.com","password":"Qwerty123!"}`
			r := httptest.NewRequest(method, path, io.NopCloser(bytes.NewBufferString(rBody)))
			w := httptest.NewRecorder()

			router := gin.New()
			Bind(router, userAuthenticateCase, validator, serverErrorHandlerFactory)

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expCode, w.Code)
			assert.JSONEq(t, tt.expBody, w.Body.String())
		})
	}
}

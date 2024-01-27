package v1_user_activate

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/art-es/blog/internal/common/api"

	"github.com/art-es/blog/internal/auth/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/api/endpoint/v1_user_activate/mock"
	mock_api "github.com/art-es/blog/internal/common/api/mock"
)

func TestEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		userActivateCase          = mock.NewMockuserActivateCase(ctrl)
		validator                 = mock_api.NewMockValidator(ctrl)
		serverErrorHandlerFactory = mock_api.NewMockServerErrorHandlerFactory(ctrl)
		serverErrorHandler        = mock_api.NewMockServerErrorHandler(ctrl)

		code       = "1eb62291-9374-4887-8c1a-96382e54fcad"
		noError    = (error)(nil)
		dummyError = errors.New("dummy error")

		expectedRequestInValidator = &request{
			Code: code,
		}
		expectedUserActivateIn = &dto.UserActivateIn{
			Code: code,
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

				userActivateCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserActivateIn)).
					Return(noError)
			},
			expCode: 200,
			expBody: `{"message":"Your account have been activated."}`,
		},
		{
			name: "Bad request: request validation failed",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Any()).
					Return(dummyError)
			},
			expCode: 400,
			expBody: `{"error":{"code":1001,"name":"Request validation failed"},"message":"dummy error"}`,
		},
		{
			name: "Bad request: expired user activation code",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userActivateCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserActivateIn)).
					Return(dto.ErrExpiredUserActivationCode)
			},
			expCode: 400,
			expBody: `{"error":{"code":2002,"name":"Expired user activation code"},"message":"Activate code is expired. Please pass the registration process again."}`,
		},
		{
			name: "Not found: user activation code not found",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userActivateCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserActivateIn)).
					Return(dto.ErrUserActivationCodeNotFound)
			},
			expCode: 404,
			expBody: `{"message":"Activation code not found."}`,
		},
		{
			name: "Not found: user not found",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userActivateCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserActivateIn)).
					Return(dto.ErrUserNotFound)
			},
			expCode: 404,
			expBody: `{"message":"Activation code not found."}`,
		},
		{
			name: "Internal server error: unexpected error in use case",
			setup: func() {
				validator.EXPECT().
					Struct(gomock.Eq(expectedRequestInValidator)).
					Return(noError)

				userActivateCase.EXPECT().
					Use(gomock.Any(), gomock.Eq(expectedUserActivateIn)).
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

			rBody := `{"code":"1eb62291-9374-4887-8c1a-96382e54fcad"}`
			r := httptest.NewRequest(method, path, io.NopCloser(bytes.NewBufferString(rBody)))
			w := httptest.NewRecorder()

			router := gin.New()
			Bind(router, userActivateCase, validator, serverErrorHandlerFactory)

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expCode, w.Code)
			assert.JSONEq(t, tt.expBody, w.Body.String())
		})
	}
}

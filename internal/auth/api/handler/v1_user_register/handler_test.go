package v1_user_register

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/api/handler/v1_user_register/mock"
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
		name              = "dummy name"
		email             = "dummy@example.com"
		password          = "1234Qwerty"
		dummyResponseBody = `{"message":"dummy response"}`
		dummyError        = errors.New("dummy error")
		noError           = (error)(nil)

		requestBody = func() io.ReadCloser {
			return testutil.ReadCloserFromJSON(map[string]any{
				"name":     name,
				"email":    email,
				"password": password,
			})
		}
		validatorRequest = &request{
			Name:     name,
			Email:    email,
			Password: password,
		}
		registerIn = &dto.RegisterIn{
			Name:     name,
			Email:    email,
			Password: password,
		}
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
				r.Body = requestBody()

				validator.EXPECT().
					Struct(gomock.Eq(validatorRequest)).
					Return(noError)

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(registerIn)).
					Return(noError)
			},
			expCode: 200,
			expBody: `{"ok":true}`,
		},
		{
			name: "email is busy",
			setup: func(r *http.Request) {
				r.Body = requestBody()

				validator.EXPECT().
					Struct(gomock.Eq(validatorRequest)).
					Return(noError)

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(registerIn)).
					Return(dto.ErrEmailIsBusy)
			},
			expCode: 400,
			expBody: `{"ok":false,"message":"email is busy"}`,
		},
		{
			name: "validation error",
			setup: func(r *http.Request) {
				validator.EXPECT().
					Struct(gomock.Any()).
					Return(dummyError)
			},
			expCode: 400,
			expBody: `{"ok":false,"message":"dummy error"}`,
		},
		{
			name: "server error",
			setup: func(r *http.Request) {
				validator.EXPECT().
					Struct(gomock.Any()).
					Return(noError)

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Any()).
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
			w := httptest.NewRecorder()
			r := httptest.NewRequest(handler.Method(), handler.Endpoint(), nil)

			if tt.setup != nil {
				tt.setup(r)
			}

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expCode, w.Code)
			assert.JSONEq(t, tt.expBody, w.Body.String())
		})
	}
}
